// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: storage/auth_provider.proto

package storage

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

// Next Tag: 9
type AuthProvider struct {
	Id         string            `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty" sql:"pk"`
	Name       string            `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty" sql:"unique"`
	Type       string            `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`
	UiEndpoint string            `protobuf:"bytes,4,opt,name=ui_endpoint,json=uiEndpoint,proto3" json:"ui_endpoint,omitempty"`
	Enabled    bool              `protobuf:"varint,5,opt,name=enabled,proto3" json:"enabled,omitempty"`
	Config     map[string]string `protobuf:"bytes,6,rep,name=config,proto3" json:"config,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// The login URL will be provided by the backend, and may not be specified in a request.
	LoginUrl  string `protobuf:"bytes,7,opt,name=login_url,json=loginUrl,proto3" json:"login_url,omitempty"`
	Validated bool   `protobuf:"varint,8,opt,name=validated,proto3" json:"validated,omitempty"` // Deprecated: Do not use.
	// UI endpoints which to allow in addition to `ui_endpoint`. I.e., if a login request
	// is coming from any of these, the auth request will use these for the callback URL,
	// not ui_endpoint.
	ExtraUiEndpoints []string `protobuf:"bytes,9,rep,name=extra_ui_endpoints,json=extraUiEndpoints,proto3" json:"extra_ui_endpoints,omitempty"`
	Active           bool     `protobuf:"varint,10,opt,name=active,proto3" json:"active,omitempty"`
	// EXPERIMENTAL.
	RequiredAttributes   []*AuthProvider_RequiredAttribute `protobuf:"bytes,11,rep,name=required_attributes,json=requiredAttributes,proto3" json:"required_attributes,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                          `json:"-"`
	XXX_unrecognized     []byte                            `json:"-"`
	XXX_sizecache        int32                             `json:"-"`
}

func (m *AuthProvider) Reset()         { *m = AuthProvider{} }
func (m *AuthProvider) String() string { return proto.CompactTextString(m) }
func (*AuthProvider) ProtoMessage()    {}
func (*AuthProvider) Descriptor() ([]byte, []int) {
	return fileDescriptor_4ed6b69aa5a381c8, []int{0}
}
func (m *AuthProvider) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AuthProvider) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AuthProvider.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AuthProvider) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AuthProvider.Merge(m, src)
}
func (m *AuthProvider) XXX_Size() int {
	return m.Size()
}
func (m *AuthProvider) XXX_DiscardUnknown() {
	xxx_messageInfo_AuthProvider.DiscardUnknown(m)
}

var xxx_messageInfo_AuthProvider proto.InternalMessageInfo

func (m *AuthProvider) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *AuthProvider) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *AuthProvider) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *AuthProvider) GetUiEndpoint() string {
	if m != nil {
		return m.UiEndpoint
	}
	return ""
}

func (m *AuthProvider) GetEnabled() bool {
	if m != nil {
		return m.Enabled
	}
	return false
}

func (m *AuthProvider) GetConfig() map[string]string {
	if m != nil {
		return m.Config
	}
	return nil
}

func (m *AuthProvider) GetLoginUrl() string {
	if m != nil {
		return m.LoginUrl
	}
	return ""
}

// Deprecated: Do not use.
func (m *AuthProvider) GetValidated() bool {
	if m != nil {
		return m.Validated
	}
	return false
}

func (m *AuthProvider) GetExtraUiEndpoints() []string {
	if m != nil {
		return m.ExtraUiEndpoints
	}
	return nil
}

func (m *AuthProvider) GetActive() bool {
	if m != nil {
		return m.Active
	}
	return false
}

func (m *AuthProvider) GetRequiredAttributes() []*AuthProvider_RequiredAttribute {
	if m != nil {
		return m.RequiredAttributes
	}
	return nil
}

func (m *AuthProvider) MessageClone() proto.Message {
	return m.Clone()
}
func (m *AuthProvider) Clone() *AuthProvider {
	if m == nil {
		return nil
	}
	cloned := new(AuthProvider)
	*cloned = *m

	if m.Config != nil {
		cloned.Config = make(map[string]string, len(m.Config))
		for k, v := range m.Config {
			cloned.Config[k] = v
		}
	}
	if m.ExtraUiEndpoints != nil {
		cloned.ExtraUiEndpoints = make([]string, len(m.ExtraUiEndpoints))
		copy(cloned.ExtraUiEndpoints, m.ExtraUiEndpoints)
	}
	if m.RequiredAttributes != nil {
		cloned.RequiredAttributes = make([]*AuthProvider_RequiredAttribute, len(m.RequiredAttributes))
		for idx, v := range m.RequiredAttributes {
			cloned.RequiredAttributes[idx] = v.Clone()
		}
	}
	return cloned
}

