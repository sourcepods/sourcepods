package repository

import "time"

// Tree is a git repository with some meta information for gitpods.
type Repository struct {
	ID            string
	Name          string
	Description   string
	Website       string
	DefaultBranch string
	Created       time.Time
	Updated       time.Time
}

type Branch struct {
	Name      string
	Sha1      string
	Type      string
	Protected bool
}
