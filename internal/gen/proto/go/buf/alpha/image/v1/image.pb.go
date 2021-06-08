// Copyright 2020-2021 Buf Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0-devel
// 	protoc        v3.17.3
// source: buf/alpha/image/v1/image.proto

package imagev1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Image is analogous to a FileDescriptorSet.
type Image struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// file matches the file field of a FileDescriptorSet.
	File []*descriptorpb.FileDescriptorProto `protobuf:"bytes,1,rep,name=file" json:"file,omitempty"`
	// bufbuild_image_extension is the ImageExtension for this image.
	//
	// The prefixed name and high tag value is used to all but guarantee there
	// will never be any conflict with Google's FileDescriptorSet definition.
	// The definition of a FileDescriptorSet has not changed in 11 years, so
	// we're not too worried about a conflict here.
	BufbuildImageExtension *ImageExtension `protobuf:"bytes,8042,opt,name=bufbuild_image_extension,json=bufbuildImageExtension" json:"bufbuild_image_extension,omitempty"`
}

func (x *Image) Reset() {
	*x = Image{}
	if protoimpl.UnsafeEnabled {
		mi := &file_buf_alpha_image_v1_image_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Image) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Image) ProtoMessage() {}

func (x *Image) ProtoReflect() protoreflect.Message {
	mi := &file_buf_alpha_image_v1_image_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Image.ProtoReflect.Descriptor instead.
func (*Image) Descriptor() ([]byte, []int) {
	return file_buf_alpha_image_v1_image_proto_rawDescGZIP(), []int{0}
}

func (x *Image) GetFile() []*descriptorpb.FileDescriptorProto {
	if x != nil {
		return x.File
	}
	return nil
}

func (x *Image) GetBufbuildImageExtension() *ImageExtension {
	if x != nil {
		return x.BufbuildImageExtension
	}
	return nil
}

// ImageExtension contains extensions to Images.
//
// The fields are not included directly on the Image so that we can both
// detect if extensions exist, which signifies this was created by buf
// and not by protoc, and so that we can add fields in a freeform manner
// without worrying about conflicts with google.protobuf.FileDescriptorSet.
type ImageExtension struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// image_import_refs are the image import references for this specific Image.
	//
	// A given FileDescriptorProto may or may not be an import depending on
	// the image context, so this information is not stored on each FileDescriptorProto.
	ImageImportRefs []*ImageImportRef `protobuf:"bytes,1,rep,name=image_import_refs,json=imageImportRefs" json:"image_import_refs,omitempty"`
	// ModuleCommitRefs are the module commit commits for this specific Image.
	//
	// The lack of a ModuleCommitRef for a given file means that the module
	// commit information is not known, and must be treated as unknown.
	ModuleCommitRefs []*ModuleCommitRef `protobuf:"bytes,3,rep,name=module_commit_refs,json=moduleCommitRefs" json:"module_commit_refs,omitempty"`
}

func (x *ImageExtension) Reset() {
	*x = ImageExtension{}
	if protoimpl.UnsafeEnabled {
		mi := &file_buf_alpha_image_v1_image_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ImageExtension) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ImageExtension) ProtoMessage() {}

func (x *ImageExtension) ProtoReflect() protoreflect.Message {
	mi := &file_buf_alpha_image_v1_image_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ImageExtension.ProtoReflect.Descriptor instead.
func (*ImageExtension) Descriptor() ([]byte, []int) {
	return file_buf_alpha_image_v1_image_proto_rawDescGZIP(), []int{1}
}

func (x *ImageExtension) GetImageImportRefs() []*ImageImportRef {
	if x != nil {
		return x.ImageImportRefs
	}
	return nil
}

func (x *ImageExtension) GetModuleCommitRefs() []*ModuleCommitRef {
	if x != nil {
		return x.ModuleCommitRefs
	}
	return nil
}

// ImageImportRef is a reference to an image import.
//
// This is a message type instead of a scalar type so that we can add
// additional information about an import reference in the future, such as
// the external location of the import.
type ImageImportRef struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// file_index is the index within the Image file array of the import.
	//
	// This signifies that file[file_index] is an import.
	// This field must be set.
	FileIndex *uint32 `protobuf:"varint,1,opt,name=file_index,json=fileIndex" json:"file_index,omitempty"`
}

func (x *ImageImportRef) Reset() {
	*x = ImageImportRef{}
	if protoimpl.UnsafeEnabled {
		mi := &file_buf_alpha_image_v1_image_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ImageImportRef) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ImageImportRef) ProtoMessage() {}

