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
		repo, _, _, err := r.repositories.Find(*args.Owner, *args.Name)
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
		}}
	}
	return nil
}

// Repositories returns a slice of repositoryResolver based on their owner.
func (r *RepositoryResolver) Repositories(args struct{ Owner string }) []*repositoryResolver {
	repos, _, _, err := r.repositories.ListByOwnerUsername(args.Owner)
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
