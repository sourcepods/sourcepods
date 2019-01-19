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

func (s *storageServer) Branches(ctx context.Context, req *BranchesRequest) (*BranchesResponse, error) {
	branches, err := s.storage.Branches(ctx, req.GetOwner(), req.GetName())
	if err != nil {
		return nil, err
	}

	res := &BranchesResponse{}
	for _, b := range branches {
		res.Branch = append(res.Branch, &BranchResponse{
			Name: b.Name,
			Sha1: b.Sha1,
			Type: b.Type,
		})
	}

	return res, nil
}

func (s *storageServer) Commit(ctx context.Context, req *CommitRequest) (*CommitResponse, error) {
	c, err := s.storage.Commit(ctx, req.GetOwner(), req.GetName(), req.GetRev())
	if err != nil {
		return nil, err
	}

	return &CommitResponse{
		Hash:           c.Hash,
		Tree:           c.Tree,
		Parent:         c.Parent,
		Message:        c.Message,
		Author:         c.Author.Name,
		AuthorEmail:    c.Author.Email,
		AuthorDate:     c.Author.Date.Unix(),
		Committer:      c.Committer.Name,
		CommitterEmail: c.Committer.Email,
		CommitterDate:  c.Committer.Date.Unix(),
	}, nil
}