func (x *ImageImportRef) ProtoReflect() protoreflect.Message {
	mi := &file_buf_alpha_image_v1_image_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ImageImportRef.ProtoReflect.Descriptor instead.
func (*ImageImportRef) Descriptor() ([]byte, []int) {
	return file_buf_alpha_image_v1_image_proto_rawDescGZIP(), []int{2}
}

func (x *ImageImportRef) GetFileIndex() uint32 {
	if x != nil && x.FileIndex != nil {
		return *x.FileIndex
	}
	return 0
}

// ModuleCommitRef is a commit to a module commit.
type ModuleCommitRef struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// file_index is the index within the Image file array.
	//
	// This signifies that file[file_index] is the file that
	// uses this module commit.
	//
	// This field must be set.
	FileIndex  *uint32 `protobuf:"varint,1,opt,name=file_index,json=fileIndex" json:"file_index,omitempty"`
	Remote     *string `protobuf:"bytes,2,opt,name=remote" json:"remote,omitempty"`
	Owner      *string `protobuf:"bytes,3,opt,name=owner" json:"owner,omitempty"`
	Repository *string `protobuf:"bytes,4,opt,name=repository" json:"repository,omitempty"`
	Commit     *string `protobuf:"bytes,5,opt,name=commit" json:"commit,omitempty"`
}

func (x *ModuleCommitRef) Reset() {
	*x = ModuleCommitRef{}
	if protoimpl.UnsafeEnabled {
		mi := &file_buf_alpha_image_v1_image_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ModuleCommitRef) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ModuleCommitRef) ProtoMessage() {}

func (x *ModuleCommitRef) ProtoReflect() protoreflect.Message {
	mi := &file_buf_alpha_image_v1_image_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ModuleCommitRef.ProtoReflect.Descriptor instead.
func (*ModuleCommitRef) Descriptor() ([]byte, []int) {
	return file_buf_alpha_image_v1_image_proto_rawDescGZIP(), []int{3}
}

func (x *ModuleCommitRef) GetFileIndex() uint32 {
	if x != nil && x.FileIndex != nil {
		return *x.FileIndex
	}
	return 0
}

func (x *ModuleCommitRef) GetRemote() string {
	if x != nil && x.Remote != nil {
		return *x.Remote
	}
	return ""
}

func (x *ModuleCommitRef) GetOwner() string {
	if x != nil && x.Owner != nil {
		return *x.Owner
	}
	return ""
}

func (x *ModuleCommitRef) GetRepository() string {
	if x != nil && x.Repository != nil {
		return *x.Repository
	}
	return ""
}

func (x *ModuleCommitRef) GetCommit() string {
	if x != nil && x.Commit != nil {
		return *x.Commit
	}
	return ""
}

var File_buf_alpha_image_v1_image_proto protoreflect.FileDescriptor

var file_buf_alpha_image_v1_image_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x2f, 0x69, 0x6d, 0x61, 0x67,
	0x65, 0x2f, 0x76, 0x31, 0x2f, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x12, 0x62, 0x75, 0x66, 0x2e, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x2e, 0x69, 0x6d, 0x61, 0x67,
	0x65, 0x2e, 0x76, 0x31, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa0, 0x01, 0x0a, 0x05, 0x49, 0x6d, 0x61, 0x67, 0x65,
	0x12, 0x38, 0x0a, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x46, 0x69, 0x6c, 0x65, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x52, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x12, 0x5d, 0x0a, 0x18, 0x62, 0x75,
	0x66, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x5f, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x65, 0x78, 0x74,
	0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0xea, 0x3e, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e,
	0x62, 0x75, 0x66, 0x2e, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x2e, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x2e,
	0x76, 0x31, 0x2e, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f,
	0x6e, 0x52, 0x16, 0x62, 0x75, 0x66, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x49, 0x6d, 0x61, 0x67, 0x65,
	0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0xd0, 0x01, 0x0a, 0x0e, 0x49, 0x6d,
	0x61, 0x67, 0x65, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x4e, 0x0a, 0x11,
	0x69, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x69, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x5f, 0x72, 0x65, 0x66,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x62, 0x75, 0x66, 0x2e, 0x61, 0x6c,
	0x70, 0x68, 0x61, 0x2e, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x6d, 0x61,
	0x67, 0x65, 0x49, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x66, 0x52, 0x0f, 0x69, 0x6d, 0x61,
	0x67, 0x65, 0x49, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x66, 0x73, 0x12, 0x51, 0x0a, 0x12,
	0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x5f, 0x72, 0x65,
	0x66, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x62, 0x75, 0x66, 0x2e, 0x61,
	0x6c, 0x70, 0x68, 0x61, 0x2e, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x6f,
	0x64, 0x75, 0x6c, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x52, 0x65, 0x66, 0x52, 0x10, 0x6d,
	0x6f, 0x64, 0x75, 0x6c, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x52, 0x65, 0x66, 0x73, 0x4a,
	0x04, 0x08, 0x02, 0x10, 0x03, 0x52, 0x15, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x5f, 0x72, 0x65,
	0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x5f, 0x72, 0x65, 0x66, 0x73, 0x22, 0x2f, 0x0a, 0x0e,
	0x49, 0x6d, 0x61, 0x67, 0x65, 0x49, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x66, 0x12, 0x1d,
	0x0a, 0x0a, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x22, 0x96, 0x01,
	0x0a, 0x0f, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x52, 0x65,
	0x66, 0x12, 0x1d, 0x0a, 0x0a, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x6e, 0x64, 0x65, 0x78,
	0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6f, 0x77, 0x6e, 0x65,
	0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x12, 0x1e,
	0x0a, 0x0a, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x12, 0x16,
	0x0a, 0x06, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x42, 0x73, 0x0a, 0x16, 0x63, 0x6f, 0x6d, 0x2e, 0x62, 0x75,
	0x66, 0x2e, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x2e, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31,
	0x42, 0x0a, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x48, 0x01, 0x5a, 0x48,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x75, 0x66, 0x62, 0x75,
	0x69, 0x6c, 0x64, 0x2f, 0x62, 0x75, 0x66, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x2f, 0x62, 0x75,
	0x66, 0x2f, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x2f, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x2f, 0x76, 0x31,
	0x3b, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x76, 0x31, 0xf8, 0x01, 0x01,
}

