package handlers

import (
	"github.com/go-openapi/runtime/middleware"

	"lim/cmd/registry"
	"lim/internal/lim-handler/restapi/operations"
	"lim/internal/lim-handler/restapi/operations/auth"
	"lim/internal/lim-handler/restapi/jwt"
)

var AuthHandlersApp registry.App

func ConfigureAuthHandlers(app registry.App, api *operations.SwaggerLIMAPI) {
	AuthHandlersApp = app
	api.AuthLoginHandler = auth.LoginHandlerFunc(loginHandlerFunc)
	api.AuthSignupHandler = auth.SignupHandlerFunc(signupHandlerFunc)
}

func loginHandlerFunc(params auth.LoginParams) middleware.Responder {
	return jwt.LoginHandlerFunc(params, AuthHandlersApp)
}

func signupHandlerFunc(params auth.SignupParams) middleware.Responder {
	return jwt.SignupHandlerFunc(params, AuthHandlersApp)
}