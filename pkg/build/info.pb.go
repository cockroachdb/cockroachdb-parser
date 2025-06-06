// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: build/info.proto

package build

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
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

// Info describes build information for this CockroachDB binary.
type Info struct {
	// go_version is the version of the Go toolchain used to compile this executable.
	GoVersion string `protobuf:"bytes,1,opt,name=go_version,json=goVersion" json:"go_version"`
	// tag is the binary version for the revision of the source code for
	// this executable.  This used to be the git tag for the build, but
	// has since moved to the contents of the `version.txt` file.
	Tag string `protobuf:"bytes,2,opt,name=tag" json:"tag"`
	// time is the time at which the build started.
	Time string `protobuf:"bytes,3,opt,name=time" json:"time"`
	// revision is the git commit identifier for the source code of this executable.
	Revision string `protobuf:"bytes,4,opt,name=revision" json:"revision"`
	// cgo_compiler is the C/C++ compiler used to build non-go dependencies.
	CgoCompiler string `protobuf:"bytes,5,opt,name=cgo_compiler,json=cgoCompiler" json:"cgo_compiler"`
	// cgo_target_triple is the platform identifier that identifies the cross-compilation target for C/C++ components.
	CgoTargetTriple string `protobuf:"bytes,10,opt,name=cgo_target_triple,json=cgoTargetTriple" json:"cgo_target_triple"`
	// platform is the platform identifiers that identifies the cross-compilation target for Go code.
	Platform string `protobuf:"bytes,6,opt,name=platform" json:"platform"`
	// distribution indicates which licensing conditions apply (OSS: full open source; CCL: includes CCL code).
	Distribution string `protobuf:"bytes,7,opt,name=distribution" json:"distribution"`
	// type indicates whether this is a development or release build.
	Type string `protobuf:"bytes,8,opt,name=type" json:"type"`
	// channel identifies through which product channel the executable was released.
	Channel string `protobuf:"bytes,9,opt,name=channel" json:"channel"`
	// env_channel identifies the product channel as overridden by the COCKROACH_CHANNEL environment variable.
	EnvChannel string `protobuf:"bytes,11,opt,name=env_channel,json=envChannel" json:"env_channel"`
	// enabled_assertions returns the value of 'CrdbTestBuild' (true iff compiled with 'crdb_test' tag)
	EnabledAssertions bool `protobuf:"varint,12,opt,name=enabled_assertions,json=enabledAssertions" json:"enabled_assertions"`
	// dependencies exists to allow tests that run against old clusters
	// to unmarshal JSON containing this field. The tag is unimportant,
	// but the field name must remain unchanged.
	//
	// alternatively, we could set jsonpb.Unmarshaler.AllowUnknownFields
	// to true in httputil.doJSONRequest, but that comes at the expense
	// of run-time type checking, which is nice to have.
	Dependencies *string `protobuf:"bytes,10000,opt,name=dependencies" json:"dependencies,omitempty"`
}

func (m *Info) Reset()         { *m = Info{} }
func (m *Info) String() string { return proto.CompactTextString(m) }
func (*Info) ProtoMessage()    {}
func (*Info) Descriptor() ([]byte, []int) {
	return fileDescriptor_4e07f274d5866c11, []int{0}
}
func (m *Info) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Info) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *Info) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Info.Merge(m, src)
}
func (m *Info) XXX_Size() int {
	return m.Size()
}
func (m *Info) XXX_DiscardUnknown() {
	xxx_messageInfo_Info.DiscardUnknown(m)
}

var xxx_messageInfo_Info proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Info)(nil), "cockroach.parser.build.Info")
}

func init() { proto.RegisterFile("build/info.proto", fileDescriptor_4e07f274d5866c11) }

