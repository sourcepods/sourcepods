package repository

import "time"

// Repository containing source code.
type Repository struct {
	ID            string
	Name          string
	Description   string
	Website       string
	DefaultBranch string
	Created       time.Time
	Updated       time.Time
}

// Branch of a Repository.
type Branch struct {
	Name      string
	Sha1      string
	Type      string
	Protected bool
}
