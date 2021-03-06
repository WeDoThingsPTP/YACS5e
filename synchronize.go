package main

import (
	"database/sql"
	"log"
	pb "github.com/WeDoThingsPTP/YACS5eproto"
	"google.golang.org/grpc/status"
)

// Synchronize...
func (server *YACS5eServer) Synchronize(stream pb.YACS5E_SynchronizeServer) error {
	var (
		user *pb.TUser
	)

	log.Println("Start synchronize task")

	// Check recaptcha

	streamIn, err := stream.Recv()
	if err != nil {
		LogUnknownError(err)
		return status.Errorf(54, "Error getting data from input stream")
	}

	// 1. Check credentials

	user, err = partialLogin(streamIn)
	if err != nil {
		log.Println("Synchronize: Error logging user:", err)
		logErr := stream.Send(&pb.TTalk{Union: &pb.TTalk_Good{Good: false}})
		if logErr != nil {
			ErrorStatus(logErr)
		}

		return err
	}

	log.Println("Synchronize for user", user.Login)

	err = stream.Send(&pb.TTalk{Union: &pb.TTalk_Good{Good: true}})
	if err != nil {
		LogUnknownError(err)
		return status.Errorf(55, "Error sending data to input stream")
	}

	// 2a. Create slice of uuids'. If after app-sync there will be any left, app does not have them.
	uuidQuery, err := db.Query("SELECT uuid FROM characters WHERE users_id=(SELECT id FROM users WHERE login=?)", user.Login)
	if err != nil {
		ErrorStatus(err)
	}

	uuidSlice := make([]string, 0)
	for uuidQuery.Next() {
		var uuid string

		err := uuidQuery.Scan(&uuid)
		if err != nil {
			ErrorStatus(err)
		}

		uuidSlice = append(uuidSlice, uuid)
	}

	// 2b. Perform char sync one-by-one

	log.Println("Synchronize: Perform char sync one-by-one")

	exchangeCharInfo := true
	for exchangeCharInfo {
		log.Println("Loop...")

		// get login, uuid
		streamIn, err := stream.Recv()
		if err != nil {
			return ErrorStatus(err)
		}

		switch ttalk := streamIn.Union.(type) {
		case *pb.TTalk_Character:

			var (
				uuid     string
				lastSync uint64
				lastMod  uint64
				data     string
			)

			// character is to be deleted
			if ttalk.Character.Delete {
				log.Println("Synchronize: (5) Character is to be deleted", uuid)
				_, err := db.Exec("DELETE FROM characters WHERE uuid=?", ttalk.Character.Uuid)
				if err != nil {
					return ErrorStatus(err)
				}
				continue
			}

			err := db.QueryRow(
				"SELECT uuid, last_sync, last_mod, data FROM characters WHERE users_id=(SELECT id FROM users WHERE login=?) AND uuid=? LIMIT 1",
				user.Login,
				ttalk.Character.Uuid,
			).Scan(&uuid, &lastSync, &lastMod, &data)

			if err == sql.ErrNoRows {
				log.Println("Synchronize: (4) character not found on server, ask for complete data")
				// 4 - not on server - receive empty uuid, send complete character
				onCharacterNotFound(stream, *user)
				break
			} else if err != nil {
				return ErrorStatus(err)
			}

			uuidSlice = RemoveFromSlice(uuidSlice, uuid)

			err = stream.Send(&pb.TTalk{Union: &pb.TTalk_Character{Character: &pb.TCharacter{Uuid: uuid, LastSync: lastSync, LastMod: lastMod}}})
			if err != nil {
				return ErrorStatus(err)
			}

			// Character is even
			if lastSync == ttalk.Character.GetLastSync() && lastMod == ttalk.Character.GetLastMod() {
				log.Println("Synchronize: (0) Character is even", uuid)
				continue
			}

			// if not even - app wants to send data
			streamIn, err = stream.Recv()
			if err != nil {
				return ErrorStatus(err)
			}
			switch ttalk := streamIn.Union.(type) {
			case *pb.TTalk_Character:
				tChar := ttalk.Character
				if tChar.GetLastSync() != 0 && tChar.GetLastMod() != 0 && tChar.Uuid != "" && len(tChar.Blob) > 0 {

					log.Println("Synchronize: app wants to insert / update character uuid: " + tChar.Uuid)

					_, err := db.Exec("INSERT INTO characters "+
						"SET uuid=?, users_id=(SELECT id FROM users WHERE login=?), last_sync=?, last_mod=?, data=? "+
						"ON DUPLICATE KEY UPDATE last_sync=?, last_mod=?, data=?",
						tChar.Uuid, user.Login, tChar.LastSync, tChar.LastMod, tChar.Blob, tChar.LastSync, tChar.LastMod, tChar.Blob)

					if err != nil {
						return ErrorStatus(err)
					}
					continue

				} else if tChar.LastMod == 0 && tChar.LastSync == 0 {
					log.Println("Synchronize: app asks for data")
					err := stream.Send(&pb.TTalk{Union: &pb.TTalk_Character{Character: &pb.TCharacter{
						Uuid:     uuid,
						LastSync: lastSync,
						LastMod:  lastMod,
						Blob:     data,
					}}})
					if err != nil {
						return ErrorStatus(err)
					}
					continue
				}
			}

			log.Println("Synchronize: Unimplemented route...")

		case *pb.TTalk_Good:
			log.Println("Synchronize: no more characters on client")
			exchangeCharInfo = false
			continue

		default:
			return status.Errorf(125, "Unexpected type")
		}

	}

	if len(uuidSlice) > 0 {
		log.Println("Synchronize: more characters on db")
		err := stream.Send(&pb.TTalk{Union: &pb.TTalk_Good{Good: true}})
		if err != nil {
			return ErrorStatus(err)
		}
	} else {
		log.Println("Synchronize: no characters on db")
		err := stream.Send(&pb.TTalk{Union: &pb.TTalk_Good{Good: false}})
		if err != nil {
			return ErrorStatus(err)
		}
	}

	for _, uuid := range uuidSlice {
		var (
			lastSync uint64
			lastMod  uint64
			data     string
		)

		err := db.QueryRow(
			"SELECT last_sync, last_mod, data FROM characters WHERE users_id=(SELECT id FROM users WHERE login=?) AND uuid=? LIMIT 1",
			user.Login,
			uuid,
		).Scan(&lastSync, &lastMod, &data)
		if err != nil {
			return ErrorStatus(err)
		}

		err = stream.Send(&pb.TTalk{Union: &pb.TTalk_Character{Character: &pb.TCharacter{
			Uuid:     uuid,
			LastSync: lastSync,
			LastMod:  lastMod,
			Blob:     data,
		}}})
		if err != nil {
			return ErrorStatus(err)
		}
	}

	err = stream.Send(&pb.TTalk{Union: &pb.TTalk_Good{Good: true}})
	if err != nil {
		return ErrorStatus(err)
	}
	log.Println("Synchronize: Complete")

	return status.Errorf(0, "Complete")
}

// 4 - not on server - receive empty uuid, send complete character
func onCharacterNotFound(stream pb.YACS5E_SynchronizeServer, user pb.TUser) error {
	err := stream.Send(&pb.TTalk{Union: &pb.TTalk_Character{Character: &pb.TCharacter{Uuid: ""}}})
	if err != nil {
		return ErrorStatus(err)
	}

	streamIn, err := stream.Recv()
	if err != nil {
		return ErrorStatus(err)
	}

	switch tCharacter := streamIn.Union.(type) {
	case *pb.TTalk_Character:
		char := tCharacter.Character
		_, err := db.Exec(
			"INSERT INTO characters "+
				"SET uuid=?, users_id=(SELECT id FROM users WHERE login=?), last_sync=?, last_mod=?, data=?",
			char.Uuid,
			user.Login,
			char.LastSync,
			char.LastMod,
			char.Blob,
		)
		if err == sql.ErrNoRows {
			log.Println("Synchronize: internar error:", err)
			return status.Errorf(2, "Unexpected error")
		} else if err != nil {
			return ErrorStatus(err)
		}
		log.Println("Synchronize: New character uuid:", char.Uuid)
	}

	return nil
}
