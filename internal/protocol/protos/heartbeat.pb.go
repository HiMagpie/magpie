// Code generated by protoc-gen-go.
// source: heartbeat.proto
// DO NOT EDIT!

/*
Package protos is a generated protocol buffer package.

It is generated from these files:
	heartbeat.proto
	login.proto
	msg.proto
	msg_ack.proto
	notify.proto

It has these top-level messages:
	Heartbeat
	Login
	Msg
	MsgAck
	Notify
*/
package protos

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

type Heartbeat struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *Heartbeat) Reset()                    { *m = Heartbeat{} }
func (m *Heartbeat) String() string            { return proto.CompactTextString(m) }
func (*Heartbeat) ProtoMessage()               {}
func (*Heartbeat) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func init() {
	proto.RegisterType((*Heartbeat)(nil), "protos.Heartbeat")
}

func init() { proto.RegisterFile("heartbeat.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 53 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0xcf, 0x48, 0x4d, 0x2c,
	0x2a, 0x49, 0x4a, 0x4d, 0x2c, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x03, 0x53, 0xc5,
	0x4a, 0xdc, 0x5c, 0x9c, 0x1e, 0x30, 0x29, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x47, 0x74, 0x04,
	0x22, 0x26, 0x00, 0x00, 0x00,
}
