// Code generated by protoc-gen-go. DO NOT EDIT.
// source: user.proto

package api

import (
	fmt "fmt"
	proto "github.com/lightbrotherV/gin-protobuf/proto"
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

type LoginReq struct {
	// 账号
	Account string `protobuf:"bytes,1,opt,name=account,proto3" json:"account,omitempty"`
	// 密码
	Password             string   `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LoginReq) Reset()         { *m = LoginReq{} }
func (m *LoginReq) String() string { return proto.CompactTextString(m) }
func (*LoginReq) ProtoMessage()    {}
func (*LoginReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_116e343673f7ffaf, []int{0}
}

func (m *LoginReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoginReq.Unmarshal(m, b)
}
func (m *LoginReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoginReq.Marshal(b, m, deterministic)
}
func (m *LoginReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoginReq.Merge(m, src)
}
func (m *LoginReq) XXX_Size() int {
	return xxx_messageInfo_LoginReq.Size(m)
}
func (m *LoginReq) XXX_DiscardUnknown() {
	xxx_messageInfo_LoginReq.DiscardUnknown(m)
}

var xxx_messageInfo_LoginReq proto.InternalMessageInfo

func (m *LoginReq) GetAccount() string {
	if m != nil {
		return m.Account
	}
	return ""
}

func (m *LoginReq) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type LoginResp struct {
	// 登录态token
	Token                string   `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LoginResp) Reset()         { *m = LoginResp{} }
func (m *LoginResp) String() string { return proto.CompactTextString(m) }
func (*LoginResp) ProtoMessage()    {}
func (*LoginResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_116e343673f7ffaf, []int{1}
}

func (m *LoginResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoginResp.Unmarshal(m, b)
}
func (m *LoginResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoginResp.Marshal(b, m, deterministic)
}
func (m *LoginResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoginResp.Merge(m, src)
}
func (m *LoginResp) XXX_Size() int {
	return xxx_messageInfo_LoginResp.Size(m)
}
func (m *LoginResp) XXX_DiscardUnknown() {
	xxx_messageInfo_LoginResp.DiscardUnknown(m)
}

var xxx_messageInfo_LoginResp proto.InternalMessageInfo

func (m *LoginResp) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

type LogoutReq struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LogoutReq) Reset()         { *m = LogoutReq{} }
func (m *LogoutReq) String() string { return proto.CompactTextString(m) }
func (*LogoutReq) ProtoMessage()    {}
func (*LogoutReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_116e343673f7ffaf, []int{2}
}

func (m *LogoutReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LogoutReq.Unmarshal(m, b)
}
func (m *LogoutReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LogoutReq.Marshal(b, m, deterministic)
}
func (m *LogoutReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogoutReq.Merge(m, src)
}
func (m *LogoutReq) XXX_Size() int {
	return xxx_messageInfo_LogoutReq.Size(m)
}
func (m *LogoutReq) XXX_DiscardUnknown() {
	xxx_messageInfo_LogoutReq.DiscardUnknown(m)
}

var xxx_messageInfo_LogoutReq proto.InternalMessageInfo

type LogoutResp struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LogoutResp) Reset()         { *m = LogoutResp{} }
func (m *LogoutResp) String() string { return proto.CompactTextString(m) }
func (*LogoutResp) ProtoMessage()    {}
func (*LogoutResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_116e343673f7ffaf, []int{3}
}

func (m *LogoutResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LogoutResp.Unmarshal(m, b)
}
func (m *LogoutResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LogoutResp.Marshal(b, m, deterministic)
}
func (m *LogoutResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogoutResp.Merge(m, src)
}
func (m *LogoutResp) XXX_Size() int {
	return xxx_messageInfo_LogoutResp.Size(m)
}
func (m *LogoutResp) XXX_DiscardUnknown() {
	xxx_messageInfo_LogoutResp.DiscardUnknown(m)
}

var xxx_messageInfo_LogoutResp proto.InternalMessageInfo

