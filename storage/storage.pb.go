// Code generated by protoc-gen-go. DO NOT EDIT.
// source: storage.proto

/*
Package storage is a generated protocol buffer package.

It is generated from these files:
	storage.proto

It has these top-level messages:
	EmptyResponse
	CreateRequest
	SetDescriptionRequest
	BranchesRequest
	BranchObject
	BranchesResponse
	TreeRequest
	TreeRespone
	TreeObjectResponse
	CommitResponse
*/
package storage

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type EmptyResponse struct {
}

func (m *EmptyResponse) Reset()                    { *m = EmptyResponse{} }
func (m *EmptyResponse) String() string            { return proto.CompactTextString(m) }
func (*EmptyResponse) ProtoMessage()               {}
func (*EmptyResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type CreateRequest struct {
	Owner string `protobuf:"bytes,1,opt,name=owner" json:"owner,omitempty"`
	Name  string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
}

func (m *CreateRequest) Reset()                    { *m = CreateRequest{} }
func (m *CreateRequest) String() string            { return proto.CompactTextString(m) }
func (*CreateRequest) ProtoMessage()               {}
func (*CreateRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *CreateRequest) GetOwner() string {
	if m != nil {
		return m.Owner
	}
	return ""
}

func (m *CreateRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type SetDescriptionRequest struct {
	Owner       string `protobuf:"bytes,1,opt,name=owner" json:"owner,omitempty"`
	Name        string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	Description string `protobuf:"bytes,3,opt,name=description" json:"description,omitempty"`
}

func (m *SetDescriptionRequest) Reset()                    { *m = SetDescriptionRequest{} }
func (m *SetDescriptionRequest) String() string            { return proto.CompactTextString(m) }
func (*SetDescriptionRequest) ProtoMessage()               {}
func (*SetDescriptionRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *SetDescriptionRequest) GetOwner() string {
	if m != nil {
		return m.Owner
	}
	return ""
}

func (m *SetDescriptionRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *SetDescriptionRequest) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

type BranchesRequest struct {
	Owner string `protobuf:"bytes,1,opt,name=owner" json:"owner,omitempty"`
	Name  string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
}

func (m *BranchesRequest) Reset()                    { *m = BranchesRequest{} }
func (m *BranchesRequest) String() string            { return proto.CompactTextString(m) }
func (*BranchesRequest) ProtoMessage()               {}
func (*BranchesRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *BranchesRequest) GetOwner() string {
	if m != nil {
		return m.Owner
	}
	return ""
}

func (m *BranchesRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type BranchObject struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Sha1 string `protobuf:"bytes,2,opt,name=sha1" json:"sha1,omitempty"`
	Type string `protobuf:"bytes,3,opt,name=type" json:"type,omitempty"`
}

func (m *BranchObject) Reset()                    { *m = BranchObject{} }
func (m *BranchObject) String() string            { return proto.CompactTextString(m) }
func (*BranchObject) ProtoMessage()               {}
func (*BranchObject) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *BranchObject) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *BranchObject) GetSha1() string {
	if m != nil {
		return m.Sha1
	}
	return ""
}

func (m *BranchObject) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

type BranchesResponse struct {
	Branch []*BranchObject `protobuf:"bytes,1,rep,name=branch" json:"branch,omitempty"`
}

func (m *BranchesResponse) Reset()                    { *m = BranchesResponse{} }
func (m *BranchesResponse) String() string            { return proto.CompactTextString(m) }
func (*BranchesResponse) ProtoMessage()               {}
func (*BranchesResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *BranchesResponse) GetBranch() []*BranchObject {
	if m != nil {
		return m.Branch
	}
	return nil
}

type TreeRequest struct {
	Owner     string `protobuf:"bytes,1,opt,name=owner" json:"owner,omitempty"`
	Name      string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	Branch    string `protobuf:"bytes,3,opt,name=branch" json:"branch,omitempty"`
	Recursive bool   `protobuf:"varint,4,opt,name=recursive" json:"recursive,omitempty"`
}

