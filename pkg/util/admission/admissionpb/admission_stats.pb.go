// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: util/admission/admissionpb/admission_stats.proto

package admissionpb

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// AdmissionWorkQueueStats is recorded for work items waiting in the admission
// work queue.
type AdmissionWorkQueueStats struct {
	// Duration spent waiting.
	WaitDurationNanos time.Duration `protobuf:"varint,1,opt,name=wait_duration_nanos,json=waitDurationNanos,proto3,casttype=time.Duration" json:"wait_duration_nanos,omitempty"`
	// String representation of admission queue kind.
	QueueKind string `protobuf:"bytes,2,opt,name=queue_kind,json=queueKind,proto3" json:"queue_kind,omitempty"`
	// Set to true if deadline was exceeded.
	DeadlineExceeded bool `protobuf:"varint,3,opt,name=deadline_exceeded,json=deadlineExceeded,proto3" json:"deadline_exceeded,omitempty"`
	// String representation of work priority.
	WorkPriority string `protobuf:"bytes,4,opt,name=work_priority,json=workPriority,proto3" json:"work_priority,omitempty"`
}

func (m *AdmissionWorkQueueStats) Reset()         { *m = AdmissionWorkQueueStats{} }
func (m *AdmissionWorkQueueStats) String() string { return proto.CompactTextString(m) }
func (*AdmissionWorkQueueStats) ProtoMessage()    {}
func (*AdmissionWorkQueueStats) Descriptor() ([]byte, []int) {
	return fileDescriptor_695e7752d0aa7c6e, []int{0}
}
func (m *AdmissionWorkQueueStats) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AdmissionWorkQueueStats) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *AdmissionWorkQueueStats) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AdmissionWorkQueueStats.Merge(m, src)
}
func (m *AdmissionWorkQueueStats) XXX_Size() int {
	return m.Size()
}
func (m *AdmissionWorkQueueStats) XXX_DiscardUnknown() {
	xxx_messageInfo_AdmissionWorkQueueStats.DiscardUnknown(m)
}

var xxx_messageInfo_AdmissionWorkQueueStats proto.InternalMessageInfo

func init() {
	proto.RegisterType((*AdmissionWorkQueueStats)(nil), "cockroach.parser.util.admission.admissionpb.AdmissionWorkQueueStats")
}

func init() {
	proto.RegisterFile("util/admission/admissionpb/admission_stats.proto", fileDescriptor_695e7752d0aa7c6e)
}

var fileDescriptor_695e7752d0aa7c6e = []byte{
	// 313 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x91, 0xc1, 0x4a, 0xf3, 0x40,
	0x14, 0x85, 0x33, 0x7f, 0x7f, 0xc4, 0x0e, 0x16, 0x6c, 0x14, 0x0c, 0x82, 0x63, 0x51, 0x17, 0x05,
	0x61, 0x22, 0xf8, 0x00, 0xd2, 0xa2, 0x2b, 0x41, 0xb4, 0x2e, 0x04, 0x37, 0x61, 0x92, 0x19, 0xd2,
	0x21, 0xed, 0xdc, 0x38, 0x99, 0x50, 0x7d, 0x0b, 0x1f, 0xab, 0x1b, 0xa1, 0xcb, 0xae, 0x44, 0xd3,
	0xb7, 0x70, 0x25, 0x13, 0x93, 0xc6, 0x8d, 0xbb, 0x33, 0xe7, 0x9e, 0xef, 0x0c, 0xdc, 0x8b, 0xcf,
	0x72, 0x23, 0x27, 0x3e, 0xe3, 0x53, 0x99, 0x65, 0x12, 0x54, 0xa3, 0xd2, 0xb0, 0xd1, 0x41, 0x66,
	0x98, 0xc9, 0x68, 0xaa, 0xc1, 0x80, 0x7b, 0x12, 0x41, 0x94, 0x68, 0x60, 0xd1, 0x98, 0x5a, 0x96,
	0xae, 0x53, 0xf4, 0x17, 0xbb, 0xbf, 0x1b, 0x43, 0x0c, 0x25, 0xe0, 0x5b, 0xf5, 0xc3, 0x1e, 0xbd,
	0x21, 0xbc, 0x37, 0xa8, 0x53, 0x0f, 0xa0, 0x93, 0xbb, 0x5c, 0xe4, 0xe2, 0xde, 0xb6, 0xbb, 0x03,
	0xbc, 0x33, 0x63, 0xd2, 0x04, 0x3c, 0xd7, 0xcc, 0xd8, 0x4f, 0x15, 0x53, 0x90, 0x79, 0xa8, 0x87,
	0xfa, 0xad, 0x61, 0xf7, 0xeb, 0xfd, 0xb0, 0x63, 0xe4, 0x54, 0xd0, 0xcb, 0x6a, 0x3c, 0xea, 0xda,
	0x74, 0xfd, 0xba, 0xb1, 0x59, 0xf7, 0x00, 0xe3, 0x27, 0x5b, 0x18, 0x24, 0x52, 0x71, 0xef, 0x5f,
	0x0f, 0xf5, 0xdb, 0xa3, 0x76, 0xe9, 0x5c, 0x4b, 0xc5, 0xdd, 0x53, 0xdc, 0xe5, 0x82, 0xf1, 0x89,
	0x54, 0x22, 0x10, 0xcf, 0x91, 0x10, 0x5c, 0x70, 0xaf, 0xd5, 0x43, 0xfd, 0xcd, 0xd1, 0x76, 0x3d,
	0xb8, 0xaa, 0x7c, 0xf7, 0x18, 0x77, 0x66, 0xa0, 0x93, 0x20, 0xd5, 0x12, 0xb4, 0x34, 0x2f, 0xde,
	0xff, 0xb2, 0x6e, 0xcb, 0x9a, 0xb7, 0x95, 0x37, 0x64, 0xf3, 0x4f, 0xe2, 0xcc, 0x0b, 0x82, 0x16,
	0x05, 0x41, 0xcb, 0x82, 0xa0, 0x8f, 0x82, 0xa0, 0xd7, 0x15, 0x71, 0x16, 0x2b, 0xe2, 0x2c, 0x57,
	0xc4, 0x79, 0xbc, 0x88, 0xa5, 0x19, 0xe7, 0x21, 0x8d, 0x60, 0xea, 0xaf, 0x17, 0xc7, 0xc3, 0x46,
	0xfb, 0x69, 0x12, 0xfb, 0x7f, 0x1f, 0x21, 0xdc, 0x28, 0x37, 0x77, 0xfe, 0x1d, 0x00, 0x00, 0xff,
	0xff, 0x8f, 0x5e, 0x5b, 0x8e, 0xa9, 0x01, 0x00, 0x00,
}

