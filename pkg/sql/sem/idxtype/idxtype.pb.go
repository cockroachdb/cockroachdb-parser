// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sql/sem/idxtype/idxtype.proto

package idxtype

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

// T represents different types of indexes that can be defined on a table.
type T int32

const (
	// FORWARD is a standard relational index, containing at most one entry for
	// each row in the table (if a partial index, a table row may not have a
	// corresponding entry in the index).
	FORWARD T = 0
	// INVERTED indexes can contain multiple entries for each row in the table,
	// which is useful for indexing collection or composite data types like JSONB,
	// ARRAY, and GEOGRAPHY.
	INVERTED T = 1
	// VECTOR indexes high-dimensional vectors to enable rapid similarity search
	// using an approximate nearest neighbor (ANN) algorithm.
	VECTOR T = 2
)

var T_name = map[int32]string{
	0: "FORWARD",
	1: "INVERTED",
	2: "VECTOR",
}

var T_value = map[string]int32{
	"FORWARD":  0,
	"INVERTED": 1,
	"VECTOR":   2,
}

func (x T) String() string {
	return proto.EnumName(T_name, int32(x))
}

func (T) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_f11c9d993fdc9c28, []int{0}
}

func init() {
	proto.RegisterEnum("cockroach.parser.sql.sem.idxtype.T", T_name, T_value)
}

func init() { proto.RegisterFile("sql/sem/idxtype/idxtype.proto", fileDescriptor_f11c9d993fdc9c28) }

var fileDescriptor_f11c9d993fdc9c28 = []byte{
	// 201 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x2d, 0x2e, 0xcc, 0xd1,
	0x2f, 0x4e, 0xcd, 0xd5, 0xcf, 0x4c, 0xa9, 0x28, 0xa9, 0x2c, 0x48, 0x85, 0xd1, 0x7a, 0x05, 0x45,
	0xf9, 0x25, 0xf9, 0x42, 0x92, 0xc9, 0xf9, 0xc9, 0xd9, 0x45, 0xf9, 0x89, 0xc9, 0x19, 0x7a, 0xc5,
	0x85, 0x39, 0x7a, 0xc5, 0xa9, 0xb9, 0x7a, 0x50, 0x05, 0x52, 0x22, 0xe9, 0xf9, 0xe9, 0xf9, 0x60,
	0x55, 0xfa, 0x20, 0x16, 0x44, 0x83, 0x96, 0x01, 0x17, 0x63, 0x88, 0x10, 0x37, 0x17, 0xbb, 0x9b,
	0x7f, 0x50, 0xb8, 0x63, 0x90, 0x8b, 0x00, 0x83, 0x10, 0x0f, 0x17, 0x87, 0xa7, 0x5f, 0x98, 0x6b,
	0x50, 0x88, 0xab, 0x8b, 0x00, 0xa3, 0x10, 0x17, 0x17, 0x5b, 0x98, 0xab, 0x73, 0x88, 0x7f, 0x90,
	0x00, 0x93, 0x14, 0x4b, 0xc7, 0x62, 0x39, 0x06, 0xa7, 0xb0, 0x13, 0x0f, 0xe5, 0x18, 0x4e, 0x3c,
	0x92, 0x63, 0xbc, 0xf0, 0x48, 0x8e, 0xf1, 0xc6, 0x23, 0x39, 0xc6, 0x07, 0x8f, 0xe4, 0x18, 0x27,
	0x3c, 0x96, 0x63, 0xb8, 0xf0, 0x58, 0x8e, 0xe1, 0xc6, 0x63, 0x39, 0x86, 0x28, 0x93, 0xf4, 0xcc,
	0x92, 0x8c, 0xd2, 0x24, 0xbd, 0xe4, 0xfc, 0x5c, 0x7d, 0xb8, 0x7b, 0x52, 0x92, 0x10, 0x6c, 0xfd,
	0x82, 0xec, 0x74, 0x7d, 0x34, 0x8f, 0x24, 0xb1, 0x81, 0x1d, 0x64, 0x0c, 0x08, 0x00, 0x00, 0xff,
	0xff, 0x16, 0x50, 0xca, 0xe4, 0xe2, 0x00, 0x00, 0x00,
}

