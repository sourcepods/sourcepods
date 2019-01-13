// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/gitpods/gitpods/internal/api/v1/restapi/operations"
	"github.com/gitpods/gitpods/internal/api/v1/restapi/operations/repositories"
	"github.com/gitpods/gitpods/internal/api/v1/restapi/operations/users"
)

//go:generate swagger generate server --target ../../v1 --name Gitpods --spec ../../../../swagger.yaml --exclude-main

func configureFlags(api *operations.GitpodsAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.GitpodsAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.RepositoriesCreateRepositoryHandler = repositories.CreateRepositoryHandlerFunc(func(params repositories.CreateRepositoryParams) middleware.Responder {
		return middleware.NotImplemented("operation repositories.CreateRepository has not yet been implemented")
	})
	api.RepositoriesGetOwnerRepositoriesHandler = repositories.GetOwnerRepositoriesHandlerFunc(func(params repositories.GetOwnerRepositoriesParams) middleware.Responder {
		return middleware.NotImplemented("operation repositories.GetOwnerRepositories has not yet been implemented")
	})
	api.RepositoriesGetRepositoryHandler = repositories.GetRepositoryHandlerFunc(func(params repositories.GetRepositoryParams) middleware.Responder {
		return middleware.NotImplemented("operation repositories.GetRepository has not yet been implemented")
	})
	api.UsersGetUserHandler = users.GetUserHandlerFunc(func(params users.GetUserParams) middleware.Responder {
		return middleware.NotImplemented("operation users.GetUser has not yet been implemented")
	})
	api.UsersGetUserMeHandler = users.GetUserMeHandlerFunc(func(params users.GetUserMeParams) middleware.Responder {
		return middleware.NotImplemented("operation users.GetUserMe has not yet been implemented")
	})
	api.UsersListUsersHandler = users.ListUsersHandlerFunc(func(params users.ListUsersParams) middleware.Responder {
		return middleware.NotImplemented("operation users.ListUsers has not yet been implemented")
	})
	api.UsersUpdateUserHandler = users.UpdateUserHandlerFunc(func(params users.UpdateUserParams) middleware.Responder {
		return middleware.NotImplemented("operation users.UpdateUser has not yet been implemented")
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