type InfoReq struct {
	// 登录态token
	Token                string   `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InfoReq) Reset()         { *m = InfoReq{} }
func (m *InfoReq) String() string { return proto.CompactTextString(m) }
func (*InfoReq) ProtoMessage()    {}
func (*InfoReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_116e343673f7ffaf, []int{4}
}

func (m *InfoReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InfoReq.Unmarshal(m, b)
}
func (m *InfoReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InfoReq.Marshal(b, m, deterministic)
}
func (m *InfoReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InfoReq.Merge(m, src)
}
func (m *InfoReq) XXX_Size() int {
	return xxx_messageInfo_InfoReq.Size(m)
}
func (m *InfoReq) XXX_DiscardUnknown() {
	xxx_messageInfo_InfoReq.DiscardUnknown(m)
}

var xxx_messageInfo_InfoReq proto.InternalMessageInfo

func (m *InfoReq) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

type InfoResp struct {
	// 账号
	Account string `protobuf:"bytes,1,opt,name=account,proto3" json:"account,omitempty"`
	// 密码
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	// 昵称
	Nickname string `protobuf:"bytes,3,opt,name=nickname,proto3" json:"nickname,omitempty"`
	// 权限
	Roles []string `protobuf:"bytes,4,rep,name=roles,proto3" json:"roles,omitempty"`
	// 头像
	Avatar               string   `protobuf:"bytes,5,opt,name=avatar,proto3" json:"avatar,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InfoResp) Reset()         { *m = InfoResp{} }
func (m *InfoResp) String() string { return proto.CompactTextString(m) }
func (*InfoResp) ProtoMessage()    {}
func (*InfoResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_116e343673f7ffaf, []int{5}
}

func (m *InfoResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InfoResp.Unmarshal(m, b)
}
func (m *InfoResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InfoResp.Marshal(b, m, deterministic)
}
func (m *InfoResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InfoResp.Merge(m, src)
}
func (m *InfoResp) XXX_Size() int {
	return xxx_messageInfo_InfoResp.Size(m)
}
func (m *InfoResp) XXX_DiscardUnknown() {
	xxx_messageInfo_InfoResp.DiscardUnknown(m)
}

var xxx_messageInfo_InfoResp proto.InternalMessageInfo

func (m *InfoResp) GetAccount() string {
	if m != nil {
		return m.Account
	}
	return ""
}

func (m *InfoResp) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *InfoResp) GetNickname() string {
	if m != nil {
		return m.Nickname
	}
	return ""
}

func (m *InfoResp) GetRoles() []string {
	if m != nil {
		return m.Roles
	}
	return nil
}

func (m *InfoResp) GetAvatar() string {
	if m != nil {
		return m.Avatar
	}
	return ""
}

