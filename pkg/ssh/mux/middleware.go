package mux

import (
	"context"

	"github.com/gliderlabs/ssh"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type (
	// MiddlewareFunc is a middleware
	MiddlewareFunc func(context.Context, HandlerFunc, ssh.Session) error
)

// RecoverWare is a Middleware that does panic recovery
func RecoverWare(logger log.Logger) MiddlewareFunc {
	return func(ctx context.Context, h HandlerFunc, sess ssh.Session) error {
		defer func() {
			if r := recover(); r != nil {
				// TODO: Logging :thinking:
				level.Error(logger).Log(
					"msg", "handler paniced",
					"err", r,
				)
			}
		}()

		return h(ctx, sess)
	}
}