func (m *TreeRequest) Reset()                    { *m = TreeRequest{} }
func (m *TreeRequest) String() string            { return proto.CompactTextString(m) }
func (*TreeRequest) ProtoMessage()               {}
func (*TreeRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *TreeRequest) GetOwner() string {
	if m != nil {
		return m.Owner
	}
	return ""
}

func (m *TreeRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *TreeRequest) GetBranch() string {
	if m != nil {
		return m.Branch
	}
	return ""
}

func (m *TreeRequest) GetRecursive() bool {
	if m != nil {
		return m.Recursive
	}
	return false
}

type TreeRespone struct {
	Objects []*TreeObjectResponse `protobuf:"bytes,1,rep,name=objects" json:"objects,omitempty"`
}

func (m *TreeRespone) Reset()                    { *m = TreeRespone{} }
func (m *TreeRespone) String() string            { return proto.CompactTextString(m) }
func (*TreeRespone) ProtoMessage()               {}
func (*TreeRespone) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *TreeRespone) GetObjects() []*TreeObjectResponse {
	if m != nil {
		return m.Objects
	}
	return nil
}

type TreeObjectResponse struct {
	Mode   string          `protobuf:"bytes,1,opt,name=Mode" json:"Mode,omitempty"`
	Type   string          `protobuf:"bytes,2,opt,name=Type" json:"Type,omitempty"`
	Object string          `protobuf:"bytes,3,opt,name=Object" json:"Object,omitempty"`
	File   string          `protobuf:"bytes,4,opt,name=File" json:"File,omitempty"`
	Commit *CommitResponse `protobuf:"bytes,5,opt,name=Commit" json:"Commit,omitempty"`
}

func (m *TreeObjectResponse) Reset()                    { *m = TreeObjectResponse{} }
func (m *TreeObjectResponse) String() string            { return proto.CompactTextString(m) }
func (*TreeObjectResponse) ProtoMessage()               {}
func (*TreeObjectResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *TreeObjectResponse) GetMode() string {
	if m != nil {
		return m.Mode
	}
	return ""
}

func (m *TreeObjectResponse) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *TreeObjectResponse) GetObject() string {
	if m != nil {
		return m.Object
	}
	return ""
}

func (m *TreeObjectResponse) GetFile() string {
	if m != nil {
		return m.File
	}
	return ""
}

func (m *TreeObjectResponse) GetCommit() *CommitResponse {
	if m != nil {
		return m.Commit
	}
	return nil
}

type CommitResponse struct {
	Hash           string `protobuf:"bytes,1,opt,name=Hash" json:"Hash,omitempty"`
	Tree           string `protobuf:"bytes,2,opt,name=Tree" json:"Tree,omitempty"`
	Parent         string `protobuf:"bytes,3,opt,name=Parent" json:"Parent,omitempty"`
	Subject        string `protobuf:"bytes,4,opt,name=Subject" json:"Subject,omitempty"`
	Author         string `protobuf:"bytes,5,opt,name=Author" json:"Author,omitempty"`
	AuthorEmail    string `protobuf:"bytes,6,opt,name=AuthorEmail" json:"AuthorEmail,omitempty"`
	AuthorDate     int64  `protobuf:"varint,7,opt,name=AuthorDate" json:"AuthorDate,omitempty"`
	Committer      string `protobuf:"bytes,8,opt,name=Committer" json:"Committer,omitempty"`
	CommitterEmail string `protobuf:"bytes,9,opt,name=CommitterEmail" json:"CommitterEmail,omitempty"`
	CommitterDate  int64  `protobuf:"varint,10,opt,name=CommitterDate" json:"CommitterDate,omitempty"`
}

func (m *CommitResponse) Reset()                    { *m = CommitResponse{} }
func (m *CommitResponse) String() string            { return proto.CompactTextString(m) }
func (*CommitResponse) ProtoMessage()               {}
func (*CommitResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *CommitResponse) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

func (m *CommitResponse) GetTree() string {
	if m != nil {
		return m.Tree
	}
	return ""
}

func (m *CommitResponse) GetParent() string {
	if m != nil {
		return m.Parent
	}
	return ""
}

