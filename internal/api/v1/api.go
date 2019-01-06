package v1

import (
	"net/http"

	"github.com/gitpods/gitpods/internal/api/v1/models"
	"github.com/gitpods/gitpods/internal/api/v1/restapi"
	"github.com/gitpods/gitpods/internal/api/v1/restapi/operations"
	"github.com/gitpods/gitpods/internal/api/v1/restapi/operations/users"
	"github.com/gitpods/gitpods/user"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

type API struct {
	Handler http.Handler
}

func New(us user.Service) (*API, error) {
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		return nil, err
	}

	gitpodsAPI := operations.NewGitpodsAPI(swaggerSpec)

	gitpodsAPI.UsersListUsersHandler = ListUsersHandler(us)
	gitpodsAPI.UsersGetUserHandler = GetUserHandler(us)
	gitpodsAPI.UsersUpdateUserHandler = UpdateUserHandlerFunc(us)

	return &API{
		Handler: gitpodsAPI.Serve(nil),
	}, nil
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

func GetUserHandler(us user.Service) users.GetUserHandlerFunc {
	return func(params users.GetUserParams) middleware.Responder {
		u, err := us.FindByUsername(params.HTTPRequest.Context(), params.Username)
		if err != nil {
			if err == user.NotFoundError {
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

func UpdateUserHandlerFunc(us user.Service) users.UpdateUserHandlerFunc {
	return func(params users.UpdateUserParams) middleware.Responder {
		old, err := us.FindByUsername(params.HTTPRequest.Context(), params.Username)
		if err != nil {
			if err == user.NotFoundError {
				message := "user not found"
				return users.NewUpdateUserNotFound().WithPayload(&models.Error{
					Message: &message,
				})
			}
			return users.NewUpdateUserDefault(http.StatusInternalServerError)
		}

		old.Name = *params.Body.Name

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
