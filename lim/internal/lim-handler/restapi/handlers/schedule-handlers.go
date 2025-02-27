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
	"lim/internal/lim-handler/restapi/operations/schedule"

	dto_models "lim/internal/lim-handler/models"
)

var ScheduleHandlersApp registry.App

func ConfigureScheduleHandlers(app registry.App, api *operations.SwaggerLIMAPI) {
	ScheduleHandlersApp = app
	api.ScheduleGetScheduleHandler = schedule.GetScheduleHandlerFunc(getScheduleHandlerFunc)
}

func getScheduleHandlerFunc(params schedule.GetScheduleParams) middleware.Responder {
	ctx := context.Background()

	// Проверять на роль и маркировать рассписание для разных действий

	currentTime := time.Now()
	oneMonthLater := currentTime.AddDate(0, 1, 0)
	trainings, err := ScheduleHandlersApp.Services.TrainingService.GetAllBetweenDateTime(ctx, currentTime, oneMonthLater)
	if err != nil {
		return middleware.Error(http.StatusInternalServerError, "Can't get schedule")
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