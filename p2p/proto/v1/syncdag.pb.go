// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: syncdag.proto

package qitmeer_p2p_v1

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/golang/protobuf/proto"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type SyncDAG struct {
	MainLocator          []*Hash     `protobuf:"bytes,1,rep,name=mainLocator,proto3" json:"mainLocator,omitempty" ssz-max:"32"`
	GraphState           *GraphState `protobuf:"bytes,2,opt,name=graphState,proto3" json:"graphState,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *SyncDAG) Reset()         { *m = SyncDAG{} }
func (m *SyncDAG) String() string { return proto.CompactTextString(m) }
func (*SyncDAG) ProtoMessage()    {}
func (*SyncDAG) Descriptor() ([]byte, []int) {
	return fileDescriptor_9cb77bcde7ac0e2c, []int{0}
}
func (m *SyncDAG) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SyncDAG) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SyncDAG.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SyncDAG) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SyncDAG.Merge(m, src)
}
func (m *SyncDAG) XXX_Size() int {
	return m.Size()
}
func (m *SyncDAG) XXX_DiscardUnknown() {
	xxx_messageInfo_SyncDAG.DiscardUnknown(m)
}

var xxx_messageInfo_SyncDAG proto.InternalMessageInfo

func (m *SyncDAG) GetMainLocator() []*Hash {
	if m != nil {
		return m.MainLocator
	}
	return nil
}

func (m *SyncDAG) GetGraphState() *GraphState {
	if m != nil {
		return m.GraphState
	}
	return nil
}

type SubDAG struct {
	SyncPoint            *Hash        `protobuf:"bytes,1,opt,name=syncPoint,proto3" json:"syncPoint,omitempty"`
	GraphState           *GraphState  `protobuf:"bytes,2,opt,name=graphState,proto3" json:"graphState,omitempty"`
	Blocks               []*BlockData `protobuf:"bytes,3,rep,name=blocks,proto3" json:"blocks,omitempty" ssz-max:"500"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *SubDAG) Reset()         { *m = SubDAG{} }
func (m *SubDAG) String() string { return proto.CompactTextString(m) }
func (*SubDAG) ProtoMessage()    {}
func (*SubDAG) Descriptor() ([]byte, []int) {
	return fileDescriptor_9cb77bcde7ac0e2c, []int{1}
}
func (m *SubDAG) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SubDAG) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SubDAG.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SubDAG) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SubDAG.Merge(m, src)
}
func (m *SubDAG) XXX_Size() int {
	return m.Size()
}
func (m *SubDAG) XXX_DiscardUnknown() {
	xxx_messageInfo_SubDAG.DiscardUnknown(m)
}

var xxx_messageInfo_SubDAG proto.InternalMessageInfo

func (m *SubDAG) GetSyncPoint() *Hash {
	if m != nil {
		return m.SyncPoint
	}
	return nil
}

func (m *SubDAG) GetGraphState() *GraphState {
	if m != nil {
		return m.GraphState
	}
	return nil
}

func (m *SubDAG) GetBlocks() []*BlockData {
	if m != nil {
		return m.Blocks
	}
	return nil
}

func init() {
	proto.RegisterType((*SyncDAG)(nil), "qitmeer.p2p.v1.SyncDAG")
	proto.RegisterType((*SubDAG)(nil), "qitmeer.p2p.v1.SubDAG")
}

func init() { proto.RegisterFile("syncdag.proto", fileDescriptor_9cb77bcde7ac0e2c) }

