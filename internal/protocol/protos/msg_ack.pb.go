// Code generated by protoc-gen-go.
// source: msg_ack.proto
// DO NOT EDIT!

package protos

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type MsgAck struct {
	MsgId            *uint64 `protobuf:"varint,1,req,name=msg_id,json=msgId" json:"msg_id,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *MsgAck) Reset()                    { *m = MsgAck{} }
func (m *MsgAck) String() string            { return proto.CompactTextString(m) }
func (*MsgAck) ProtoMessage()               {}
func (*MsgAck) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{0} }

func (m *MsgAck) GetMsgId() uint64 {
	if m != nil && m.MsgId != nil {
		return *m.MsgId
	}
	return 0
}

func init() {
	proto.RegisterType((*MsgAck)(nil), "protos.MsgAck")
}

func init() { proto.RegisterFile("msg_ack.proto", fileDescriptor3) }

var fileDescriptor3 = []byte{
	// 73 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0xcd, 0x2d, 0x4e, 0x8f,
	0x4f, 0x4c, 0xce, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x03, 0x53, 0xc5, 0x4a, 0xf2,
	0x5c, 0x6c, 0xbe, 0xc5, 0xe9, 0x8e, 0xc9, 0xd9, 0x42, 0xa2, 0x5c, 0x6c, 0x20, 0x25, 0x99, 0x29,
	0x12, 0x8c, 0x0a, 0x4c, 0x1a, 0x2c, 0x41, 0xac, 0x40, 0x9e, 0x67, 0x0a, 0x20, 0x00, 0x00, 0xff,
	0xff, 0x03, 0xab, 0x9b, 0xb9, 0x38, 0x00, 0x00, 0x00,
}