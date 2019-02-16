// Code generated by protoc-gen-go. DO NOT EDIT.
// source: rexecd.proto

package rexecd

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
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

type KeyType int32

const (
	KeyType_ssh_rsa             KeyType = 0
	KeyType_rsa_sha2_256        KeyType = 1
	KeyType_rsa_sha2_512        KeyType = 2
	KeyType_ssh_dss             KeyType = 3
	KeyType_ecdsa_sha2_nistp256 KeyType = 4
	KeyType_ecdsa_sha2_nistp384 KeyType = 5
	KeyType_ecdsa_sha2_nistp521 KeyType = 6
	KeyType_ssh_ed25519         KeyType = 7
)

var KeyType_name = map[int32]string{
	0: "ssh_rsa",
	1: "rsa_sha2_256",
	2: "rsa_sha2_512",
	3: "ssh_dss",
	4: "ecdsa_sha2_nistp256",
	5: "ecdsa_sha2_nistp384",
	6: "ecdsa_sha2_nistp521",
	7: "ssh_ed25519",
}

var KeyType_value = map[string]int32{
	"ssh_rsa":             0,
	"rsa_sha2_256":        1,
	"rsa_sha2_512":        2,
	"ssh_dss":             3,
	"ecdsa_sha2_nistp256": 4,
	"ecdsa_sha2_nistp384": 5,
	"ecdsa_sha2_nistp521": 6,
	"ssh_ed25519":         7,
}

func (x KeyType) String() string {
	return proto.EnumName(KeyType_name, int32(x))
}

func (KeyType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_a5a97cdcf024307f, []int{0}
}

