// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.27.1
// source: sso/v1/sso.proto

package sso

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RegisterUserRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Email         string                 `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	Password      string                 `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RegisterUserRequest) Reset() {
	*x = RegisterUserRequest{}
	mi := &file_sso_v1_sso_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RegisterUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterUserRequest) ProtoMessage() {}

func (x *RegisterUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_sso_v1_sso_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterUserRequest.ProtoReflect.Descriptor instead.
func (*RegisterUserRequest) Descriptor() ([]byte, []int) {
	return file_sso_v1_sso_proto_rawDescGZIP(), []int{0}
}

func (x *RegisterUserRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *RegisterUserRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type RegisterUserResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        int64                  `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RegisterUserResponse) Reset() {
	*x = RegisterUserResponse{}
	mi := &file_sso_v1_sso_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RegisterUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterUserResponse) ProtoMessage() {}

func (x *RegisterUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_sso_v1_sso_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterUserResponse.ProtoReflect.Descriptor instead.
func (*RegisterUserResponse) Descriptor() ([]byte, []int) {
	return file_sso_v1_sso_proto_rawDescGZIP(), []int{1}
}

func (x *RegisterUserResponse) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type LoginUserRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Email         string                 `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	Password      string                 `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LoginUserRequest) Reset() {
	*x = LoginUserRequest{}
	mi := &file_sso_v1_sso_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LoginUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginUserRequest) ProtoMessage() {}

func (x *LoginUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_sso_v1_sso_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginUserRequest.ProtoReflect.Descriptor instead.
func (*LoginUserRequest) Descriptor() ([]byte, []int) {
	return file_sso_v1_sso_proto_rawDescGZIP(), []int{2}
}

func (x *LoginUserRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *LoginUserRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type LoginUserResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Token         string                 `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	UserId        int64                  `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LoginUserResponse) Reset() {
	*x = LoginUserResponse{}
	mi := &file_sso_v1_sso_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LoginUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginUserResponse) ProtoMessage() {}

func (x *LoginUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_sso_v1_sso_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginUserResponse.ProtoReflect.Descriptor instead.
func (*LoginUserResponse) Descriptor() ([]byte, []int) {
	return file_sso_v1_sso_proto_rawDescGZIP(), []int{3}
}

func (x *LoginUserResponse) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *LoginUserResponse) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type VerifyTokenRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Token         string                 `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *VerifyTokenRequest) Reset() {
	*x = VerifyTokenRequest{}
	mi := &file_sso_v1_sso_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *VerifyTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifyTokenRequest) ProtoMessage() {}

func (x *VerifyTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_sso_v1_sso_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VerifyTokenRequest.ProtoReflect.Descriptor instead.
func (*VerifyTokenRequest) Descriptor() ([]byte, []int) {
	return file_sso_v1_sso_proto_rawDescGZIP(), []int{4}
}

func (x *VerifyTokenRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type VerifyTokenResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        int64                  `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *VerifyTokenResponse) Reset() {
	*x = VerifyTokenResponse{}
	mi := &file_sso_v1_sso_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *VerifyTokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifyTokenResponse) ProtoMessage() {}

func (x *VerifyTokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_sso_v1_sso_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VerifyTokenResponse.ProtoReflect.Descriptor instead.
func (*VerifyTokenResponse) Descriptor() ([]byte, []int) {
	return file_sso_v1_sso_proto_rawDescGZIP(), []int{5}
}

func (x *VerifyTokenResponse) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

var File_sso_v1_sso_proto protoreflect.FileDescriptor

const file_sso_v1_sso_proto_rawDesc = "" +
	"\n" +
	"\x10sso/v1/sso.proto\x12\x03sso\x1a\x1cgoogle/api/annotations.proto\x1a\x1bgoogle/protobuf/empty.proto\x1a\x17validate/validate.proto\"a\n" +
	"\x13RegisterUserRequest\x12\"\n" +
	"\x05email\x18\x01 \x01(\tB\f\xfaB\tr\a\x10\x05\x18\xfe\x01`\x01R\x05email\x12&\n" +
	"\bpassword\x18\x02 \x01(\tB\n" +
	"\xfaB\ar\x05\x10\x06\x18\x80\x01R\bpassword\"/\n" +
	"\x14RegisterUserResponse\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\x03R\x06userId\"^\n" +
	"\x10LoginUserRequest\x12\"\n" +
	"\x05email\x18\x01 \x01(\tB\f\xfaB\tr\a\x10\x05\x18\xfe\x01`\x01R\x05email\x12&\n" +
	"\bpassword\x18\x02 \x01(\tB\n" +
	"\xfaB\ar\x05\x10\x06\x18\x80\x01R\bpassword\"B\n" +
	"\x11LoginUserResponse\x12\x14\n" +
	"\x05token\x18\x01 \x01(\tR\x05token\x12\x17\n" +
	"\auser_id\x18\x02 \x01(\x03R\x06userId\"*\n" +
	"\x12VerifyTokenRequest\x12\x14\n" +
	"\x05token\x18\x01 \x01(\tR\x05token\".\n" +
	"\x13VerifyTokenResponse\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\x03R\x06userId2\x8b\x02\n" +
	"\x03SSO\x12\\\n" +
	"\fRegisterUser\x12\x18.sso.RegisterUserRequest\x1a\x19.sso.RegisterUserResponse\"\x17\x82\xd3\xe4\x93\x02\x11:\x01*\"\f/v1/register\x12P\n" +
	"\tLoginUser\x12\x15.sso.LoginUserRequest\x1a\x16.sso.LoginUserResponse\"\x14\x82\xd3\xe4\x93\x02\x0e:\x01*\"\t/v1/login\x12T\n" +
	"\vVerifyToken\x12\x17.sso.VerifyTokenRequest\x1a\x18.sso.VerifyTokenResponse\"\x12\x82\xd3\xe4\x93\x02\f\x12\n" +
	"/v1/verifyB1Z/github.com/iamvkosarev/sso/pkg/proto/sso/v1;ssob\x06proto3"

var (
	file_sso_v1_sso_proto_rawDescOnce sync.Once
	file_sso_v1_sso_proto_rawDescData []byte
)

func file_sso_v1_sso_proto_rawDescGZIP() []byte {
	file_sso_v1_sso_proto_rawDescOnce.Do(func() {
		file_sso_v1_sso_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_sso_v1_sso_proto_rawDesc), len(file_sso_v1_sso_proto_rawDesc)))
	})
	return file_sso_v1_sso_proto_rawDescData
}

