package api

import (
	"context"
	"net/http"

	"github.com/satori/go.uuid"
)

const contextReqID = "request_id"

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

func GetRequestID(ctx context.Context) string {
	value := ctx.Value(contextReqID)
	if reqID, ok := value.(string); ok {
		return reqID
	}
	return ""
}
