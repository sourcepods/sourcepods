package resolver

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/gitpods/gitpods/repository"
	"github.com/gitpods/gitpods/session"
	"github.com/gitpods/gitpods/user"
	graphql "github.com/neelance/graphql-go"
)

// RepositoryResolver communicates with the service to interact with repositories.
type RepositoryResolver struct {
	repositories repository.Service
	users        user.Service
}

type graphqlRepository struct {
	ID            graphql.ID
	Name          string
	Description   string
	Website       string
	DefaultBranch string
	Private       bool
	Bare          bool
	Created       time.Time
	Updated       time.Time

	Stars int
	Forks int

	IssueStats       graphqlIssueStats
	PullRequestStats graphqlPullRequestStats
}

type graphqlIssueStats struct {
	Total  int32
	Open   int32
	Closed int32
}

type graphqlPullRequestStats struct {
	Total  int32
	Open   int32
	Closed int32
}

func newGraphqlRepository(repo *repository.Repository, stats *repository.Stats) *graphqlRepository {
	r := &graphqlRepository{
		ID:            graphql.ID(repo.ID),
		Name:          repo.Name,
		Description:   repo.Description,
		Website:       repo.Website,
		DefaultBranch: repo.DefaultBranch,
		Private:       repo.Private,
		Bare:          repo.Bare,
		Created:       repo.Created,
		Updated:       repo.Updated,
	}

	if stats != nil {
		r.Stars = stats.Stars
		r.Forks = stats.Forks
		r.IssueStats = graphqlIssueStats{
			Total:  int32(stats.IssueTotalCount),
			Open:   int32(stats.IssueOpenCount),
			Closed: int32(stats.IssueClosedCount),
		}
		r.PullRequestStats = graphqlPullRequestStats{
			Total:  int32(stats.PullRequestTotalCount),
			Open:   int32(stats.PullRequestOpenCount),
			Closed: int32(stats.PullRequestClosedCount),
		}
	}

	return r
}

// NewRepository returns a new RepositoryResolver.
func NewRepository(rs repository.Service, us user.Service) *RepositoryResolver {
	return &RepositoryResolver{
		repositories: rs,
		users:        us,
	}
}

type repositoryArgs struct {
	Owner string
	Name  string
}

// Repository returns a repositoryResolver based on an ID or Owner and Name.
func (r *RepositoryResolver) Repository(ctx context.Context, args repositoryArgs) (*repositoryResolver, error) {
	repo, stats, _, err := r.repositories.Find(ctx, args.Owner, args.Name)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &repositoryResolver{repository: newGraphqlRepository(repo, stats)}, nil
}

// Repositories returns a slice of repositoryResolver based on their owner.
func (r *RepositoryResolver) Repositories(ctx context.Context, args struct{ Owner string }) []*repositoryResolver {
	repos, stats, _, err := r.repositories.List(ctx, args.Owner)
	if err != nil {
		log.Println(err)
		return nil
	}

	var res []*repositoryResolver
	for i := range repos {
		res = append(res, &repositoryResolver{repository: newGraphqlRepository(repos[i], stats[i])})
	}

	return res
}

type newRepository struct {
	Name        string
	Description *string
	Website     *string
	Private     bool
}

func (r *RepositoryResolver) CreateRepository(ctx context.Context, args struct{ Repository newRepository }) (*repositoryResolver, error) {
	sessUser := session.GetSessionUser(ctx)
	sessOwner := sessUser

	description := ""
	if args.Repository.Description != nil {
		description = strings.TrimSpace(*args.Repository.Description)
	}

	website := ""
	if args.Repository.Website != nil {
		website = strings.TrimSpace(*args.Repository.Website)
	}

	repo := &repository.Repository{
		Name:          strings.TrimSpace(args.Repository.Name),
		Description:   description,
		Website:       website,
		DefaultBranch: "master",
		Private:       args.Repository.Private,
		Bare:          true,
	}

	owner, err := r.users.Find(ctx, sessOwner.ID)
	if err != nil {
		return nil, err
	}

	repo, err = r.repositories.Create(ctx, owner.Username, repo)
	if err != nil {
		return nil, err
	}

	return &repositoryResolver{
		repository: newGraphqlRepository(repo, nil),
		owner:      newGraphqlUser(owner),
	}, nil
}

type repositoryResolver struct {
	repository *graphqlRepository
	owner      *graphqlUser
}

func (r *repositoryResolver) ID() graphql.ID {
	return r.repository.ID
}

func (r *repositoryResolver) Name() string {
	return r.repository.Name
}

func (r *repositoryResolver) Description() string {
	return r.repository.Description
}

func (r *repositoryResolver) Website() string {
	return r.repository.Website
}

func (r *repositoryResolver) DefaultBranch() string {
	return r.repository.DefaultBranch
}

func (r *repositoryResolver) Private() bool {
	return r.repository.Private
}

func (r *repositoryResolver) Bare() bool {
	return r.repository.Bare
}

func (r *repositoryResolver) CreatedAt() int32 {
	return int32(r.repository.Created.Unix())
}

func (r *repositoryResolver) UpdatedAt() int32 {
	return int32(r.repository.Updated.Unix())
}

func (r *repositoryResolver) Stars() int32 {
	return int32(r.repository.Stars)
}

func (r *repositoryResolver) Forks() int32 {
	return int32(r.repository.Forks)
}

func (r *repositoryResolver) IssueStats() *issueStatsResolver {
	return &issueStatsResolver{
		total:  r.repository.IssueStats.Total,
		open:   r.repository.IssueStats.Open,
		closed: r.repository.IssueStats.Closed,
	}
}

func (r *repositoryResolver) PullRequestStats() *pullRequestStatsResolver {
	return &pullRequestStatsResolver{
		total:  r.repository.PullRequestStats.Total,
		open:   r.repository.PullRequestStats.Open,
		closed: r.repository.PullRequestStats.Closed,
	}
}

type issueStatsResolver struct {
	total  int32
	open   int32
	closed int32
}

func (r *issueStatsResolver) Total() int32 {
	return r.total
}

func (r *issueStatsResolver) Open() int32 {
	return r.open
}

func (r *issueStatsResolver) Closed() int32 {
	return r.closed
}

type pullRequestStatsResolver struct {
	total  int32
	open   int32
	closed int32
}

func (r *pullRequestStatsResolver) Total() int32 {
	return r.total
}

func (r *pullRequestStatsResolver) Open() int32 {
	return r.open
}

func (r *pullRequestStatsResolver) Closed() int32 {
	return r.closed
}

func (r *repositoryResolver) Owner() *userResolver {
	return &userResolver{user: r.owner}
}