type UpdateReq struct {
	// 账号
	Account string `protobuf:"bytes,1,opt,name=account,proto3" json:"account,omitempty"`
	// 密码
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	// 昵称
	Nickname string `protobuf:"bytes,3,opt,name=nickname,proto3" json:"nickname,omitempty"`
	// 权限
	Roles []string `protobuf:"bytes,4,rep,name=roles,proto3" json:"roles,omitempty"`
	// 头像
	Avatar               string   `protobuf:"bytes,5,opt,name=avatar,proto3" json:"avatar,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateReq) Reset()         { *m = UpdateReq{} }
func (m *UpdateReq) String() string { return proto.CompactTextString(m) }
func (*UpdateReq) ProtoMessage()    {}
func (*UpdateReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_116e343673f7ffaf, []int{6}
}

func (m *UpdateReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateReq.Unmarshal(m, b)
}
func (m *UpdateReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateReq.Marshal(b, m, deterministic)
}
func (m *UpdateReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateReq.Merge(m, src)
}
func (m *UpdateReq) XXX_Size() int {
	return xxx_messageInfo_UpdateReq.Size(m)
}
func (m *UpdateReq) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateReq.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateReq proto.InternalMessageInfo

func (m *UpdateReq) GetAccount() string {
	if m != nil {
		return m.Account
	}
	return ""
}

func (m *UpdateReq) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *UpdateReq) GetNickname() string {
	if m != nil {
		return m.Nickname
	}
	return ""
}

func (m *UpdateReq) GetRoles() []string {
	if m != nil {
		return m.Roles
	}
	return nil
}

func (m *UpdateReq) GetAvatar() string {
	if m != nil {
		return m.Avatar
	}
	return ""
}

type UpdateResp struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateResp) Reset()         { *m = UpdateResp{} }
func (m *UpdateResp) String() string { return proto.CompactTextString(m) }
func (*UpdateResp) ProtoMessage()    {}
func (*UpdateResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_116e343673f7ffaf, []int{7}
}

func (m *UpdateResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateResp.Unmarshal(m, b)
}
func (m *UpdateResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateResp.Marshal(b, m, deterministic)
}
func (m *UpdateResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateResp.Merge(m, src)
}
func (m *UpdateResp) XXX_Size() int {
	return xxx_messageInfo_UpdateResp.Size(m)
}
func (m *UpdateResp) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateResp.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateResp proto.InternalMessageInfo

func init() {
	proto.RegisterType((*LoginReq)(nil), "most.simple.mcd.LoginReq")
	proto.RegisterType((*LoginResp)(nil), "most.simple.mcd.LoginResp")
	proto.RegisterType((*LogoutReq)(nil), "most.simple.mcd.LogoutReq")
	proto.RegisterType((*LogoutResp)(nil), "most.simple.mcd.LogoutResp")
	proto.RegisterType((*InfoReq)(nil), "most.simple.mcd.InfoReq")
	proto.RegisterType((*InfoResp)(nil), "most.simple.mcd.InfoResp")
	proto.RegisterType((*UpdateReq)(nil), "most.simple.mcd.UpdateReq")
	proto.RegisterType((*UpdateResp)(nil), "most.simple.mcd.UpdateResp")
}

func init() {
	proto.RegisterFile("user.proto", fileDescriptor_116e343673f7ffaf)
}

var fileDescriptor_116e343673f7ffaf = []byte{
	// 312 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x92, 0x31, 0x4f, 0xc3, 0x30,
	0x10, 0x85, 0x95, 0xb6, 0x49, 0x93, 0xa3, 0x12, 0x92, 0x85, 0x90, 0x1b, 0x06, 0x8a, 0x27, 0xa6,
	0x0c, 0x30, 0x22, 0x21, 0x60, 0x43, 0x62, 0x8a, 0xd4, 0x85, 0xcd, 0x24, 0x2e, 0x8a, 0x9a, 0xd8,
	0xae, 0xcf, 0x81, 0xdf, 0x00, 0xfc, 0x69, 0x14, 0x27, 0x29, 0x12, 0x24, 0x4b, 0xc5, 0xf8, 0xe9,
	0xbd, 0xdc, 0xbb, 0xdc, 0x33, 0x40, 0x8d, 0xc2, 0x24, 0xda, 0x28, 0xab, 0xc8, 0x71, 0xa5, 0xd0,
	0x26, 0x58, 0x54, 0xba, 0x14, 0x49, 0x95, 0xe5, 0xec, 0x0e, 0xc2, 0x27, 0xf5, 0x5a, 0xc8, 0x54,
	0xec, 0x08, 0x85, 0x39, 0xcf, 0x32, 0x55, 0x4b, 0x4b, 0xbd, 0x95, 0x77, 0x19, 0xa5, 0x3d, 0x92,
	0x18, 0x42, 0xcd, 0x11, 0xdf, 0x95, 0xc9, 0xe9, 0xc4, 0x49, 0x7b, 0x66, 0x17, 0x10, 0x75, 0x13,
	0x50, 0x93, 0x13, 0xf0, 0xad, 0xda, 0x0a, 0xd9, 0x0d, 0x68, 0x81, 0x1d, 0x39, 0x8b, 0xaa, 0x6d,
	0x2a, 0x76, 0x6c, 0x01, 0xd0, 0x03, 0x6a, 0x76, 0x0e, 0xf3, 0x47, 0xb9, 0x51, 0x4d, 0xfc, 0xf0,
	0xb7, 0x9f, 0x1e, 0x84, 0xad, 0x03, 0xf5, 0x61, 0x1b, 0x36, 0x9a, 0x2c, 0xb2, 0xad, 0xe4, 0x95,
	0xa0, 0xd3, 0x56, 0xeb, 0xb9, 0x09, 0x35, 0xaa, 0x14, 0x48, 0x67, 0xab, 0x69, 0x13, 0xea, 0x80,
	0x9c, 0x42, 0xc0, 0xdf, 0xb8, 0xe5, 0x86, 0xfa, 0xce, 0xdf, 0x11, 0xfb, 0xf2, 0x20, 0x5a, 0xeb,
	0x9c, 0x5b, 0x71, 0xf0, 0xbd, 0xfe, 0x71, 0x9b, 0x05, 0x40, 0xbf, 0x0c, 0xea, 0xab, 0x8f, 0x09,
	0xcc, 0xd6, 0x28, 0x0c, 0xb9, 0x05, 0xbf, 0x6c, 0x0a, 0x21, 0xcb, 0xe4, 0x57, 0xdb, 0x49, 0x5f,
	0x75, 0x1c, 0x8f, 0x49, 0xa8, 0xc9, 0x3d, 0x04, 0xa5, 0x2b, 0x88, 0x0c, 0xba, 0xda, 0x1a, 0xe3,
	0xb3, 0x51, 0x0d, 0x35, 0xb9, 0x81, 0x59, 0x21, 0x37, 0x8a, 0xd0, 0x3f, 0xa6, 0xae, 0xec, 0x78,
	0x39, 0xa2, 0xb4, 0xf9, 0xb5, 0xfb, 0xad, 0x81, 0xfc, 0xfd, 0xf1, 0x07, 0xf2, 0x7f, 0x6e, 0xf1,
	0xe0, 0x3f, 0x4f, 0xb9, 0x2e, 0x5e, 0x02, 0xf7, 0xe8, 0xaf, 0xbf, 0x03, 0x00, 0x00, 0xff, 0xff,
	0x0e, 0x68, 0xeb, 0x3c, 0x02, 0x03, 0x00, 0x00,
}
