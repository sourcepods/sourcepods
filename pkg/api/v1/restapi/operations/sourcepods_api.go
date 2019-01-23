// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"net/http"
	"strings"

	errors "github.com/go-openapi/errors"
	loads "github.com/go-openapi/loads"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	security "github.com/go-openapi/runtime/security"
	spec "github.com/go-openapi/spec"
	strfmt "github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/sourcepods/sourcepods/pkg/api/v1/restapi/operations/repositories"
	"github.com/sourcepods/sourcepods/pkg/api/v1/restapi/operations/users"
)

// NewSourcepodsAPI creates a new Sourcepods instance
func NewSourcepodsAPI(spec *loads.Document) *SourcepodsAPI {
	return &SourcepodsAPI{
		handlers:            make(map[string]map[string]http.Handler),
		formats:             strfmt.Default,
		defaultConsumes:     "application/json",
		defaultProduces:     "application/json",
		customConsumers:     make(map[string]runtime.Consumer),
		customProducers:     make(map[string]runtime.Producer),
		ServerShutdown:      func() {},
		spec:                spec,
		ServeError:          errors.ServeError,
		BasicAuthenticator:  security.BasicAuth,
		APIKeyAuthenticator: security.APIKeyAuth,
		BearerAuthenticator: security.BearerAuth,
		JSONConsumer:        runtime.JSONConsumer(),
		JSONProducer:        runtime.JSONProducer(),
		RepositoriesCreateRepositoryHandler: repositories.CreateRepositoryHandlerFunc(func(params repositories.CreateRepositoryParams) middleware.Responder {
			return middleware.NotImplemented("operation RepositoriesCreateRepository has not yet been implemented")
		}),
		RepositoriesGetOwnerRepositoriesHandler: repositories.GetOwnerRepositoriesHandlerFunc(func(params repositories.GetOwnerRepositoriesParams) middleware.Responder {
			return middleware.NotImplemented("operation RepositoriesGetOwnerRepositories has not yet been implemented")
		}),
		RepositoriesGetRepositoryHandler: repositories.GetRepositoryHandlerFunc(func(params repositories.GetRepositoryParams) middleware.Responder {
			return middleware.NotImplemented("operation RepositoriesGetRepository has not yet been implemented")
		}),
		RepositoriesGetRepositoryBranchesHandler: repositories.GetRepositoryBranchesHandlerFunc(func(params repositories.GetRepositoryBranchesParams) middleware.Responder {
			return middleware.NotImplemented("operation RepositoriesGetRepositoryBranches has not yet been implemented")
		}),
		UsersGetUserHandler: users.GetUserHandlerFunc(func(params users.GetUserParams) middleware.Responder {
			return middleware.NotImplemented("operation UsersGetUser has not yet been implemented")
		}),
		UsersGetUserMeHandler: users.GetUserMeHandlerFunc(func(params users.GetUserMeParams) middleware.Responder {
			return middleware.NotImplemented("operation UsersGetUserMe has not yet been implemented")
		}),
		UsersListUsersHandler: users.ListUsersHandlerFunc(func(params users.ListUsersParams) middleware.Responder {
			return middleware.NotImplemented("operation UsersListUsers has not yet been implemented")
		}),
		UsersUpdateUserHandler: users.UpdateUserHandlerFunc(func(params users.UpdateUserParams) middleware.Responder {
			return middleware.NotImplemented("operation UsersUpdateUser has not yet been implemented")
		}),
	}
}

