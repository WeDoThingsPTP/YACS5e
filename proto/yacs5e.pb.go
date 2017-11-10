// Code generated by protoc-gen-go. DO NOT EDIT.
// source: yacs5e.proto

/*
Package yacs5e is a generated protocol buffer package.

It is generated from these files:
	yacs5e.proto

It has these top-level messages:
	TUser
	TCharacter
	Empty
*/
package yacs5e

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type TUser struct {
	Login       string `protobuf:"bytes,1,opt,name=login" json:"login,omitempty"`
	Password    string `protobuf:"bytes,2,opt,name=password" json:"password,omitempty"`
	RespToken   string `protobuf:"bytes,3,opt,name=respToken" json:"respToken,omitempty"`
	VisibleName string `protobuf:"bytes,4,opt,name=visibleName" json:"visibleName,omitempty"`
}

func (m *TUser) Reset()                    { *m = TUser{} }
func (m *TUser) String() string            { return proto.CompactTextString(m) }
func (*TUser) ProtoMessage()               {}
func (*TUser) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *TUser) GetLogin() string {
	if m != nil {
		return m.Login
	}
	return ""
}

func (m *TUser) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *TUser) GetRespToken() string {
	if m != nil {
		return m.RespToken
	}
	return ""
}

func (m *TUser) GetVisibleName() string {
	if m != nil {
		return m.VisibleName
	}
	return ""
}

type TCharacter struct {
	// Types that are valid to be assigned to Union:
	//	*TCharacter_User
	//	*TCharacter_Blob
	//	*TCharacter_Id
	Union isTCharacter_Union `protobuf_oneof:"union"`
}

func (m *TCharacter) Reset()                    { *m = TCharacter{} }
func (m *TCharacter) String() string            { return proto.CompactTextString(m) }
func (*TCharacter) ProtoMessage()               {}
func (*TCharacter) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type isTCharacter_Union interface {
	isTCharacter_Union()
}

type TCharacter_User struct {
	User *TUser `protobuf:"bytes,5,opt,name=user,oneof"`
}
type TCharacter_Blob struct {
	Blob []byte `protobuf:"bytes,6,opt,name=blob,proto3,oneof"`
}
type TCharacter_Id struct {
	Id int32 `protobuf:"varint,7,opt,name=id,oneof"`
}

func (*TCharacter_User) isTCharacter_Union() {}
func (*TCharacter_Blob) isTCharacter_Union() {}
func (*TCharacter_Id) isTCharacter_Union()   {}

func (m *TCharacter) GetUnion() isTCharacter_Union {
	if m != nil {
		return m.Union
	}
	return nil
}

func (m *TCharacter) GetUser() *TUser {
	if x, ok := m.GetUnion().(*TCharacter_User); ok {
		return x.User
	}
	return nil
}

func (m *TCharacter) GetBlob() []byte {
	if x, ok := m.GetUnion().(*TCharacter_Blob); ok {
		return x.Blob
	}
	return nil
}

func (m *TCharacter) GetId() int32 {
	if x, ok := m.GetUnion().(*TCharacter_Id); ok {
		return x.Id
	}
	return 0
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*TCharacter) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _TCharacter_OneofMarshaler, _TCharacter_OneofUnmarshaler, _TCharacter_OneofSizer, []interface{}{
		(*TCharacter_User)(nil),
		(*TCharacter_Blob)(nil),
		(*TCharacter_Id)(nil),
	}
}

func _TCharacter_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*TCharacter)
	// union
	switch x := m.Union.(type) {
	case *TCharacter_User:
		b.EncodeVarint(5<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.User); err != nil {
			return err
		}
	case *TCharacter_Blob:
		b.EncodeVarint(6<<3 | proto.WireBytes)
		b.EncodeRawBytes(x.Blob)
	case *TCharacter_Id:
		b.EncodeVarint(7<<3 | proto.WireVarint)
		b.EncodeVarint(uint64(x.Id))
	case nil:
	default:
		return fmt.Errorf("TCharacter.Union has unexpected type %T", x)
	}
	return nil
}

func _TCharacter_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*TCharacter)
	switch tag {
	case 5: // union.user
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(TUser)
		err := b.DecodeMessage(msg)
		m.Union = &TCharacter_User{msg}
		return true, err
	case 6: // union.blob
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeRawBytes(true)
		m.Union = &TCharacter_Blob{x}
		return true, err
	case 7: // union.id
		if wire != proto.WireVarint {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeVarint()
		m.Union = &TCharacter_Id{int32(x)}
		return true, err
	default:
		return false, nil
	}
}

