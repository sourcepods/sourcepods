package repository

import "time"

// Repository is a git repository with some meta information for gitpods.
type Repository struct {
	ID            string
	Name          string
	Description   string
	Website       string
	DefaultBranch string
	Private       bool
	Bare          bool
	Created       time.Time
	Updated       time.Time
}

type Stats struct {
	Stars int
	Forks int

	IssueTotalCount        int
	IssueOpenCount         int
	IssueClosedCount       int
	PullRequestTotalCount  int
	PullRequestOpenCount   int
	PullRequestClosedCount int
}

type Owner struct {
	ID string
}
