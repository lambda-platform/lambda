// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.3
// source: lambda.proto

package proto

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

type RoleParams struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *RoleParams) Reset() {
	*x = RoleParams{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lambda_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RoleParams) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RoleParams) ProtoMessage() {}

func (x *RoleParams) ProtoReflect() protoreflect.Message {
	mi := &file_lambda_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RoleParams.ProtoReflect.Descriptor instead.
func (*RoleParams) Descriptor() ([]byte, []int) {
	return file_lambda_proto_rawDescGZIP(), []int{0}
}

func (x *RoleParams) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type SchemaParams struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     int32  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Schema string `protobuf:"bytes,2,opt,name=schema,proto3" json:"schema,omitempty"`
	Type   string `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`
}

func (x *SchemaParams) Reset() {
	*x = SchemaParams{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lambda_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SchemaParams) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SchemaParams) ProtoMessage() {}

func (x *SchemaParams) ProtoReflect() protoreflect.Message {
	mi := &file_lambda_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SchemaParams.ProtoReflect.Descriptor instead.
func (*SchemaParams) Descriptor() ([]byte, []int) {
	return file_lambda_proto_rawDescGZIP(), []int{1}
}

func (x *SchemaParams) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *SchemaParams) GetSchema() string {
	if x != nil {
		return x.Schema
	}
	return ""
}

func (x *SchemaParams) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

type TableOption struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Table  string   `protobuf:"bytes,1,opt,name=table,proto3" json:"table,omitempty"`
	Key    string   `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Field  string   `protobuf:"bytes,3,opt,name=field,proto3" json:"field,omitempty"`
	Values []string `protobuf:"bytes,4,rep,name=values,proto3" json:"values,omitempty"`
}

func (x *TableOption) Reset() {
	*x = TableOption{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lambda_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TableOption) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TableOption) ProtoMessage() {}

func (x *TableOption) ProtoReflect() protoreflect.Message {
	mi := &file_lambda_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TableOption.ProtoReflect.Descriptor instead.
func (*TableOption) Descriptor() ([]byte, []int) {
	return file_lambda_proto_rawDescGZIP(), []int{2}
}

func (x *TableOption) GetTable() string {
	if x != nil {
		return x.Table
	}
	return ""
}

func (x *TableOption) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *TableOption) GetField() string {
	if x != nil {
		return x.Field
	}
	return ""
}

func (x *TableOption) GetValues() []string {
	if x != nil {
		return x.Values
	}
	return nil
}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status bool `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lambda_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_lambda_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_lambda_proto_rawDescGZIP(), []int{3}
}

func (x *Response) GetStatus() bool {
	if x != nil {
		return x.Status
	}
	return false
}

type IntRows struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Rows []*IntRow `protobuf:"bytes,1,rep,name=rows,proto3" json:"rows,omitempty"`
}

func (x *IntRows) Reset() {
	*x = IntRows{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lambda_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IntRows) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IntRows) ProtoMessage() {}

