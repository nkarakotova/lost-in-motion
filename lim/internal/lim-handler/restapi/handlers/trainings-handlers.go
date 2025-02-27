package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"lim/cmd/registry"
	"lim/internal/lim-core/errors/servicesErrors"
	"lim/internal/lim-core/models"
	"lim/internal/lim-handler/restapi/jwt"
	"lim/internal/lim-handler/restapi/operations"
	"lim/internal/lim-handler/restapi/operations/trainings"
)

var TrainingsHandlersApp registry.App

func ConfigureTrainingsHandlers(app registry.App, api *operations.SwaggerLIMAPI) {
	TrainingsHandlersApp = app
	api.TrainingsCreateTrainingHandler = trainings.CreateTrainingHandlerFunc(createTrainingHandlerFunc)
	api.TrainingsDeleteTrainingHandler = trainings.DeleteTrainingHandlerFunc(deleteTrainingHandlerFunc)
}

func createTrainingHandlerFunc(params trainings.CreateTrainingParams, principal interface{}) middleware.Responder {
	_, role, err := jwt.GetAuthenticatedUser(params.HTTPRequest)
	if err != nil {
		return middleware.Error(http.StatusUnauthorized, "Authorization error")
	}
	if role != "admin" {
		return middleware.Error(http.StatusForbidden, "No rights")
	}

	ctx := context.Background()

	coach, _ := ScheduleHandlersApp.Services.CoachService.GetByName(ctx, *params.TrainingInfo.CoachName)
	coachID := coach.ID

	hall, _ := ScheduleHandlersApp.Services.HallService.GetByNumber(ctx, *params.TrainingInfo.HallNumber)
	hallID := hall.ID

	coach_id := coachID
	hall_id := hallID
	name := params.TrainingInfo.Name
	date_time := params.TrainingInfo.DateTime
	places_num := params.TrainingInfo.PlacesNum

	training := models.Training{
		CoachID:   coach_id,
		HallID:    hall_id,
		Name:      *name,
		DateTime:  time.Time(*date_time),
		PlacesNum: *places_num,
	}

	err = TrainingsHandlersApp.Services.TrainingService.Create(ctx, &training)
	if err != nil && err == servicesErrors.IncorrectTrainingTime {
		return middleware.Error(http.StatusUnprocessableEntity, "Incorrect training time")
	} else if err != nil && err == servicesErrors.BusyDateTime {
		return middleware.Error(http.StatusConflict, "Busy time")
	} else if err != nil && err == servicesErrors.CoachDoesNotExists {
		return middleware.Error(http.StatusNotFound, "Coach does not exists")
	} else if err != nil && err == servicesErrors.HallDoesNotExists {
		return middleware.Error(http.StatusNotFound, "Hall does not exists")
	} else if err != nil {
		return middleware.Error(http.StatusInternalServerError, "Can't create training")
	}

	return middleware.ResponderFunc(func(rw http.ResponseWriter, p runtime.Producer) {
		rw.WriteHeader(http.StatusOK)
	})
}

func deleteTrainingHandlerFunc(params trainings.DeleteTrainingParams, principal interface{}) middleware.Responder {
	_, role, err := jwt.GetAuthenticatedUser(params.HTTPRequest)
	if err != nil {
		return middleware.Error(http.StatusUnauthorized, "Authorization error")
	}
	if role != "admin" {
		return middleware.Error(http.StatusForbidden, "No rights")
	}

	ctx := context.Background()

	training_id := params.TrainingID.TrainingID

	err = TrainingsHandlersApp.Services.TrainingService.Delete(ctx, *training_id)
	if err != nil && err == servicesErrors.TrainingDoesNotExists {
		return middleware.Error(http.StatusNotFound, "Training does not exists")
	} else if err != nil {
		return middleware.Error(http.StatusInternalServerError, "Can't delete training")
	}

	return middleware.ResponderFunc(func(rw http.ResponseWriter, p runtime.Producer) {
		rw.WriteHeader(http.StatusOK)
	})
}
