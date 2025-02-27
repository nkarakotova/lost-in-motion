package handlers

import (
	"net/http"
	"context"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/runtime"

	"lim/cmd/registry"
	"lim/internal/lim-handler/restapi/operations"
	"lim/internal/lim-handler/restapi/operations/halls"
	"lim/internal/lim-handler/restapi/jwt"
	"lim/internal/lim-core/models"
	"lim/internal/lim-core/services"
	"lim/internal/lim-core/errors/servicesErrors"

	dto_models "lim/internal/lim-handler/models"
)

var HallsHandlersApp registry.App

func ConfigureHallsHandlers(app registry.App, api *operations.SwaggerLIMAPI) {
	HallsHandlersApp = app
	api.HallsCreateHallHandler = halls.CreateHallHandlerFunc(createHallHandlerFunc)
	api.HallsGetHallsHandler = halls.GetHallsHandlerFunc(getHallsHandlerFunc)
}

func createHallHandlerFunc(params halls.CreateHallParams, principal interface{}) middleware.Responder {
	_, role, err := jwt.GetAuthenticatedUser(params.HTTPRequest)
	if err != nil {
		return middleware.Error(http.StatusUnauthorized, "Authorization error")
	}
	if role != "admin" {
		return middleware.Error(http.StatusForbidden, "No rights")
	}

	ctx := context.Background()

	number := params.HallInfo.Number

	hall := models.Hall{
		Number: *number,
	}

	err = HallsHandlersApp.Services.HallService.Create(ctx, &hall)
	if err != nil && err == servicesErrors.HallAlreadyExists {
		return middleware.Error(http.StatusConflict, "Hall alredy exists")
	} else if err != nil {
		return middleware.Error(http.StatusInternalServerError, "Сan't create hall")
	}

	return middleware.ResponderFunc(func(rw http.ResponseWriter, p runtime.Producer) {
		rw.WriteHeader(http.StatusOK)
	})
}

func getHallsHandlerFunc(params halls.GetHallsParams, principal interface{}) middleware.Responder {
	_, role, err := jwt.GetAuthenticatedUser(params.HTTPRequest)
	if err != nil {
		return middleware.Error(http.StatusUnauthorized, "Authorization error")
	}
	if role != "admin" {
		return middleware.Error(http.StatusForbidden, "No rights")
	}

	ctx := context.Background()

	date, err := time.Parse("2006-01-02", params.Date.String())
	if err != nil {
		return middleware.Error(http.StatusInternalServerError, "Can't parse date")
	}

	slots := make([]time.Time, services.LastTrainingTime-services.FirstTrainingTime)
	slots[0] = date.Add(time.Duration(services.FirstTrainingTime) * time.Hour)

	for idx := 1; idx < len(slots); idx++ {
		slots[idx] = slots[idx-1].Add(time.Hour)
	}

	type HallResponse struct {
		Timestamp string            `json:"timestamp"`
		Hall      []dto_models.Hall `json:"hall"`
	}

	hallTimeDTO := make([]HallResponse, len(slots))

	for i, t := range slots {
		hour, _, _ := t.Clock()
		halls, err := HallsHandlersApp.Services.HallService.GetFreeOnDateTime(ctx, time.Date(date.Year(), date.Month(), date.Day(), hour, 0, 0, 0, time.UTC))
		if err != nil {
			return middleware.Error(http.StatusInternalServerError, "Can't get halls")
		}

		hallsDTO := make([]dto_models.Hall, len(halls))
		j := 0
		for _, hall := range halls {
			hallsDTO[j] = dto_models.Hall{
				HallID: &hall.ID,
				Number: &hall.Number,
			}
			j++
		}

		hallTimeDTO[i] = HallResponse{
			Timestamp: t.GoString(),
			Hall:      hallsDTO,
		}
	}

	return middleware.ResponderFunc(func(rw http.ResponseWriter, p runtime.Producer) {
		rw.WriteHeader(http.StatusOK)
		_ = p.Produce(rw, hallTimeDTO)
	})
}