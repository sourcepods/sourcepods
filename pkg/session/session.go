package session

import (
	"time"
)

const (
	defaultExpiry = 24 * time.Hour
)

type (
	// User only has an ID and a username
	// which is enough to find what you need from stores.
	User struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	}

	// Session has an ID, expiry and a User.
	Session struct {
		ID     string    `json:"id"`
		Expiry time.Time `json:"expiry"`
		User   User      `json:"user"`
	}
)
