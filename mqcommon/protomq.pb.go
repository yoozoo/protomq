// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protomq.proto

package mqcommon

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

var E_Topic = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.MessageOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         55008,
	Name:          "protomq.topic",
	Tag:           "bytes,55008,opt,name=topic",
	Filename:      "protomq.proto",
}

func init() {
	proto.RegisterExtension(E_Topic)
}

func init() { proto.RegisterFile("protomq.proto", fileDescriptor_06c3109ad938d537) }

var fileDescriptor_06c3109ad938d537 = []byte{
	// 126 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2d, 0x28, 0xca, 0x2f,
	0xc9, 0xcf, 0x2d, 0xd4, 0x03, 0xd3, 0x42, 0xec, 0x50, 0xae, 0x94, 0x42, 0x7a, 0x7e, 0x7e, 0x7a,
	0x4e, 0xaa, 0x3e, 0x98, 0x9f, 0x54, 0x9a, 0xa6, 0x9f, 0x92, 0x5a, 0x9c, 0x5c, 0x94, 0x59, 0x50,
	0x92, 0x5f, 0x04, 0x51, 0x6a, 0x65, 0xce, 0xc5, 0x5a, 0x92, 0x5f, 0x90, 0x99, 0x2c, 0x24, 0xaf,
	0x07, 0x51, 0xab, 0x07, 0x53, 0xab, 0xe7, 0x9b, 0x5a, 0x5c, 0x9c, 0x98, 0x9e, 0xea, 0x5f, 0x50,
	0x92, 0x99, 0x9f, 0x57, 0x2c, 0xf1, 0x60, 0x2d, 0xb3, 0x02, 0xa3, 0x06, 0x67, 0x10, 0x44, 0xbd,
	0x13, 0x57, 0x14, 0x47, 0x6e, 0x61, 0x72, 0x7e, 0x6e, 0x6e, 0x7e, 0x5e, 0x12, 0x1b, 0x58, 0x8f,
	0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0x60, 0xd1, 0x6b, 0x29, 0x87, 0x00, 0x00, 0x00,
}