func _TCharacter_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*TCharacter)
	// union
	switch x := m.Union.(type) {
	case *TCharacter_User:
		s := proto.Size(x.User)
		n += proto.SizeVarint(5<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *TCharacter_Blob:
		n += proto.SizeVarint(6<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(len(x.Blob)))
		n += len(x.Blob)
	case *TCharacter_Id:
		n += proto.SizeVarint(7<<3 | proto.WireVarint)
		n += proto.SizeVarint(uint64(x.Id))
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type Empty struct {
}

func (m *Empty) Reset()                    { *m = Empty{} }
func (m *Empty) String() string            { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()               {}
func (*Empty) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func init() {
	proto.RegisterType((*TUser)(nil), "TUser")
	proto.RegisterType((*TCharacter)(nil), "TCharacter")
	proto.RegisterType((*Empty)(nil), "Empty")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for YACS5E service

type YACS5EClient interface {
	// ERROR CODES:
	// 100: UNKNOWN ERROR
	// 101: INVALID LOGIN
	// 102: INVALID PASSWORD
	// 103: USER EXISTS
	Registration(ctx context.Context, in *TUser, opts ...grpc.CallOption) (*Empty, error)
	// ERROR CODES:
	// 110: UNKNOWN ERROR
	// 111: INVALID CREDENTIALS
	Login(ctx context.Context, in *TUser, opts ...grpc.CallOption) (*Empty, error)
	// ERROR CODES:
	// 120: UNKNOWN ERROR
	// 121: INVALID CREDENTIALS
	// 122: CHARACTER DOES NOT EXISTS
	GetCharacter(ctx context.Context, opts ...grpc.CallOption) (YACS5E_GetCharacterClient, error)
}

type yACS5EClient struct {
	cc *grpc.ClientConn
}

func NewYACS5EClient(cc *grpc.ClientConn) YACS5EClient {
	return &yACS5EClient{cc}
}

func (c *yACS5EClient) Registration(ctx context.Context, in *TUser, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/YACS5e/Registration", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *yACS5EClient) Login(ctx context.Context, in *TUser, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/YACS5e/Login", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *yACS5EClient) GetCharacter(ctx context.Context, opts ...grpc.CallOption) (YACS5E_GetCharacterClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_YACS5E_serviceDesc.Streams[0], c.cc, "/YACS5e/GetCharacter", opts...)
	if err != nil {
		return nil, err
	}
	x := &yACS5EGetCharacterClient{stream}
	return x, nil
}

type YACS5E_GetCharacterClient interface {
	Send(*TCharacter) error
	Recv() (*TCharacter, error)
	grpc.ClientStream
}

type yACS5EGetCharacterClient struct {
	grpc.ClientStream
}

func (x *yACS5EGetCharacterClient) Send(m *TCharacter) error {
	return x.ClientStream.SendMsg(m)
}

func (x *yACS5EGetCharacterClient) Recv() (*TCharacter, error) {
	m := new(TCharacter)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for YACS5E service

type YACS5EServer interface {
	// ERROR CODES:
	// 100: UNKNOWN ERROR
	// 101: INVALID LOGIN
	// 102: INVALID PASSWORD
	// 103: USER EXISTS
	Registration(context.Context, *TUser) (*Empty, error)
	// ERROR CODES:
	// 110: UNKNOWN ERROR
	// 111: INVALID CREDENTIALS
	Login(context.Context, *TUser) (*Empty, error)
	// ERROR CODES:
	// 120: UNKNOWN ERROR
	// 121: INVALID CREDENTIALS
	// 122: CHARACTER DOES NOT EXISTS
	GetCharacter(YACS5E_GetCharacterServer) error
}

func RegisterYACS5EServer(s *grpc.Server, srv YACS5EServer) {
	s.RegisterService(&_YACS5E_serviceDesc, srv)
}

func _YACS5E_Registration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TUser)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(YACS5EServer).Registration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/YACS5e/Registration",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(YACS5EServer).Registration(ctx, req.(*TUser))
	}
	return interceptor(ctx, in, info, handler)
}

func _YACS5E_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TUser)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(YACS5EServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/YACS5e/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(YACS5EServer).Login(ctx, req.(*TUser))
	}
	return interceptor(ctx, in, info, handler)
}

func _YACS5E_GetCharacter_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(YACS5EServer).GetCharacter(&yACS5EGetCharacterServer{stream})
}

