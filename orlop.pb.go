// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: orlop.proto

package orlop

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
	reflect "reflect"
	strings "strings"
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

// Redirect represents a redirection to a new location
type Redirect struct {
	// Location to redirect to
	Location string `protobuf:"bytes,1,opt,name=location,proto3" json:"location,omitempty"`
}

func (m *Redirect) Reset()      { *m = Redirect{} }
func (*Redirect) ProtoMessage() {}
func (*Redirect) Descriptor() ([]byte, []int) {
	return fileDescriptor_dd37fa353e556da6, []int{0}
}
func (m *Redirect) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Redirect) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Redirect.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Redirect) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Redirect.Merge(m, src)
}
func (m *Redirect) XXX_Size() int {
	return m.Size()
}
func (m *Redirect) XXX_DiscardUnknown() {
	xxx_messageInfo_Redirect.DiscardUnknown(m)
}

var xxx_messageInfo_Redirect proto.InternalMessageInfo

func (m *Redirect) GetLocation() string {
	if m != nil {
		return m.Location
	}
	return ""
}

// ErrorMessage represents an error message
type ErrorMessage struct {
	// Code description
	Code int32 `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	// Error description
	Error string `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
	// Message description
	Message string `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
}

func (m *ErrorMessage) Reset()      { *m = ErrorMessage{} }
func (*ErrorMessage) ProtoMessage() {}
func (*ErrorMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_dd37fa353e556da6, []int{1}
}
func (m *ErrorMessage) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ErrorMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ErrorMessage.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ErrorMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ErrorMessage.Merge(m, src)
}
func (m *ErrorMessage) XXX_Size() int {
	return m.Size()
}
func (m *ErrorMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_ErrorMessage.DiscardUnknown(m)
}

var xxx_messageInfo_ErrorMessage proto.InternalMessageInfo

func (m *ErrorMessage) GetCode() int32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *ErrorMessage) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func (m *ErrorMessage) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*Redirect)(nil), "orlop.Redirect")
	proto.RegisterType((*ErrorMessage)(nil), "orlop.ErrorMessage")
}

func init() { proto.RegisterFile("orlop.proto", fileDescriptor_dd37fa353e556da6) }

var fileDescriptor_dd37fa353e556da6 = []byte{
	// 215 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xce, 0x2f, 0xca, 0xc9,
	0x2f, 0xd0, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x73, 0x94, 0xd4, 0xb8, 0x38, 0x82,
	0x52, 0x53, 0x32, 0x8b, 0x52, 0x93, 0x4b, 0x84, 0xa4, 0xb8, 0x38, 0x72, 0xf2, 0x93, 0x13, 0x4b,
	0x32, 0xf3, 0xf3, 0x24, 0x18, 0x15, 0x18, 0x35, 0x38, 0x83, 0xe0, 0x7c, 0xa5, 0x20, 0x2e, 0x1e,
	0xd7, 0xa2, 0xa2, 0xfc, 0x22, 0xdf, 0xd4, 0xe2, 0xe2, 0xc4, 0xf4, 0x54, 0x21, 0x21, 0x2e, 0x96,
	0xe4, 0xfc, 0x94, 0x54, 0xb0, 0x3a, 0xd6, 0x20, 0x30, 0x5b, 0x48, 0x84, 0x8b, 0x35, 0x15, 0xa4,
	0x46, 0x82, 0x09, 0xac, 0x19, 0xc2, 0x11, 0x92, 0xe0, 0x62, 0xcf, 0x85, 0x68, 0x92, 0x60, 0x06,
	0x8b, 0xc3, 0xb8, 0x4e, 0x41, 0x17, 0x1e, 0xca, 0x31, 0xdc, 0x78, 0x28, 0xc7, 0xf0, 0xe1, 0xa1,
	0x1c, 0x63, 0xc3, 0x23, 0x39, 0xc6, 0x15, 0x8f, 0xe4, 0x18, 0x4f, 0x3c, 0x92, 0x63, 0xbc, 0xf0,
	0x48, 0x8e, 0xf1, 0xc1, 0x23, 0x39, 0xc6, 0x17, 0x8f, 0xe4, 0x18, 0x3e, 0x3c, 0x92, 0x63, 0x9c,
	0xf0, 0x58, 0x8e, 0xe1, 0xc2, 0x63, 0x39, 0x86, 0x1b, 0x8f, 0xe5, 0x18, 0xa2, 0x64, 0xd2, 0xf3,
	0xf5, 0xb2, 0x53, 0x4b, 0x92, 0x33, 0xf4, 0x92, 0xf3, 0x73, 0xf5, 0x73, 0x32, 0x93, 0xf4, 0xc1,
	0x3e, 0xb1, 0x06, 0x93, 0x49, 0x6c, 0x60, 0xdf, 0x19, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x08,
	0x11, 0x7c, 0x43, 0xec, 0x00, 0x00, 0x00,
}

func (this *Redirect) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Redirect)
	if !ok {
		that2, ok := that.(Redirect)
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
	if this.Location != that1.Location {
		return false
	}
	return true
}
func (this *ErrorMessage) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*ErrorMessage)
	if !ok {
		that2, ok := that.(ErrorMessage)
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
	if this.Code != that1.Code {
		return false
	}
	if this.Error != that1.Error {
		return false
	}
	if this.Message != that1.Message {
		return false
	}
	return true
}
func (this *Redirect) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 5)
	s = append(s, "&orlop.Redirect{")
	s = append(s, "Location: "+fmt.Sprintf("%#v", this.Location)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *ErrorMessage) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 7)
	s = append(s, "&orlop.ErrorMessage{")
	s = append(s, "Code: "+fmt.Sprintf("%#v", this.Code)+",\n")
	s = append(s, "Error: "+fmt.Sprintf("%#v", this.Error)+",\n")
	s = append(s, "Message: "+fmt.Sprintf("%#v", this.Message)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringOrlop(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}
func (m *Redirect) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Redirect) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Redirect) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Location) > 0 {
		i -= len(m.Location)
		copy(dAtA[i:], m.Location)
		i = encodeVarintOrlop(dAtA, i, uint64(len(m.Location)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *ErrorMessage) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ErrorMessage) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ErrorMessage) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Message) > 0 {
		i -= len(m.Message)
		copy(dAtA[i:], m.Message)
		i = encodeVarintOrlop(dAtA, i, uint64(len(m.Message)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Error) > 0 {
		i -= len(m.Error)
		copy(dAtA[i:], m.Error)
		i = encodeVarintOrlop(dAtA, i, uint64(len(m.Error)))
		i--
		dAtA[i] = 0x12
	}
	if m.Code != 0 {
		i = encodeVarintOrlop(dAtA, i, uint64(m.Code))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintOrlop(dAtA []byte, offset int, v uint64) int {
	offset -= sovOrlop(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Redirect) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Location)
	if l > 0 {
		n += 1 + l + sovOrlop(uint64(l))
	}
	return n
}

func (m *ErrorMessage) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Code != 0 {
		n += 1 + sovOrlop(uint64(m.Code))
	}
	l = len(m.Error)
	if l > 0 {
		n += 1 + l + sovOrlop(uint64(l))
	}
	l = len(m.Message)
	if l > 0 {
		n += 1 + l + sovOrlop(uint64(l))
	}
	return n
}

func sovOrlop(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozOrlop(x uint64) (n int) {
	return sovOrlop(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *Redirect) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Redirect{`,
		`Location:` + fmt.Sprintf("%v", this.Location) + `,`,
		`}`,
	}, "")
	return s
}
func (this *ErrorMessage) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&ErrorMessage{`,
		`Code:` + fmt.Sprintf("%v", this.Code) + `,`,
		`Error:` + fmt.Sprintf("%v", this.Error) + `,`,
		`Message:` + fmt.Sprintf("%v", this.Message) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringOrlop(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *Redirect) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowOrlop
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
			return fmt.Errorf("proto: Redirect: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Redirect: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Location", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrlop
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
				return ErrInvalidLengthOrlop
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthOrlop
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Location = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipOrlop(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthOrlop
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthOrlop
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
func (m *ErrorMessage) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowOrlop
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
			return fmt.Errorf("proto: ErrorMessage: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ErrorMessage: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Code", wireType)
			}
			m.Code = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrlop
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Code |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Error", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrlop
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
				return ErrInvalidLengthOrlop
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthOrlop
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Error = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Message", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrlop
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
				return ErrInvalidLengthOrlop
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthOrlop
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Message = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipOrlop(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthOrlop
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthOrlop
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
func skipOrlop(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowOrlop
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
					return 0, ErrIntOverflowOrlop
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
					return 0, ErrIntOverflowOrlop
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
				return 0, ErrInvalidLengthOrlop
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupOrlop
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthOrlop
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthOrlop        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowOrlop          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupOrlop = fmt.Errorf("proto: unexpected end of group")
)