func (x *IntRows) ProtoReflect() protoreflect.Message {
	mi := &file_lambda_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IntRows.ProtoReflect.Descriptor instead.
func (*IntRows) Descriptor() ([]byte, []int) {
	return file_lambda_proto_rawDescGZIP(), []int{4}
}

func (x *IntRows) GetRows() []*IntRow {
	if x != nil {
		return x.Rows
	}
	return nil
}

type IntRow struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   int32 `protobuf:"varint,1,opt,name=key,proto3" json:"key,omitempty"`
	Value int32 `protobuf:"varint,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *IntRow) Reset() {
	*x = IntRow{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lambda_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IntRow) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IntRow) ProtoMessage() {}

func (x *IntRow) ProtoReflect() protoreflect.Message {
	mi := &file_lambda_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IntRow.ProtoReflect.Descriptor instead.
func (*IntRow) Descriptor() ([]byte, []int) {
	return file_lambda_proto_rawDescGZIP(), []int{5}
}

func (x *IntRow) GetKey() int32 {
	if x != nil {
		return x.Key
	}
	return 0
}

func (x *IntRow) GetValue() int32 {
	if x != nil {
		return x.Value
	}
	return 0
}

type StringRows struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Rows []*StringRow `protobuf:"bytes,1,rep,name=rows,proto3" json:"rows,omitempty"`
}

func (x *StringRows) Reset() {
	*x = StringRows{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lambda_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StringRows) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StringRows) ProtoMessage() {}

func (x *StringRows) ProtoReflect() protoreflect.Message {
	mi := &file_lambda_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StringRows.ProtoReflect.Descriptor instead.
func (*StringRows) Descriptor() ([]byte, []int) {
	return file_lambda_proto_rawDescGZIP(), []int{6}
}

func (x *StringRows) GetRows() []*StringRow {
	if x != nil {
		return x.Rows
	}
	return nil
}

type StringRow struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *StringRow) Reset() {
	*x = StringRow{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lambda_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StringRow) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StringRow) ProtoMessage() {}

func (x *StringRow) ProtoReflect() protoreflect.Message {
	mi := &file_lambda_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StringRow.ProtoReflect.Descriptor instead.
func (*StringRow) Descriptor() ([]byte, []int) {
	return file_lambda_proto_rawDescGZIP(), []int{7}
}

func (x *StringRow) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *StringRow) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

var File_lambda_proto protoreflect.FileDescriptor

var file_lambda_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x6c, 0x61, 0x6d, 0x62, 0x64, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x6c, 0x61, 0x6d, 0x62, 0x64, 0x61, 0x22, 0x1c, 0x0a, 0x0a, 0x52, 0x6f, 0x6c, 0x65, 0x50, 0x61,
	0x72, 0x61, 0x6d, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x69, 0x64, 0x22, 0x4a, 0x0a, 0x0c, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x50, 0x61,
	0x72, 0x61, 0x6d, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x12, 0x12, 0x0a, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x22, 0x63, 0x0a, 0x0b, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x14, 0x0a, 0x05, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x74, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x66, 0x69, 0x65, 0x6c, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x16, 0x0a,
	0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x73, 0x22, 0x22, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x2d, 0x0a, 0x07, 0x49, 0x6e, 0x74,
	0x52, 0x6f, 0x77, 0x73, 0x12, 0x22, 0x0a, 0x04, 0x72, 0x6f, 0x77, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x6c, 0x61, 0x6d, 0x62, 0x64, 0x61, 0x2e, 0x49, 0x6e, 0x74, 0x52,
	0x6f, 0x77, 0x52, 0x04, 0x72, 0x6f, 0x77, 0x73, 0x22, 0x30, 0x0a, 0x06, 0x49, 0x6e, 0x74, 0x52,
	0x6f, 0x77, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x33, 0x0a, 0x0a, 0x53, 0x74,
	0x72, 0x69, 0x6e, 0x67, 0x52, 0x6f, 0x77, 0x73, 0x12, 0x25, 0x0a, 0x04, 0x72, 0x6f, 0x77, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x6c, 0x61, 0x6d, 0x62, 0x64, 0x61, 0x2e,
	0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x52, 0x6f, 0x77, 0x52, 0x04, 0x72, 0x6f, 0x77, 0x73, 0x22,
	0x33, 0x0a, 0x09, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x52, 0x6f, 0x77, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x32, 0xa6, 0x02, 0x0a, 0x06, 0x4c, 0x61, 0x6d, 0x62, 0x64, 0x61, 0x12,
	0x34, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x49, 0x6e, 0x74, 0x44, 0x61, 0x74, 0x61, 0x12, 0x13, 0x2e,
	0x6c, 0x61, 0x6d, 0x62, 0x64, 0x61, 0x2e, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x4f, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x1a, 0x0f, 0x2e, 0x6c, 0x61, 0x6d, 0x62, 0x64, 0x61, 0x2e, 0x49, 0x6e, 0x74, 0x52,
	0x6f, 0x77, 0x73, 0x22, 0x00, 0x12, 0x3a, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x53, 0x74, 0x72, 0x69,
	0x6e, 0x67, 0x44, 0x61, 0x74, 0x61, 0x12, 0x13, 0x2e, 0x6c, 0x61, 0x6d, 0x62, 0x64, 0x61, 0x2e,
	0x54, 0x61, 0x62, 0x6c, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x12, 0x2e, 0x6c, 0x61,
	0x6d, 0x62, 0x64, 0x61, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x52, 0x6f, 0x77, 0x73, 0x22,
	0x00, 0x12, 0x39, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x44, 0x61,
	0x74, 0x61, 0x12, 0x14, 0x2e, 0x6c, 0x61, 0x6d, 0x62, 0x64, 0x61, 0x2e, 0x53, 0x63, 0x68, 0x65,
	0x6d, 0x61, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x1a, 0x10, 0x2e, 0x6c, 0x61, 0x6d, 0x62, 0x64,
	0x61, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x35, 0x0a, 0x0b,
	0x47, 0x65, 0x74, 0x52, 0x6f, 0x6c, 0x65, 0x44, 0x61, 0x74, 0x61, 0x12, 0x12, 0x2e, 0x6c, 0x61,
	0x6d, 0x62, 0x64, 0x61, 0x2e, 0x52, 0x6f, 0x6c, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x1a,
	0x10, 0x2e, 0x6c, 0x61, 0x6d, 0x62, 0x64, 0x61, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x12, 0x38, 0x0a, 0x0e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x4d, 0x79, 0x53,
	0x43, 0x48, 0x45, 0x4d, 0x41, 0x12, 0x12, 0x2e, 0x6c, 0x61, 0x6d, 0x62, 0x64, 0x61, 0x2e, 0x52,
	0x6f, 0x6c, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x1a, 0x10, 0x2e, 0x6c, 0x61, 0x6d, 0x62,
	0x64, 0x61, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x0a, 0x5a,
	0x08, 0x2e, 0x2f, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_lambda_proto_rawDescOnce sync.Once
	file_lambda_proto_rawDescData = file_lambda_proto_rawDesc
)

func file_lambda_proto_rawDescGZIP() []byte {
	file_lambda_proto_rawDescOnce.Do(func() {
		file_lambda_proto_rawDescData = protoimpl.X.CompressGZIP(file_lambda_proto_rawDescData)
	})
	return file_lambda_proto_rawDescData
}

var file_lambda_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_lambda_proto_goTypes = []interface{}{
	(*RoleParams)(nil),   // 0: lambda.RoleParams
	(*SchemaParams)(nil), // 1: lambda.SchemaParams
	(*TableOption)(nil),  // 2: lambda.TableOption
	(*Response)(nil),     // 3: lambda.Response
	(*IntRows)(nil),      // 4: lambda.IntRows
	(*IntRow)(nil),       // 5: lambda.IntRow
	(*StringRows)(nil),   // 6: lambda.StringRows
	(*StringRow)(nil),    // 7: lambda.StringRow
}
var file_lambda_proto_depIdxs = []int32{
	5, // 0: lambda.IntRows.rows:type_name -> lambda.IntRow
	7, // 1: lambda.StringRows.rows:type_name -> lambda.StringRow
	2, // 2: lambda.Lambda.GetIntData:input_type -> lambda.TableOption
	2, // 3: lambda.Lambda.GetStringData:input_type -> lambda.TableOption
	1, // 4: lambda.Lambda.GetSchemaData:input_type -> lambda.SchemaParams
	0, // 5: lambda.Lambda.GetRoleData:input_type -> lambda.RoleParams
	0, // 6: lambda.Lambda.UploadMySCHEMA:input_type -> lambda.RoleParams
	4, // 7: lambda.Lambda.GetIntData:output_type -> lambda.IntRows
	6, // 8: lambda.Lambda.GetStringData:output_type -> lambda.StringRows
	3, // 9: lambda.Lambda.GetSchemaData:output_type -> lambda.Response
	3, // 10: lambda.Lambda.GetRoleData:output_type -> lambda.Response
	3, // 11: lambda.Lambda.UploadMySCHEMA:output_type -> lambda.Response
	7, // [7:12] is the sub-list for method output_type
	2, // [2:7] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_lambda_proto_init() }
func file_lambda_proto_init() {
	if File_lambda_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_lambda_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RoleParams); i {
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
		file_lambda_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SchemaParams); i {
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
		file_lambda_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TableOption); i {
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
		file_lambda_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
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
		file_lambda_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IntRows); i {
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
		file_lambda_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IntRow); i {
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
		file_lambda_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StringRows); i {
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
		file_lambda_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StringRow); i {
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
			RawDescriptor: file_lambda_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_lambda_proto_goTypes,
		DependencyIndexes: file_lambda_proto_depIdxs,
		MessageInfos:      file_lambda_proto_msgTypes,
	}.Build()
	File_lambda_proto = out.File
	file_lambda_proto_rawDesc = nil
	file_lambda_proto_goTypes = nil
	file_lambda_proto_depIdxs = nil
}