type RegisterHostRequest struct {
	Fqdn                 string   `protobuf:"bytes,1,opt,name=fqdn,proto3" json:"fqdn,omitempty"`
	Port                 string   `protobuf:"bytes,2,opt,name=port,proto3" json:"port,omitempty"`
	PrivateKey           []byte   `protobuf:"bytes,3,opt,name=private_key,json=privateKey,proto3" json:"private_key,omitempty"`
	PublicKey            []byte   `protobuf:"bytes,4,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
	KeyType              KeyType  `protobuf:"varint,5,opt,name=key_type,json=keyType,proto3,enum=rexecd.KeyType" json:"key_type,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RegisterHostRequest) Reset()         { *m = RegisterHostRequest{} }
func (m *RegisterHostRequest) String() string { return proto.CompactTextString(m) }
func (*RegisterHostRequest) ProtoMessage()    {}
func (*RegisterHostRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a5a97cdcf024307f, []int{0}
}

func (m *RegisterHostRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterHostRequest.Unmarshal(m, b)
}
func (m *RegisterHostRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterHostRequest.Marshal(b, m, deterministic)
}
func (m *RegisterHostRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterHostRequest.Merge(m, src)
}
func (m *RegisterHostRequest) XXX_Size() int {
	return xxx_messageInfo_RegisterHostRequest.Size(m)
}
func (m *RegisterHostRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterHostRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterHostRequest proto.InternalMessageInfo

func (m *RegisterHostRequest) GetFqdn() string {
	if m != nil {
		return m.Fqdn
	}
	return ""
}

func (m *RegisterHostRequest) GetPort() string {
	if m != nil {
		return m.Port
	}
	return ""
}

func (m *RegisterHostRequest) GetPrivateKey() []byte {
	if m != nil {
		return m.PrivateKey
	}
	return nil
}

func (m *RegisterHostRequest) GetPublicKey() []byte {
	if m != nil {
		return m.PublicKey
	}
	return nil
}

func (m *RegisterHostRequest) GetKeyType() KeyType {
	if m != nil {
		return m.KeyType
	}
	return KeyType_ssh_rsa
}

type HostConnect struct {
	Fqdn                 string   `protobuf:"bytes,1,opt,name=fqdn,proto3" json:"fqdn,omitempty"`
	Port                 string   `protobuf:"bytes,2,opt,name=port,proto3" json:"port,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HostConnect) Reset()         { *m = HostConnect{} }
func (m *HostConnect) String() string { return proto.CompactTextString(m) }
func (*HostConnect) ProtoMessage()    {}
func (*HostConnect) Descriptor() ([]byte, []int) {
	return fileDescriptor_a5a97cdcf024307f, []int{1}
}

func (m *HostConnect) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HostConnect.Unmarshal(m, b)
}
func (m *HostConnect) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HostConnect.Marshal(b, m, deterministic)
}
func (m *HostConnect) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HostConnect.Merge(m, src)
}
func (m *HostConnect) XXX_Size() int {
	return xxx_messageInfo_HostConnect.Size(m)
}
func (m *HostConnect) XXX_DiscardUnknown() {
	xxx_messageInfo_HostConnect.DiscardUnknown(m)
}

var xxx_messageInfo_HostConnect proto.InternalMessageInfo

func (m *HostConnect) GetFqdn() string {
	if m != nil {
		return m.Fqdn
	}
	return ""
}

func (m *HostConnect) GetPort() string {
	if m != nil {
		return m.Port
	}
	return ""
}

type RegisterHostResponse struct {
	Id                   int64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RegisterHostResponse) Reset()         { *m = RegisterHostResponse{} }
func (m *RegisterHostResponse) String() string { return proto.CompactTextString(m) }
func (*RegisterHostResponse) ProtoMessage()    {}
func (*RegisterHostResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a5a97cdcf024307f, []int{2}
}

func (m *RegisterHostResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterHostResponse.Unmarshal(m, b)
}
func (m *RegisterHostResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterHostResponse.Marshal(b, m, deterministic)
}
func (m *RegisterHostResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterHostResponse.Merge(m, src)
}
func (m *RegisterHostResponse) XXX_Size() int {
	return xxx_messageInfo_RegisterHostResponse.Size(m)
}
func (m *RegisterHostResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterHostResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterHostResponse proto.InternalMessageInfo

func (m *RegisterHostResponse) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

type RegisterUserRequest struct {
	Username             string   `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	FirstName            string   `protobuf:"bytes,2,opt,name=first_name,json=firstName,proto3" json:"first_name,omitempty"`
	LastName             string   `protobuf:"bytes,3,opt,name=last_name,json=lastName,proto3" json:"last_name,omitempty"`
	Admin                bool     `protobuf:"varint,4,opt,name=admin,proto3" json:"admin,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RegisterUserRequest) Reset()         { *m = RegisterUserRequest{} }
func (m *RegisterUserRequest) String() string { return proto.CompactTextString(m) }
func (*RegisterUserRequest) ProtoMessage()    {}
func (*RegisterUserRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a5a97cdcf024307f, []int{3}
}

func (m *RegisterUserRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterUserRequest.Unmarshal(m, b)
}
func (m *RegisterUserRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterUserRequest.Marshal(b, m, deterministic)
}
func (m *RegisterUserRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterUserRequest.Merge(m, src)
}
func (m *RegisterUserRequest) XXX_Size() int {
	return xxx_messageInfo_RegisterUserRequest.Size(m)
}
func (m *RegisterUserRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterUserRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterUserRequest proto.InternalMessageInfo

func (m *RegisterUserRequest) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *RegisterUserRequest) GetFirstName() string {
	if m != nil {
		return m.FirstName
	}
	return ""
}

func (m *RegisterUserRequest) GetLastName() string {
	if m != nil {
		return m.LastName
	}
	return ""
}

func (m *RegisterUserRequest) GetAdmin() bool {
	if m != nil {
		return m.Admin
	}
	return false
}

type RegisterUserResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RegisterUserResponse) Reset()         { *m = RegisterUserResponse{} }
func (m *RegisterUserResponse) String() string { return proto.CompactTextString(m) }
func (*RegisterUserResponse) ProtoMessage()    {}
func (*RegisterUserResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a5a97cdcf024307f, []int{4}
}

func (m *RegisterUserResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterUserResponse.Unmarshal(m, b)
}
func (m *RegisterUserResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterUserResponse.Marshal(b, m, deterministic)
}
func (m *RegisterUserResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterUserResponse.Merge(m, src)
}
func (m *RegisterUserResponse) XXX_Size() int {
	return xxx_messageInfo_RegisterUserResponse.Size(m)
}
func (m *RegisterUserResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterUserResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterUserResponse proto.InternalMessageInfo

type CommandRequest struct {
	Id                   int64             `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Cmd                  string            `protobuf:"bytes,2,opt,name=cmd,proto3" json:"cmd,omitempty"`
	Username             string            `protobuf:"bytes,3,opt,name=username,proto3" json:"username,omitempty"`
	PrivateKey           []byte            `protobuf:"bytes,4,opt,name=private_key,json=privateKey,proto3" json:"private_key,omitempty"`
	HostConnect          []*HostConnect    `protobuf:"bytes,5,rep,name=host_connect,json=hostConnect,proto3" json:"host_connect,omitempty"`
	Env                  map[string]string `protobuf:"bytes,7,rep,name=env,proto3" json:"env,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Distribute           bool              `protobuf:"varint,8,opt,name=distribute,proto3" json:"distribute,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *CommandRequest) Reset()         { *m = CommandRequest{} }
func (m *CommandRequest) String() string { return proto.CompactTextString(m) }
func (*CommandRequest) ProtoMessage()    {}
func (*CommandRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a5a97cdcf024307f, []int{5}
}

func (m *CommandRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommandRequest.Unmarshal(m, b)
}
func (m *CommandRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommandRequest.Marshal(b, m, deterministic)
}
func (m *CommandRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommandRequest.Merge(m, src)
}
func (m *CommandRequest) XXX_Size() int {
	return xxx_messageInfo_CommandRequest.Size(m)
}
func (m *CommandRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CommandRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CommandRequest proto.InternalMessageInfo

func (m *CommandRequest) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *CommandRequest) GetCmd() string {
	if m != nil {
		return m.Cmd
	}
	return ""
}

func (m *CommandRequest) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *CommandRequest) GetPrivateKey() []byte {
	if m != nil {
		return m.PrivateKey
	}
	return nil
}

func (m *CommandRequest) GetHostConnect() []*HostConnect {
	if m != nil {
		return m.HostConnect
	}
	return nil
}

func (m *CommandRequest) GetEnv() map[string]string {
	if m != nil {
		return m.Env
	}
	return nil
}

func (m *CommandRequest) GetDistribute() bool {
	if m != nil {
		return m.Distribute
	}
	return false
}

type CommandResponse struct {
	Id                   int64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Stdout               []byte   `protobuf:"bytes,2,opt,name=stdout,proto3" json:"stdout,omitempty"`
	Stderr               []byte   `protobuf:"bytes,3,opt,name=stderr,proto3" json:"stderr,omitempty"`
	Timestamp            int64    `protobuf:"varint,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	ExitCode             int64    `protobuf:"varint,5,opt,name=exit_code,json=exitCode,proto3" json:"exit_code,omitempty"`
	ErrorMsg             string   `protobuf:"bytes,6,opt,name=error_msg,json=errorMsg,proto3" json:"error_msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CommandResponse) Reset()         { *m = CommandResponse{} }
func (m *CommandResponse) String() string { return proto.CompactTextString(m) }
func (*CommandResponse) ProtoMessage()    {}
func (*CommandResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a5a97cdcf024307f, []int{6}
}

func (m *CommandResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommandResponse.Unmarshal(m, b)
}
func (m *CommandResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommandResponse.Marshal(b, m, deterministic)
}
func (m *CommandResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommandResponse.Merge(m, src)
}
func (m *CommandResponse) XXX_Size() int {
	return xxx_messageInfo_CommandResponse.Size(m)
}
func (m *CommandResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CommandResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CommandResponse proto.InternalMessageInfo

func (m *CommandResponse) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *CommandResponse) GetStdout() []byte {
	if m != nil {
		return m.Stdout
	}
	return nil
}

func (m *CommandResponse) GetStderr() []byte {
	if m != nil {
		return m.Stderr
	}
	return nil
}

func (m *CommandResponse) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *CommandResponse) GetExitCode() int64 {
	if m != nil {
		return m.ExitCode
	}
	return 0
}

func (m *CommandResponse) GetErrorMsg() string {
	if m != nil {
		return m.ErrorMsg
	}
	return ""
}

func init() {
	proto.RegisterEnum("rexecd.KeyType", KeyType_name, KeyType_value)
	proto.RegisterType((*RegisterHostRequest)(nil), "rexecd.RegisterHostRequest")
	proto.RegisterType((*HostConnect)(nil), "rexecd.HostConnect")
	proto.RegisterType((*RegisterHostResponse)(nil), "rexecd.RegisterHostResponse")
	proto.RegisterType((*RegisterUserRequest)(nil), "rexecd.RegisterUserRequest")
	proto.RegisterType((*RegisterUserResponse)(nil), "rexecd.RegisterUserResponse")
	proto.RegisterType((*CommandRequest)(nil), "rexecd.CommandRequest")
	proto.RegisterMapType((map[string]string)(nil), "rexecd.CommandRequest.EnvEntry")
	proto.RegisterType((*CommandResponse)(nil), "rexecd.CommandResponse")
}

func init() { proto.RegisterFile("rexecd.proto", fileDescriptor_a5a97cdcf024307f) }

var fileDescriptor_a5a97cdcf024307f = []byte{
	// 634 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x54, 0xcf, 0x4f, 0xdb, 0x30,
	0x14, 0x26, 0x4d, 0x7f, 0xbe, 0x56, 0x10, 0x19, 0x04, 0x51, 0x61, 0xa3, 0xea, 0x61, 0xaa, 0x38,
	0xa0, 0x35, 0xac, 0x88, 0xed, 0xb0, 0x0b, 0x42, 0x9a, 0x84, 0xb6, 0x83, 0xb5, 0x9d, 0xa3, 0x50,
	0x3f, 0xa8, 0x05, 0xf9, 0x81, 0xed, 0x56, 0xe4, 0xbc, 0x3f, 0x65, 0xa7, 0x5d, 0xf7, 0x27, 0xed,
	0x2f, 0x99, 0x6c, 0x27, 0x90, 0xd2, 0x1c, 0x76, 0xf3, 0xfb, 0x3e, 0xfb, 0xe5, 0xf3, 0xf7, 0x3d,
	0x07, 0x06, 0x02, 0x9f, 0x70, 0xce, 0x4e, 0x33, 0x91, 0xaa, 0x94, 0xb4, 0x6d, 0x35, 0xfe, 0xed,
	0xc0, 0x2e, 0xc5, 0x3b, 0x2e, 0x15, 0x8a, 0x2f, 0xa9, 0x54, 0x14, 0x1f, 0x97, 0x28, 0x15, 0x21,
	0xd0, 0xbc, 0x7d, 0x64, 0x89, 0xef, 0x8c, 0x9c, 0x49, 0x8f, 0x9a, 0xb5, 0xc6, 0xb2, 0x54, 0x28,
	0xbf, 0x61, 0x31, 0xbd, 0x26, 0xc7, 0xd0, 0xcf, 0x04, 0x5f, 0x45, 0x0a, 0xc3, 0x7b, 0xcc, 0x7d,
	0x77, 0xe4, 0x4c, 0x06, 0x14, 0x0a, 0xe8, 0x1a, 0x73, 0xf2, 0x06, 0x20, 0x5b, 0xde, 0x3c, 0xf0,
	0xb9, 0xe1, 0x9b, 0x86, 0xef, 0x59, 0x44, 0xd3, 0x27, 0xd0, 0xbd, 0xc7, 0x3c, 0x54, 0x79, 0x86,
	0x7e, 0x6b, 0xe4, 0x4c, 0xb6, 0x83, 0x9d, 0xd3, 0x42, 0xe8, 0x35, 0xe6, 0xdf, 0xf3, 0x0c, 0x69,
	0xe7, 0xde, 0x2e, 0xc6, 0x33, 0xe8, 0x6b, 0x89, 0x97, 0x69, 0x92, 0xe0, 0xfc, 0xbf, 0x25, 0x8e,
	0xdf, 0xc1, 0xde, 0xfa, 0x0d, 0x65, 0x96, 0x26, 0x12, 0xc9, 0x36, 0x34, 0x38, 0x33, 0xa7, 0x5d,
	0xda, 0xe0, 0x6c, 0xfc, 0xb3, 0x62, 0xc5, 0x0f, 0x89, 0xa2, 0xb4, 0x62, 0x08, 0xdd, 0xa5, 0x44,
	0x91, 0x44, 0x31, 0x16, 0xdf, 0x7a, 0xae, 0xf5, 0xed, 0x6e, 0xb9, 0x90, 0x2a, 0x34, 0xac, 0xfd,
	0x6a, 0xcf, 0x20, 0xdf, 0x34, 0x7d, 0x08, 0xbd, 0x87, 0xa8, 0x64, 0x5d, 0x7b, 0x56, 0x03, 0x86,
	0xdc, 0x83, 0x56, 0xc4, 0x62, 0x9e, 0x18, 0x53, 0xba, 0xd4, 0x16, 0xe3, 0xfd, 0x17, 0xb5, 0x56,
	0x84, 0x55, 0x3b, 0xfe, 0xd3, 0x80, 0xed, 0xcb, 0x34, 0x8e, 0xa3, 0x84, 0x95, 0xc2, 0x5e, 0x5d,
	0x80, 0x78, 0xe0, 0xce, 0x63, 0x56, 0xa8, 0xd0, 0xcb, 0x35, 0xe9, 0xee, 0x2b, 0xe9, 0xaf, 0x92,
	0x6b, 0x6e, 0x24, 0x77, 0x0e, 0x83, 0x45, 0x2a, 0x55, 0x38, 0xb7, 0x7e, 0xfb, 0xad, 0x91, 0x3b,
	0xe9, 0x07, 0xbb, 0x65, 0x3c, 0x95, 0x28, 0x68, 0x7f, 0x51, 0xc9, 0x65, 0x0a, 0x2e, 0x26, 0x2b,
	0xbf, 0x63, 0xb6, 0x1f, 0x97, 0xdb, 0xd7, 0xb5, 0x9f, 0x5e, 0x25, 0xab, 0xab, 0x44, 0x89, 0x9c,
	0xea, 0xbd, 0xe4, 0x2d, 0x00, 0xe3, 0x52, 0x09, 0x7e, 0xb3, 0x54, 0xe8, 0x77, 0x8d, 0x1f, 0x15,
	0x64, 0x78, 0x0e, 0xdd, 0xf2, 0x80, 0xbe, 0xa5, 0xd6, 0x6b, 0x93, 0xd0, 0x4b, 0x6d, 0xe4, 0x2a,
	0x7a, 0x58, 0x96, 0xfe, 0xdb, 0xe2, 0x53, 0xe3, 0xc2, 0xd1, 0xd3, 0xbd, 0xf3, 0xfc, 0xe1, 0xfa,
	0xd8, 0xc9, 0x3e, 0xb4, 0xa5, 0x62, 0xe9, 0xd2, 0x0e, 0xcd, 0x80, 0x16, 0x55, 0x81, 0xa3, 0x10,
	0xc5, 0x50, 0x17, 0x15, 0x39, 0x82, 0x9e, 0xe2, 0x31, 0x4a, 0x15, 0xc5, 0x99, 0x71, 0xcd, 0xa5,
	0x2f, 0x80, 0x4e, 0x1c, 0x9f, 0xb8, 0x36, 0x8d, 0xd9, 0x81, 0x76, 0x69, 0x57, 0x03, 0x97, 0x29,
	0x33, 0xe3, 0x80, 0x42, 0xa4, 0x22, 0x8c, 0xe5, 0x9d, 0xdf, 0xb6, 0x79, 0x18, 0xe0, 0xab, 0xbc,
	0x3b, 0xf9, 0xe5, 0x40, 0xa7, 0x18, 0x79, 0xd2, 0x87, 0x8e, 0x94, 0x8b, 0x50, 0xc8, 0xc8, 0xdb,
	0x22, 0x1e, 0x0c, 0x84, 0x8c, 0x42, 0xb9, 0x88, 0x82, 0x30, 0x98, 0x9d, 0x7b, 0xce, 0x1a, 0x32,
	0x9b, 0x06, 0x5e, 0xa3, 0x3c, 0xc0, 0xa4, 0xf4, 0x5c, 0x72, 0x00, 0xbb, 0x38, 0x67, 0xe5, 0x86,
	0x84, 0x4b, 0x95, 0xe9, 0x73, 0xcd, 0x3a, 0xe2, 0xec, 0xe2, 0x83, 0xd7, 0xaa, 0x23, 0x66, 0xc1,
	0xd4, 0x6b, 0x93, 0x1d, 0xe8, 0xeb, 0xbe, 0xc8, 0x82, 0xd9, 0x6c, 0xfa, 0xd1, 0xeb, 0x04, 0x7f,
	0x1d, 0x68, 0x53, 0x93, 0x28, 0xf9, 0x0c, 0x9d, 0xc2, 0x5b, 0xb2, 0x5f, 0x9f, 0xf2, 0xf0, 0x60,
	0x03, 0x2f, 0xa6, 0x79, 0xeb, 0xbd, 0x43, 0xae, 0x61, 0x50, 0x7d, 0x97, 0xe4, 0xb0, 0xdc, 0x5c,
	0xf3, 0x3f, 0x1a, 0x1e, 0xd5, 0x93, 0x65, 0xbb, 0x6a, 0x33, 0xfd, 0x6c, 0x36, 0x9b, 0x55, 0x5e,
	0xf4, 0x66, 0xb3, 0xb5, 0x97, 0xb6, 0x75, 0xd3, 0x36, 0xff, 0xc8, 0xb3, 0x7f, 0x01, 0x00, 0x00,
	0xff, 0xff, 0x9a, 0xae, 0xcc, 0x93, 0x33, 0x05, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// RexecdClient is the client API for Rexecd service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RexecdClient interface {
	Command(ctx context.Context, in *CommandRequest, opts ...grpc.CallOption) (Rexecd_CommandClient, error)
	RegisterHost(ctx context.Context, in *RegisterHostRequest, opts ...grpc.CallOption) (*RegisterHostResponse, error)
	RegisterUser(ctx context.Context, in *RegisterUserRequest, opts ...grpc.CallOption) (*RegisterUserResponse, error)
}

type rexecdClient struct {
	cc *grpc.ClientConn
}

func NewRexecdClient(cc *grpc.ClientConn) RexecdClient {
	return &rexecdClient{cc}
}

func (c *rexecdClient) Command(ctx context.Context, in *CommandRequest, opts ...grpc.CallOption) (Rexecd_CommandClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Rexecd_serviceDesc.Streams[0], "/rexecd.Rexecd/Command", opts...)
	if err != nil {
		return nil, err
	}
	x := &rexecdCommandClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Rexecd_CommandClient interface {
	Recv() (*CommandResponse, error)
	grpc.ClientStream
}

type rexecdCommandClient struct {
	grpc.ClientStream
}

func (x *rexecdCommandClient) Recv() (*CommandResponse, error) {
	m := new(CommandResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *rexecdClient) RegisterHost(ctx context.Context, in *RegisterHostRequest, opts ...grpc.CallOption) (*RegisterHostResponse, error) {
	out := new(RegisterHostResponse)
	err := c.cc.Invoke(ctx, "/rexecd.Rexecd/RegisterHost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rexecdClient) RegisterUser(ctx context.Context, in *RegisterUserRequest, opts ...grpc.CallOption) (*RegisterUserResponse, error) {
	out := new(RegisterUserResponse)
	err := c.cc.Invoke(ctx, "/rexecd.Rexecd/RegisterUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RexecdServer is the server API for Rexecd service.
type RexecdServer interface {
	Command(*CommandRequest, Rexecd_CommandServer) error
	RegisterHost(context.Context, *RegisterHostRequest) (*RegisterHostResponse, error)
	RegisterUser(context.Context, *RegisterUserRequest) (*RegisterUserResponse, error)
}

func RegisterRexecdServer(s *grpc.Server, srv RexecdServer) {
	s.RegisterService(&_Rexecd_serviceDesc, srv)
}

func _Rexecd_Command_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(CommandRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(RexecdServer).Command(m, &rexecdCommandServer{stream})
}

type Rexecd_CommandServer interface {
	Send(*CommandResponse) error
	grpc.ServerStream
}

type rexecdCommandServer struct {
	grpc.ServerStream
}

func (x *rexecdCommandServer) Send(m *CommandResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Rexecd_RegisterHost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterHostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RexecdServer).RegisterHost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rexecd.Rexecd/RegisterHost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RexecdServer).RegisterHost(ctx, req.(*RegisterHostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rexecd_RegisterUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RexecdServer).RegisterUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rexecd.Rexecd/RegisterUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RexecdServer).RegisterUser(ctx, req.(*RegisterUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Rexecd_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rexecd.Rexecd",
	HandlerType: (*RexecdServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterHost",
			Handler:    _Rexecd_RegisterHost_Handler,
		},
		{
			MethodName: "RegisterUser",
			Handler:    _Rexecd_RegisterUser_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Command",
			Handler:       _Rexecd_Command_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "rexecd.proto",
}
