package handlers

import (
	"net/http"
	"context"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"lim/cmd/registry"
	"lim/internal/lim-handler/restapi/operations"
	"lim/internal/lim-handler/restapi/operations/clients"
	"lim/internal/lim-handler/restapi/jwt"
	"lim/internal/lim-core/errors/servicesErrors"

	dto_models "lim/internal/lim-handler/models"
)

var ClientsHandlersApp registry.App

func ConfigureClientsHandlers(app registry.App, api *operations.SwaggerLIMAPI) {
	ClientsHandlersApp = app
	api.ClientsChangePasswordHandler = clients.ChangePasswordHandlerFunc(changePasswordHandlerFunc)
	api.ClientsGetTrainingsByClientHandler = clients.GetTrainingsByClientHandlerFunc(getTrainingsByClientHandlerFunc)
}

func changePasswordHandlerFunc(params clients.ChangePasswordParams, principal interface{}) middleware.Responder {
	telephone, role, err := jwt.GetAuthenticatedUser(params.HTTPRequest)
	if err != nil {
		return middleware.Error(http.StatusUnauthorized, "Authorization error")
	}
	if role != "client" {
		return middleware.Error(http.StatusForbidden, "No rights")
	}

	ctx := context.Background()

	client, err := ClientsHandlersApp.Services.ClientService.GetByTelephone(ctx, telephone)
	if err != nil && err == servicesErrors.ClientDoesNotExists {
		return middleware.Error(http.StatusNotFound, "Client does not exists")
	} else if err != nil {
		return middleware.Error(http.StatusInternalServerError, "Can't get client by telephone")
	}

	password := params.NewPassword.Password

	err = ClientsHandlersApp.Services.ClientService.ChangePassword(ctx, client.ID, *password)
	if err != nil {
		return middleware.Error(http.StatusInternalServerError, "Can't change password")
	}

	return middleware.ResponderFunc(func(rw http.ResponseWriter, p runtime.Producer) {
		rw.WriteHeader(http.StatusOK)
	})
}

func getTrainingsByClientHandlerFunc(params clients.GetTrainingsByClientParams, principal interface{}) middleware.Responder {
	telephone, role, err := jwt.GetAuthenticatedUser(params.HTTPRequest)
	if err != nil {
		return middleware.Error(http.StatusUnauthorized, "Authorization error")
	}
	if role != "client" {
		return middleware.Error(http.StatusForbidden, "No rights")
	}

	ctx := context.Background()

	client, err := ClientsHandlersApp.Services.ClientService.GetByTelephone(ctx, telephone)
	if err != nil && err == servicesErrors.ClientDoesNotExists {
		return middleware.Error(http.StatusNotFound, "Client does not exists")
	} else if err != nil {
		return middleware.Error(http.StatusInternalServerError, "Can't get client by telephone")
	}

	currentTime := time.Now()
	oneMonthLater := currentTime.AddDate(0, 1, 0)
	trainings, err := TrainingsHandlersApp.Services.TrainingService.GetAllByClientBetweenDateTime(ctx, client.ID, currentTime, oneMonthLater)
	if err != nil {
		return middleware.Error(http.StatusInternalServerError, "Can't get trainings by client")
	}

	trainingDTO := make([]dto_models.Training, len(trainings))
	for i, training := range trainings {
		copyTime := strfmt.DateTime(training.DateTime)

		coach, _ := ScheduleHandlersApp.Services.CoachService.GetByID(ctx, training.CoachID)
		coachName := coach.Name

		hall, _ := ScheduleHandlersApp.Services.HallService.GetByID(ctx, training.HallID)
		hallNumber := hall.Number

		trainingDTO[i] = dto_models.Training{
			TrainingID: &training.ID,
			CoachName:  &coachName,
			HallNumber: &hallNumber,
			Name:       &training.Name,
			DateTime:   &copyTime,
			PlacesNum:  &training.PlacesNum,
		}
	}

	return middleware.ResponderFunc(func(rw http.ResponseWriter, p runtime.Producer) {
		rw.WriteHeader(http.StatusOK)
		_ = p.Produce(rw, trainingDTO)
	})
}