func (m *CommitResponse) GetSubject() string {
	if m != nil {
		return m.Subject
	}
	return ""
}

func (m *CommitResponse) GetAuthor() string {
	if m != nil {
		return m.Author
	}
	return ""
}

func (m *CommitResponse) GetAuthorEmail() string {
	if m != nil {
		return m.AuthorEmail
	}
	return ""
}

func (m *CommitResponse) GetAuthorDate() int64 {
	if m != nil {
		return m.AuthorDate
	}
	return 0
}

func (m *CommitResponse) GetCommitter() string {
	if m != nil {
		return m.Committer
	}
	return ""
}

func (m *CommitResponse) GetCommitterEmail() string {
	if m != nil {
		return m.CommitterEmail
	}
	return ""
}

func (m *CommitResponse) GetCommitterDate() int64 {
	if m != nil {
		return m.CommitterDate
	}
	return 0
}

func init() {
	proto.RegisterType((*EmptyResponse)(nil), "storage.EmptyResponse")
	proto.RegisterType((*CreateRequest)(nil), "storage.CreateRequest")
	proto.RegisterType((*SetDescriptionRequest)(nil), "storage.SetDescriptionRequest")
	proto.RegisterType((*BranchesRequest)(nil), "storage.BranchesRequest")
	proto.RegisterType((*BranchObject)(nil), "storage.BranchObject")
	proto.RegisterType((*BranchesResponse)(nil), "storage.BranchesResponse")
	proto.RegisterType((*TreeRequest)(nil), "storage.TreeRequest")
	proto.RegisterType((*TreeRespone)(nil), "storage.TreeRespone")
	proto.RegisterType((*TreeObjectResponse)(nil), "storage.TreeObjectResponse")
	proto.RegisterType((*CommitResponse)(nil), "storage.CommitResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Storage service

type StorageClient interface {
	Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*EmptyResponse, error)
	SetDescriptions(ctx context.Context, in *SetDescriptionRequest, opts ...grpc.CallOption) (*EmptyResponse, error)
	Branches(ctx context.Context, in *BranchesRequest, opts ...grpc.CallOption) (*BranchesResponse, error)
	Tree(ctx context.Context, in *TreeRequest, opts ...grpc.CallOption) (*TreeRespone, error)
}

type storageClient struct {
	cc *grpc.ClientConn
}

func NewStorageClient(cc *grpc.ClientConn) StorageClient {
	return &storageClient{cc}
}

func (c *storageClient) Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := grpc.Invoke(ctx, "/storage.Storage/Create", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storageClient) SetDescriptions(ctx context.Context, in *SetDescriptionRequest, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := grpc.Invoke(ctx, "/storage.Storage/SetDescriptions", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storageClient) Branches(ctx context.Context, in *BranchesRequest, opts ...grpc.CallOption) (*BranchesResponse, error) {
	out := new(BranchesResponse)
	err := grpc.Invoke(ctx, "/storage.Storage/Branches", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storageClient) Tree(ctx context.Context, in *TreeRequest, opts ...grpc.CallOption) (*TreeRespone, error) {
	out := new(TreeRespone)
	err := grpc.Invoke(ctx, "/storage.Storage/Tree", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Storage service

type StorageServer interface {
	Create(context.Context, *CreateRequest) (*EmptyResponse, error)
	SetDescriptions(context.Context, *SetDescriptionRequest) (*EmptyResponse, error)
	Branches(context.Context, *BranchesRequest) (*BranchesResponse, error)
	Tree(context.Context, *TreeRequest) (*TreeRespone, error)
}

func RegisterStorageServer(s *grpc.Server, srv StorageServer) {
	s.RegisterService(&_Storage_serviceDesc, srv)
}

func _Storage_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/storage.Storage/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServer).Create(ctx, req.(*CreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Storage_SetDescriptions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetDescriptionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServer).SetDescriptions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/storage.Storage/SetDescriptions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServer).SetDescriptions(ctx, req.(*SetDescriptionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Storage_Branches_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BranchesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServer).Branches(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/storage.Storage/Branches",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServer).Branches(ctx, req.(*BranchesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Storage_Tree_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TreeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServer).Tree(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/storage.Storage/Tree",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServer).Tree(ctx, req.(*TreeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Storage_serviceDesc = grpc.ServiceDesc{
	ServiceName: "storage.Storage",
	HandlerType: (*StorageServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _Storage_Create_Handler,
		},
		{
			MethodName: "SetDescriptions",
			Handler:    _Storage_SetDescriptions_Handler,
		},
		{
			MethodName: "Branches",
			Handler:    _Storage_Branches_Handler,
		},
		{
			MethodName: "Tree",
			Handler:    _Storage_Tree_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "storage.proto",
}

func init() { proto.RegisterFile("storage.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 535 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x94, 0xcd, 0x6e, 0x13, 0x31,
	0x10, 0xc7, 0xb5, 0x49, 0x9a, 0x8f, 0x09, 0x21, 0x68, 0xd4, 0x16, 0x53, 0x50, 0xb5, 0x5a, 0x21,
	0x94, 0x0b, 0x45, 0x04, 0x21, 0x81, 0x38, 0xa0, 0xd2, 0x14, 0x01, 0x12, 0x02, 0x6d, 0xfa, 0x02,
	0x9b, 0x74, 0x44, 0x16, 0x75, 0xd7, 0xc1, 0x76, 0x40, 0x3d, 0xf3, 0x14, 0xbc, 0x05, 0x8f, 0x88,
	0xfc, 0xb9, 0xbb, 0xa1, 0x3d, 0xe4, 0x36, 0xf3, 0xf7, 0x8c, 0xe7, 0xe7, 0xb1, 0xc7, 0x30, 0x92,
	0x8a, 0x8b, 0xec, 0x1b, 0x9d, 0xac, 0x05, 0x57, 0x1c, 0x7b, 0xce, 0x4d, 0xc6, 0x30, 0x3a, 0x2f,
	0xd6, 0xea, 0x3a, 0x25, 0xb9, 0xe6, 0xa5, 0xa4, 0xe4, 0x35, 0x8c, 0xce, 0x04, 0x65, 0x8a, 0x52,
	0xfa, 0xb1, 0x21, 0xa9, 0x70, 0x1f, 0xf6, 0xf8, 0xaf, 0x92, 0x04, 0x8b, 0xe2, 0x68, 0x32, 0x48,
	0xad, 0x83, 0x08, 0x9d, 0x32, 0x2b, 0x88, 0xb5, 0x8c, 0x68, 0xec, 0x64, 0x09, 0x07, 0x73, 0x52,
	0x33, 0x92, 0x4b, 0x91, 0xaf, 0x55, 0xce, 0xcb, 0x9d, 0xb7, 0xc0, 0x18, 0x86, 0x97, 0x55, 0x3e,
	0x6b, 0x9b, 0xa5, 0xba, 0x94, 0xbc, 0x81, 0xf1, 0x3b, 0x91, 0x95, 0xcb, 0x15, 0xc9, 0xdd, 0x09,
	0x3f, 0xc1, 0x1d, 0x9b, 0xfc, 0x65, 0xf1, 0x9d, 0x96, 0x2a, 0xc4, 0x44, 0x35, 0x04, 0x84, 0x8e,
	0x5c, 0x65, 0xcf, 0x7d, 0x9e, 0xb6, 0xb5, 0xa6, 0xae, 0xd7, 0xe4, 0x78, 0x8c, 0x9d, 0x9c, 0xc2,
	0xbd, 0x0a, 0xc4, 0x36, 0x0f, 0x9f, 0x42, 0x77, 0x61, 0x34, 0x16, 0xc5, 0xed, 0xc9, 0x70, 0x7a,
	0x70, 0xe2, 0xdb, 0x5e, 0x2f, 0x9b, 0xba, 0xa0, 0xa4, 0x80, 0xe1, 0x85, 0xa0, 0xdd, 0x3b, 0x8d,
	0x87, 0xa1, 0x8e, 0x25, 0x72, 0x1e, 0x3e, 0x82, 0x81, 0xa0, 0xe5, 0x46, 0xc8, 0xfc, 0x27, 0xb1,
	0x4e, 0x1c, 0x4d, 0xfa, 0x69, 0x25, 0x24, 0x33, 0x5f, 0x4e, 0xd3, 0x12, 0xbe, 0x84, 0x1e, 0x37,
	0x3c, 0xd2, 0xd1, 0x3e, 0x0c, 0xb4, 0x3a, 0xcc, 0xb1, 0xba, 0xa3, 0xa5, 0x3e, 0x36, 0xf9, 0x13,
	0x01, 0xfe, 0xbf, 0xae, 0x31, 0x3f, 0xf3, 0xcb, 0xd0, 0x4a, 0x6d, 0x6b, 0xed, 0x42, 0xb7, 0xcd,
	0xa1, 0x6b, 0x5b, 0xa3, 0xdb, 0x4c, 0x8f, 0x5e, 0x5d, 0xc5, 0xfb, 0xfc, 0xca, 0x52, 0x0f, 0x52,
	0x63, 0xe3, 0x33, 0xe8, 0x9e, 0xf1, 0xa2, 0xc8, 0x15, 0xdb, 0x8b, 0xa3, 0xc9, 0x70, 0x7a, 0x3f,
	0x00, 0x5a, 0x39, 0xc0, 0xb9, 0xb0, 0xe4, 0x6f, 0x0b, 0xee, 0x36, 0x97, 0xf4, 0xbe, 0x1f, 0x32,
	0xb9, 0xf2, 0x5c, 0xda, 0x36, 0x5c, 0x82, 0x2a, 0x2e, 0x41, 0x86, 0xeb, 0x6b, 0x26, 0xa8, 0x0c,
	0x5c, 0xd6, 0x43, 0x06, 0xbd, 0xf9, 0xc6, 0x02, 0x5b, 0x34, 0xef, 0xea, 0x8c, 0xd3, 0x8d, 0x5a,
	0x71, 0x61, 0xe8, 0x06, 0xa9, 0xf3, 0xf4, 0x1b, 0xb6, 0xd6, 0x79, 0x91, 0xe5, 0x57, 0xac, 0x6b,
	0xdf, 0x70, 0x4d, 0xc2, 0x63, 0x00, 0xeb, 0xce, 0x32, 0x45, 0xac, 0x17, 0x47, 0x93, 0x76, 0x5a,
	0x53, 0xf4, 0x35, 0xda, 0x53, 0x28, 0x12, 0xac, 0x6f, 0xf2, 0x2b, 0x01, 0x9f, 0xf8, 0x33, 0x2a,
	0x72, 0x25, 0x06, 0x26, 0x64, 0x4b, 0xc5, 0xc7, 0x30, 0x0a, 0x8a, 0x29, 0x04, 0xa6, 0x50, 0x53,
	0x9c, 0xfe, 0x6e, 0x41, 0x6f, 0x6e, 0xbb, 0x8a, 0xaf, 0xa0, 0x6b, 0x67, 0x1f, 0x0f, 0xab, 0x4e,
	0xd7, 0x3f, 0x83, 0xa3, 0x4a, 0x6f, 0xfc, 0x1a, 0xf8, 0x11, 0xc6, 0xcd, 0xd1, 0x97, 0x78, 0x1c,
	0x42, 0x6f, 0xfc, 0x14, 0x6e, 0xdd, 0xea, 0x2d, 0xf4, 0xfd, 0x5c, 0x21, 0xdb, 0x9a, 0x9f, 0x30,
	0xf3, 0x47, 0x0f, 0x6e, 0x58, 0x71, 0x1b, 0x4c, 0xed, 0xed, 0xe2, 0x7e, 0xe3, 0x39, 0xfb, 0xc4,
	0x6d, 0xd5, 0xcc, 0xc2, 0xa2, 0x6b, 0xbe, 0xc5, 0x17, 0xff, 0x02, 0x00, 0x00, 0xff, 0xff, 0x4b,
	0xbd, 0xbd, 0xa3, 0x27, 0x05, 0x00, 0x00,
}
