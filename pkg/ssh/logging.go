package ssh

import (
	"context"
	"net"
	"time"

	"github.com/gliderlabs/ssh"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/sourcepods/sourcepods/pkg/ssh/mux"
)

type contextType string

var contextLogger = contextType("logger")

// loggerWare logs connections, and injects `logger` into the context.
func loggerWare(logger log.Logger) mux.MiddlewareFunc {
	return func(ctx context.Context, next mux.HandlerFunc, sess ssh.Session) error {
		sessID := ctx.Value(ssh.ContextKeySessionID).(string)
		start := time.Now()

		level.Debug(logger).Log(
			"msg", "new session",
			"user", sess.User(),
			"session-id", sessID,
			"remote-addr", sess.RemoteAddr().String(),
			"remote-addr", ctx.Value(ssh.ContextKeyRemoteAddr).(net.Addr).String(),
			"cmd", ctx.Value(mux.ContextCmd).(string),
		)
		ctx = context.WithValue(ctx, contextLogger, logger)

		err := next(ctx, sess)

		level.Debug(logger).Log(
			"err", err,
			"msg", "session closed",
			"user", sess.User(),
			"session-length", time.Now().Sub(start),
		)
		return err
	}

}