// RequiredAttribute allows to specify a set of attributes which ALL are required to be returned
// by the auth provider.
// If any attribute is missing within the external claims of the token issued by Central, the
// authentication request to this IdP is considered failed.
type AuthProvider_RequiredAttribute struct {
	AttributeKey         string   `protobuf:"bytes,1,opt,name=attribute_key,json=attributeKey,proto3" json:"attribute_key,omitempty"`
	AttributeValue       string   `protobuf:"bytes,2,opt,name=attribute_value,json=attributeValue,proto3" json:"attribute_value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AuthProvider_RequiredAttribute) Reset()         { *m = AuthProvider_RequiredAttribute{} }
func (m *AuthProvider_RequiredAttribute) String() string { return proto.CompactTextString(m) }
func (*AuthProvider_RequiredAttribute) ProtoMessage()    {}
func (*AuthProvider_RequiredAttribute) Descriptor() ([]byte, []int) {
	return fileDescriptor_4ed6b69aa5a381c8, []int{0, 1}
}
func (m *AuthProvider_RequiredAttribute) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AuthProvider_RequiredAttribute) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AuthProvider_RequiredAttribute.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AuthProvider_RequiredAttribute) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AuthProvider_RequiredAttribute.Merge(m, src)
}
func (m *AuthProvider_RequiredAttribute) XXX_Size() int {
	return m.Size()
}
func (m *AuthProvider_RequiredAttribute) XXX_DiscardUnknown() {
	xxx_messageInfo_AuthProvider_RequiredAttribute.DiscardUnknown(m)
}

var xxx_messageInfo_AuthProvider_RequiredAttribute proto.InternalMessageInfo

func (m *AuthProvider_RequiredAttribute) GetAttributeKey() string {
	if m != nil {
		return m.AttributeKey
	}
	return ""
}

func (m *AuthProvider_RequiredAttribute) GetAttributeValue() string {
	if m != nil {
		return m.AttributeValue
	}
	return ""
}

func (m *AuthProvider_RequiredAttribute) MessageClone() proto.Message {
	return m.Clone()
}
func (m *AuthProvider_RequiredAttribute) Clone() *AuthProvider_RequiredAttribute {
	if m == nil {
		return nil
	}
	cloned := new(AuthProvider_RequiredAttribute)
	*cloned = *m

	return cloned
}

func init() {
	proto.RegisterType((*AuthProvider)(nil), "storage.AuthProvider")
	proto.RegisterMapType((map[string]string)(nil), "storage.AuthProvider.ConfigEntry")
	proto.RegisterType((*AuthProvider_RequiredAttribute)(nil), "storage.AuthProvider.RequiredAttribute")
}

func init() { proto.RegisterFile("storage/auth_provider.proto", fileDescriptor_4ed6b69aa5a381c8) }

var fileDescriptor_4ed6b69aa5a381c8 = []byte{
	// 456 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x92, 0xcf, 0x6e, 0xd3, 0x40,
	0x10, 0xc6, 0xb1, 0x93, 0x26, 0xf1, 0x24, 0x40, 0x18, 0x2a, 0xb4, 0xa4, 0x28, 0x31, 0x01, 0xa9,
	0x39, 0x20, 0x57, 0x02, 0x0e, 0xb4, 0xb7, 0x06, 0xf5, 0xc4, 0x05, 0xad, 0x54, 0x84, 0xb8, 0x58,
	0x9b, 0x78, 0x71, 0x57, 0x31, 0x5e, 0x67, 0xbd, 0x1b, 0x35, 0x6f, 0xc2, 0x23, 0x71, 0xe4, 0xca,
	0xa5, 0x42, 0xe1, 0x0d, 0xfa, 0x04, 0xc8, 0x9b, 0xcd, 0x1f, 0x41, 0x6f, 0x33, 0xbf, 0x6f, 0x66,
	0x67, 0x3e, 0xed, 0xc0, 0x51, 0xa9, 0xa5, 0x62, 0x29, 0x3f, 0x61, 0x46, 0x5f, 0xc5, 0x85, 0x92,
	0x0b, 0x91, 0x70, 0x15, 0x15, 0x4a, 0x6a, 0x89, 0x4d, 0x27, 0xf6, 0x0e, 0x53, 0x99, 0x4a, 0xcb,
	0x4e, 0xaa, 0x68, 0x2d, 0x0f, 0x7f, 0xd5, 0xa1, 0x73, 0x6e, 0xf4, 0xd5, 0x47, 0xd7, 0x85, 0xcf,
	0xc0, 0x17, 0x09, 0xf1, 0x42, 0x6f, 0x14, 0x8c, 0x3b, 0xb7, 0x37, 0x83, 0x56, 0x39, 0xcf, 0xce,
	0x86, 0xc5, 0x6c, 0x48, 0x7d, 0x91, 0xe0, 0x4b, 0xa8, 0xe7, 0xec, 0x1b, 0x27, 0xbe, 0xd5, 0xbb,
	0xb7, 0x37, 0x83, 0x8e, 0xd5, 0x4d, 0x2e, 0xe6, 0x86, 0x0f, 0xa9, 0x55, 0x11, 0xa1, 0xae, 0x97,
	0x05, 0x27, 0xb5, 0xaa, 0x8a, 0xda, 0x18, 0x07, 0xd0, 0x36, 0x22, 0xe6, 0x79, 0x52, 0x48, 0x91,
	0x6b, 0x52, 0xb7, 0x12, 0x18, 0x71, 0xe1, 0x08, 0x12, 0x68, 0xf2, 0x9c, 0x4d, 0x32, 0x9e, 0x90,
	0x83, 0xd0, 0x1b, 0xb5, 0xe8, 0x26, 0xc5, 0x53, 0x68, 0x4c, 0x65, 0xfe, 0x55, 0xa4, 0xa4, 0x11,
	0xd6, 0x46, 0xed, 0xd7, 0xcf, 0x23, 0xe7, 0x29, 0xda, 0xdf, 0x3c, 0x7a, 0x6f, 0x6b, 0x2e, 0x72,
	0xad, 0x96, 0xd4, 0x35, 0xe0, 0x11, 0x04, 0x99, 0x4c, 0x45, 0x1e, 0x1b, 0x95, 0x91, 0xa6, 0x9d,
	0xd9, 0xb2, 0xe0, 0x52, 0x65, 0x18, 0x42, 0xb0, 0x60, 0x99, 0x48, 0x98, 0xe6, 0x09, 0x69, 0x55,
	0x33, 0xc7, 0x3e, 0xf1, 0xe8, 0x0e, 0xe2, 0x2b, 0x40, 0x7e, 0xad, 0x15, 0x8b, 0xf7, 0x56, 0x2f,
	0x49, 0x10, 0xd6, 0x46, 0x01, 0xed, 0x5a, 0xe5, 0x72, 0x6b, 0xa0, 0xc4, 0x27, 0xd0, 0x60, 0x53,
	0x2d, 0x16, 0x9c, 0x80, 0x35, 0xe0, 0x32, 0xfc, 0x0c, 0x8f, 0x15, 0x9f, 0x1b, 0xa1, 0x78, 0x12,
	0x33, 0xad, 0x95, 0x98, 0x18, 0xcd, 0x4b, 0xd2, 0xb6, 0x66, 0x8e, 0xef, 0x36, 0x43, 0x5d, 0xc3,
	0xf9, 0xa6, 0x9e, 0xa2, 0xfa, 0x17, 0x95, 0xbd, 0x53, 0x68, 0xef, 0xb9, 0xc6, 0x2e, 0xd4, 0x66,
	0x7c, 0xb9, 0xfe, 0x3c, 0x5a, 0x85, 0x78, 0x08, 0x07, 0x0b, 0x96, 0x19, 0xf7, 0x61, 0x74, 0x9d,
	0x9c, 0xf9, 0xef, 0xbc, 0x1e, 0x83, 0x47, 0xff, 0xcd, 0xc0, 0x17, 0x70, 0x7f, 0xbb, 0x60, 0xbc,
	0x7b, 0xaa, 0xb3, 0x85, 0x1f, 0xf8, 0x12, 0x8f, 0xe1, 0xe1, 0xae, 0x68, 0xff, 0xf5, 0x07, 0x5b,
	0xfc, 0xa9, 0xa2, 0xe3, 0xb7, 0x3f, 0x56, 0x7d, 0xef, 0xe7, 0xaa, 0xef, 0xfd, 0x5e, 0xf5, 0xbd,
	0xef, 0x7f, 0xfa, 0xf7, 0xe0, 0xa9, 0x90, 0x51, 0xa9, 0xd9, 0x74, 0xa6, 0xe4, 0xf5, 0xfa, 0x00,
	0x37, 0xee, 0xbf, 0x6c, 0xee, 0x74, 0xd2, 0xb0, 0xfc, 0xcd, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff,
	0x26, 0x13, 0x36, 0xf2, 0xd6, 0x02, 0x00, 0x00,
}

func (m *AuthProvider) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AuthProvider) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AuthProvider) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.RequiredAttributes) > 0 {
		for iNdEx := len(m.RequiredAttributes) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.RequiredAttributes[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintAuthProvider(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x5a
		}
	}
	if m.Active {
		i--
		if m.Active {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x50
	}
	if len(m.ExtraUiEndpoints) > 0 {
		for iNdEx := len(m.ExtraUiEndpoints) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.ExtraUiEndpoints[iNdEx])
			copy(dAtA[i:], m.ExtraUiEndpoints[iNdEx])
			i = encodeVarintAuthProvider(dAtA, i, uint64(len(m.ExtraUiEndpoints[iNdEx])))
			i--
			dAtA[i] = 0x4a
		}
	}
	if m.Validated {
		i--
		if m.Validated {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x40
	}
	if len(m.LoginUrl) > 0 {
		i -= len(m.LoginUrl)
		copy(dAtA[i:], m.LoginUrl)
		i = encodeVarintAuthProvider(dAtA, i, uint64(len(m.LoginUrl)))
		i--
		dAtA[i] = 0x3a
	}
	if len(m.Config) > 0 {
		for k := range m.Config {
			v := m.Config[k]
			baseI := i
			i -= len(v)
			copy(dAtA[i:], v)
			i = encodeVarintAuthProvider(dAtA, i, uint64(len(v)))
			i--
			dAtA[i] = 0x12
			i -= len(k)
			copy(dAtA[i:], k)
			i = encodeVarintAuthProvider(dAtA, i, uint64(len(k)))
			i--
			dAtA[i] = 0xa
			i = encodeVarintAuthProvider(dAtA, i, uint64(baseI-i))
			i--
			dAtA[i] = 0x32
		}
	}
	if m.Enabled {
		i--
		if m.Enabled {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x28
	}
	if len(m.UiEndpoint) > 0 {
		i -= len(m.UiEndpoint)
		copy(dAtA[i:], m.UiEndpoint)
		i = encodeVarintAuthProvider(dAtA, i, uint64(len(m.UiEndpoint)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Type) > 0 {
		i -= len(m.Type)
		copy(dAtA[i:], m.Type)
		i = encodeVarintAuthProvider(dAtA, i, uint64(len(m.Type)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintAuthProvider(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Id) > 0 {
		i -= len(m.Id)
		copy(dAtA[i:], m.Id)
		i = encodeVarintAuthProvider(dAtA, i, uint64(len(m.Id)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *AuthProvider_RequiredAttribute) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AuthProvider_RequiredAttribute) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AuthProvider_RequiredAttribute) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.AttributeValue) > 0 {
		i -= len(m.AttributeValue)
		copy(dAtA[i:], m.AttributeValue)
		i = encodeVarintAuthProvider(dAtA, i, uint64(len(m.AttributeValue)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.AttributeKey) > 0 {
		i -= len(m.AttributeKey)
		copy(dAtA[i:], m.AttributeKey)
		i = encodeVarintAuthProvider(dAtA, i, uint64(len(m.AttributeKey)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintAuthProvider(dAtA []byte, offset int, v uint64) int {
	offset -= sovAuthProvider(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *AuthProvider) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovAuthProvider(uint64(l))
	}
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovAuthProvider(uint64(l))
	}
	l = len(m.Type)
	if l > 0 {
		n += 1 + l + sovAuthProvider(uint64(l))
	}
	l = len(m.UiEndpoint)
	if l > 0 {
		n += 1 + l + sovAuthProvider(uint64(l))
	}
	if m.Enabled {
		n += 2
	}
	if len(m.Config) > 0 {
		for k, v := range m.Config {
			_ = k
			_ = v
			mapEntrySize := 1 + len(k) + sovAuthProvider(uint64(len(k))) + 1 + len(v) + sovAuthProvider(uint64(len(v)))
			n += mapEntrySize + 1 + sovAuthProvider(uint64(mapEntrySize))
		}
	}
	l = len(m.LoginUrl)
	if l > 0 {
		n += 1 + l + sovAuthProvider(uint64(l))
	}
	if m.Validated {
		n += 2
	}
	if len(m.ExtraUiEndpoints) > 0 {
		for _, s := range m.ExtraUiEndpoints {
			l = len(s)
			n += 1 + l + sovAuthProvider(uint64(l))
		}
	}
	if m.Active {
		n += 2
	}
	if len(m.RequiredAttributes) > 0 {
		for _, e := range m.RequiredAttributes {
			l = e.Size()
			n += 1 + l + sovAuthProvider(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *AuthProvider_RequiredAttribute) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.AttributeKey)
	if l > 0 {
		n += 1 + l + sovAuthProvider(uint64(l))
	}
	l = len(m.AttributeValue)
	if l > 0 {
		n += 1 + l + sovAuthProvider(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovAuthProvider(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozAuthProvider(x uint64) (n int) {
	return sovAuthProvider(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *AuthProvider) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAuthProvider
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
			return fmt.Errorf("proto: AuthProvider: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AuthProvider: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAuthProvider
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
				return ErrInvalidLengthAuthProvider
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAuthProvider
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAuthProvider
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
				return ErrInvalidLengthAuthProvider
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAuthProvider
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAuthProvider
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
				return ErrInvalidLengthAuthProvider
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAuthProvider
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Type = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UiEndpoint", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAuthProvider
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
				return ErrInvalidLengthAuthProvider
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAuthProvider
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.UiEndpoint = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Enabled", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAuthProvider
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
			m.Enabled = bool(v != 0)
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Config", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAuthProvider
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
				return ErrInvalidLengthAuthProvider
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAuthProvider
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Config == nil {
				m.Config = make(map[string]string)
			}
			var mapkey string
			var mapvalue string
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowAuthProvider
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
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowAuthProvider
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthAuthProvider
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey < 0 {
						return ErrInvalidLengthAuthProvider
					}
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var stringLenmapvalue uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowAuthProvider
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapvalue |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapvalue := int(stringLenmapvalue)
					if intStringLenmapvalue < 0 {
						return ErrInvalidLengthAuthProvider
					}
					postStringIndexmapvalue := iNdEx + intStringLenmapvalue
					if postStringIndexmapvalue < 0 {
						return ErrInvalidLengthAuthProvider
					}
					if postStringIndexmapvalue > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = string(dAtA[iNdEx:postStringIndexmapvalue])
					iNdEx = postStringIndexmapvalue
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipAuthProvider(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if (skippy < 0) || (iNdEx+skippy) < 0 {
						return ErrInvalidLengthAuthProvider
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Config[mapkey] = mapvalue
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LoginUrl", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAuthProvider
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
				return ErrInvalidLengthAuthProvider
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAuthProvider
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.LoginUrl = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Validated", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAuthProvider
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
			m.Validated = bool(v != 0)
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExtraUiEndpoints", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAuthProvider
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
				return ErrInvalidLengthAuthProvider
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAuthProvider
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ExtraUiEndpoints = append(m.ExtraUiEndpoints, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 10:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Active", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAuthProvider
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
			m.Active = bool(v != 0)
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RequiredAttributes", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAuthProvider
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
				return ErrInvalidLengthAuthProvider
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAuthProvider
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RequiredAttributes = append(m.RequiredAttributes, &AuthProvider_RequiredAttribute{})
			if err := m.RequiredAttributes[len(m.RequiredAttributes)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipAuthProvider(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthAuthProvider
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
func (m *AuthProvider_RequiredAttribute) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAuthProvider
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
			return fmt.Errorf("proto: RequiredAttribute: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RequiredAttribute: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AttributeKey", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAuthProvider
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
				return ErrInvalidLengthAuthProvider
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAuthProvider
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AttributeKey = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AttributeValue", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAuthProvider
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
				return ErrInvalidLengthAuthProvider
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAuthProvider
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AttributeValue = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipAuthProvider(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthAuthProvider
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
func skipAuthProvider(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowAuthProvider
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
					return 0, ErrIntOverflowAuthProvider
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
					return 0, ErrIntOverflowAuthProvider
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
				return 0, ErrInvalidLengthAuthProvider
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupAuthProvider
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthAuthProvider
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthAuthProvider        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowAuthProvider          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupAuthProvider = fmt.Errorf("proto: unexpected end of group")
)