var fileDescriptor_4e07f274d5866c11 = []byte{
	// 386 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x92, 0xc1, 0x6a, 0x1b, 0x31,
	0x10, 0x40, 0x77, 0x6b, 0xb7, 0xb1, 0xc7, 0x86, 0x34, 0xa2, 0x14, 0xd1, 0x83, 0x6a, 0x1a, 0x4a,
	0x4d, 0x0f, 0xde, 0x42, 0xbf, 0xa0, 0xc9, 0xa9, 0x50, 0x7a, 0x08, 0xa1, 0x87, 0x5e, 0x96, 0x5d,
	0xed, 0x58, 0x16, 0x59, 0x6b, 0x84, 0x56, 0x36, 0xe4, 0x2f, 0xf2, 0x59, 0x3e, 0xe6, 0x98, 0x53,
	0x69, 0xec, 0x1f, 0x29, 0x2b, 0xaf, 0x5d, 0xf9, 0x26, 0xde, 0x7b, 0x68, 0x24, 0x18, 0x78, 0x5d,
	0xae, 0x74, 0x5d, 0x65, 0xda, 0xcc, 0x69, 0x66, 0x1d, 0x79, 0x62, 0xe7, 0x92, 0xe4, 0x9d, 0xa3,
	0x42, 0x2e, 0x66, 0xc1, 0xbd, 0x7b, 0xa3, 0x48, 0x51, 0x70, 0x59, 0x7b, 0xda, 0x67, 0x1f, 0x9e,
	0x7b, 0xd0, 0xff, 0x6e, 0xe6, 0xc4, 0x2e, 0x01, 0x14, 0xe5, 0x6b, 0x74, 0x8d, 0x26, 0xc3, 0xd3,
	0x49, 0x3a, 0x1d, 0x5e, 0xf5, 0x37, 0x7f, 0xde, 0x27, 0x37, 0x43, 0x45, 0xbf, 0xf6, 0x98, 0xbd,
	0x85, 0x9e, 0x2f, 0x14, 0x7f, 0x11, 0xd9, 0x16, 0x30, 0x0e, 0x7d, 0xaf, 0x97, 0xc8, 0x7b, 0x91,
	0x08, 0x84, 0x4d, 0x60, 0xe0, 0x70, 0xad, 0xc3, 0xa5, 0xfd, 0xc8, 0x1e, 0x29, 0xfb, 0x04, 0x63,
	0xa9, 0x28, 0x97, 0xb4, 0xb4, 0xba, 0x46, 0xc7, 0x5f, 0x46, 0xd5, 0x48, 0x2a, 0xba, 0xee, 0x04,
	0xfb, 0x02, 0x17, 0x6d, 0xe8, 0x0b, 0xa7, 0xd0, 0xe7, 0xde, 0x69, 0x5b, 0x23, 0x87, 0xa8, 0x3e,
	0x97, 0x8a, 0x6e, 0x83, 0xbd, 0x0d, 0xb2, 0x1d, 0x6e, 0xeb, 0xc2, 0xcf, 0xc9, 0x2d, 0xf9, 0xab,
	0x78, 0xf8, 0x81, 0xb2, 0x29, 0x8c, 0x2b, 0xdd, 0x78, 0xa7, 0xcb, 0x95, 0x6f, 0x9f, 0x78, 0x16,
	0x55, 0x27, 0x26, 0x7c, 0xf1, 0xde, 0x22, 0x1f, 0x9c, 0x7c, 0xf1, 0xde, 0x22, 0x13, 0x70, 0x26,
	0x17, 0x85, 0x31, 0x58, 0xf3, 0x61, 0x24, 0x0f, 0x90, 0x7d, 0x84, 0x11, 0x9a, 0x75, 0x7e, 0x68,
	0x46, 0x51, 0x03, 0x68, 0xd6, 0xd7, 0x5d, 0xf6, 0x15, 0x18, 0x9a, 0xa2, 0xac, 0xb1, 0xca, 0x8b,
	0xa6, 0x41, 0xd7, 0x4e, 0x6d, 0xf8, 0x78, 0x92, 0x4e, 0x07, 0x5d, 0x7d, 0xd1, 0xf9, 0x6f, 0x47,
	0xcd, 0x2e, 0x61, 0x5c, 0xa1, 0x45, 0x53, 0xa1, 0x91, 0x1a, 0x1b, 0xfe, 0xf0, 0xb3, 0xbd, 0xfd,
	0xe6, 0x04, 0x5e, 0xfd, 0xd8, 0x3c, 0x8b, 0x64, 0xb3, 0x15, 0xe9, 0xe3, 0x56, 0xa4, 0x4f, 0x5b,
	0x91, 0xfe, 0xdd, 0x8a, 0xf4, 0x61, 0x27, 0x92, 0xc7, 0x9d, 0x48, 0x9e, 0x76, 0x22, 0xf9, 0xfd,
	0x59, 0x69, 0xbf, 0x58, 0x95, 0x33, 0x49, 0xcb, 0xec, 0xb8, 0x37, 0x55, 0xf9, 0xff, 0x9c, 0xd9,
	0x3b, 0x95, 0x85, 0x3d, 0xfa, 0x17, 0x00, 0x00, 0xff, 0xff, 0x7a, 0x99, 0x80, 0x4a, 0x6b, 0x02,
	0x00, 0x00,
}

