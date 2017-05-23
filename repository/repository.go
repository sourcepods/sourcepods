package repository

import "time"

// Repository is a git repository with some meta information for gitpods.
type Repository struct {
	ID            string    `json:"id"`   //valid:"required,uuidv4"
	Name          string    `json:"name"` //valid:"required"
	Description   string    `json:"description"`
	Website       string    `json:"website"`
	DefaultBranch string    `json:"default_branch"`
	Private       bool      `json:"private"`
	Bare          bool      `json:"bare"`
	Created       time.Time `json:"created_at"`
	Updated       time.Time `json:"updated_at"`
}

type RepositoryAggregate struct {
	*Repository
	Stars int `json:"stars"`
	Forks int `json:"forks"`
}
