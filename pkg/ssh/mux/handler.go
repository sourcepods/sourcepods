package mux

import (
	"context"
	"fmt"

	"github.com/gliderlabs/ssh"
	"github.com/pkg/errors"
)

var (
	// ErrExists returned if handler is already registered
	ErrExists = errors.New("service already exists")

	noCommandHandler = HandlerFunc(func(ctx context.Context, sess ssh.Session) error {
		fmt.Fprintf(sess, "Welcome %s\n", sess.User())
		sess.Exit(0)
		return nil
	})

	unknownCommandHandler = HandlerFunc(func(ctx context.Context, sess ssh.Session) error {
		fmt.Fprintf(sess, "Welcome %s\n", sess.User())
		fmt.Fprintf(sess, "  Unknown command given\n")
		sess.Exit(1)
		return nil
	})
)

type (

	// HandlerFunc is a loose service function
	HandlerFunc func(context.Context, ssh.Session) error

	// Handler defines an executable handler
	Handler interface {
		// Execute is the function called by mux.Handle()
		Execute(ctx context.Context, sess ssh.Session) error
	}
)

// Execute makes HandlerFunc implement Handler
func (h HandlerFunc) Execute(ctx context.Context, sess ssh.Session) error {
	return h(ctx, sess)
}
