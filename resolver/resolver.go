package resolver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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
			"createdAt": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.DateTime),
				Description: "The time the user was first created",
			},
			"updatedAt": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.DateTime),
				Description: "The time the user was updated last",
			},
		},
	})

	gCommit := graphql.NewObject(graphql.ObjectConfig{
		Name: "Commit",
		Fields: graphql.Fields{
			"sha1": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"parent": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"message": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"author": &graphql.Field{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "Author",
					Fields: graphql.Fields{
						"name": &graphql.Field{
							Type: graphql.NewNonNull(graphql.String),
						},
						"email": &graphql.Field{
							Type: graphql.NewNonNull(graphql.String),
						},
						"date": &graphql.Field{
							Type: graphql.NewNonNull(graphql.DateTime),
						},
					},
				}),
			},
			"committer": &graphql.Field{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "Committer",
					Fields: graphql.Fields{
						"name": &graphql.Field{
							Type: graphql.NewNonNull(graphql.String),
						},
						"email": &graphql.Field{
							Type: graphql.NewNonNull(graphql.String),
						},
						"date": &graphql.Field{
							Type: graphql.NewNonNull(graphql.DateTime),
						},
					},
				}),
			},
		},
	})

	gBranch := graphql.NewObject(graphql.ObjectConfig{
		Name: "Branch",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"sha1": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"type": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"protected": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Boolean),
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
			"defaultBranch": &graphql.Field{
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
			"createdAt": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.DateTime),
				Description: "The time the repository was first created",
			},
			"updatedAt": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.DateTime),
				Description: "The time the repository was updated last",
			},
			"owner": &graphql.Field{
				Type:        graphql.NewNonNull(gUser),
				Description: "The owner of this repository",
				Resolve:     h.ResolveRepositoryOwner(),
			},
			"branches": &graphql.Field{
				Type:        graphql.NewList(gBranch),
				Description: "", // TODO
				Resolve:     h.ResolveRepositoryBranches(),
			},
			"commit": &graphql.Field{
				Type:        graphql.NewNonNull(gCommit),
				Description: "", // TODO
				Resolve:     h.ResolveRepositoryCommit(),
				Args: graphql.FieldConfigArgument{
					"rev": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "", // TODO
					},
				},
			},
		},
	})

	query := graphql.NewObject(graphql.ObjectConfig{
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
				Type:    gRepository,
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
	})

	updatedUser := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "UpdatedUser",
		Fields: graphql.InputObjectConfigFieldMap{
			"name": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
	})

	createRepository := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "newRepository",
		Fields: graphql.InputObjectConfigFieldMap{
			"name": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"description": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"website": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
		},
	})

	mutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"updateUser": &graphql.Field{
				Type: gUser,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
					"user": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(updatedUser),
					},
				},
				Resolve: h.MutateUpdateUser(),
			},
			"createRepository": &graphql.Field{
				Type: gRepository,
				Args: graphql.FieldConfigArgument{
					"repository": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(createRepository),
					},
				},
				Resolve: h.MutateCreateRepository(),
			},
		},
	})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    query,
		Mutation: mutation,
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
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables"`
	}{}

	const megabyte = 1048576
	if err := json.NewDecoder(io.LimitReader(r.Body, megabyte)).Decode(&payload); err != nil {
		http.Error(w, "can't decode json body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	res := graphql.Do(graphql.Params{
		Schema:         h.schema,
		RequestString:  payload.Query,
		VariableValues: payload.Variables,
		Context:        ctx,
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
	Created  time.Time `json:"createdAt"`
	Updated  time.Time `json:"updatedAt"`
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

func (h *handler) MutateUpdateUser() graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		id, ok := p.Args["id"].(string)
		if !ok {
			return nil, fmt.Errorf("can't retreive id from arguments")
		}
		uu, ok := p.Args["user"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("can't retreive user from arguments")
		}
		uuName, ok := uu["name"].(string)
		if !ok {
			return nil, fmt.Errorf("can't retreive user's name from arguments")
		}

		sessUser := session.GetSessionUser(p.Context)
		if sessUser.ID != id {
			return nil, fmt.Errorf("not allowed to update other users")
		}

		u, err := h.users.Find(p.Context, id)
		if err != nil {
			return nil, err
		}

		u.Name = strings.TrimSpace(uuName)

		updated, err := h.users.Update(p.Context, u)
		if err != nil {
			return nil, err // TODO
		}

		return userResponse{
			ID:       updated.ID,
			Email:    updated.Email,
			Username: updated.Username,
			Name:     updated.Name,
			Created:  updated.Created,
			Updated:  updated.Updated,
		}, nil
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
	DefaultBranch string    `json:"defaultBranch"`
	Private       bool      `json:"private"`
	Bare          bool      `json:"bare"`
	Created       time.Time `json:"createdAt"`
	Updated       time.Time `json:"updatedAt"`

	Owner string
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

			Owner: owner,
		}, nil
	}
}

