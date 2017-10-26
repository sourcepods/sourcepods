package resolver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gitpods/gitpods/repository"
	"github.com/gitpods/gitpods/session"
	"github.com/gitpods/gitpods/user"
	"github.com/graphql-go/graphql"
	"github.com/opentracing/opentracing-go"
)

type handler struct {
	schema       graphql.Schema
	users        user.Service
	repositories repository.Service
}

func Handler(repositories repository.Service, users user.Service) http.Handler {
	h := &handler{
		repositories: repositories,
		users:        users,
	}

	gUser := graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.ID),
				Description: "The user's ID",
			},
			"email": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The user's email address",
			},
			"username": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The user's username",
			},
			"name": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The user's name",
			},
			"created_at": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.DateTime),
				Description: "The time the user was first created",
			},
			"updated_at": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.DateTime),
				Description: "The time the user was updated last",
			},
		},
	})

	gRepository := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Repository",
		Description: "Repository of code and more information, think of it like a software project",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.ID),
				Description: "The repository's ID",
			},
			"name": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The repository's name",
			},
			"description": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The repository's description",
			},
			"website": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The repository's website",
			},
			"default_branch": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The user's name",
			},
			"private": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Boolean),
				Description: "True when the repository is private and not public",
			},
			"bare": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Boolean),
				Description: "True when the repository is bare",
			},
			"created_at": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.DateTime),
				Description: "The time the repository was first created",
			},
			"updated_at": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.DateTime),
				Description: "The time the repository was updated last",
			},
			"owner": &graphql.Field{
				Type:        graphql.NewNonNull(gUser),
				Description: "The owner of this repository",
				Resolve:     h.ResolveRepositoryOwner(),
			},
		},
	})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"me": &graphql.Field{
					Type:    gUser,
					Resolve: h.ResolveMe(),
				},
				"user": &graphql.Field{
					Type:    gUser,
					Resolve: h.ResolveUser(),
					Args: graphql.FieldConfigArgument{
						"username": &graphql.ArgumentConfig{
							Type:        graphql.NewNonNull(graphql.String),
							Description: "The username used to search for the user",
						},
					},
				},
				"users": &graphql.Field{
					Type:    graphql.NewNonNull(graphql.NewList(gUser)),
					Resolve: h.ResolveUsers(),
				},
				"repository": &graphql.Field{
					Type:    graphql.NewNonNull(gRepository),
					Resolve: h.ResolveRepository(),
					Args: graphql.FieldConfigArgument{
						"owner": &graphql.ArgumentConfig{
							Type:        graphql.NewNonNull(graphql.String),
							Description: "The username of the repository's owner",
						},
						"name": &graphql.ArgumentConfig{
							Type:        graphql.NewNonNull(graphql.String),
							Description: "The name of the repository",
						},
					},
				},
				"repositories": &graphql.Field{
					Type:    graphql.NewNonNull(graphql.NewList(gRepository)),
					Resolve: h.ResolveRepositories(),
					Args: graphql.FieldConfigArgument{
						"owner": &graphql.ArgumentConfig{
							Type:        graphql.NewNonNull(graphql.String),
							Description: "The username of the repository's owner",
						},
					},
				},
			},
		}),
	})
	if err != nil {
		panic(err) // TODO
	}
	h.schema = schema

	return h
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "graphql.ServeHTTP")
	defer span.Finish()

	var payload = struct {
		Query string `json:"query"`
	}{}

	const megabyte = 1048576
	if err := json.NewDecoder(io.LimitReader(r.Body, megabyte)).Decode(&payload); err != nil {
		http.Error(w, "can't decode json body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	res := graphql.Do(graphql.Params{
		Schema:        h.schema,
		RequestString: payload.Query,
		Context:       ctx,
	})

	if err := json.NewEncoder(w).Encode(res); err != nil {
		panic(err) // TODO
	}
}

