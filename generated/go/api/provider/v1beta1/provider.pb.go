// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: api/provider/v1beta1/provider.proto

package v1beta1

import (
	v1beta1 "github.com/tak-sh/tak/generated/go/api/metadata/v1beta1"
	v1beta11 "github.com/tak-sh/tak/generated/go/api/script/v1beta1"
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

// Provides the instructions necessary for a user to add an Account to
// tak.
type Provider struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metadata *v1beta1.Metadata `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Spec     *Spec             `protobuf:"bytes,2,opt,name=spec,proto3" json:"spec,omitempty"`
}

func (x *Provider) Reset() {
	*x = Provider{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_provider_v1beta1_provider_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Provider) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Provider) ProtoMessage() {}

func (x *Provider) ProtoReflect() protoreflect.Message {
	mi := &file_api_provider_v1beta1_provider_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Provider.ProtoReflect.Descriptor instead.
func (*Provider) Descriptor() ([]byte, []int) {
	return file_api_provider_v1beta1_provider_proto_rawDescGZIP(), []int{0}
}

func (x *Provider) GetMetadata() *v1beta1.Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *Provider) GetSpec() *Spec {
	if x != nil {
		return x.Spec
	}
	return nil
}

type Spec struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// A script to handle the login flow for an account. This must take MFA into consideration
	// when writing.
	Login *v1beta11.Script `protobuf:"bytes,1,opt,name=login,proto3" json:"login,omitempty"`
	// A script to handle downloading transactions file from the account.
	DownloadTransactions *v1beta11.Script `protobuf:"bytes,2,opt,name=download_transactions,proto3" json:"download_transactions,omitempty"`
}

func (x *Spec) Reset() {
	*x = Spec{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_provider_v1beta1_provider_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Spec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Spec) ProtoMessage() {}

func (x *Spec) ProtoReflect() protoreflect.Message {
	mi := &file_api_provider_v1beta1_provider_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Spec.ProtoReflect.Descriptor instead.
func (*Spec) Descriptor() ([]byte, []int) {
	return file_api_provider_v1beta1_provider_proto_rawDescGZIP(), []int{1}
}

func (x *Spec) GetLogin() *v1beta11.Script {
	if x != nil {
		return x.Login
	}
	return nil
}

func (x *Spec) GetDownloadTransactions() *v1beta11.Script {
	if x != nil {
		return x.DownloadTransactions
	}
	return nil
}

var File_api_provider_v1beta1_provider_proto protoreflect.FileDescriptor

var file_api_provider_v1beta1_provider_proto_rawDesc = []byte{
	0x0a, 0x23, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x2f, 0x76,
	0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2f, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1b, 0x74, 0x61, 0x6b, 0x2e, 0x73, 0x68, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74,
	0x61, 0x31, 0x1a, 0x1f, 0x61, 0x70, 0x69, 0x2f, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x2f, 0x76,
	0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2f, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x23, 0x61, 0x70, 0x69, 0x2f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x84, 0x01, 0x0a, 0x08, 0x50, 0x72, 0x6f,
	0x76, 0x69, 0x64, 0x65, 0x72, 0x12, 0x41, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x74, 0x61, 0x6b, 0x2e, 0x73, 0x68,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x76, 0x31,
	0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x08,
	0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x35, 0x0a, 0x04, 0x73, 0x70, 0x65, 0x63,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x74, 0x61, 0x6b, 0x2e, 0x73, 0x68, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x62,
	0x65, 0x74, 0x61, 0x31, 0x2e, 0x53, 0x70, 0x65, 0x63, 0x52, 0x04, 0x73, 0x70, 0x65, 0x63, 0x22,
	0x98, 0x01, 0x0a, 0x04, 0x53, 0x70, 0x65, 0x63, 0x12, 0x37, 0x0a, 0x05, 0x6c, 0x6f, 0x67, 0x69,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x74, 0x61, 0x6b, 0x2e, 0x73, 0x68,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x2e, 0x76, 0x31, 0x62, 0x65,
	0x74, 0x61, 0x31, 0x2e, 0x53, 0x63, 0x72, 0x69, 0x70, 0x74, 0x52, 0x05, 0x6c, 0x6f, 0x67, 0x69,
	0x6e, 0x12, 0x57, 0x0a, 0x15, 0x64, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x74, 0x72,
	0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x21, 0x2e, 0x74, 0x61, 0x6b, 0x2e, 0x73, 0x68, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x63,
	0x72, 0x69, 0x70, 0x74, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x53, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x52, 0x15, 0x64, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x74, 0x72,
	0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x42, 0xfa, 0x01, 0x0a, 0x1f, 0x63,
	0x6f, 0x6d, 0x2e, 0x74, 0x61, 0x6b, 0x2e, 0x73, 0x68, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72,
	0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x42, 0x0d,
	0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a,
	0x37, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x61, 0x6b, 0x2d,
	0x73, 0x68, 0x2f, 0x74, 0x61, 0x6b, 0x2f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64,
	0x2f, 0x67, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72,
	0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0xa2, 0x02, 0x04, 0x54, 0x53, 0x41, 0x50, 0xaa,
	0x02, 0x1b, 0x54, 0x61, 0x6b, 0x2e, 0x53, 0x68, 0x2e, 0x41, 0x70, 0x69, 0x2e, 0x50, 0x72, 0x6f,
	0x76, 0x69, 0x64, 0x65, 0x72, 0x2e, 0x56, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0xca, 0x02, 0x1b,
	0x54, 0x61, 0x6b, 0x5c, 0x53, 0x68, 0x5c, 0x41, 0x70, 0x69, 0x5c, 0x50, 0x72, 0x6f, 0x76, 0x69,
	0x64, 0x65, 0x72, 0x5c, 0x56, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0xe2, 0x02, 0x27, 0x54, 0x61,
	0x6b, 0x5c, 0x53, 0x68, 0x5c, 0x41, 0x70, 0x69, 0x5c, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65,
	0x72, 0x5c, 0x56, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x1f, 0x54, 0x61, 0x6b, 0x3a, 0x3a, 0x53, 0x68, 0x3a,
	0x3a, 0x41, 0x70, 0x69, 0x3a, 0x3a, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x3a, 0x3a,
	0x56, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_provider_v1beta1_provider_proto_rawDescOnce sync.Once
	file_api_provider_v1beta1_provider_proto_rawDescData = file_api_provider_v1beta1_provider_proto_rawDesc
)

func file_api_provider_v1beta1_provider_proto_rawDescGZIP() []byte {
	file_api_provider_v1beta1_provider_proto_rawDescOnce.Do(func() {
		file_api_provider_v1beta1_provider_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_provider_v1beta1_provider_proto_rawDescData)
	})
	return file_api_provider_v1beta1_provider_proto_rawDescData
}

var file_api_provider_v1beta1_provider_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_api_provider_v1beta1_provider_proto_goTypes = []any{
	(*Provider)(nil),         // 0: tak.sh.api.provider.v1beta1.Provider
	(*Spec)(nil),             // 1: tak.sh.api.provider.v1beta1.Spec
	(*v1beta1.Metadata)(nil), // 2: tak.sh.api.metadata.v1beta1.Metadata
	(*v1beta11.Script)(nil),  // 3: tak.sh.api.script.v1beta1.Script
}
var file_api_provider_v1beta1_provider_proto_depIdxs = []int32{
	2, // 0: tak.sh.api.provider.v1beta1.Provider.metadata:type_name -> tak.sh.api.metadata.v1beta1.Metadata
	1, // 1: tak.sh.api.provider.v1beta1.Provider.spec:type_name -> tak.sh.api.provider.v1beta1.Spec
	3, // 2: tak.sh.api.provider.v1beta1.Spec.login:type_name -> tak.sh.api.script.v1beta1.Script
	3, // 3: tak.sh.api.provider.v1beta1.Spec.download_transactions:type_name -> tak.sh.api.script.v1beta1.Script
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_api_provider_v1beta1_provider_proto_init() }
func file_api_provider_v1beta1_provider_proto_init() {
	if File_api_provider_v1beta1_provider_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_provider_v1beta1_provider_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Provider); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_provider_v1beta1_provider_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*Spec); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_provider_v1beta1_provider_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_provider_v1beta1_provider_proto_goTypes,
		DependencyIndexes: file_api_provider_v1beta1_provider_proto_depIdxs,
		MessageInfos:      file_api_provider_v1beta1_provider_proto_msgTypes,
	}.Build()
	File_api_provider_v1beta1_provider_proto = out.File
	file_api_provider_v1beta1_provider_proto_rawDesc = nil
	file_api_provider_v1beta1_provider_proto_goTypes = nil
	file_api_provider_v1beta1_provider_proto_depIdxs = nil
}