func (m *AdmissionWorkQueueStats) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AdmissionWorkQueueStats) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AdmissionWorkQueueStats) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.WorkPriority) > 0 {
		i -= len(m.WorkPriority)
		copy(dAtA[i:], m.WorkPriority)
		i = encodeVarintAdmissionStats(dAtA, i, uint64(len(m.WorkPriority)))
		i--
		dAtA[i] = 0x22
	}
	if m.DeadlineExceeded {
		i--
		if m.DeadlineExceeded {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x18
	}
	if len(m.QueueKind) > 0 {
		i -= len(m.QueueKind)
		copy(dAtA[i:], m.QueueKind)
		i = encodeVarintAdmissionStats(dAtA, i, uint64(len(m.QueueKind)))
		i--
		dAtA[i] = 0x12
	}
	if m.WaitDurationNanos != 0 {
		i = encodeVarintAdmissionStats(dAtA, i, uint64(m.WaitDurationNanos))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintAdmissionStats(dAtA []byte, offset int, v uint64) int {
	offset -= sovAdmissionStats(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *AdmissionWorkQueueStats) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.WaitDurationNanos != 0 {
		n += 1 + sovAdmissionStats(uint64(m.WaitDurationNanos))
	}
	l = len(m.QueueKind)
	if l > 0 {
		n += 1 + l + sovAdmissionStats(uint64(l))
	}
	if m.DeadlineExceeded {
		n += 2
	}
	l = len(m.WorkPriority)
	if l > 0 {
		n += 1 + l + sovAdmissionStats(uint64(l))
	}
	return n
}

func sovAdmissionStats(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozAdmissionStats(x uint64) (n int) {
	return sovAdmissionStats(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *AdmissionWorkQueueStats) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAdmissionStats
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
			return fmt.Errorf("proto: AdmissionWorkQueueStats: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AdmissionWorkQueueStats: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field WaitDurationNanos", wireType)
			}
			m.WaitDurationNanos = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAdmissionStats
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.WaitDurationNanos |= time.Duration(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field QueueKind", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAdmissionStats
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
				return ErrInvalidLengthAdmissionStats
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAdmissionStats
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.QueueKind = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field DeadlineExceeded", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAdmissionStats
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
			m.DeadlineExceeded = bool(v != 0)
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field WorkPriority", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAdmissionStats
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
				return ErrInvalidLengthAdmissionStats
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAdmissionStats
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.WorkPriority = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipAdmissionStats(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthAdmissionStats
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
func skipAdmissionStats(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowAdmissionStats
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
					return 0, ErrIntOverflowAdmissionStats
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
					return 0, ErrIntOverflowAdmissionStats
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
				return 0, ErrInvalidLengthAdmissionStats
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupAdmissionStats
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthAdmissionStats
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthAdmissionStats        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowAdmissionStats          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupAdmissionStats = fmt.Errorf("proto: unexpected end of group")
)

