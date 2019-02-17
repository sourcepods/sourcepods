package mux

import (
	"context"
	"regexp"
	"strings"

	"github.com/gliderlabs/ssh"
	"github.com/pkg/errors"
)

var (
	// ContextCmd points to the command-line in the Context
	ContextCmd = contextType("mux-cmd")
	// ContextPattern points to the pattern that matched the Handler
	ContextPattern = contextType("mux-pattern")
	// ContextArguments points to the arguments matched by the pattern
	ContextArguments = contextType("mux-args")
	// ContextHandlerName points to the handlers name. Good for logging and/or tracing
	ContextHandlerName = contextType("mux-handler-name")
)

type (
	contextType string
	// Muxer is a simple SSH Handler muxer
	Muxer interface {
		// Use, inserts one or more middleware(s).
		//  They are executed in FIFE (First In, First Executed).
		//  This means the last one in the list is the one executing
		//  closest to the Handler.
		Use(mws ...MiddlewareFunc)

		// AddHandler to the Muxer
		//  `pattern` is a regexp Pattern, submatches will be injected
		//  into the context under ContextArgument
		AddHandler(pattern string, name string, h Handler) error

		// Handle a given session. This should be given to gliderlabs/ssh.Server.Handler()
		Handle() ssh.Handler
	}

	reHandler struct {
		handler Handler
		regexp  *regexp.Regexp
		name    string
	}

	muxer struct {
		handlers map[string]reHandler
		mws      []MiddlewareFunc
	}
)

// New creates a new Muxer
func New() Muxer {
	return &muxer{
		handlers: make(map[string]reHandler),
	}
}

func (s *muxer) Use(mws ...MiddlewareFunc) {
	// We prepend the middlewares instead, so that
	//  Use(Foo, Bar) results in Foo(Bar(Handler))
	//  being ran, instead of Bar(Foo(Handler))
	// NOTE: This is probably somewhat inefficient
	//  (but only ran at setup, so who cares)
	for _, mw := range mws {
		s.mws = append([]MiddlewareFunc{mw}, s.mws...)
	}
}

func (s *muxer) AddHandler(pattern, name string, h Handler) error {
	if _, ok := s.handlers[pattern]; ok {
		return ErrExists
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return errors.Wrap(err, "invalid pattern")
	}

	s.handlers[pattern] = reHandler{
		handler: h,
		regexp:  re,
		name:    name,
	}
	return nil
}

func (s *muxer) Handle() ssh.Handler {
	return func(sess ssh.Session) {
		cmd := strings.Join(sess.Command(), " ")
		ctx := context.WithValue(sess.Context(), ContextCmd, cmd)

		errorCode := 0
		ctx, handler := s.match(ctx, cmd)
		err := s.wrapMiddlewares(ctx, handler, sess)
		if err != nil {
			// Any logging of this error should be done in a Middleware
			errorCode = 1
			if exitStatus, ok := err.(*ExitStatus); ok {
				errorCode = exitStatus.Code
			}
		}
		sess.Exit(errorCode)
		return
	}
}

func (s *muxer) match(ctx context.Context, cmd string) (context.Context, Handler) {
	ctx = context.WithValue(ctx, ContextHandlerName, "ssh.Handler.Unknown")
	if len(cmd) == 0 {
		return ctx, noCommandHandler
	}
	for _, handler := range s.handlers {
		matches := handler.regexp.FindStringSubmatch(cmd)
		if len(matches) == 0 {
			continue
		}
		ctx = context.WithValue(ctx, ContextHandlerName, handler.name)
		ctx = context.WithValue(ctx, ContextArguments, matches[1:])
		return ctx, handler.handler
	}

	return ctx, unknownCommandHandler
}

func (s *muxer) wrapMiddlewares(ctx context.Context, serv Handler, sess ssh.Session) error {
	f := func(ctx context.Context, sess ssh.Session) error {
		return serv.Execute(ctx, sess)
	}

	for _, mw := range s.mws {
		f = func(f2 HandlerFunc, mw2 MiddlewareFunc) HandlerFunc {
			return func(ctx context.Context, sess ssh.Session) error {
				return mw2(ctx, f2, sess)
			}
		}(f, mw)
	}

	return f(ctx, sess)
}
