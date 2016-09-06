// Code generated by protoc-gen-go.
// source: common.proto
// DO NOT EDIT!

/*
Package meta is a generated protocol buffer package.

It is generated from these files:
	common.proto
	gate.proto
	master.proto
	notice.proto
	store.proto

It has these top-level messages:
	ResponseHeader
	Message
	UserFindUserRequest
	UserFindUserResponse
	UserAddFriendRequest
	UserAddFriendResponse
	UserCreateGroupRequest
	UserCreateGroupResponse
	UserRegisterRequest
	UserRegisterResponse
	UpdateUserInfoRequest
	UpdateUserInfoResponse
	HeartbeatRequest
	HeartbeatResponse
	UserLoginRequest
	UserLoginResponse
	UserLogoutRequest
	UserLogoutResponse
	UploadImageRequest
	UploadImageResponse
	DownloadImageRequest
	DownloadImageResponse
	NoticeRequest
	NoticeResponse
	NewIDRequest
	NewIDResponse
	SubscribeRequest
	SubscribeResponse
	UnsubscribeRequest
	UnsubscribeResponse
	PushRequest
	PushResponse
	FindUserRequest
	FindUserResponse
	AddFriendRequest
	AddFriendResponse
	RegisterRequest
	RegisterResponse
	AuthRequest
	AuthResponse
	CreateGroupRequest
	CreateGroupResponse
	NewMessageRequest
	NewMessageResponse
*/
package meta

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ResponseHeader struct {
	Code int32  `protobuf:"varint,1,opt,name=Code,json=code" json:"Code,omitempty"`
	Msg  string `protobuf:"bytes,2,opt,name=Msg,json=msg" json:"Msg,omitempty"`
}

func (m *ResponseHeader) Reset()                    { *m = ResponseHeader{} }
func (m *ResponseHeader) String() string            { return proto.CompactTextString(m) }
func (*ResponseHeader) ProtoMessage()               {}
func (*ResponseHeader) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Message struct {
	ID   int64  `protobuf:"varint,1,opt,name=ID,json=iD" json:"ID,omitempty"`
	From int64  `protobuf:"varint,2,opt,name=From,json=from" json:"From,omitempty"`
	To   int64  `protobuf:"varint,3,opt,name=To,json=to" json:"To,omitempty"`
	Body string `protobuf:"bytes,4,opt,name=Body,json=body" json:"Body,omitempty"`
}

func (m *Message) Reset()                    { *m = Message{} }
func (m *Message) String() string            { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()               {}
func (*Message) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func init() {
	proto.RegisterType((*ResponseHeader)(nil), "candy.meta.ResponseHeader")
	proto.RegisterType((*Message)(nil), "candy.meta.Message")
}

func init() { proto.RegisterFile("common.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 173 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x24, 0xce, 0xb1, 0x0e, 0x82, 0x30,
	0x10, 0x80, 0xe1, 0x40, 0x2b, 0xc6, 0x8b, 0x21, 0xa6, 0x53, 0x47, 0xc2, 0xc4, 0xc4, 0x62, 0xe2,
	0x03, 0x20, 0x31, 0x3a, 0x30, 0xd8, 0x38, 0xb9, 0x15, 0x7a, 0x12, 0x87, 0x72, 0xa4, 0x65, 0xe9,
	0xdb, 0x9b, 0xd6, 0xfd, 0xcb, 0x9f, 0x1f, 0x8e, 0x13, 0x59, 0x4b, 0x4b, 0xbb, 0x3a, 0xda, 0x48,
	0xc0, 0xa4, 0x17, 0x13, 0x5a, 0x8b, 0x9b, 0xae, 0x2f, 0x50, 0x2a, 0xf4, 0x2b, 0x2d, 0x1e, 0xef,
	0xa8, 0x0d, 0x3a, 0x21, 0x80, 0x5f, 0xc9, 0xa0, 0xcc, 0xaa, 0xac, 0xd9, 0x29, 0x3e, 0x91, 0x41,
	0x71, 0x02, 0x36, 0xf8, 0x59, 0xe6, 0x55, 0xd6, 0x1c, 0x14, 0xb3, 0x7e, 0xae, 0x9f, 0xb0, 0x1f,
	0xd0, 0x7b, 0x3d, 0xa3, 0x28, 0x21, 0x7f, 0xf4, 0x89, 0x33, 0x95, 0x7f, 0xfb, 0x18, 0xb8, 0x39,
	0xb2, 0x49, 0x33, 0xc5, 0x3f, 0x8e, 0x6c, 0x34, 0x2f, 0x92, 0xec, 0x6f, 0x36, 0x8a, 0xa6, 0x23,
	0x13, 0x24, 0x4f, 0x45, 0x3e, 0x92, 0x09, 0x5d, 0xf1, 0xe6, 0x71, 0x69, 0x2c, 0xd2, 0xe5, 0xf9,
	0x17, 0x00, 0x00, 0xff, 0xff, 0xbf, 0x0d, 0x85, 0xfd, 0xb5, 0x00, 0x00, 0x00,
}
