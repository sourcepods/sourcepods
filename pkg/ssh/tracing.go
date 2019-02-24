package ssh

import (
	"context"

	"github.com/gliderlabs/ssh"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/sourcepods/sourcepods/pkg/ssh/mux"
)

// tracerWare traces connections, and injects
func tracerWare() mux.MiddlewareFunc {
	return func(ctx context.Context, next mux.HandlerFunc, sess ssh.Session) error {
		sessID := sess.Context().(ssh.Context).SessionID()

		span, ctx := opentracing.StartSpanFromContext(ctx, ctx.Value(mux.ContextHandlerName).(string))
		span.SetTag("remote-addr", sess.RemoteAddr().String())
		span.SetTag("user", sess.User())
		span.SetTag("session-id", sessID)
		args, ok := ctx.Value(mux.ContextArguments).([]string)
		if ok && len(args) > 0 {
			span.SetTag("args", ctx.Value(mux.ContextArguments).([]string))
		}
		defer span.Finish()

		err := next(ctx, sess)

		if err != nil {
			span.SetTag("error", true)
			span.LogKV("event", "error", "message", err)
		}
		return err
	}
}
