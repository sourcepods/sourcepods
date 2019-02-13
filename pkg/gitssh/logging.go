package gitssh

import (
	"fmt"

	"github.com/gliderlabs/ssh"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	opentracing "github.com/opentracing/opentracing-go"
)

func logHandler(next ssh.Handler, logger log.Logger) ssh.Handler {
	return func(s ssh.Session) {
		defer s.Close()
		span, _ := opentracing.StartSpanFromContext(s.Context(), "ssh.MainHandler")
		span.SetTag("remote-addr", s.RemoteAddr().String())
		span.SetTag("user", s.User())
		defer span.Finish()

		level.Info(logger).Log(
			"msg", "new session",
			"user", s.User(),
			"remote-addr", s.RemoteAddr().String(),
			"command", fmt.Sprintf("%v", s.Command()),
		)

		next(s)

		level.Info(logger).Log(
			"msg", "session closed",
			"user", s.User(),
		)
	}
}