type userResponse struct {
	ID       string    `json:"id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	Name     string    `json:"name"`
	Created  time.Time `json:"created_at"`
	Updated  time.Time `json:"updated_at"`
}

func (h *handler) ResolveMe() graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		sessUser := session.GetSessionUser(p.Context)

		u, err := h.users.FindByUsername(p.Context, sessUser.Username)
		if err != nil {
			return nil, err // TODO
		}

		return userResponse{
			ID:       u.ID,
			Email:    u.Email,
			Username: u.Username,
			Name:     u.Name,
			Created:  u.Created,
			Updated:  u.Updated,
		}, err
	}
}

func (h *handler) ResolveUser() graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		username, ok := p.Args["username"].(string)
		if !ok {
			return nil, fmt.Errorf("can't retreive username from arguments")
		}

		u, err := h.users.FindByUsername(p.Context, username)
		if err != nil {
			return nil, err // TODO
		}

		return userResponse{
			ID:       u.ID,
			Email:    u.Email,
			Username: u.Username,
			Name:     u.Name,
			Created:  u.Created,
			Updated:  u.Updated,
		}, err
	}
}
func (h *handler) ResolveUsers() graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		us, err := h.users.FindAll(p.Context)
		if err != nil {
			return nil, err // TODO
		}

		var res []userResponse
		for _, u := range us {
			res = append(res, userResponse{
				ID:       u.ID,
				Email:    u.Email,
				Username: u.Username,
				Name:     u.Name,
				Created:  u.Created,
				Updated:  u.Updated,
			})
		}

		return res, nil
	}
}

type repositoryResponse struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Website       string    `json:"website"`
	DefaultBranch string    `json:"default_branch"`
	Private       bool      `json:"private"`
	Bare          bool      `json:"bare"`
	Created       time.Time `json:"created_at"`
	Updated       time.Time `json:"updated_at"`
}

func (h *handler) ResolveRepository() graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		name, ok := p.Args["name"].(string)
		if !ok {
			return nil, fmt.Errorf("can't retreive name from arguments")
		}
		owner, ok := p.Args["owner"].(string)
		if !ok {
			return nil, fmt.Errorf("can't retreive owner's username from arguments")
		}

		r, _, _, err := h.repositories.Find(p.Context, owner, name)
		if err != nil {
			return nil, err // TODO
		}

		return repositoryResponse{
			ID:            r.ID,
			Name:          r.Name,
			Description:   r.Description,
			Website:       r.Website,
			DefaultBranch: r.DefaultBranch,
			Private:       r.Private,
			Bare:          r.Bare,
			Created:       r.Created,
			Updated:       r.Updated,
		}, nil
	}
}

func (h *handler) ResolveRepositories() graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		owner, ok := p.Args["owner"].(string)
		if !ok {
			return nil, fmt.Errorf("can't retreive owner's username from arguments")
		}

		rs, _, _, err := h.repositories.List(p.Context, owner)
		if err != nil {
			return nil, err // TODO
		}

		var res []repositoryResponse
		for _, r := range rs {
			res = append(res, repositoryResponse{
				ID:            r.ID,
				Name:          r.Name,
				Description:   r.Description,
				Website:       r.Website,
				DefaultBranch: r.DefaultBranch,
				Private:       r.Private,
				Bare:          r.Bare,
				Created:       r.Created,
				Updated:       r.Updated,
			})
		}
		return res, nil
	}
}

func (h *handler) ResolveRepositoryOwner() graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		r, ok := p.Source.(repositoryResponse)
		if !ok {
			return nil, fmt.Errorf("can't retreive repository from source")
		}

		owner, err := h.users.FindRepositoryOwner(p.Context, r.ID)
		if err != nil {
			return nil, err // TODO
		}

		return userResponse{
			ID:       owner.ID,
			Email:    owner.Email,
			Username: owner.Username,
			Name:     owner.Name,
			Created:  owner.Created,
			Updated:  owner.Updated,
		}, err
	}
}
