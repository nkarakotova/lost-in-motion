package services

import (
	"time"
	"context"
	"lim/internal/lim-core/models"
)

type TrainingService interface {
	Create(ctx context.Context, training *models.Training) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*models.Training, error)
	GetAllByClient(ctx context.Context, id uint64) ([]models.Training, error)
	GetAllByCoachOnDate(ctx context.Context, id uint64, date time.Time) ([]models.Training, error)
	GetAllByDateTime(ctx context.Context, dateTime time.Time) ([]models.Training, error)
	GetAllBetweenDateTime(ctx context.Context, start time.Time, end time.Time) ([]models.Training, error)
	GetAllByClientBetweenDateTime(ctx context.Context, id uint64, start time.Time, end time.Time) ([]models.Training, error)
}

var FirstTrainingTime = 10
var LastTrainingTime = 22
