package storage

import (
	"strings"

	"google.golang.org/grpc/status"

	opentracing "github.com/opentracing/opentracing-go"
	"gitlab.com/gitlab-org/gitaly/streamio"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type sshService struct {
	storage Storage
}

func (s sshService) UploadPack(stream SSH_UploadPackServer) error {
	span, ctx := opentracing.StartSpanFromContext(stream.Context(), "storage.SSH.UploadPack")
	defer span.Finish()

	req, err := stream.Recv()
	if err != nil {
		return grpc.Errorf(codes.Internal, "recv: %v", err)
	}

	if req.GetId() == "" {
		return grpc.Errorf(codes.FailedPrecondition, "no repo id given")
	}
	span.SetTag("repo_path", req.GetId())

	id := strings.Replace(req.GetId(), "/", "", -1)
	span.SetTag("repo_hash", id)

	repo, err := s.storage.GetRepository(ctx, id)
	if err != nil {
		return grpc.Errorf(codes.FailedPrecondition, "repo does not exist")
	}

	if req.GetStdin() != nil {
		return grpc.Errorf(codes.FailedPrecondition, "stdin given on first message")
	}

	ec, err := repo.UploadPack(ctx,
		streamio.NewReader(func() ([]byte, error) {
			request, err := stream.Recv()
			return request.GetStdin(), err
		}),
		streamio.NewWriter(func(p []byte) error {
			span.LogEvent(string(p))
			return stream.Send(&GREResponse{Stdout: p})
		}),
		streamio.NewWriter(func(p []byte) error {
			span.LogEvent(string(p))
			return stream.Send(&GREResponse{Stderr: p})
		}),
	)

	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}
	if ec != 0 {
		return stream.Send(&GREResponse{ExitCode: &GREExitCode{ExitCode: ec}})
	}

	return stream.Send(&GREResponse{ExitCode: &GREExitCode{ExitCode: 0}})
}

func (s sshService) ReceivePack(stream SSH_ReceivePackServer) error {
	span, ctx := opentracing.StartSpanFromContext(stream.Context(), "storage.SSH.ReceivePack")
	defer span.Finish()

	req, err := stream.Recv()
	if err != nil {
		return grpc.Errorf(codes.Internal, "recv: %v", err)
	}

	if req.GetId() == "" {
		return grpc.Errorf(codes.FailedPrecondition, "no repo id given")
	}
	span.SetTag("repo_path", req.GetId())

	id := strings.Replace(req.GetId(), "/", "", -1)
	span.SetTag("repo_hash", id)

	repo, err := s.storage.GetRepository(ctx, id)
	if err != nil {
		return grpc.Errorf(codes.FailedPrecondition, "repo does not exist")
	}

	if req.GetStdin() != nil {
		return grpc.Errorf(codes.FailedPrecondition, "stdin given on first message")
	}

	ec, err := repo.ReceivePack(ctx,
		streamio.NewReader(func() ([]byte, error) {
			request, err := stream.Recv()
			span.LogEvent("stdin: " + string(request.GetStdin()))
			return request.GetStdin(), err
		}),
		streamio.NewWriter(func(p []byte) error {
			span.LogEvent("stdout: " + string(p))
			return stream.Send(&GREResponse{Stdout: p})
		}),
		streamio.NewWriter(func(p []byte) error {
			span.LogEvent("stderr: " + string(p))
			return stream.Send(&GREResponse{Stderr: p})
		}),
	)

	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}
	if ec != 0 {
		return stream.Send(&GREResponse{ExitCode: &GREExitCode{ExitCode: ec}})
	}

	return stream.Send(&GREResponse{ExitCode: &GREExitCode{ExitCode: 0}})
}
