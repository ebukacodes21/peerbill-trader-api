// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v5.28.3
// source: rpc_reset.proto

package pb

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

type ResetPasswordRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token    string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *ResetPasswordRequest) Reset() {
	*x = ResetPasswordRequest{}
	mi := &file_rpc_reset_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ResetPasswordRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResetPasswordRequest) ProtoMessage() {}

func (x *ResetPasswordRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_reset_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResetPasswordRequest.ProtoReflect.Descriptor instead.
func (*ResetPasswordRequest) Descriptor() ([]byte, []int) {
	return file_rpc_reset_proto_rawDescGZIP(), []int{0}
}

func (x *ResetPasswordRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *ResetPasswordRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type ResetPasswordResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *ResetPasswordResponse) Reset() {
	*x = ResetPasswordResponse{}
	mi := &file_rpc_reset_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ResetPasswordResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResetPasswordResponse) ProtoMessage() {}

func (x *ResetPasswordResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_reset_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResetPasswordResponse.ProtoReflect.Descriptor instead.
func (*ResetPasswordResponse) Descriptor() ([]byte, []int) {
	return file_rpc_reset_proto_rawDescGZIP(), []int{1}
}

func (x *ResetPasswordResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_rpc_reset_proto protoreflect.FileDescriptor

var file_rpc_reset_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x72, 0x70, 0x63, 0x5f, 0x72, 0x65, 0x73, 0x65, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x02, 0x70, 0x62, 0x22, 0x48, 0x0a, 0x14, 0x52, 0x65, 0x73, 0x65, 0x74, 0x50, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a,
	0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f,
	0x6b, 0x65, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x22,
	0x31, 0x0a, 0x15, 0x52, 0x65, 0x73, 0x65, 0x74, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x42, 0x30, 0x5a, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x65, 0x62, 0x75, 0x6b, 0x61, 0x63, 0x6f, 0x64, 0x65, 0x73, 0x32, 0x31, 0x2f, 0x70, 0x65,
	0x65, 0x72, 0x62, 0x69, 0x6c, 0x6c, 0x2d, 0x74, 0x72, 0x61, 0x64, 0x65, 0x72, 0x2d, 0x61, 0x70,
	0x69, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rpc_reset_proto_rawDescOnce sync.Once
	file_rpc_reset_proto_rawDescData = file_rpc_reset_proto_rawDesc
)

func file_rpc_reset_proto_rawDescGZIP() []byte {
	file_rpc_reset_proto_rawDescOnce.Do(func() {
		file_rpc_reset_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc_reset_proto_rawDescData)
	})
	return file_rpc_reset_proto_rawDescData
}

var file_rpc_reset_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_rpc_reset_proto_goTypes = []any{
	(*ResetPasswordRequest)(nil),  // 0: pb.ResetPasswordRequest
	(*ResetPasswordResponse)(nil), // 1: pb.ResetPasswordResponse
}
var file_rpc_reset_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_rpc_reset_proto_init() }
func file_rpc_reset_proto_init() {
	if File_rpc_reset_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_rpc_reset_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_rpc_reset_proto_goTypes,
		DependencyIndexes: file_rpc_reset_proto_depIdxs,
		MessageInfos:      file_rpc_reset_proto_msgTypes,
	}.Build()
	File_rpc_reset_proto = out.File
	file_rpc_reset_proto_rawDesc = nil
	file_rpc_reset_proto_goTypes = nil
	file_rpc_reset_proto_depIdxs = nil
}
