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

func (s *storageServer) SetDescriptions(ctx context.Context, req *SetDescriptionRequest) (*EmptyResponse, error) {
	return &EmptyResponse{}, s.storage.SetDescription(ctx, req.GetOwner(), req.GetName(), req.GetDescription())
}

func (s *storageServer) Tree(ctx context.Context, req *TreeRequest) (*TreeRespone, error) {
	objects, err := s.storage.Tree(ctx, req.GetOwner(), req.GetName(), req.GetBranch(), req.GetRecursive())
	if err != nil {
		return nil, err
	}

	res := &TreeRespone{}
	for _, object := range objects {
		res.Objects = append(res.Objects, &TreeObjectResponse{
			Mode:   object.Mode,
			Type:   object.Type,
			Object: object.Object,
			File:   object.File,
			Commit: &CommitResponse{
				Hash:           object.Commit.Hash,
				Tree:           object.Commit.Tree,
				Parent:         object.Commit.Parent,
				Subject:        object.Commit.Subject,
				Author:         object.Commit.Author,
				AuthorEmail:    object.Commit.AuthorEmail,
				AuthorDate:     object.Commit.AuthorDate.Unix(),
				Committer:      object.Commit.Committer,
				CommitterEmail: object.Commit.CommitterEmail,
				CommitterDate:  object.Commit.CommitterDate.Unix(),
			},
		})
	}

	return res, nil
}
