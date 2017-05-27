package user

import "time"

// User of GitPods.
type User struct {
	ID       string
	Email    string
	Username string
	Name     string
	Password string
	Created  time.Time
	Updated  time.Time
}