type YACS5E_GetCharacterServer interface {
	Send(*TCharacter) error
	Recv() (*TCharacter, error)
	grpc.ServerStream
}

type yACS5EGetCharacterServer struct {
	grpc.ServerStream
}

func (x *yACS5EGetCharacterServer) Send(m *TCharacter) error {
	return x.ServerStream.SendMsg(m)
}

func (x *yACS5EGetCharacterServer) Recv() (*TCharacter, error) {
	m := new(TCharacter)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _YACS5E_serviceDesc = grpc.ServiceDesc{
	ServiceName: "YACS5e",
	HandlerType: (*YACS5EServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Registration",
			Handler:    _YACS5E_Registration_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _YACS5E_Login_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetCharacter",
			Handler:       _YACS5E_GetCharacter_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "yacs5e.proto",
}

func init() { proto.RegisterFile("yacs5e.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 266 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x90, 0x41, 0x4b, 0xc3, 0x30,
	0x1c, 0x47, 0x97, 0xda, 0xb4, 0xee, 0xdf, 0x1e, 0x24, 0x0c, 0x0c, 0x65, 0x48, 0xe9, 0xa9, 0x07,
	0x29, 0x32, 0xd9, 0x07, 0xd0, 0x21, 0xee, 0x20, 0x1e, 0x6a, 0x3d, 0x88, 0xa7, 0x74, 0xfd, 0x33,
	0x83, 0x5d, 0x53, 0x92, 0x4c, 0xe9, 0xb7, 0x97, 0x45, 0xb1, 0x63, 0xc7, 0xf7, 0x5e, 0x20, 0xbf,
	0x04, 0xe2, 0x41, 0x6c, 0xcc, 0x12, 0x8b, 0x5e, 0x2b, 0xab, 0xb2, 0x01, 0x68, 0xf5, 0x6a, 0x50,
	0xb3, 0x19, 0xd0, 0x56, 0x6d, 0x65, 0xc7, 0x49, 0x4a, 0xf2, 0x69, 0xf9, 0x0b, 0x2c, 0x81, 0xf3,
	0x5e, 0x18, 0xf3, 0xad, 0x74, 0xc3, 0x3d, 0x17, 0xfe, 0x99, 0xcd, 0x61, 0xaa, 0xd1, 0xf4, 0x95,
	0xfa, 0xc4, 0x8e, 0x9f, 0xb9, 0x38, 0x0a, 0x96, 0x42, 0xf4, 0x25, 0x8d, 0xac, 0x5b, 0x7c, 0x16,
	0x3b, 0xe4, 0xbe, 0xeb, 0xc7, 0x2a, 0x7b, 0x07, 0xa8, 0x56, 0x1f, 0x42, 0x8b, 0x8d, 0x45, 0xcd,
	0xe6, 0xe0, 0xef, 0x0d, 0x6a, 0x4e, 0x53, 0x92, 0x47, 0x8b, 0xa0, 0x70, 0xab, 0xd6, 0x93, 0xd2,
	0x59, 0x36, 0x03, 0xbf, 0x6e, 0x55, 0xcd, 0x83, 0x94, 0xe4, 0xf1, 0xc1, 0x1e, 0x88, 0x5d, 0x80,
	0x27, 0x1b, 0x1e, 0xa6, 0x24, 0xa7, 0xeb, 0x49, 0xe9, 0xc9, 0xe6, 0x3e, 0x04, 0xba, 0xef, 0xa4,
	0xea, 0xb2, 0x10, 0xe8, 0xc3, 0xae, 0xb7, 0xc3, 0x42, 0x41, 0xf0, 0x76, 0xb7, 0x7a, 0x59, 0x22,
	0xbb, 0x82, 0xb8, 0xc4, 0xad, 0x34, 0x56, 0x0b, 0x2b, 0x55, 0xc7, 0xfe, 0xee, 0x48, 0x82, 0xc2,
	0x9d, 0x64, 0x97, 0x40, 0x9f, 0xdc, 0xa3, 0x4f, 0xc3, 0x35, 0xc4, 0x8f, 0x68, 0xc7, 0xa9, 0x51,
	0x31, 0xee, 0x4e, 0x8e, 0x21, 0x27, 0x37, 0xa4, 0x0e, 0xdc, 0xc7, 0xde, 0xfe, 0x04, 0x00, 0x00,
	0xff, 0xff, 0x1e, 0x82, 0x1a, 0x0f, 0x68, 0x01, 0x00, 0x00,
}