/*SourcepodsAPI This is the API for SourcePods - git in the cloud. */
type SourcepodsAPI struct {
	spec            *loads.Document
	context         *middleware.Context
	handlers        map[string]map[string]http.Handler
	formats         strfmt.Registry
	customConsumers map[string]runtime.Consumer
	customProducers map[string]runtime.Producer
	defaultConsumes string
	defaultProduces string
	Middleware      func(middleware.Builder) http.Handler

	// BasicAuthenticator generates a runtime.Authenticator from the supplied basic auth function.
	// It has a default implemention in the security package, however you can replace it for your particular usage.
	BasicAuthenticator func(security.UserPassAuthentication) runtime.Authenticator
	// APIKeyAuthenticator generates a runtime.Authenticator from the supplied token auth function.
	// It has a default implemention in the security package, however you can replace it for your particular usage.
	APIKeyAuthenticator func(string, string, security.TokenAuthentication) runtime.Authenticator
	// BearerAuthenticator generates a runtime.Authenticator from the supplied bearer token auth function.
	// It has a default implemention in the security package, however you can replace it for your particular usage.
	BearerAuthenticator func(string, security.ScopedTokenAuthentication) runtime.Authenticator

	// JSONConsumer registers a consumer for a "application/json" mime type
	JSONConsumer runtime.Consumer

	// JSONProducer registers a producer for a "application/json" mime type
	JSONProducer runtime.Producer

	// RepositoriesCreateRepositoryHandler sets the operation handler for the create repository operation
	RepositoriesCreateRepositoryHandler repositories.CreateRepositoryHandler
	// RepositoriesGetOwnerRepositoriesHandler sets the operation handler for the get owner repositories operation
	RepositoriesGetOwnerRepositoriesHandler repositories.GetOwnerRepositoriesHandler
	// RepositoriesGetRepositoryHandler sets the operation handler for the get repository operation
	RepositoriesGetRepositoryHandler repositories.GetRepositoryHandler
	// RepositoriesGetRepositoryBranchesHandler sets the operation handler for the get repository branches operation
	RepositoriesGetRepositoryBranchesHandler repositories.GetRepositoryBranchesHandler
	// UsersGetUserHandler sets the operation handler for the get user operation
	UsersGetUserHandler users.GetUserHandler
	// UsersGetUserMeHandler sets the operation handler for the get user me operation
	UsersGetUserMeHandler users.GetUserMeHandler
	// UsersListUsersHandler sets the operation handler for the list users operation
	UsersListUsersHandler users.ListUsersHandler
	// UsersUpdateUserHandler sets the operation handler for the update user operation
	UsersUpdateUserHandler users.UpdateUserHandler

	// ServeError is called when an error is received, there is a default handler
	// but you can set your own with this
	ServeError func(http.ResponseWriter, *http.Request, error)

	// ServerShutdown is called when the HTTP(S) server is shut down and done
	// handling all active connections and does not accept connections any more
	ServerShutdown func()

	// Custom command line argument groups with their descriptions
	CommandLineOptionsGroups []swag.CommandLineOptionsGroup

	// User defined logger function.
	Logger func(string, ...interface{})
}

// SetDefaultProduces sets the default produces media type
func (o *SourcepodsAPI) SetDefaultProduces(mediaType string) {
	o.defaultProduces = mediaType
}

// SetDefaultConsumes returns the default consumes media type
func (o *SourcepodsAPI) SetDefaultConsumes(mediaType string) {
	o.defaultConsumes = mediaType
}

// SetSpec sets a spec that will be served for the clients.
func (o *SourcepodsAPI) SetSpec(spec *loads.Document) {
	o.spec = spec
}

// DefaultProduces returns the default produces media type
func (o *SourcepodsAPI) DefaultProduces() string {
	return o.defaultProduces
}

// DefaultConsumes returns the default consumes media type
func (o *SourcepodsAPI) DefaultConsumes() string {
	return o.defaultConsumes
}

// Formats returns the registered string formats
func (o *SourcepodsAPI) Formats() strfmt.Registry {
	return o.formats
}

// RegisterFormat registers a custom format validator
func (o *SourcepodsAPI) RegisterFormat(name string, format strfmt.Format, validator strfmt.Validator) {
	o.formats.Add(name, format, validator)
}

// Validate validates the registrations in the SourcepodsAPI
func (o *SourcepodsAPI) Validate() error {
	var unregistered []string

	if o.JSONConsumer == nil {
		unregistered = append(unregistered, "JSONConsumer")
	}

	if o.JSONProducer == nil {
		unregistered = append(unregistered, "JSONProducer")
	}

	if o.RepositoriesCreateRepositoryHandler == nil {
		unregistered = append(unregistered, "repositories.CreateRepositoryHandler")
	}

	if o.RepositoriesGetOwnerRepositoriesHandler == nil {
		unregistered = append(unregistered, "repositories.GetOwnerRepositoriesHandler")
	}

	if o.RepositoriesGetRepositoryHandler == nil {
		unregistered = append(unregistered, "repositories.GetRepositoryHandler")
	}

	if o.RepositoriesGetRepositoryBranchesHandler == nil {
		unregistered = append(unregistered, "repositories.GetRepositoryBranchesHandler")
	}

	if o.UsersGetUserHandler == nil {
		unregistered = append(unregistered, "users.GetUserHandler")
	}

	if o.UsersGetUserMeHandler == nil {
		unregistered = append(unregistered, "users.GetUserMeHandler")
	}

	if o.UsersListUsersHandler == nil {
		unregistered = append(unregistered, "users.ListUsersHandler")
	}

	if o.UsersUpdateUserHandler == nil {
		unregistered = append(unregistered, "users.UpdateUserHandler")
	}

	if len(unregistered) > 0 {
		return fmt.Errorf("missing registration: %s", strings.Join(unregistered, ", "))
	}

	return nil
}

