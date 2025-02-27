package handlers

import (
	"net/http"
	"context"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/runtime"

	"lim/cmd/registry"
	"lim/internal/lim-handler/restapi/operations"
	"lim/internal/lim-handler/restapi/operations/coaches"
	"lim/internal/lim-handler/restapi/jwt"
	"lim/internal/lim-core/models"
	"lim/internal/lim-core/errors/servicesErrors"

	dto_models "lim/internal/lim-handler/models"
)

var CoachesHandlersApp registry.App

func ConfigureCoachesHandlers(app registry.App, api *operations.SwaggerLIMAPI) {
	CoachesHandlersApp = app
	api.CoachesCreateCoachHandler = coaches.CreateCoachHandlerFunc(createCoachHandlerFunc)
	api.CoachesGetCoachesHandler = coaches.GetCoachesHandlerFunc(getCoachesHandlerFunc)
}

func createCoachHandlerFunc(params coaches.CreateCoachParams, principal interface{}) middleware.Responder {

	_, role, err := jwt.GetAuthenticatedUser(params.HTTPRequest)
	if err != nil {
		return middleware.Error(http.StatusUnauthorized, "Authorization error")
	}
	if role != "admin" {
		return middleware.Error(http.StatusForbidden, "No rights")
	}

	ctx := context.Background()

	name := params.CoachInfo.Name

	coach := models.Coach{
		Name: *name,
	}

	err = CoachesHandlersApp.Services.CoachService.Create(ctx, &coach)
	if err != nil && err == servicesErrors.CoachAlreadyExists {
		return middleware.Error(http.StatusConflict, "Coach alredy exists")
	} else if err != nil {
		return middleware.Error(http.StatusInternalServerError, "Can't create coach")
	}

	return middleware.ResponderFunc(func(rw http.ResponseWriter, p runtime.Producer) {
		rw.WriteHeader(http.StatusOK)
	})
}

func getCoachesHandlerFunc(params coaches.GetCoachesParams, principal interface{}) middleware.Responder {

	CoachesHandlersApp.Logger.Info("GET", "getCoaches")

	_, role, err := jwt.GetAuthenticatedUser(params.HTTPRequest)
	if err != nil {
		return middleware.Error(http.StatusUnauthorized, "Authorization error")
	}
	if role != "admin" {
		return middleware.Error(http.StatusForbidden, "No rights")
	}

	ctx := context.Background()

	coaches, err := CoachesHandlersApp.Services.CoachService.GetAll(ctx)
	if err != nil {
		return middleware.Error(http.StatusInternalServerError, "Can't get coach")
	}

	coachDTO := make([]dto_models.Coach, len(coaches))
	for i, coach := range coaches {
		coachDTO[i] = dto_models.Coach{
			CoachID: &coach.ID,
			Name:    &coach.Name,
		}
	}

	return middleware.ResponderFunc(func(rw http.ResponseWriter, p runtime.Producer) {
		rw.WriteHeader(http.StatusOK)
		_ = p.Produce(rw, coachDTO)
	})
}