var (
	file_buf_alpha_image_v1_image_proto_rawDescOnce sync.Once
	file_buf_alpha_image_v1_image_proto_rawDescData = file_buf_alpha_image_v1_image_proto_rawDesc
)

func file_buf_alpha_image_v1_image_proto_rawDescGZIP() []byte {
	file_buf_alpha_image_v1_image_proto_rawDescOnce.Do(func() {
		file_buf_alpha_image_v1_image_proto_rawDescData = protoimpl.X.CompressGZIP(file_buf_alpha_image_v1_image_proto_rawDescData)
	})
	return file_buf_alpha_image_v1_image_proto_rawDescData
}

var file_buf_alpha_image_v1_image_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_buf_alpha_image_v1_image_proto_goTypes = []interface{}{
	(*Image)(nil),                            // 0: buf.alpha.image.v1.Image
	(*ImageExtension)(nil),                   // 1: buf.alpha.image.v1.ImageExtension
	(*ImageImportRef)(nil),                   // 2: buf.alpha.image.v1.ImageImportRef
	(*ModuleCommitRef)(nil),                  // 3: buf.alpha.image.v1.ModuleCommitRef
	(*descriptorpb.FileDescriptorProto)(nil), // 4: google.protobuf.FileDescriptorProto
}
var file_buf_alpha_image_v1_image_proto_depIdxs = []int32{
	4, // 0: buf.alpha.image.v1.Image.file:type_name -> google.protobuf.FileDescriptorProto
	1, // 1: buf.alpha.image.v1.Image.bufbuild_image_extension:type_name -> buf.alpha.image.v1.ImageExtension
	2, // 2: buf.alpha.image.v1.ImageExtension.image_import_refs:type_name -> buf.alpha.image.v1.ImageImportRef
	3, // 3: buf.alpha.image.v1.ImageExtension.module_commit_refs:type_name -> buf.alpha.image.v1.ModuleCommitRef
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_buf_alpha_image_v1_image_proto_init() }
func file_buf_alpha_image_v1_image_proto_init() {
	if File_buf_alpha_image_v1_image_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_buf_alpha_image_v1_image_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Image); i {
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
		file_buf_alpha_image_v1_image_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ImageExtension); i {
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
		file_buf_alpha_image_v1_image_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ImageImportRef); i {
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
		file_buf_alpha_image_v1_image_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ModuleCommitRef); i {
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
			RawDescriptor: file_buf_alpha_image_v1_image_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_buf_alpha_image_v1_image_proto_goTypes,
		DependencyIndexes: file_buf_alpha_image_v1_image_proto_depIdxs,
		MessageInfos:      file_buf_alpha_image_v1_image_proto_msgTypes,
	}.Build()
	File_buf_alpha_image_v1_image_proto = out.File
	file_buf_alpha_image_v1_image_proto_rawDesc = nil
	file_buf_alpha_image_v1_image_proto_goTypes = nil
	file_buf_alpha_image_v1_image_proto_depIdxs = nil
}
