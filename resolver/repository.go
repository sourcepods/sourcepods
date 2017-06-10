package resolver

import (
	"log"
	"time"

	"github.com/gitpods/gitpods/repository"
	graphql "github.com/neelance/graphql-go"
)

// RepositoryResolver communicates with the service to interact with repositories.
type RepositoryResolver struct {
	repositories repository.Service
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

// NewRepository returns a new RepositoryResolver.
func NewRepository(rs repository.Service) *RepositoryResolver {
	return &RepositoryResolver{repositories: rs}
}

type repositoryArgs struct {
	ID    *graphql.ID
	Owner *string
	Name  *string
}

// Repository returns a repositoryResolver based on an ID or Owner and Name.
func (r *RepositoryResolver) Repository(args repositoryArgs) *repositoryResolver {
	if args.ID != nil { // TODO
		return nil
	}
	if args.Owner != nil && args.Name != nil {
		repo, stats, _, err := r.repositories.Find(*args.Owner, *args.Name)
		if err != nil {
			log.Println(err)
			return nil
		}

		return &repositoryResolver{repository: &graphqlRepository{
			ID:            graphql.ID(repo.ID),
			Name:          repo.Name,
			Description:   repo.Description,
			Website:       repo.Website,
			DefaultBranch: repo.DefaultBranch,
			Private:       repo.Private,
			Bare:          repo.Bare,
			Created:       repo.Created,
			Updated:       repo.Updated,

			Stars: stats.Stars,
			Forks: stats.Forks,

			IssueStats: graphqlIssueStats{
				Total:  int32(stats.IssueTotalCount),
				Open:   int32(stats.IssueOpenCount),
				Closed: int32(stats.IssueClosedCount),
			},
			PullRequestStats: graphqlPullRequestStats{
				Total:  int32(stats.PullRequestTotalCount),
				Open:   int32(stats.PullRequestOpenCount),
				Closed: int32(stats.PullRequestClosedCount),
			},
		}}
	}
	return nil
}

// Repositories returns a slice of repositoryResolver based on their owner.
func (r *RepositoryResolver) Repositories(args struct{ Owner string }) []*repositoryResolver {
	repos, stats, _, err := r.repositories.ListByOwnerUsername(args.Owner)
	if err != nil {
		log.Println(err)
		return nil
	}

	var res []*repositoryResolver
	for i := range repos {
		res = append(res, &repositoryResolver{repository: &graphqlRepository{
			ID:            graphql.ID(repos[i].ID),
			Name:          repos[i].Name,
			Description:   repos[i].Description,
			Website:       repos[i].Website,
			DefaultBranch: repos[i].DefaultBranch,
			Private:       repos[i].Private,
			Bare:          repos[i].Bare,
			Created:       repos[i].Created,
			Updated:       repos[i].Updated,

			Stars: stats[i].Stars,
			Forks: stats[i].Forks,
		}})
	}
	return res
}

type repositoryResolver struct {
	repository *graphqlRepository
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
