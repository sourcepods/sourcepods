package gitssh

import (
	"fmt"
	"time"

	"github.com/gliderlabs/ssh"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	opentracing "github.com/opentracing/opentracing-go"
)

// logHandler logs connections, and injects `span-ctx` and `logger` into the context.
func logHandler(next ssh.Handler, logger log.Logger) ssh.Handler {
	return func(s ssh.Session) {
		sessID := s.Context().(ssh.Context).SessionID()
		defer s.Close()
		span, spanCtx := opentracing.StartSpanFromContext(s.Context(), "ssh.Handler")
		s.Context().(ssh.Context).SetValue("span-ctx", spanCtx)
		span.SetTag("remote-addr", s.RemoteAddr().String())
		span.SetTag("user", s.User())
		span.SetTag("session-id", sessID)
		defer span.Finish()

		start := time.Now()

		level.Info(logger).Log(
			"msg", "new session",
			"user", s.User(),
			"session-id", sessID,
			"remote-addr", s.RemoteAddr().String(),
			"command", fmt.Sprintf("%v", s.Command()),
		)
		s.Context().(ssh.Context).SetValue("logger", logger)

		next(s)

		level.Info(logger).Log(
			"msg", "session closed",
			"user", s.User(),
			"session-length", time.Now().Sub(start),
		)
	}
}
