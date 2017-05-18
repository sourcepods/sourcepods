package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/metrics"
	"github.com/gorilla/sessions"
	"github.com/pressly/chi"
)

var (
	JsonUnauthorized   = []byte(`{"message":"Unauthorized"}`) // 401
	JsonNotFound       = []byte(`{"error":"Not Found"}`)      // 404
	JsonBadCredentials = []byte(`{"message":"Bad credentials"}`)
)

type RouterStore struct {
	AuthorizeStore AuthorizeStore
	CookieStore    sessions.Store
	UsersStore     UsersStore
}

type RouterMetrics struct {
	LoginAttempts metrics.Counter
}

func NewRouter(logger log.Logger, metrics RouterMetrics, store *RouterStore) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/authorize", Authorize(logger, metrics.LoginAttempts, store.CookieStore, store.AuthorizeStore))

	apiAuthRouter := NewAuthRouter(logger, metrics, store)
	r.With(Authorized(logger, store.CookieStore)).Mount("/", apiAuthRouter)

	return r
}

func NewAuthRouter(logger log.Logger, metrics RouterMetrics, store *RouterStore) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/user", User(logger, store.UsersStore))

	users := &UsersAPI{Logger: logger, Store: store.UsersStore}
	r.Mount("/users", users.Routes())

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		jsonResponseBytes(w, JsonNotFound, http.StatusNotFound)
	})

	return r
}

func LoggerMiddleware(logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			level.Debug(logger).Log(
				"duration", time.Since(start),
				"method", r.Method,
				"path", r.URL.Path,
			)
		})
	}
}

func jsonResponse(w http.ResponseWriter, v interface{}, code int) {
	data, err := json.Marshal(v)
	if err != nil {
		http.Error(w, "failed to marshal to json", http.StatusInternalServerError)
		return
	}
	jsonResponseBytes(w, data, code)
}

func jsonResponseBytes(w http.ResponseWriter, payload []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(payload)
}
