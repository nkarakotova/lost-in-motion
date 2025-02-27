// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"
	"fmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/charmbracelet/log"
	"github.com/rs/cors"

	"lim/internal/lim-handler/restapi/operations"
	"lim/cmd/registry"
	"lim/internal/lim-handler/restapi/handlers"
	"lim/internal/lim-handler/restapi/jwt"
)

//go:generate swagger generate server --target ../../lim-handler --name SwaggerLIM --spec ../../../docs/swagger.yaml --principal interface{}

func configureFlags(api *operations.SwaggerLIMAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.SwaggerLIMAPI) http.Handler {
	api.ServeError = errors.ServeError

	app := registry.App{}
	errConfig := app.Config.ParseConfig("config.json", "./cmd/app")
	if errConfig != nil {
		log.Fatal("error opening config file")
	}
	err := app.Run()
	if err != nil {
		log.Fatal(fmt.Errorf("%s", err.Error()))
	}

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()

	api.JWTAuthAuth = func(token string) (interface{}, error) {
		_, role, err := jwt.VerifyToken(token)
		if err != nil {
			return nil, errors.Unauthenticated("invalid token")
		}
		return role, nil
	}

	handlers.ConfigureAuthHandlers(app, api)
	handlers.ConfigureAssignmentsHandlers(app, api)
	handlers.ConfigureClientsHandlers(app, api)
	handlers.ConfigureCoachesHandlers(app, api)
	handlers.ConfigureHallsHandlers(app, api)
	handlers.ConfigureScheduleHandlers(app, api)
	handlers.ConfigureTrainingsHandlers(app, api)

	api.PreServerShutdown = func() {}
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
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	corsHandler := cors.New(cors.Options{
		AllowedMethods: []string{"HEAD", "PATCH", "PUT", "GET", "POST", "DELETE", "OPTIONS"},
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"Content-Type", "X-Requested-With", "token"},
	})

	return corsHandler.Handler(handler)
}
