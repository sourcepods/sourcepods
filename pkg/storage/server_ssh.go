package storage

import (
	"strings"

	context "golang.org/x/net/context"
	"google.golang.org/grpc/status"

	opentracing "github.com/opentracing/opentracing-go"
	"gitlab.com/gitlab-org/gitaly/streamio"
	"google.golang.org/grpc/codes"
)

type sshService struct {
	storage Storage
}

// NOTE: #namingThings. And this does more than it should... for "simplicity"
func validateRepoGRERequest(ctx context.Context, s Storage, req *GRERequest, errIn error) (Repository, error) {
	if errIn != nil {
		return nil, status.Errorf(codes.Internal, errIn.Error())
	}
	if len(req.GetId()) == 0 {
		return nil, status.Errorf(codes.FailedPrecondition, "no repo id given")
	}
	if len(req.GetStdin()) != 0 {
		return nil, status.Errorf(codes.FailedPrecondition, "stdin given on first message")
	}
	id := strings.Replace(req.GetId(), "/", "", -1)
	repo, err := s.GetRepository(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "repo does not exist")
	}
	return repo, nil
}

func (s sshService) UploadPack(stream SSH_UploadPackServer) error {
	span, ctx := opentracing.StartSpanFromContext(stream.Context(), "storage.SSH.UploadPack")
	defer span.Finish()

	req, err := stream.Recv()
	repo, err := validateRepoGRERequest(ctx, s.storage, req, err)
	if err != nil {
		return err
	}

	span.SetTag("repo_hash", repo.GetID())

	ec, err := repo.UploadPack(ctx,
		streamio.NewReader(func() ([]byte, error) {
			request, err := stream.Recv()
			return request.GetStdin(), err
		}),
		streamio.NewWriter(func(p []byte) error {
			return stream.Send(&GREResponse{Stdout: p})
		}),
		streamio.NewWriter(func(p []byte) error {
			return stream.Send(&GREResponse{Stderr: p})
		}),
	)

	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}
	return stream.Send(&GREResponse{ExitCode: &GREExitCode{ExitCode: ec}})
}

func (s sshService) ReceivePack(stream SSH_ReceivePackServer) error {
	span, ctx := opentracing.StartSpanFromContext(stream.Context(), "storage.SSH.ReceivePack")
	defer span.Finish()

	req, err := stream.Recv()
	repo, err := validateRepoGRERequest(ctx, s.storage, req, err)
	if err != nil {
		return err
	}

	span.SetTag("repo_hash", repo.GetID())

	ec, err := repo.ReceivePack(ctx,
		streamio.NewReader(func() ([]byte, error) {
			request, err := stream.Recv()
			return request.GetStdin(), err
		}),
		streamio.NewWriter(func(p []byte) error {
			return stream.Send(&GREResponse{Stdout: p})
		}),
		streamio.NewWriter(func(p []byte) error {
			return stream.Send(&GREResponse{Stderr: p})
		}),
	)

	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}
	return stream.Send(&GREResponse{ExitCode: &GREExitCode{ExitCode: ec}})
}
