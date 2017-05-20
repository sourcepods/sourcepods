package session

import (
	"time"
)

const (
	DefaultSessionDuration = 24 * time.Hour
)

type SessionUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type Session struct {
	ID     string      `json:"id"`
	Expiry time.Time   `json:"expiry"`
	User   SessionUser `json:"user"`
}
