package session

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type ctxKey int

const (
	// CookieName is the name to store the cookie in the browser with.
	CookieName                = "_gitpods_session"
	cookieUserID       ctxKey = iota
	cookieUserUsername ctxKey = iota
)

// Authorized users will have a user information in the next handlers.
func Authorized(s Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(CookieName)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintln(w, err)
				return // TODO
			}

			if cookie.Value == "" {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintln(w, err)
				return // TODO
			}

			if time.Now().Before(cookie.Expires) {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintln(w, "your cookie is expired")
				return // TODO
			}

			session, err := s.FindSession(cookie.Value)
			if err != nil {
				println(err)
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintln(w, err)
				return // TODO
			}

			ctx := context.WithValue(r.Context(), cookieUserID, session.User.ID)
			ctx = context.WithValue(ctx, cookieUserUsername, session.User.Username)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// GetSessionUser from the http.Request
func GetSessionUser(r *http.Request) *User {
	ctx := r.Context()
	return &User{
		ID:       ctx.Value(cookieUserID).(string),
		Username: ctx.Value(cookieUserUsername).(string),
	}
}
