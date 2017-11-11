package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	pb "github.com/krzysztofromanowski94/YACS5e-cloud/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
	"log"
	"runtime"
	"strconv"
	"strings"
)

var (
	sqlDBName                = "{{ sql_dbname }}"
	sqlHostname              = "{{ sql_hostname }}"
	sqlPassword              = "{{ sql_password }}"
	sqlPort                  = "{{ sql_port }}"
	sqlUser                  = "{{ sql_user }}"
	sqlMaxOpenConnections, _ = strconv.ParseInt("{{ sql_max_open_connections }}", 10, 64)
	db                       *sql.DB
	dataSourceName           = sqlUser + ":" + sqlPassword +
		"@tcp(" + sqlHostname + ":" + sqlPort + ")/" + sqlDBName
)

// rpc Registration (User) returns (Empty)
// ERROR CODES:
// 100: UNKNOWN ERROR
// 101: INVALID LOGIN
// 102: INVALID PASSWORD
// 103: USER EXISTS
func (server *YACS5eServer) Registration(ctx context.Context, user *pb.TUser) (*pb.Empty, error) {

	// Here should be checking if recaptcha is right

	_, err := db.Exec("INSERT INTO users VALUES (NULL, ?, ?, ?)", user.Login, user.Password, user.VisibleName)
	if err != nil {
		switch strErr := err.Error(); {

		case strings.Contains(strErr, "Error 1062"):
			returnStr := fmt.Sprint("User ", user.Login, " exists.")
			return &pb.Empty{}, status.Errorf(103, returnStr)

		default:
			LogUnknownError(err)
			return &pb.Empty{}, status.Errorf(100, "Unknown error: ", err)
		}
	}

	returnStr := fmt.Sprint("Registered user: ", user.Login)
	return &pb.Empty{}, status.Errorf(0, returnStr)
}

// rpc Login (User) returns (Empty)
// ERROR CODES:
// 110: UNKNOWN ERROR
// 111: INVALID CREDENTIALS
// 112: USER NOT FOUND
func (server *YACS5eServer) Login(ctx context.Context, user *pb.TUser) (*pb.Empty, error) {

	// Here should be checking if recaptcha is right

	row, err := db.Query("SELECT login, visible_name FROM users WHERE login=? AND password=? LIMIT 1", user.Login, user.Password)

	if err != nil {
		LogUnknownError(err)
		returnStr := fmt.Sprint("UNKNOWN ERROR: ", err)
		return &pb.Empty{}, status.Errorf(110, returnStr)
	}

	for row.Next() {
		var (
			login       string
			visibleName string
		)

		err := row.Scan(&login, &visibleName)
		if err != nil {
			LogUnknownError(err)
			returnStr := fmt.Sprint("UNKNOWN ERROR: ", err)
			return &pb.Empty{}, status.Errorf(110, returnStr)
		}

		return &pb.Empty{}, status.Errorf(0, "User found")
	}

	returnStr := fmt.Sprint("User ", user.Login, " not found")
	return &pb.Empty{}, status.Errorf(112, returnStr)
}

