// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v5.28.3
// source: rpc_create_order.proto

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

type CreateOrderRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EscrowAddress string  `protobuf:"bytes,1,opt,name=escrow_address,json=escrowAddress,proto3" json:"escrow_address,omitempty"`
	Crypto        string  `protobuf:"bytes,2,opt,name=crypto,proto3" json:"crypto,omitempty"`
	Fiat          string  `protobuf:"bytes,3,opt,name=fiat,proto3" json:"fiat,omitempty"`
	FiatAmount    float32 `protobuf:"fixed32,4,opt,name=fiat_amount,json=fiatAmount,proto3" json:"fiat_amount,omitempty"`
	CryptoAmount  float32 `protobuf:"fixed32,5,opt,name=crypto_amount,json=cryptoAmount,proto3" json:"crypto_amount,omitempty"`
	Username      string  `protobuf:"bytes,6,opt,name=username,proto3" json:"username,omitempty"`
	Rate          float32 `protobuf:"fixed32,7,opt,name=rate,proto3" json:"rate,omitempty"`
	UserAddress   string  `protobuf:"bytes,8,opt,name=user_address,json=userAddress,proto3" json:"user_address,omitempty"`
	OrderType     string  `protobuf:"bytes,9,opt,name=order_type,json=orderType,proto3" json:"order_type,omitempty"`
	BankName      *string `protobuf:"bytes,10,opt,name=bank_name,json=bankName,proto3,oneof" json:"bank_name,omitempty"`
	AccountNumber *string `protobuf:"bytes,11,opt,name=account_number,json=accountNumber,proto3,oneof" json:"account_number,omitempty"`
	AccountHolder *string `protobuf:"bytes,12,opt,name=account_holder,json=accountHolder,proto3,oneof" json:"account_holder,omitempty"`
}

func (x *CreateOrderRequest) Reset() {
	*x = CreateOrderRequest{}
	mi := &file_rpc_create_order_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateOrderRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateOrderRequest) ProtoMessage() {}

func (x *CreateOrderRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_create_order_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateOrderRequest.ProtoReflect.Descriptor instead.
func (*CreateOrderRequest) Descriptor() ([]byte, []int) {
	return file_rpc_create_order_proto_rawDescGZIP(), []int{0}
}

func (x *CreateOrderRequest) GetEscrowAddress() string {
	if x != nil {
		return x.EscrowAddress
	}
	return ""
}

func (x *CreateOrderRequest) GetCrypto() string {
	if x != nil {
		return x.Crypto
	}
	return ""
}

func (x *CreateOrderRequest) GetFiat() string {
	if x != nil {
		return x.Fiat
	}
	return ""
}

func (x *CreateOrderRequest) GetFiatAmount() float32 {
	if x != nil {
		return x.FiatAmount
	}
	return 0
}

func (x *CreateOrderRequest) GetCryptoAmount() float32 {
	if x != nil {
		return x.CryptoAmount
	}
	return 0
}

func (x *CreateOrderRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *CreateOrderRequest) GetRate() float32 {
	if x != nil {
		return x.Rate
	}
	return 0
}

func (x *CreateOrderRequest) GetUserAddress() string {
	if x != nil {
		return x.UserAddress
	}
	return ""
}

func (x *CreateOrderRequest) GetOrderType() string {
	if x != nil {
		return x.OrderType
	}
	return ""
}

func (x *CreateOrderRequest) GetBankName() string {
	if x != nil && x.BankName != nil {
		return *x.BankName
	}
	return ""
}

func (x *CreateOrderRequest) GetAccountNumber() string {
	if x != nil && x.AccountNumber != nil {
		return *x.AccountNumber
	}
	return ""
}

func (x *CreateOrderRequest) GetAccountHolder() string {
	if x != nil && x.AccountHolder != nil {
		return *x.AccountHolder
	}
	return ""
}

type CreateOrderResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Order *Order `protobuf:"bytes,1,opt,name=order,proto3" json:"order,omitempty"`
}

func (x *CreateOrderResponse) Reset() {
	*x = CreateOrderResponse{}
	mi := &file_rpc_create_order_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateOrderResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateOrderResponse) ProtoMessage() {}

