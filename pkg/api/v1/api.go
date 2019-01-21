package v1

import (
	"net/http"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/sourcepods/sourcepods/pkg/api/v1/models"
	"github.com/sourcepods/sourcepods/pkg/api/v1/restapi"
	"github.com/sourcepods/sourcepods/pkg/api/v1/restapi/operations"
	"github.com/sourcepods/sourcepods/pkg/api/v1/restapi/operations/repositories"
	"github.com/sourcepods/sourcepods/pkg/api/v1/restapi/operations/users"
	"github.com/sourcepods/sourcepods/pkg/session"
	"github.com/sourcepods/sourcepods/pkg/sourcepods/repository"
	"github.com/sourcepods/sourcepods/pkg/sourcepods/user"
)

// API has the http.Handler for the OpenAPI implementation
type API struct {
	Handler http.Handler
}

// New creates a new API that adds our own Handler implementations
func New(rs repository.Service, us user.Service) (*API, error) {
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		return nil, err
	}

	sourcepodsAPI := operations.NewSourcepodsAPI(swaggerSpec)

	sourcepodsAPI.Middleware = func(b middleware.Builder) http.Handler {
		return middleware.Spec("", nil, sourcepodsAPI.Context().RoutesHandler(b))
	}

	sourcepodsAPI.RepositoriesCreateRepositoryHandler = CreateRepositoryHandler(rs)
	sourcepodsAPI.RepositoriesGetOwnerRepositoriesHandler = GetOwnerRepositoriesHandler(rs)
	sourcepodsAPI.RepositoriesGetRepositoryBranchesHandler = GetRepositoryBranchesHandler(rs)
	sourcepodsAPI.RepositoriesGetRepositoryHandler = GetRepositoryHandler(rs)
	sourcepodsAPI.UsersGetUserHandler = GetUserHandler(us)
	sourcepodsAPI.UsersGetUserMeHandler = GetUserMeHandler(us)
	sourcepodsAPI.UsersListUsersHandler = ListUsersHandler(us)
	sourcepodsAPI.UsersUpdateUserHandler = UpdateUserHandler(us)

	return &API{
		Handler: sourcepodsAPI.Serve(nil),
	}, nil
}

func convertRepository(r *repository.Repository) *models.Repository {
	return &models.Repository{
		ID:            strfmt.UUID(r.ID),
		Name:          &r.Name,
		Description:   r.Description,
		DefaultBranch: r.DefaultBranch,
		Website:       r.Website,
		CreatedAt:     strfmt.DateTime(r.Created),
		UpdatedAt:     strfmt.DateTime(r.Updated),

		//Owner: nil, // TODO: Include via query parameter if wanted
	}
}

//CreateRepositoryHandler creates a new repository from given input
func CreateRepositoryHandler(rs repository.Service) repositories.CreateRepositoryHandlerFunc {
	return func(params repositories.CreateRepositoryParams) middleware.Responder {
		ctx := params.HTTPRequest.Context()
		owner := session.GetSessionUser(ctx)

		r, err := rs.Create(ctx, owner.Username, &repository.Repository{
			Name:        *params.NewRepository.Name,
			Description: params.NewRepository.Description,
			Website:     params.NewRepository.Website,
		})
		if err != nil {
			if v, ok := err.(repository.ValidationErrors); ok {
				message := "The given repository input is invalid"
				payload := &models.ValidationError{
					Message: &message,
				}
				for _, verr := range v.Errors {
					payload.Errors = append(payload.Errors, &models.ValidationErrorErrorsItems0{
						Field:   verr.Field,
						Message: verr.Error.Error(),
					})
				}
				return repositories.NewCreateRepositoryUnprocessableEntity().WithPayload(payload)
			}

			return repositories.NewCreateRepositoryDefault(http.StatusInternalServerError)
		}

		return repositories.NewCreateRepositoryOK().WithPayload(convertRepository(r))
	}
}

//GetOwnerRepositoriesHandler gets a repository by the owner's username
func GetOwnerRepositoriesHandler(rs repository.Service) repositories.GetOwnerRepositoriesHandlerFunc {
	return func(params repositories.GetOwnerRepositoriesParams) middleware.Responder {
		list, _, err := rs.List(params.HTTPRequest.Context(), params.Owner)
		if err != nil {
			if err == repository.ErrOwnerNotFound {
				message := "owner not found"
				return repositories.NewGetOwnerRepositoriesNotFound().WithPayload(&models.Error{
					Message: &message,
				})
			}
			return repositories.NewGetOwnerRepositoriesDefault(http.StatusInternalServerError)
		}

		var payload []*models.Repository
		for _, r := range list {
			payload = append(payload, convertRepository(r))
		}

		return repositories.NewGetOwnerRepositoriesOK().WithPayload(payload)
	}
}

