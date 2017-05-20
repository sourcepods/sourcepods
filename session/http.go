package session

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	// CookieName is the name used to store the cookie on the client
	CookieName         = "_gitpods_session"
	CookieUserID       = "user_id"
	CookieUserUsername = "user_username"
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

			ctx := context.WithValue(r.Context(), CookieUserID, session.User.ID)
			ctx = context.WithValue(ctx, CookieUserUsername, session.User.Username)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func GetSessionUser(r *http.Request) *SessionUser {
	ctx := r.Context()
	return &SessionUser{
		ID:       ctx.Value(CookieUserID).(string),
		Username: ctx.Value(CookieUserUsername).(string),
	}
}