func (x *CreateOrderResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_create_order_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateOrderResponse.ProtoReflect.Descriptor instead.
func (*CreateOrderResponse) Descriptor() ([]byte, []int) {
	return file_rpc_create_order_proto_rawDescGZIP(), []int{1}
}

func (x *CreateOrderResponse) GetOrder() *Order {
	if x != nil {
		return x.Order
	}
	return nil
}

var File_rpc_create_order_proto protoreflect.FileDescriptor

var file_rpc_create_order_proto_rawDesc = []byte{
	0x0a, 0x16, 0x72, 0x70, 0x63, 0x5f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x5f, 0x6f, 0x72, 0x64,
	0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x1a, 0x0b, 0x6f, 0x72,
	0x64, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xcd, 0x03, 0x0a, 0x12, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x25, 0x0a, 0x0e, 0x65, 0x73, 0x63, 0x72, 0x6f, 0x77, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x65, 0x73, 0x63, 0x72, 0x6f, 0x77,
	0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x72, 0x79, 0x70, 0x74,
	0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x63, 0x72, 0x79, 0x70, 0x74, 0x6f, 0x12,
	0x12, 0x0a, 0x04, 0x66, 0x69, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x66,
	0x69, 0x61, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x66, 0x69, 0x61, 0x74, 0x5f, 0x61, 0x6d, 0x6f, 0x75,
	0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x02, 0x52, 0x0a, 0x66, 0x69, 0x61, 0x74, 0x41, 0x6d,
	0x6f, 0x75, 0x6e, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x63, 0x72, 0x79, 0x70, 0x74, 0x6f, 0x5f, 0x61,
	0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x02, 0x52, 0x0c, 0x63, 0x72, 0x79,
	0x70, 0x74, 0x6f, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65,
	0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65,
	0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x72, 0x61, 0x74, 0x65, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x02, 0x52, 0x04, 0x72, 0x61, 0x74, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x75, 0x73, 0x65,
	0x72, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x75, 0x73, 0x65, 0x72, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x1d, 0x0a, 0x0a,
	0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x54, 0x79, 0x70, 0x65, 0x12, 0x20, 0x0a, 0x09, 0x62,
	0x61, 0x6e, 0x6b, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00,
	0x52, 0x08, 0x62, 0x61, 0x6e, 0x6b, 0x4e, 0x61, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x12, 0x2a, 0x0a,
	0x0e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18,
	0x0b, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x0d, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x88, 0x01, 0x01, 0x12, 0x2a, 0x0a, 0x0e, 0x61, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x68, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x18, 0x0c, 0x20, 0x01, 0x28,
	0x09, 0x48, 0x02, 0x52, 0x0d, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x48, 0x6f, 0x6c, 0x64,
	0x65, 0x72, 0x88, 0x01, 0x01, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x62, 0x61, 0x6e, 0x6b, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x42, 0x11, 0x0a, 0x0f, 0x5f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f,
	0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x42, 0x11, 0x0a, 0x0f, 0x5f, 0x61, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x5f, 0x68, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x22, 0x36, 0x0a, 0x13, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x1f, 0x0a, 0x05, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x09, 0x2e, 0x70, 0x62, 0x2e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x52, 0x05, 0x6f, 0x72, 0x64, 0x65,
	0x72, 0x42, 0x18, 0x5a, 0x16, 0x70, 0x65, 0x65, 0x72, 0x62, 0x69, 0x6c, 0x6c, 0x2d, 0x74, 0x72,
	0x61, 0x64, 0x65, 0x72, 0x2d, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_rpc_create_order_proto_rawDescOnce sync.Once
	file_rpc_create_order_proto_rawDescData = file_rpc_create_order_proto_rawDesc
)

func file_rpc_create_order_proto_rawDescGZIP() []byte {
	file_rpc_create_order_proto_rawDescOnce.Do(func() {
		file_rpc_create_order_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc_create_order_proto_rawDescData)
	})
	return file_rpc_create_order_proto_rawDescData
}

var file_rpc_create_order_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_rpc_create_order_proto_goTypes = []any{
	(*CreateOrderRequest)(nil),  // 0: pb.CreateOrderRequest
	(*CreateOrderResponse)(nil), // 1: pb.CreateOrderResponse
	(*Order)(nil),               // 2: pb.Order
}
var file_rpc_create_order_proto_depIdxs = []int32{
	2, // 0: pb.CreateOrderResponse.order:type_name -> pb.Order
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_rpc_create_order_proto_init() }
func file_rpc_create_order_proto_init() {
	if File_rpc_create_order_proto != nil {
		return
	}
	file_order_proto_init()
	file_rpc_create_order_proto_msgTypes[0].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_rpc_create_order_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_rpc_create_order_proto_goTypes,
		DependencyIndexes: file_rpc_create_order_proto_depIdxs,
		MessageInfos:      file_rpc_create_order_proto_msgTypes,
	}.Build()
	File_rpc_create_order_proto = out.File
	file_rpc_create_order_proto_rawDesc = nil
	file_rpc_create_order_proto_goTypes = nil
	file_rpc_create_order_proto_depIdxs = nil
}