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

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.1.0
// - protoc             v3.17.3
// source: buf/alpha/registry/v1alpha1/repository_commit.proto

package registryv1alpha1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// RepositoryCommitServiceClient is the client API for RepositoryCommitService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RepositoryCommitServiceClient interface {
	// ListRepositoryCommitsByBranch lists the repository commits associated
	// with a repository branch on a repository, ordered by their create time.
	ListRepositoryCommitsByBranch(ctx context.Context, in *ListRepositoryCommitsByBranchRequest, opts ...grpc.CallOption) (*ListRepositoryCommitsByBranchResponse, error)
	// GetRepositoryCommitByReference returns the repository commit matching
	// the provided reference, if it exists.
	GetRepositoryCommitByReference(ctx context.Context, in *GetRepositoryCommitByReferenceRequest, opts ...grpc.CallOption) (*GetRepositoryCommitByReferenceResponse, error)
	// GetRepositoryCommitBySequenceID returns the repository commit matching
	// the provided sequence ID and branch, if it exists.
	GetRepositoryCommitBySequenceID(ctx context.Context, in *GetRepositoryCommitBySequenceIDRequest, opts ...grpc.CallOption) (*GetRepositoryCommitBySequenceIDResponse, error)
}

type repositoryCommitServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRepositoryCommitServiceClient(cc grpc.ClientConnInterface) RepositoryCommitServiceClient {
	return &repositoryCommitServiceClient{cc}
}

func (c *repositoryCommitServiceClient) ListRepositoryCommitsByBranch(ctx context.Context, in *ListRepositoryCommitsByBranchRequest, opts ...grpc.CallOption) (*ListRepositoryCommitsByBranchResponse, error) {
	out := new(ListRepositoryCommitsByBranchResponse)
	err := c.cc.Invoke(ctx, "/buf.alpha.registry.v1alpha1.RepositoryCommitService/ListRepositoryCommitsByBranch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *repositoryCommitServiceClient) GetRepositoryCommitByReference(ctx context.Context, in *GetRepositoryCommitByReferenceRequest, opts ...grpc.CallOption) (*GetRepositoryCommitByReferenceResponse, error) {
	out := new(GetRepositoryCommitByReferenceResponse)
	err := c.cc.Invoke(ctx, "/buf.alpha.registry.v1alpha1.RepositoryCommitService/GetRepositoryCommitByReference", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *repositoryCommitServiceClient) GetRepositoryCommitBySequenceID(ctx context.Context, in *GetRepositoryCommitBySequenceIDRequest, opts ...grpc.CallOption) (*GetRepositoryCommitBySequenceIDResponse, error) {
	out := new(GetRepositoryCommitBySequenceIDResponse)
	err := c.cc.Invoke(ctx, "/buf.alpha.registry.v1alpha1.RepositoryCommitService/GetRepositoryCommitBySequenceID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RepositoryCommitServiceServer is the server API for RepositoryCommitService service.
// All implementations should embed UnimplementedRepositoryCommitServiceServer
// for forward compatibility
type RepositoryCommitServiceServer interface {
	// ListRepositoryCommitsByBranch lists the repository commits associated
	// with a repository branch on a repository, ordered by their create time.
	ListRepositoryCommitsByBranch(context.Context, *ListRepositoryCommitsByBranchRequest) (*ListRepositoryCommitsByBranchResponse, error)
	// GetRepositoryCommitByReference returns the repository commit matching
	// the provided reference, if it exists.
	GetRepositoryCommitByReference(context.Context, *GetRepositoryCommitByReferenceRequest) (*GetRepositoryCommitByReferenceResponse, error)
	// GetRepositoryCommitBySequenceID returns the repository commit matching
	// the provided sequence ID and branch, if it exists.
	GetRepositoryCommitBySequenceID(context.Context, *GetRepositoryCommitBySequenceIDRequest) (*GetRepositoryCommitBySequenceIDResponse, error)
}

// UnimplementedRepositoryCommitServiceServer should be embedded to have forward compatible implementations.
type UnimplementedRepositoryCommitServiceServer struct {
}

func (UnimplementedRepositoryCommitServiceServer) ListRepositoryCommitsByBranch(context.Context, *ListRepositoryCommitsByBranchRequest) (*ListRepositoryCommitsByBranchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListRepositoryCommitsByBranch not implemented")
}
func (UnimplementedRepositoryCommitServiceServer) GetRepositoryCommitByReference(context.Context, *GetRepositoryCommitByReferenceRequest) (*GetRepositoryCommitByReferenceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRepositoryCommitByReference not implemented")
}
func (UnimplementedRepositoryCommitServiceServer) GetRepositoryCommitBySequenceID(context.Context, *GetRepositoryCommitBySequenceIDRequest) (*GetRepositoryCommitBySequenceIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRepositoryCommitBySequenceID not implemented")
}

// UnsafeRepositoryCommitServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RepositoryCommitServiceServer will
// result in compilation errors.
type UnsafeRepositoryCommitServiceServer interface {
	mustEmbedUnimplementedRepositoryCommitServiceServer()
}

func RegisterRepositoryCommitServiceServer(s grpc.ServiceRegistrar, srv RepositoryCommitServiceServer) {
	s.RegisterService(&RepositoryCommitService_ServiceDesc, srv)
}

func _RepositoryCommitService_ListRepositoryCommitsByBranch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRepositoryCommitsByBranchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RepositoryCommitServiceServer).ListRepositoryCommitsByBranch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/buf.alpha.registry.v1alpha1.RepositoryCommitService/ListRepositoryCommitsByBranch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RepositoryCommitServiceServer).ListRepositoryCommitsByBranch(ctx, req.(*ListRepositoryCommitsByBranchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RepositoryCommitService_GetRepositoryCommitByReference_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRepositoryCommitByReferenceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RepositoryCommitServiceServer).GetRepositoryCommitByReference(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/buf.alpha.registry.v1alpha1.RepositoryCommitService/GetRepositoryCommitByReference",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RepositoryCommitServiceServer).GetRepositoryCommitByReference(ctx, req.(*GetRepositoryCommitByReferenceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RepositoryCommitService_GetRepositoryCommitBySequenceID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRepositoryCommitBySequenceIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RepositoryCommitServiceServer).GetRepositoryCommitBySequenceID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/buf.alpha.registry.v1alpha1.RepositoryCommitService/GetRepositoryCommitBySequenceID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RepositoryCommitServiceServer).GetRepositoryCommitBySequenceID(ctx, req.(*GetRepositoryCommitBySequenceIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RepositoryCommitService_ServiceDesc is the grpc.ServiceDesc for RepositoryCommitService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RepositoryCommitService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "buf.alpha.registry.v1alpha1.RepositoryCommitService",
	HandlerType: (*RepositoryCommitServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListRepositoryCommitsByBranch",
			Handler:    _RepositoryCommitService_ListRepositoryCommitsByBranch_Handler,
		},
		{
			MethodName: "GetRepositoryCommitByReference",
			Handler:    _RepositoryCommitService_GetRepositoryCommitByReference_Handler,
		},
		{
			MethodName: "GetRepositoryCommitBySequenceID",
			Handler:    _RepositoryCommitService_GetRepositoryCommitBySequenceID_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "buf/alpha/registry/v1alpha1/repository_commit.proto",
}
