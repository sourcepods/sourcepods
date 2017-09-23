package storage

import (
	"context"

	grpcopentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"google.golang.org/grpc"
)

type storageServer struct {
	storage Storage
}

func NewStorageServer(storage Storage) *grpc.Server {
	var opts []grpc.ServerOption
	opts = append(opts, grpc.UnaryInterceptor(grpcopentracing.UnaryServerInterceptor()))

	s := grpc.NewServer(opts...)

	RegisterStorageServer(s, &storageServer{storage: storage})
	return s
}

func (s *storageServer) Create(ctx context.Context, req *CreateRequest) (*EmptyResponse, error) {
	return &EmptyResponse{}, s.storage.Create(ctx, req.GetOwner(), req.GetName())
}

func (s *storageServer) Descriptions(ctx context.Context, req *DescriptionRequest) (*EmptyResponse, error) {
	return &EmptyResponse{}, s.storage.Description(ctx, req.GetOwner(), req.GetName(), req.GetDescription())
}
