// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sql/lex/encode.proto

package lex

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
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
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// BytesEncodeFormat is the configuration for bytes to string conversions.
type BytesEncodeFormat int32

const (
	// BytesEncodeHex uses the hex format: e'abc\n'::BYTES::STRING -> '\x61626312'.
	// This is the default, for compatibility with PostgreSQL.
	BytesEncodeHex BytesEncodeFormat = 0
	// BytesEncodeEscape uses the escaped format: e'abc\n'::BYTES::STRING -> 'abc\012'.
	BytesEncodeEscape BytesEncodeFormat = 1
	// BytesEncodeBase64 uses base64 encoding.
	BytesEncodeBase64 BytesEncodeFormat = 2
)

var BytesEncodeFormat_name = map[int32]string{
	0: "BytesEncodeHex",
	1: "BytesEncodeEscape",
	2: "BytesEncodeBase64",
}

var BytesEncodeFormat_value = map[string]int32{
	"BytesEncodeHex":    0,
	"BytesEncodeEscape": 1,
	"BytesEncodeBase64": 2,
}

func (BytesEncodeFormat) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_3d1d8cf249b1b776, []int{0}
}

func init() {
	proto.RegisterEnum("cockroach.sql.sessiondatapb.BytesEncodeFormat", BytesEncodeFormat_name, BytesEncodeFormat_value)
}

func init() { proto.RegisterFile("sql/lex/encode.proto", fileDescriptor_3d1d8cf249b1b776) }

var fileDescriptor_3d1d8cf249b1b776 = []byte{
	// 196 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x29, 0x2e, 0xcc, 0xd1,
	0xcf, 0x49, 0xad, 0xd0, 0x4f, 0xcd, 0x4b, 0xce, 0x4f, 0x49, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9,
	0x17, 0x92, 0x4e, 0xce, 0x4f, 0xce, 0x2e, 0xca, 0x4f, 0x4c, 0xce, 0xd0, 0x2b, 0x2e, 0xcc, 0xd1,
	0x2b, 0x4e, 0x2d, 0x2e, 0xce, 0xcc, 0xcf, 0x4b, 0x49, 0x2c, 0x49, 0x2c, 0x48, 0x92, 0x12, 0x49,
	0xcf, 0x4f, 0xcf, 0x07, 0xab, 0xd3, 0x07, 0xb1, 0x20, 0x5a, 0xb4, 0xe2, 0xb9, 0x04, 0x9d, 0x2a,
	0x4b, 0x52, 0x8b, 0x5d, 0xc1, 0xe6, 0xb8, 0xe5, 0x17, 0xe5, 0x26, 0x96, 0x08, 0x09, 0x71, 0xf1,
	0x21, 0x09, 0x7a, 0xa4, 0x56, 0x08, 0x30, 0x08, 0x89, 0xa2, 0x28, 0x74, 0x2d, 0x4e, 0x4e, 0x2c,
	0x48, 0x15, 0x60, 0x44, 0x13, 0x76, 0x4a, 0x2c, 0x4e, 0x35, 0x33, 0x11, 0x60, 0x92, 0xe2, 0xe8,
	0x58, 0x2c, 0xc7, 0xb0, 0x62, 0x89, 0x1c, 0x83, 0x93, 0xea, 0x89, 0x87, 0x72, 0x0c, 0x27, 0x1e,
	0xc9, 0x31, 0x5e, 0x78, 0x24, 0xc7, 0x78, 0xe3, 0x91, 0x1c, 0xe3, 0x83, 0x47, 0x72, 0x8c, 0x13,
	0x1e, 0xcb, 0x31, 0x5c, 0x78, 0x2c, 0xc7, 0x70, 0xe3, 0xb1, 0x1c, 0x43, 0x14, 0x73, 0x4e, 0x6a,
	0x45, 0x12, 0x1b, 0xd8, 0x39, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x9f, 0x9b, 0x2d, 0x80,
	0xd9, 0x00, 0x00, 0x00,
}