//GetRepositoryBranchesHandler gets all branches of a repository
func GetRepositoryBranchesHandler(rs repository.Service) repositories.GetRepositoryBranchesHandlerFunc {
	return func(params repositories.GetRepositoryBranchesParams) middleware.Responder {
		branches, err := rs.Branches(params.HTTPRequest.Context(), params.Owner, params.Name)
		if err != nil {
			if err == repository.ErrRepositoryNotFound {
				message := "repository not found"
				return repositories.NewGetRepositoryBranchesNotFound().WithPayload(&models.Error{
					Message: &message,
				})
			}

			return repositories.NewGetRepositoryBranchesDefault(http.StatusInternalServerError)
		}

		var payload []*models.Branch

		for _, b := range branches {
			payload = append(payload, &models.Branch{
				Name: b.Name,
				Sha1: b.Sha1,
				Type: b.Type,
			})
		}

		return repositories.NewGetRepositoryBranchesOK().WithPayload(payload)
	}
}

//GetRepositoryHandler gets a repository by name and the owner's username
func GetRepositoryHandler(rs repository.Service) repositories.GetRepositoryHandlerFunc {
	return func(params repositories.GetRepositoryParams) middleware.Responder {
		r, _, err := rs.Find(params.HTTPRequest.Context(), params.Owner, params.Name)
		if err != nil {
			if err == repository.ErrRepositoryNotFound {
				message := "repository not found"
				return repositories.NewGetRepositoryNotFound().WithPayload(&models.Error{
					Message: &message,
				})
			}
			return repositories.NewGetRepositoryDefault(http.StatusInternalServerError)
		}

		return repositories.NewGetRepositoryOK().WithPayload(convertRepository(r))
	}
}

func convertUser(u *user.User) *models.User {
	return &models.User{
		ID:        strfmt.UUID(u.ID),
		Email:     strfmt.Email(u.Email),
		Username:  &u.Username,
		Name:      u.Name,
		CreatedAt: strfmt.DateTime(u.Created),
		UpdatedAt: strfmt.DateTime(u.Updated),
	}
}

// ListUsersHandler gets a list of users from the user.Service and returns a API response
func ListUsersHandler(us user.Service) users.ListUsersHandlerFunc {
	return func(params users.ListUsersParams) middleware.Responder {
		list, err := us.FindAll(params.HTTPRequest.Context())
		if err != nil {
			return users.NewListUsersDefault(http.StatusInternalServerError)
		}

		var payload []*models.User

		for _, u := range list {
			payload = append(payload, convertUser(u))
		}

		return users.NewListUsersOK().WithPayload(payload)
	}
}

// GetUserMeHandler gets the currently authenticated user
func GetUserMeHandler(us user.Service) users.GetUserMeHandlerFunc {
	return func(params users.GetUserMeParams) middleware.Responder {
		sessUser := session.GetSessionUser(params.HTTPRequest.Context())

		u, err := us.FindByUsername(params.HTTPRequest.Context(), sessUser.Username)
		if err != nil {
			return users.NewGetUserMeDefault(http.StatusInternalServerError)
		}

		return users.NewGetUserMeOK().WithPayload(convertUser(u))
	}
}

// GetUserHandler gets a user from the user.Service and returns a API response
func GetUserHandler(us user.Service) users.GetUserHandlerFunc {
	return func(params users.GetUserParams) middleware.Responder {
		u, err := us.FindByUsername(params.HTTPRequest.Context(), params.Username)
		if err != nil {
			if err == user.ErrNotFound {
				message := "user not found"
				return users.NewGetUserNotFound().WithPayload(&models.Error{
					Message: &message,
				})
			}
			return users.NewGetUserDefault(http.StatusInternalServerError)
		}

		return users.NewGetUserOK().WithPayload(convertUser(u))
	}
}

// UpdateUserHandler receives a updated user and returns a API response after updating
func UpdateUserHandler(us user.Service) users.UpdateUserHandlerFunc {
	return func(params users.UpdateUserParams) middleware.Responder {
		old, err := us.FindByUsername(params.HTTPRequest.Context(), params.Username)
		if err != nil {
			if err == user.ErrNotFound {
				message := "user not found"
				return users.NewUpdateUserNotFound().WithPayload(&models.Error{
					Message: &message,
				})
			}
			return users.NewUpdateUserDefault(http.StatusInternalServerError)
		}

		old.Name = *params.UpdatedUser.Name

		updated, err := us.Update(params.HTTPRequest.Context(), old)
		if err != nil {
			message := "updated user is invalid"
			return users.NewUpdateUserUnprocessableEntity().WithPayload(&models.ValidationError{
				Message: &message,
			})
		}

		return users.NewUpdateUserOK().WithPayload(convertUser(updated))
	}
}
