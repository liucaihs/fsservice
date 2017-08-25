// Code generated by protoc-gen-go. DO NOT EDIT.
// source: trojan.proto

/*
Package command is a generated protocol buffer package.

It is generated from these files:
	trojan.proto

It has these top-level messages:
	PageOption
	User
*/
package command

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

type PageOption struct {
	PageSize   int32 `protobuf:"varint,1,opt,name=page_size,json=pageSize" json:"page_size,omitempty"`
	PerPageNum int32 `protobuf:"varint,2,opt,name=per_page_num,json=perPageNum" json:"per_page_num,omitempty"`
}

func (m *PageOption) Reset()                    { *m = PageOption{} }
func (m *PageOption) String() string            { return proto.CompactTextString(m) }
func (*PageOption) ProtoMessage()               {}
func (*PageOption) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *PageOption) GetPageSize() int32 {
	if m != nil {
		return m.PageSize
	}
	return 0
}

func (m *PageOption) GetPerPageNum() int32 {
	if m != nil {
		return m.PerPageNum
	}
	return 0
}

type User struct {
	Imei  string `protobuf:"bytes,1,opt,name=imei" json:"imei,omitempty"`
	Imsi  string `protobuf:"bytes,2,opt,name=imsi" json:"imsi,omitempty"`
	Iccid string `protobuf:"bytes,3,opt,name=iccid" json:"iccid,omitempty"`
}

func (m *User) Reset()                    { *m = User{} }
func (m *User) String() string            { return proto.CompactTextString(m) }
func (*User) ProtoMessage()               {}
func (*User) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *User) GetImei() string {
	if m != nil {
		return m.Imei
	}
	return ""
}

func (m *User) GetImsi() string {
	if m != nil {
		return m.Imsi
	}
	return ""
}

func (m *User) GetIccid() string {
	if m != nil {
		return m.Iccid
	}
	return ""
}

func init() {
	proto.RegisterType((*PageOption)(nil), "command.PageOption")
	proto.RegisterType((*User)(nil), "command.User")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Command service

type CommandClient interface {
	// list users on pages
	ListOnline(ctx context.Context, in *PageOption, opts ...grpc.CallOption) (Command_ListOnlineClient, error)
}

type commandClient struct {
	cc *grpc.ClientConn
}

func NewCommandClient(cc *grpc.ClientConn) CommandClient {
	return &commandClient{cc}
}

func (c *commandClient) ListOnline(ctx context.Context, in *PageOption, opts ...grpc.CallOption) (Command_ListOnlineClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Command_serviceDesc.Streams[0], c.cc, "/command.Command/ListOnline", opts...)
	if err != nil {
		return nil, err
	}
	x := &commandListOnlineClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Command_ListOnlineClient interface {
	Recv() (*User, error)
	grpc.ClientStream
}

type commandListOnlineClient struct {
	grpc.ClientStream
}

func (x *commandListOnlineClient) Recv() (*User, error) {
	m := new(User)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Command service

type CommandServer interface {
	// list users on pages
	ListOnline(*PageOption, Command_ListOnlineServer) error
}

func RegisterCommandServer(s *grpc.Server, srv CommandServer) {
	s.RegisterService(&_Command_serviceDesc, srv)
}

func _Command_ListOnline_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(PageOption)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(CommandServer).ListOnline(m, &commandListOnlineServer{stream})
}

type Command_ListOnlineServer interface {
	Send(*User) error
	grpc.ServerStream
}

type commandListOnlineServer struct {
	grpc.ServerStream
}

func (x *commandListOnlineServer) Send(m *User) error {
	return x.ServerStream.SendMsg(m)
}

var _Command_serviceDesc = grpc.ServiceDesc{
	ServiceName: "command.Command",
	HandlerType: (*CommandServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ListOnline",
			Handler:       _Command_ListOnline_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "trojan.proto",
}

func init() { proto.RegisterFile("trojan.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 200 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x44, 0x8f, 0xbf, 0x4f, 0x87, 0x40,
	0x0c, 0x47, 0x45, 0x41, 0xa4, 0xc1, 0xa5, 0x3a, 0x10, 0x5d, 0x08, 0x93, 0x13, 0x31, 0xea, 0xee,
	0xa0, 0x9b, 0x46, 0x0c, 0xc6, 0x99, 0x9c, 0xd0, 0x90, 0x1a, 0xef, 0x47, 0xee, 0x8e, 0x85, 0xbf,
	0xde, 0x50, 0x8c, 0xdf, 0xed, 0xee, 0xb5, 0xfd, 0xf4, 0x15, 0xca, 0xe8, 0xed, 0xb7, 0x32, 0xad,
	0xf3, 0x36, 0x5a, 0xcc, 0x47, 0xab, 0xb5, 0x32, 0x53, 0xf3, 0x02, 0xf0, 0xae, 0x66, 0xea, 0x5c,
	0x64, 0x6b, 0xf0, 0x1a, 0x0a, 0xa7, 0x66, 0x1a, 0x02, 0xaf, 0x54, 0x25, 0x75, 0x72, 0x93, 0xf5,
	0x67, 0x1b, 0xf8, 0xe0, 0x95, 0xb0, 0x86, 0xd2, 0x91, 0x1f, 0xa4, 0xc1, 0x2c, 0xba, 0x3a, 0x96,
	0x3a, 0x38, 0xf2, 0x5b, 0xc2, 0xdb, 0xa2, 0x9b, 0x67, 0x48, 0x3f, 0x03, 0x79, 0x44, 0x48, 0x59,
	0x13, 0x4b, 0x42, 0xd1, 0xcb, 0x7b, 0x67, 0x81, 0x65, 0x4a, 0x58, 0x60, 0xbc, 0x84, 0x8c, 0xc7,
	0x91, 0xa7, 0xea, 0x44, 0xe0, 0xfe, 0xb9, 0x7b, 0x84, 0xfc, 0x69, 0xb7, 0xc3, 0x07, 0x80, 0x57,
	0x0e, 0xb1, 0x33, 0x3f, 0x6c, 0x08, 0x2f, 0xda, 0x3f, 0xeb, 0xf6, 0xa0, 0x7c, 0x75, 0xfe, 0x0f,
	0xb7, 0xd5, 0xcd, 0xd1, 0x6d, 0xf2, 0x75, 0x2a, 0x37, 0xde, 0xff, 0x06, 0x00, 0x00, 0xff, 0xff,
	0xc7, 0x92, 0x8a, 0x13, 0xf3, 0x00, 0x00, 0x00,
}