func (h *handler) MutateCreateRepository() graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		repoArgs := p.Args["repository"].(map[string]interface{})

		name, ok := repoArgs["name"].(string)
		if !ok {
			return nil, fmt.Errorf("can't retreive name from arguments")
		}

		description, _ := repoArgs["description"].(string)
		website, _ := repoArgs["website"].(string)

		u, err := h.users.Find(p.Context, session.GetSessionUser(p.Context).ID)
		if err != nil {
			return nil, err // TODO
		}

		r, err := h.repositories.Create(p.Context, u.Username, &repository.Repository{
			Name:        strings.TrimSpace(name),
			Description: strings.TrimSpace(description),
			Website:     strings.TrimSpace(website),
		})
		if err != nil {
			return nil, err
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

				Owner: owner,
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

type branchResponse struct {
	Name      string `json:"name"`
	Sha1      string `json:"sha1"`
	Type      string `json:"type"`
	Protected bool   `json:"protected"`
}

func (h *handler) ResolveRepositoryBranches() graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		r, ok := p.Source.(repositoryResponse)
		if !ok {
			return nil, fmt.Errorf("can't retreive repository from source") // TODO
		}

		bs, err := h.repositories.Branches(p.Context, r.Owner, r.Name)
		if err != nil {
			return nil, err
		}

		var res []branchResponse
		for _, b := range bs {
			res = append(res, branchResponse{
				Name:      b.Name,
				Sha1:      b.Sha1,
				Type:      b.Type,
				Protected: b.Protected,
			})
		}

		return res, nil
	}
}

type (
	commitResponse struct {
		Sha1      string                  `json:"sha1"`
		Message   string                  `json:"message"`
		Parent    string                  `json:"parent"`
		Author    commitAuthorResponse    `json:"author"`
		Committer commitCommitterResponse `json:"committer"`
	}

	commitAuthorResponse struct {
		Name  string    `json:"name"`
		Email string    `json:"email"`
		Date  time.Time `json:"date"`
	}

	commitCommitterResponse struct {
		Name  string    `json:"name"`
		Email string    `json:"email"`
		Date  time.Time `json:"date"`
	}
)

func (h *handler) ResolveRepositoryCommit() graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		r, ok := p.Source.(repositoryResponse)
		if !ok {
			return nil, fmt.Errorf("can't retreive repository from source") // TODO
		}

		var rev string
		if p.Args["rev"] != nil {
			rev, ok = p.Args["rev"].(string)
			if !ok {
				return nil, fmt.Errorf("can't retreive rev from arguments")
			}
		} else {
			// If no rev was given use the default branch
			rev = r.DefaultBranch
		}

		commit, err := h.repositories.Commit(p.Context, r.Owner, r.Name, rev)
		if err != nil {
			return nil, err // TODO
		}

		return commitResponse{
			Sha1:    commit.Hash,
			Parent:  commit.Parent,
			Message: commit.Message,
			Author: commitAuthorResponse{
				Name:  commit.Author,
				Email: commit.AuthorEmail,
				Date:  commit.AuthorDate,
			},
			Committer: commitCommitterResponse{
				Name:  commit.Committer,
				Email: commit.CommitterEmail,
				Date:  commit.CommitterDate,
			},
		}, nil
	}
}
