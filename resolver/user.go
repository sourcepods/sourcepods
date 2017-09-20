package resolver

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gitpods/gitpods/repository"
	"github.com/gitpods/gitpods/session"
	"github.com/gitpods/gitpods/user"
	"github.com/neelance/graphql-go"
)

// UserResolver communicates with the service to interact with repositories.
type UserResolver struct {
	repositories repository.Service
	users        user.Service
}

type graphqlUser struct {
	ID       graphql.ID
	Email    string
	Username string
	Name     string
	Created  time.Time
	Updated  time.Time
}

func newGraphqlUser(u *user.User) *graphqlUser {
	return &graphqlUser{
		ID:       graphql.ID(u.ID),
		Email:    u.Email,
		Username: u.Username,
		Name:     u.Name,
		Created:  u.Created,
		Updated:  u.Updated,
	}
}

// NewUser returns a new UserResolver.
func NewUser(rs repository.Service, us user.Service) *UserResolver {
	return &UserResolver{repositories: rs, users: us}
}

// Me returns a userResolver based on the authenticated user which is retrieved from the context.
func (r *UserResolver) Me(ctx context.Context) *userResolver {
	sessionUser := session.GetSessionUser(ctx)

	u, err := r.users.FindByUsername(ctx, sessionUser.Username)
	if err != nil {
		log.Println(err)
		return nil
	}

	return &userResolver{rs: r.repositories, user: newGraphqlUser(u)}
}

type userArgs struct {
	ID       *graphql.ID
	Username *string
}

// User returns a userResolver based on an ID and Username.
func (r *UserResolver) User(args userArgs) *userResolver {
	//if args.ID != nil { TODO
	//	for _, user := range users {
	//		if user.ID == *args.ID {
	//			return &userResolver{rs: r.repositories, user}
	//		}
	//	}
	//}
	if args.Username != nil {
		u, err := r.users.FindByUsername(context.TODO(), *args.Username)
		if err != nil {
			log.Println(err)
			return nil
		}

		return &userResolver{rs: r.repositories, user: newGraphqlUser(u)}
	}
	return nil
}

// Users returns a slice of userResolver.
func (r *UserResolver) Users() []*userResolver {
	var uResolvers []*userResolver

	users, err := r.users.FindAll()
	if err != nil {
		log.Println(err)
		return nil
	}

	for _, u := range users {
		uResolvers = append(uResolvers, &userResolver{rs: r.repositories, user: newGraphqlUser(u)})
	}

	return uResolvers
}

type updateUserArgs struct {
	ID   graphql.ID
	User updatedUser
}

type updatedUser struct {
	Name string
}

func (r *UserResolver) UpdateUser(ctx context.Context, args updateUserArgs) (*userResolver, error) {
	sessUser := session.GetSessionUser(ctx)
	if sessUser.ID != string(args.ID) {
		return nil, fmt.Errorf("not allowed to update other users")
	}

	u, err := r.users.Find(string(args.ID))
	if err != nil {
		return nil, err
	}

	u.Name = strings.TrimSpace(args.User.Name)

	u, err = r.users.Update(u)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("updating user failed")
	}

	return &userResolver{rs: r.repositories, user: newGraphqlUser(u)}, nil
}

type userResolver struct {
	rs   repository.Service
	user *graphqlUser
}

func (r *userResolver) ID() graphql.ID {
	return r.user.ID
}

func (r *userResolver) Email() string {
	return r.user.Email
}

func (r *userResolver) Username() string {
	return r.user.Username
}

func (r *userResolver) Name() string {
	return r.user.Name
}

func (r *userResolver) CreatedAt() int32 {
	return int32(r.user.Created.Unix())
}

func (r *userResolver) UpdatedAt() int32 {
	return int32(r.user.Updated.Unix())
}

func (r *userResolver) Repositories() []*repositoryResolver {
	repos, stats, _, err := r.rs.List(&repository.Owner{Username: r.user.Username})
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