var fileDescriptor_9cb77bcde7ac0e2c = []byte{
	// 305 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2d, 0xae, 0xcc, 0x4b,
	0x4e, 0x49, 0x4c, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x2b, 0xcc, 0x2c, 0xc9, 0x4d,
	0x4d, 0x2d, 0xd2, 0x2b, 0x30, 0x2a, 0xd0, 0x2b, 0x33, 0x94, 0xd2, 0x4d, 0xcf, 0x2c, 0xc9, 0x28,
	0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0x4f, 0xcf, 0x4f, 0xcf, 0xd7, 0x07, 0x2b, 0x4b, 0x2a, 0x4d,
	0x03, 0xf3, 0xc0, 0x1c, 0x30, 0x0b, 0xa2, 0x5d, 0x8a, 0x2f, 0x37, 0xb5, 0xb8, 0x38, 0x31, 0x3d,
	0xb5, 0x18, 0xca, 0x17, 0x48, 0x2f, 0x4a, 0x2c, 0xc8, 0x28, 0x2e, 0x49, 0x2c, 0x49, 0x85, 0x8a,
	0xf0, 0xa7, 0xa7, 0x96, 0x24, 0xe5, 0xe4, 0x27, 0x67, 0x43, 0x95, 0x28, 0xf5, 0x33, 0x72, 0xb1,
	0x07, 0x57, 0xe6, 0x25, 0xbb, 0x38, 0xba, 0x0b, 0x79, 0x70, 0x71, 0xe7, 0x26, 0x66, 0xe6, 0xf9,
	0xe4, 0x27, 0x27, 0x96, 0xe4, 0x17, 0x49, 0x30, 0x2a, 0x30, 0x6b, 0x70, 0x1b, 0x89, 0xe8, 0xa1,
	0xba, 0x49, 0xcf, 0x23, 0xb1, 0x38, 0xc3, 0x49, 0xe0, 0xd3, 0x3d, 0x79, 0x9e, 0xe2, 0xe2, 0x2a,
	0xdd, 0xdc, 0xc4, 0x0a, 0x2b, 0x25, 0x63, 0x23, 0xa5, 0x20, 0x64, 0xad, 0x42, 0x56, 0x5c, 0x5c,
	0x60, 0xab, 0x83, 0x41, 0x56, 0x4b, 0x30, 0x29, 0x30, 0x6a, 0x70, 0x1b, 0x49, 0xa1, 0x1b, 0xe4,
	0x0e, 0x57, 0x11, 0x84, 0xa4, 0x5a, 0x69, 0x1f, 0x23, 0x17, 0x5b, 0x70, 0x69, 0x12, 0xc8, 0x41,
	0x46, 0x5c, 0x9c, 0xa0, 0xf0, 0x09, 0xc8, 0xcf, 0xcc, 0x2b, 0x91, 0x60, 0x04, 0x9b, 0x82, 0xd5,
	0x39, 0x41, 0x08, 0x65, 0x94, 0x58, 0x2d, 0xe4, 0xc2, 0xc5, 0x06, 0x09, 0x1c, 0x09, 0x66, 0xb0,
	0xdf, 0x25, 0xd1, 0xf5, 0x39, 0x81, 0x64, 0x5d, 0x12, 0x4b, 0x12, 0x9d, 0x04, 0x3f, 0xdd, 0x93,
	0xe7, 0x85, 0x07, 0x80, 0xa9, 0x81, 0x81, 0x52, 0x10, 0x54, 0xaf, 0x93, 0xc0, 0x89, 0x47, 0x72,
	0x8c, 0x17, 0x1e, 0xc9, 0x31, 0x3e, 0x78, 0x24, 0xc7, 0x38, 0xe3, 0xb1, 0x1c, 0x43, 0x12, 0x1b,
	0x38, 0xac, 0x8d, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0xa8, 0x46, 0xa2, 0x42, 0xee, 0x01, 0x00,
	0x00,
}

func (m *SyncDAG) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SyncDAG) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SyncDAG) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.GraphState != nil {
		{
			size, err := m.GraphState.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSyncdag(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.MainLocator) > 0 {
		for iNdEx := len(m.MainLocator) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.MainLocator[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintSyncdag(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *SubDAG) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SubDAG) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SubDAG) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Blocks) > 0 {
		for iNdEx := len(m.Blocks) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Blocks[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintSyncdag(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if m.GraphState != nil {
		{
			size, err := m.GraphState.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSyncdag(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.SyncPoint != nil {
		{
			size, err := m.SyncPoint.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSyncdag(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintSyncdag(dAtA []byte, offset int, v uint64) int {
	offset -= sovSyncdag(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *SyncDAG) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.MainLocator) > 0 {
		for _, e := range m.MainLocator {
			l = e.Size()
			n += 1 + l + sovSyncdag(uint64(l))
		}
	}
	if m.GraphState != nil {
		l = m.GraphState.Size()
		n += 1 + l + sovSyncdag(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *SubDAG) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.SyncPoint != nil {
		l = m.SyncPoint.Size()
		n += 1 + l + sovSyncdag(uint64(l))
	}
	if m.GraphState != nil {
		l = m.GraphState.Size()
		n += 1 + l + sovSyncdag(uint64(l))
	}
	if len(m.Blocks) > 0 {
		for _, e := range m.Blocks {
			l = e.Size()
			n += 1 + l + sovSyncdag(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovSyncdag(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozSyncdag(x uint64) (n int) {
	return sovSyncdag(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *SyncDAG) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSyncdag
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
			return fmt.Errorf("proto: SyncDAG: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SyncDAG: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MainLocator", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSyncdag
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSyncdag
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSyncdag
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MainLocator = append(m.MainLocator, &Hash{})
			if err := m.MainLocator[len(m.MainLocator)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GraphState", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSyncdag
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSyncdag
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSyncdag
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.GraphState == nil {
				m.GraphState = &GraphState{}
			}
			if err := m.GraphState.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSyncdag(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthSyncdag
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthSyncdag
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *SubDAG) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSyncdag
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
			return fmt.Errorf("proto: SubDAG: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SubDAG: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SyncPoint", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSyncdag
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSyncdag
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSyncdag
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.SyncPoint == nil {
				m.SyncPoint = &Hash{}
			}
			if err := m.SyncPoint.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GraphState", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSyncdag
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSyncdag
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSyncdag
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.GraphState == nil {
				m.GraphState = &GraphState{}
			}
			if err := m.GraphState.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Blocks", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSyncdag
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSyncdag
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSyncdag
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Blocks = append(m.Blocks, &BlockData{})
			if err := m.Blocks[len(m.Blocks)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSyncdag(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthSyncdag
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthSyncdag
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipSyncdag(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowSyncdag
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
					return 0, ErrIntOverflowSyncdag
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
					return 0, ErrIntOverflowSyncdag
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
				return 0, ErrInvalidLengthSyncdag
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupSyncdag
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthSyncdag
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthSyncdag        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowSyncdag          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupSyncdag = fmt.Errorf("proto: unexpected end of group")
)
