// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.28.0
// source: api/auth/messages.proto

package auth_proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type TokenVerificationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *TokenVerificationRequest) Reset() {
	*x = TokenVerificationRequest{}
	mi := &file_api_auth_messages_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TokenVerificationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TokenVerificationRequest) ProtoMessage() {}

func (x *TokenVerificationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_auth_messages_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TokenVerificationRequest.ProtoReflect.Descriptor instead.
func (*TokenVerificationRequest) Descriptor() ([]byte, []int) {
	return file_api_auth_messages_proto_rawDescGZIP(), []int{0}
}

func (x *TokenVerificationRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type TokenVerificationResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Verified       bool   `protobuf:"varint,1,opt,name=verified,proto3" json:"verified,omitempty"`
	UserId         string `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	ExpirationTime int64  `protobuf:"varint,3,opt,name=expiration_time,json=expirationTime,proto3" json:"expiration_time,omitempty"`
}

func (x *TokenVerificationResponse) Reset() {
	*x = TokenVerificationResponse{}
	mi := &file_api_auth_messages_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TokenVerificationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TokenVerificationResponse) ProtoMessage() {}

func (x *TokenVerificationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_auth_messages_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TokenVerificationResponse.ProtoReflect.Descriptor instead.
func (*TokenVerificationResponse) Descriptor() ([]byte, []int) {
	return file_api_auth_messages_proto_rawDescGZIP(), []int{1}
}

func (x *TokenVerificationResponse) GetVerified() bool {
	if x != nil {
		return x.Verified
	}
	return false
}

func (x *TokenVerificationResponse) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *TokenVerificationResponse) GetExpirationTime() int64 {
	if x != nil {
		return x.ExpirationTime
	}
	return 0
}

var File_api_auth_messages_proto protoreflect.FileDescriptor

var file_api_auth_messages_proto_rawDesc = []byte{
	0x0a, 0x17, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x61, 0x75, 0x74, 0x68, 0x5f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x30, 0x0a, 0x18, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x56, 0x65,
	0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x79, 0x0a, 0x19, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x64,
	0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x27, 0x0a, 0x0f, 0x65, 0x78, 0x70,
	0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x0e, 0x65, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69,
	0x6d, 0x65, 0x42, 0x14, 0x5a, 0x12, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x75,
	0x74, 0x68, 0x2d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_auth_messages_proto_rawDescOnce sync.Once
	file_api_auth_messages_proto_rawDescData = file_api_auth_messages_proto_rawDesc
)

func file_api_auth_messages_proto_rawDescGZIP() []byte {
	file_api_auth_messages_proto_rawDescOnce.Do(func() {
		file_api_auth_messages_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_auth_messages_proto_rawDescData)
	})
	return file_api_auth_messages_proto_rawDescData
}

var file_api_auth_messages_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_api_auth_messages_proto_goTypes = []any{
	(*TokenVerificationRequest)(nil),  // 0: auth_proto.TokenVerificationRequest
	(*TokenVerificationResponse)(nil), // 1: auth_proto.TokenVerificationResponse
}
var file_api_auth_messages_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_auth_messages_proto_init() }
func file_api_auth_messages_proto_init() {
	if File_api_auth_messages_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_auth_messages_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_auth_messages_proto_goTypes,
		DependencyIndexes: file_api_auth_messages_proto_depIdxs,
		MessageInfos:      file_api_auth_messages_proto_msgTypes,
	}.Build()
	File_api_auth_messages_proto = out.File
	file_api_auth_messages_proto_rawDesc = nil
	file_api_auth_messages_proto_goTypes = nil
	file_api_auth_messages_proto_depIdxs = nil
}