func (m *Info) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Info) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Info) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Dependencies != nil {
		i -= len(*m.Dependencies)
		copy(dAtA[i:], *m.Dependencies)
		i = encodeVarintInfo(dAtA, i, uint64(len(*m.Dependencies)))
		i--
		dAtA[i] = 0x4
		i--
		dAtA[i] = 0xf1
		i--
		dAtA[i] = 0x82
	}
	i--
	if m.EnabledAssertions {
		dAtA[i] = 1
	} else {
		dAtA[i] = 0
	}
	i--
	dAtA[i] = 0x60
	i -= len(m.EnvChannel)
	copy(dAtA[i:], m.EnvChannel)
	i = encodeVarintInfo(dAtA, i, uint64(len(m.EnvChannel)))
	i--
	dAtA[i] = 0x5a
	i -= len(m.CgoTargetTriple)
	copy(dAtA[i:], m.CgoTargetTriple)
	i = encodeVarintInfo(dAtA, i, uint64(len(m.CgoTargetTriple)))
	i--
	dAtA[i] = 0x52
	i -= len(m.Channel)
	copy(dAtA[i:], m.Channel)
	i = encodeVarintInfo(dAtA, i, uint64(len(m.Channel)))
	i--
	dAtA[i] = 0x4a
	i -= len(m.Type)
	copy(dAtA[i:], m.Type)
	i = encodeVarintInfo(dAtA, i, uint64(len(m.Type)))
	i--
	dAtA[i] = 0x42
	i -= len(m.Distribution)
	copy(dAtA[i:], m.Distribution)
	i = encodeVarintInfo(dAtA, i, uint64(len(m.Distribution)))
	i--
	dAtA[i] = 0x3a
	i -= len(m.Platform)
	copy(dAtA[i:], m.Platform)
	i = encodeVarintInfo(dAtA, i, uint64(len(m.Platform)))
	i--
	dAtA[i] = 0x32
	i -= len(m.CgoCompiler)
	copy(dAtA[i:], m.CgoCompiler)
	i = encodeVarintInfo(dAtA, i, uint64(len(m.CgoCompiler)))
	i--
	dAtA[i] = 0x2a
	i -= len(m.Revision)
	copy(dAtA[i:], m.Revision)
	i = encodeVarintInfo(dAtA, i, uint64(len(m.Revision)))
	i--
	dAtA[i] = 0x22
	i -= len(m.Time)
	copy(dAtA[i:], m.Time)
	i = encodeVarintInfo(dAtA, i, uint64(len(m.Time)))
	i--
	dAtA[i] = 0x1a
	i -= len(m.Tag)
	copy(dAtA[i:], m.Tag)
	i = encodeVarintInfo(dAtA, i, uint64(len(m.Tag)))
	i--
	dAtA[i] = 0x12
	i -= len(m.GoVersion)
	copy(dAtA[i:], m.GoVersion)
	i = encodeVarintInfo(dAtA, i, uint64(len(m.GoVersion)))
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintInfo(dAtA []byte, offset int, v uint64) int {
	offset -= sovInfo(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Info) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.GoVersion)
	n += 1 + l + sovInfo(uint64(l))
	l = len(m.Tag)
	n += 1 + l + sovInfo(uint64(l))
	l = len(m.Time)
	n += 1 + l + sovInfo(uint64(l))
	l = len(m.Revision)
	n += 1 + l + sovInfo(uint64(l))
	l = len(m.CgoCompiler)
	n += 1 + l + sovInfo(uint64(l))
	l = len(m.Platform)
	n += 1 + l + sovInfo(uint64(l))
	l = len(m.Distribution)
	n += 1 + l + sovInfo(uint64(l))
	l = len(m.Type)
	n += 1 + l + sovInfo(uint64(l))
	l = len(m.Channel)
	n += 1 + l + sovInfo(uint64(l))
	l = len(m.CgoTargetTriple)
	n += 1 + l + sovInfo(uint64(l))
	l = len(m.EnvChannel)
	n += 1 + l + sovInfo(uint64(l))
	n += 2
	if m.Dependencies != nil {
		l = len(*m.Dependencies)
		n += 3 + l + sovInfo(uint64(l))
	}
	return n
}

func sovInfo(x uint64) (n int) {
	return int((uint32(math_bits.Len64(x|1)+6) * 37) >> 8)
}
func sozInfo(x uint64) (n int) {
	return sovInfo(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Info) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowInfo
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Info: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Info: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GoVersion", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInfo
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthInfo
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthInfo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.GoVersion = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Tag", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInfo
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthInfo
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthInfo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Tag = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Time", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInfo
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthInfo
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthInfo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Time = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Revision", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInfo
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthInfo
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthInfo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Revision = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CgoCompiler", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInfo
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthInfo
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthInfo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CgoCompiler = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Platform", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInfo
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthInfo
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthInfo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Platform = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Distribution", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInfo
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthInfo
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthInfo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Distribution = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInfo
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthInfo
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthInfo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Type = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Channel", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInfo
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthInfo
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthInfo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Channel = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CgoTargetTriple", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInfo
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthInfo
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthInfo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CgoTargetTriple = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EnvChannel", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInfo
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthInfo
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthInfo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.EnvChannel = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 12:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EnabledAssertions", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInfo
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.EnabledAssertions = bool(v != 0)
		case 10000:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Dependencies", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInfo
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthInfo
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthInfo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(dAtA[iNdEx:postIndex])
			m.Dependencies = &s
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipInfo(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthInfo
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipInfo(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowInfo
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowInfo
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowInfo
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthInfo
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupInfo
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthInfo
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthInfo        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowInfo          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupInfo = fmt.Errorf("proto: unexpected end of group")
)

