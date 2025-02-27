package handlers

import (
	"net/http"
	"context"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/runtime"

	"lim/cmd/registry"
	"lim/internal/lim-handler/restapi/operations"
	"lim/internal/lim-handler/restapi/jwt"
	"lim/internal/lim-handler/restapi/operations/assignments"
	"lim/internal/lim-core/errors/servicesErrors"
)

var AssignmentsHandlersApp registry.App

func ConfigureAssignmentsHandlers(app registry.App, api *operations.SwaggerLIMAPI) {
	AssignmentsHandlersApp = app
	api.AssignmentsCreateAssignmentHandler = assignments.CreateAssignmentHandlerFunc(createAssignmentHandlerFunc)
	api.AssignmentsDeleteAssignmentHandler = assignments.DeleteAssignmentHandlerFunc(deleteAssignmentHandlerFunc)
}

func createAssignmentHandlerFunc(params assignments.CreateAssignmentParams, principal interface{}) middleware.Responder {
	telephone, role, err := jwt.GetAuthenticatedUser(params.HTTPRequest)
	if err != nil {
		return middleware.Error(http.StatusUnauthorized, "Authorization error")
	}
	if role != "client" {
		return middleware.Error(http.StatusForbidden, "No rights")
	}

	ctx := context.Background()

	client, err := AssignmentsHandlersApp.Services.ClientService.GetByTelephone(ctx, telephone)
	if err != nil && err == servicesErrors.ClientDoesNotExists {
		return middleware.Error(http.StatusNotFound, "Client does not exists")
	} else if err != nil {
		return middleware.Error(http.StatusInternalServerError, "Can't get client by telephone")
	}

	training_id := params.TrainingID.TrainingID

	err = AssignmentsHandlersApp.Services.ClientService.CreateAssignment(ctx, client.ID, *training_id)
	if err != nil && err == servicesErrors.ClientDoesNotExists {
		return middleware.Error(http.StatusNotFound, "Client does not exists")
	} else if err != nil && err == servicesErrors.TrainingDoesNotExists {
		return middleware.Error(http.StatusNotFound, "Training does not exists")
	} else if err != nil && err == servicesErrors.NoAvailablePlacesNum {
		return middleware.Error(http.StatusConflict, "No available places num")
	} else if err != nil && err == servicesErrors.AssignmentOnThisTimeAlreadyExists {
		return middleware.Error(http.StatusConflict, "Assignment on this time alredy exists")
	} else if err != nil {
		return middleware.Error(http.StatusInternalServerError, "Can't create assignment")
	}

	return middleware.ResponderFunc(func(rw http.ResponseWriter, p runtime.Producer) {
		rw.WriteHeader(http.StatusOK)
	})
}

func deleteAssignmentHandlerFunc(params assignments.DeleteAssignmentParams, principal interface{}) middleware.Responder {
	telephone, role, err := jwt.GetAuthenticatedUser(params.HTTPRequest)
	if err != nil {
		return middleware.Error(http.StatusUnauthorized, "Unauthorized")
	}
	if role != "client" {
		return middleware.Error(http.StatusForbidden, "Access denied")
	}

	ctx := context.Background()

	client, err := AssignmentsHandlersApp.Services.ClientService.GetByTelephone(ctx, telephone)
	if err != nil && err == servicesErrors.ClientDoesNotExists {
		return middleware.Error(http.StatusNotFound, "Client does not exists")
	} else if err != nil {
		return middleware.Error(http.StatusInternalServerError, "Can't get client by telephone")
	}

	training_id := params.TrainingID.TrainingID

	err = AssignmentsHandlersApp.Services.ClientService.DeleteAssignment(ctx, client.ID, *training_id)
	if err != nil && err == servicesErrors.AssignmentDoesNotExists {
		return middleware.Error(http.StatusNotFound, "Assignment does not exists")
	} else if err != nil {
		return middleware.Error(http.StatusInternalServerError, "Can't delete assignment")
	}

	return middleware.ResponderFunc(func(rw http.ResponseWriter, p runtime.Producer) {
		rw.WriteHeader(http.StatusOK)
	})
}