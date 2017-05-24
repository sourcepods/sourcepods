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

type Stats struct {
	Stars            int               `json:"stars"`
	Forks            int               `json:"forks"`
	IssueStats       *IssueStats       `json:"issue_stats"`
	PullRequestStats *PullRequestStats `json:"pull_request_stats"`
}

type IssueStats struct {
	TotalCount  int `json:"total_count"`
	OpenCount   int `json:"open_count"`
	ClosedCount int `json:"closed_count"`
}

type PullRequestStats struct {
	TotalCount  int `json:"total_count"`
	OpenCount   int `json:"open_count"`
	ClosedCount int `json:"closed_count"`
}