// ERROR CODES:
// 120: UNKNOWN ERROR
// 54: ERROR GETTING DATA FROM STREAM
// 124: ERROR SENDING DATA TO STREAM
// 125: INCORRECT FLOW
// 126: USER DON'T HAVE THIS CHARACTER
func (server *YACS5eServer) Synchronize(stream pb.YACS5E_SynchronizeServer) error {

	var (
		user *pb.TUser
	)

	// Check recaptcha

	//streamIn, err := stream.Recv()
	//if err != nil {
	//	LogUnknownError(err)
	//	returnStr := fmt.Sprint("ERROR GETTING DATA FROM INPUT STREAM: ", err)
	//	return status.Errorf(54, returnStr)
	//}

	// 1. Check credentials

	streamIn := &pb.TTalk{&pb.TTalk_User{&pb.TUser{"testUser", "testPass", "", "", 0}}}

	user, err := partialLogin(streamIn)
	if err != nil {
		return err
	}

	//err = stream.Send(&pb.TTalk{&pb.TTalk_Good{true}})
	//if err != nil {
	//	LogUnknownError(err)
	//	returnStr := fmt.Sprint("ERROR SENDING DATA FROM INPUT STREAM: ", err)
	//	return status.Errorf(55, returnStr)
	//}

	// 2. Get characters timestamp from client

	//var (
	//	clientTimestampList = make([]*pb.TCharacter, 0)
	//)

	//gettingTimestamps := true
	//for gettingTimestamps{
	//
	//	streamIn, err := stream.Recv()
	//	if err != nil {
	//		LogUnknownError(err)
	//	}
	//
	//	switch ttalk := streamIn.(type) {
	//
	//	case *pb.TTalk_Character:
	//		if ttalk.Character.Timestamp != 0 {
	//			clientTimestampList = append(clientTimestampList, ttalk.Character)
	//		} else {
	//			gettingTimestamps = false
	//			break
	//	}
	//
	//	default:
	//		return status.Errorf(125, "Expected type TTalk_Character")
	//	}
	//}

	// 3a. Get timestamps from database

	var (
		timestampsDB []*pb.TCharacter
	)

	log.Println(user)
	timestampsQuery, err := db.Query("SELECT timestamp, uuid, data FROM characters WHERE users_id=?", user.Id)
	if err != nil {
		switch strErr := err.Error(); {
		case strings.Contains(strErr, "Error 1062"):
			// User don't have any characters in database
			break
		default:
			LogUnknownError(err)
			returnStr := fmt.Sprint("UNKNOWN ERROR:", err)
			return status.Errorf(120, returnStr)
		}
	}

	for timestampsQuery.Next() {
		var (
			timestamp uint64
			uuid      []byte
			blob      []byte
		)
		err = timestampsQuery.Scan(
			&timestamp,
			&uuid,
			&blob,
		)
		if err != nil {
			LogUnknownError(err)
			returnStr := fmt.Sprint("UNKNOWN ERROR:", err)
			return status.Errorf(120, returnStr)
		}
		timestampsDB = append(timestampsDB, &pb.TCharacter{Uuid: uuid, Timestamp: timestamp, Blob: blob})
	}

	// 3b. Decide what characters need to be updated on server

	log.Println(timestampsDB)

	return status.Errorf(0, "")
}

func init() {
	log.Println("Connecting to: ", dataSourceName)

	var err error

	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal("Preparing the database connection caused ERROR: ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Connecting to database caused ERROR: ", err)
	}

	db.SetMaxOpenConns(int(sqlMaxOpenConnections))

	log.Println("Connection estabilished")

	testServer := newServer()

	log.Println(testServer.Synchronize(nil))

}

func partialLogin(tTalk *pb.TTalk) (user *pb.TUser, err error) {
	switch tCharacterUnion := tTalk.Union.(type) {

	case *pb.TTalk_User:
		err := db.QueryRow(
			"SELECT id FROM users WHERE login=? AND password=? LIMIT 1",
			tCharacterUnion.User.Login,
			tCharacterUnion.User.Password,
		).Scan(&tCharacterUnion.User.Id)

		switch err {
		case sql.ErrNoRows:
			// User does not exists
			return nil, status.Errorf(52, "INVALID CREDENTIALS")

		case nil:
			// User exists
			return tCharacterUnion.User, nil
			break

		default:
			LogUnknownError(err)
			returnStr := fmt.Sprint("UNKNOWN ERROR: ", err)
			return nil, status.Errorf(51, returnStr)
		}

	default:
		return nil, status.Errorf(53, "EXPECTED TYPE IS TTalk_User")
	}

	return nil, status.Errorf(2, "UNEXPECTED RETURN AT PARTIAL LOGIN")
}

func LogUnknownError(err error) {
	if err != nil {
		pc, fn, line, _ := runtime.Caller(1)
		log.Printf("ERROR in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), fn, line, err)
	}
}

type YACS5eServer struct {
}

//type CharacterList struct {
//	Characters []CharacterInfo `json:"characters"`
//}

//type CharacterInfo struct {
//	ID          int    `json:"id"`
//	Name        string `json:"name"`
//	Description string `json:"description"`
//}

func newServer() *YACS5eServer {
	return new(YACS5eServer)
}
