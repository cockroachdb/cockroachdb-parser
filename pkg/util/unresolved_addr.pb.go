// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: util/unresolved_addr.proto

package util

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

// UnresolvedAddr is an unresolved version of net.Addr.
type UnresolvedAddr struct {
	NetworkField string `protobuf:"bytes,1,opt,name=network_field,json=networkField" json:"network_field"`
	AddressField string `protobuf:"bytes,2,opt,name=address_field,json=addressField" json:"address_field"`
}

func (m *UnresolvedAddr) Reset()      { *m = UnresolvedAddr{} }
func (*UnresolvedAddr) ProtoMessage() {}
func (*UnresolvedAddr) Descriptor() ([]byte, []int) {
	return fileDescriptor_e843f4480e4927e4, []int{0}
}
func (m *UnresolvedAddr) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *UnresolvedAddr) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *UnresolvedAddr) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UnresolvedAddr.Merge(m, src)
}
func (m *UnresolvedAddr) XXX_Size() int {
	return m.Size()
}
func (m *UnresolvedAddr) XXX_DiscardUnknown() {
	xxx_messageInfo_UnresolvedAddr.DiscardUnknown(m)
}

var xxx_messageInfo_UnresolvedAddr proto.InternalMessageInfo

func init() {
	proto.RegisterType((*UnresolvedAddr)(nil), "cockroach.parser.util.UnresolvedAddr")
}

func init() { proto.RegisterFile("util/unresolved_addr.proto", fileDescriptor_e843f4480e4927e4) }

var fileDescriptor_e843f4480e4927e4 = []byte{
	// 215 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x2a, 0x2d, 0xc9, 0xcc,
	0xd1, 0x2f, 0xcd, 0x2b, 0x4a, 0x2d, 0xce, 0xcf, 0x29, 0x4b, 0x4d, 0x89, 0x4f, 0x4c, 0x49, 0x29,
	0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x4b, 0xce, 0x4f, 0xce, 0x2e, 0xca, 0x4f, 0x4c,
	0xce, 0xd0, 0x03, 0xa9, 0x92, 0x12, 0x49, 0xcf, 0x4f, 0xcf, 0x07, 0x4b, 0xe9, 0x83, 0x58, 0x10,
	0x55, 0x4a, 0x05, 0x5c, 0x7c, 0xa1, 0x70, 0xed, 0x8e, 0x29, 0x29, 0x45, 0x42, 0x9a, 0x5c, 0xbc,
	0x79, 0xa9, 0x25, 0xe5, 0xf9, 0x45, 0xd9, 0xf1, 0x69, 0x99, 0xa9, 0x39, 0x29, 0x12, 0x8c, 0x0a,
	0x8c, 0x1a, 0x9c, 0x4e, 0x2c, 0x27, 0xee, 0xc9, 0x33, 0x04, 0xf1, 0x40, 0xa5, 0xdc, 0x40, 0x32,
	0x20, 0xa5, 0x20, 0x0b, 0x53, 0x8b, 0x8b, 0xa1, 0x4a, 0x99, 0x90, 0x95, 0x42, 0xa5, 0xc0, 0x4a,
	0xad, 0x38, 0x66, 0x2c, 0x90, 0x67, 0x78, 0xb1, 0x40, 0x9e, 0xd1, 0xc9, 0xfb, 0xc4, 0x43, 0x39,
	0x86, 0x13, 0x8f, 0xe4, 0x18, 0x2f, 0x3c, 0x92, 0x63, 0xbc, 0xf1, 0x48, 0x8e, 0xf1, 0xc1, 0x23,
	0x39, 0xc6, 0x09, 0x8f, 0xe5, 0x18, 0x2e, 0x3c, 0x96, 0x63, 0xb8, 0xf1, 0x58, 0x8e, 0x21, 0x4a,
	0x33, 0x3d, 0xb3, 0x24, 0xa3, 0x34, 0x49, 0x2f, 0x39, 0x3f, 0x57, 0x1f, 0xee, 0x89, 0x94, 0x24,
	0x04, 0x5b, 0xbf, 0x20, 0x3b, 0x5d, 0x1f, 0xe4, 0x29, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0xc2,
	0x21, 0xae, 0xea, 0x01, 0x01, 0x00, 0x00,
}

func (this *UnresolvedAddr) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*UnresolvedAddr)
	if !ok {
		that2, ok := that.(UnresolvedAddr)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.NetworkField != that1.NetworkField {
		return false
	}
	if this.AddressField != that1.AddressField {
		return false
	}
	return true
}
func (m *UnresolvedAddr) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UnresolvedAddr) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *UnresolvedAddr) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	i -= len(m.AddressField)
	copy(dAtA[i:], m.AddressField)
	i = encodeVarintUnresolvedAddr(dAtA, i, uint64(len(m.AddressField)))
	i--
	dAtA[i] = 0x12
	i -= len(m.NetworkField)
	copy(dAtA[i:], m.NetworkField)
	i = encodeVarintUnresolvedAddr(dAtA, i, uint64(len(m.NetworkField)))
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintUnresolvedAddr(dAtA []byte, offset int, v uint64) int {
	offset -= sovUnresolvedAddr(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *UnresolvedAddr) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.NetworkField)
	n += 1 + l + sovUnresolvedAddr(uint64(l))
	l = len(m.AddressField)
	n += 1 + l + sovUnresolvedAddr(uint64(l))
	return n
}

func sovUnresolvedAddr(x uint64) (n int) {
	return int((uint32(math_bits.Len64(x|1)+6) * 37) >> 8)
}
func sozUnresolvedAddr(x uint64) (n int) {
	return sovUnresolvedAddr(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *UnresolvedAddr) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowUnresolvedAddr
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
			return fmt.Errorf("proto: UnresolvedAddr: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: UnresolvedAddr: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NetworkField", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUnresolvedAddr
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
				return ErrInvalidLengthUnresolvedAddr
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthUnresolvedAddr
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NetworkField = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AddressField", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUnresolvedAddr
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
				return ErrInvalidLengthUnresolvedAddr
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthUnresolvedAddr
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AddressField = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipUnresolvedAddr(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthUnresolvedAddr
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
func skipUnresolvedAddr(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowUnresolvedAddr
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
					return 0, ErrIntOverflowUnresolvedAddr
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
					return 0, ErrIntOverflowUnresolvedAddr
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
				return 0, ErrInvalidLengthUnresolvedAddr
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupUnresolvedAddr
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthUnresolvedAddr
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthUnresolvedAddr        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowUnresolvedAddr          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupUnresolvedAddr = fmt.Errorf("proto: unexpected end of group")
)