// ServeErrorFor gets a error handler for a given operation id
func (o *SourcepodsAPI) ServeErrorFor(operationID string) func(http.ResponseWriter, *http.Request, error) {
	return o.ServeError
}

// AuthenticatorsFor gets the authenticators for the specified security schemes
func (o *SourcepodsAPI) AuthenticatorsFor(schemes map[string]spec.SecurityScheme) map[string]runtime.Authenticator {

	return nil

}

// Authorizer returns the registered authorizer
func (o *SourcepodsAPI) Authorizer() runtime.Authorizer {

	return nil

}

// ConsumersFor gets the consumers for the specified media types
func (o *SourcepodsAPI) ConsumersFor(mediaTypes []string) map[string]runtime.Consumer {

	result := make(map[string]runtime.Consumer)
	for _, mt := range mediaTypes {
		switch mt {

		case "application/json":
			result["application/json"] = o.JSONConsumer

		}

		if c, ok := o.customConsumers[mt]; ok {
			result[mt] = c
		}
	}
	return result

}

// ProducersFor gets the producers for the specified media types
func (o *SourcepodsAPI) ProducersFor(mediaTypes []string) map[string]runtime.Producer {

	result := make(map[string]runtime.Producer)
	for _, mt := range mediaTypes {
		switch mt {

		case "application/json":
			result["application/json"] = o.JSONProducer

		}

		if p, ok := o.customProducers[mt]; ok {
			result[mt] = p
		}
	}
	return result

}

// HandlerFor gets a http.Handler for the provided operation method and path
func (o *SourcepodsAPI) HandlerFor(method, path string) (http.Handler, bool) {
	if o.handlers == nil {
		return nil, false
	}
	um := strings.ToUpper(method)
	if _, ok := o.handlers[um]; !ok {
		return nil, false
	}
	if path == "/" {
		path = ""
	}
	h, ok := o.handlers[um][path]
	return h, ok
}

// Context returns the middleware context for the sourcepods API
func (o *SourcepodsAPI) Context() *middleware.Context {
	if o.context == nil {
		o.context = middleware.NewRoutableContext(o.spec, o, nil)
	}

	return o.context
}

func (o *SourcepodsAPI) initHandlerCache() {
	o.Context() // don't care about the result, just that the initialization happened

	if o.handlers == nil {
		o.handlers = make(map[string]map[string]http.Handler)
	}

	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/repositories"] = repositories.NewCreateRepository(o.context, o.RepositoriesCreateRepositoryHandler)

	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/repositories/{owner}"] = repositories.NewGetOwnerRepositories(o.context, o.RepositoriesGetOwnerRepositoriesHandler)

	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/repositories/{owner}/{name}"] = repositories.NewGetRepository(o.context, o.RepositoriesGetRepositoryHandler)

	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/repositories/{owner}/{name}/branches"] = repositories.NewGetRepositoryBranches(o.context, o.RepositoriesGetRepositoryBranchesHandler)

	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/users/{username}"] = users.NewGetUser(o.context, o.UsersGetUserHandler)

	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/users/me"] = users.NewGetUserMe(o.context, o.UsersGetUserMeHandler)

	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/users"] = users.NewListUsers(o.context, o.UsersListUsersHandler)

	if o.handlers["PATCH"] == nil {
		o.handlers["PATCH"] = make(map[string]http.Handler)
	}
	o.handlers["PATCH"]["/users/{username}"] = users.NewUpdateUser(o.context, o.UsersUpdateUserHandler)

}

// Serve creates a http handler to serve the API over HTTP
// can be used directly in http.ListenAndServe(":8000", api.Serve(nil))
func (o *SourcepodsAPI) Serve(builder middleware.Builder) http.Handler {
	o.Init()

	if o.Middleware != nil {
		return o.Middleware(builder)
	}
	return o.context.APIHandler(builder)
}

// Init allows you to just initialize the handler cache, you can then recompose the middleware as you see fit
func (o *SourcepodsAPI) Init() {
	if len(o.handlers) == 0 {
		o.initHandlerCache()
	}
}

// RegisterConsumer allows you to add (or override) a consumer for a media type.
func (o *SourcepodsAPI) RegisterConsumer(mediaType string, consumer runtime.Consumer) {
	o.customConsumers[mediaType] = consumer
}

// RegisterProducer allows you to add (or override) a producer for a media type.
func (o *SourcepodsAPI) RegisterProducer(mediaType string, producer runtime.Producer) {
	o.customProducers[mediaType] = producer
}
