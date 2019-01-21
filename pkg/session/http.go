package session

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/jsonapi"
	"github.com/opentracing/opentracing-go"
)

type ctxKey int

const (
	// CookieName is the name to store the cookie in the browser with.
	CookieName                = "_sourcepods_session"
	CookieUserID       ctxKey = iota
	CookieUserUsername ctxKey = iota
)

var (
	errUnauthorized = []*jsonapi.ErrorObject{{
		Title:  http.StatusText(http.StatusUnauthorized),
		Detail: "Your Cookie is not valid",
		Status: fmt.Sprintf("%d", http.StatusUnauthorized),
	}}
)

// Authorized users will have a user information in the next handlers.
func Authorized(s Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			span, ctx := opentracing.StartSpanFromContext(r.Context(), "session.Handler.Authorized")
			defer span.Finish()

			cookie, err := r.Cookie(CookieName)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				jsonapi.MarshalErrors(w, errUnauthorized)
				return
			}

			if cookie.Value == "" {
				w.WriteHeader(http.StatusUnauthorized)
				jsonapi.MarshalErrors(w, errUnauthorized)
				return
			}

			if time.Now().Before(cookie.Expires) {
				w.WriteHeader(http.StatusUnauthorized)
				jsonapi.MarshalErrors(w, errUnauthorized)
				return
			}

			session, err := s.Find(ctx, cookie.Value)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				jsonapi.MarshalErrors(w, errUnauthorized)
				return
			}

			ctx = context.WithValue(ctx, CookieUserID, session.User.ID)
			ctx = context.WithValue(ctx, CookieUserUsername, session.User.Username)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// GetSessionUser from the http.Request
func GetSessionUser(ctx context.Context) *User {
	return &User{
		ID:       ctx.Value(CookieUserID).(string),
		Username: ctx.Value(CookieUserUsername).(string),
	}
}

func NewHandler(s Service) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/logout", logout(s))

	return r
}

func logout(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(CookieName)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			jsonapi.MarshalErrors(w, errUnauthorized)
			return
		}

		if err := s.Delete(r.Context(), cookie.Value); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
				Title:  http.StatusText(http.StatusInternalServerError),
				Detail: "Your session couldn't be terminated",
				Status: fmt.Sprintf("%d", http.StatusInternalServerError),
			}})
			return
		}

		// TODO: Fix the url
		http.Redirect(w, r, "https://try.gitpods.io/", http.StatusFound)
	}
}