var file_sso_v1_sso_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_sso_v1_sso_proto_goTypes = []any{
	(*RegisterUserRequest)(nil),  // 0: sso.RegisterUserRequest
	(*RegisterUserResponse)(nil), // 1: sso.RegisterUserResponse
	(*LoginUserRequest)(nil),     // 2: sso.LoginUserRequest
	(*LoginUserResponse)(nil),    // 3: sso.LoginUserResponse
	(*VerifyTokenRequest)(nil),   // 4: sso.VerifyTokenRequest
	(*VerifyTokenResponse)(nil),  // 5: sso.VerifyTokenResponse
}
var file_sso_v1_sso_proto_depIdxs = []int32{
	0, // 0: sso.SSO.RegisterUser:input_type -> sso.RegisterUserRequest
	2, // 1: sso.SSO.LoginUser:input_type -> sso.LoginUserRequest
	4, // 2: sso.SSO.VerifyToken:input_type -> sso.VerifyTokenRequest
	1, // 3: sso.SSO.RegisterUser:output_type -> sso.RegisterUserResponse
	3, // 4: sso.SSO.LoginUser:output_type -> sso.LoginUserResponse
	5, // 5: sso.SSO.VerifyToken:output_type -> sso.VerifyTokenResponse
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_sso_v1_sso_proto_init() }
func file_sso_v1_sso_proto_init() {
	if File_sso_v1_sso_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_sso_v1_sso_proto_rawDesc), len(file_sso_v1_sso_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_sso_v1_sso_proto_goTypes,
		DependencyIndexes: file_sso_v1_sso_proto_depIdxs,
		MessageInfos:      file_sso_v1_sso_proto_msgTypes,
	}.Build()
	File_sso_v1_sso_proto = out.File
	file_sso_v1_sso_proto_goTypes = nil
	file_sso_v1_sso_proto_depIdxs = nil
}
