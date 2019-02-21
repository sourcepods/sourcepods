package ssh

import (
	"context"
	"errors"

	"github.com/gliderlabs/ssh"
	"github.com/sourcepods/sourcepods/pkg/ssh/mux"
	"github.com/sourcepods/sourcepods/pkg/storage"
)

type gitStorage struct {
	client *storage.Client
}

func (gs gitStorage) parsePath(ctx context.Context) (string, error) {
	args, ok := ctx.Value(mux.ContextArguments).([]string)
	if !ok || len(args) < 1 || len(args[0]) == 0 {
		return "", errors.New("no path given")
	}
	return args[0], nil
}

func (gs *gitStorage) UploadPack(ctx context.Context, sess ssh.Session) error {
	path, err := gs.parsePath(ctx)
	if err != nil {
		return err
	}

	ec, err := gs.client.UploadPack(ctx, path, sess, sess, sess.Stderr())
	if err != nil {
		return err
	}
	if ec != 0 {
		return mux.ExitStatus{Code: int(ec), Err: err}
	}
	return nil

}

func (gs *gitStorage) ReceivePack(ctx context.Context, sess ssh.Session) error {
	path, err := gs.parsePath(ctx)
	if err != nil {
		return err
	}

	ec, err := gs.client.ReceivePack(ctx, path, sess, sess, sess.Stderr())
	if err != nil {
		return err
	}
	if ec != 0 {
		return mux.ExitStatus{Code: int(ec), Err: err}
	}
	return nil
}
