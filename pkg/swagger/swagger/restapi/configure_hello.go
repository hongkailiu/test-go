// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	operations2 "github.com/hongkailiu/test-go/pkg/experimental/swagger/swagger/restapi/operations"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
)

//go:generate swagger generate server --target ../swagger/swagger --name hello --spec ../swagger/swagger/swagger.yml --exclude-main

func configureFlags(api *operations2.HelloAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations2.HelloAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.TxtProducer = runtime.TextProducer()

	api.GetUserUserIDHandler = operations2.GetUserUserIDHandlerFunc(func(params operations2.GetUserUserIDParams) middleware.Responder {
		return middleware.NotImplemented("operation .GetUserUserID has not yet been implemented")
	})
	api.GetUsersHandler = operations2.GetUsersHandlerFunc(func(params operations2.GetUsersParams) middleware.Responder {
		return middleware.NotImplemented("operation .GetUsers has not yet been implemented")
	})
	api.GetHostnameHandler = operations2.GetHostnameHandlerFunc(func(params operations2.GetHostnameParams) middleware.Responder {
		return middleware.NotImplemented("operation .GetHostname has not yet been implemented")
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
