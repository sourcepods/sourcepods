package storage

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	empty "github.com/golang/protobuf/ptypes/empty"
	grpcopentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"google.golang.org/grpc"
)

// NewStorageServer returns a grpc.Server serving Storage
func NewStorageServer(storage Storage) *grpc.Server {
	var opts []grpc.ServerOption
	opts = append(opts, grpc.UnaryInterceptor(grpcopentracing.UnaryServerInterceptor()))
	opts = append(opts, grpc.StreamInterceptor(grpcopentracing.StreamServerInterceptor()))

	s := grpc.NewServer(opts...)

	RegisterRepositoryServer(s, &repositoryServer{storage: storage})
	RegisterBranchServer(s, &branchesServer{storage: storage})
	RegisterCommitServer(s, &commitServer{storage: storage})
	RegisterSSHServer(s, &sshService{storage: storage})

	return s
}

type repositoryServer struct {
	storage Storage
}

func (s *repositoryServer) Create(ctx context.Context, req *CreateRequest) (*empty.Empty, error) {
	return &empty.Empty{}, s.storage.Create(ctx, req.GetId())
}

func (s *repositoryServer) SetDescriptions(ctx context.Context, req *SetDescriptionRequest) (*empty.Empty, error) {
	repo, err := s.storage.GetRepository(ctx, req.GetId())
	if err != nil {
		return nil, grpc.Errorf(codes.NotFound, "%v", err)
	}
	return &empty.Empty{}, repo.SetDescription(ctx, req.GetDescription())
}

type branchesServer struct {
	storage Storage
}

func (s *branchesServer) List(ctx context.Context, req *BranchesRequest) (*BranchesResponse, error) {
	repo, err := s.storage.GetRepository(ctx, req.GetId())
	if err != nil {
		return nil, grpc.Errorf(codes.NotFound, "%v", err)
	}
	branches, err := repo.ListBranches(ctx)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", err)
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

type commitServer struct {
	storage Storage
}

func (s *commitServer) Get(ctx context.Context, req *CommitRequest) (*CommitResponse, error) {
	repo, err := s.storage.GetRepository(ctx, req.GetId())
	if err != nil {
		return nil, grpc.Errorf(codes.NotFound, "%v", err)
	}
	c, err := repo.GetCommit(ctx, req.GetRef())
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", err)
	}

	return commitToResponse(c), nil

}

func (s *commitServer) Count(ctx context.Context, req *CommitCountRequest) (*CommitCountResponse, error) {
	if len(req.GetRef()) == 0 {
		return nil, status.Errorf(codes.FailedPrecondition, "no ref given")
	}
	repo, err := s.storage.GetRepository(ctx, req.GetId())
	if err != nil {
		return nil, grpc.Errorf(codes.NotFound, "%v", err)
	}
	count, err := repo.CountCommits(ctx, req.GetRef())
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", err)
	}

	return &CommitCountResponse{Count: count}, nil
}

func validateCommitListRequest(req *CommitListRequest) error {
	if len(req.GetRef()) == 0 {
		return status.Errorf(codes.FailedPrecondition, "no ref given")
	}
	if req.GetLimit() == 0 {
		return status.Errorf(codes.InvalidArgument, "limit can not be 0")
	}

	return nil
}

func (s *commitServer) List(req *CommitListRequest, stream Commit_ListServer) error {
	ctx := stream.Context()
	if err := validateCommitListRequest(req); err != nil {
		return err
	}
	repo, err := s.storage.GetRepository(ctx, req.GetId())
	if err != nil {
		return status.Errorf(codes.NotFound, "%v", err)
	}
	commits, err := repo.ListCommits(ctx, req.GetRef(), req.GetLimit(), req.GetSkip())
	if err != nil {
		return status.Errorf(codes.Internal, "%v", err)
	}
	commitResp := make([]*CommitResponse, len(commits), len(commits))
	for i, commit := range commits {
		commitResp[i] = commitToResponse(commit)
	}
	if err = stream.Send(&CommitListResponse{Commits: commitResp}); err != nil {
		return status.Errorf(codes.Internal, "stream.Send: %v", err)
	}

	return nil
}

func commitToResponse(c Commit) *CommitResponse {
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
	}
}

func (s *repositoryServer) Tree(ctx context.Context, req *TreeRequest) (*TreeResponse, error) {
	repo, err := s.storage.GetRepository(ctx, req.GetId())
	if err != nil {
		return nil, grpc.Errorf(codes.NotFound, "%v", err)
	}
	entries, err := repo.Tree(ctx, req.GetRef(), req.GetPath())
	if err != nil {
		return nil, err
	}

	var treeEntryRes []*TreeEntryResponse
	for _, e := range entries {
		treeEntryRes = append(treeEntryRes, &TreeEntryResponse{
			Mode:   e.Mode,
			Type:   e.Type,
			Object: e.Object,
			Path:   e.Path,
		})
	}

	return &TreeResponse{TreeEntries: treeEntryRes}, nil
}
