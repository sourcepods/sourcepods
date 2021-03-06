package api

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/satori/go.uuid"
)

type contextValue string

const contextReqID contextValue = "request_id"

//NewRequestID creates a middleware that passes X-Request-ID via context or creates a new ID
func NewRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.NewV4().String()
		}

		r = r.WithContext(context.WithValue(r.Context(), contextReqID, reqID))

		next.ServeHTTP(w, r)
	})
}

//GetRequestID returns the request ID from context.Context
func GetRequestID(ctx context.Context) string {
	value := ctx.Value(contextReqID)
	if reqID, ok := value.(string); ok {
		return reqID
	}
	return ""
}

//NewRequestLogger creates a middleware that logs all HTTP request information
func NewRequestLogger(logger log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)

			level.Debug(logger).Log(
				"request", GetRequestID(r.Context()),
				"proto", r.Proto,
				"method", r.Method,
				"status", ww.Status(),
				"path", r.URL.Path,
				"duration", time.Since(start),
				"bytes", ww.BytesWritten(),
			)
		})
	}
}
