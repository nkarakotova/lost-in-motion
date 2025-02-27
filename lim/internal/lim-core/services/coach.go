package services

import (
	"time"
	"context"
	"lim/internal/lim-core/models"
)

type CoachService interface {
	Create(ctx context.Context, coach *models.Coach) error
	GetByID(ctx context.Context, id uint64) (*models.Coach, error)
	GetByName(ctx context.Context, name string) (*models.Coach, error)
	GetAll(ctx context.Context) ([]models.Coach, error)
	GetFreeTimeOnDate(ctx context.Context, id uint64, date time.Time) ([]time.Time, error)
}
