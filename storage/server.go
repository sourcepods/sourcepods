package storage

import (
	"context"

	"google.golang.org/grpc"
)

type storageServer struct {
	storage Storage
}

func NewStorageServer(server *grpc.Server, storage Storage) *grpc.Server {
	RegisterStorageServer(server, &storageServer{storage: storage})
	return server
}

func (s *storageServer) Create(ctx context.Context, req *CreateRequest) (*EmptyResponse, error) {
	return &EmptyResponse{}, s.storage.Create(req.GetOwner(), req.GetName())
}

func (s *storageServer) Descriptions(ctx context.Context, req *DescriptionRequest) (*EmptyResponse, error) {
	return &EmptyResponse{}, s.storage.Description(req.GetOwner(), req.GetName(), req.GetDescription())
}
