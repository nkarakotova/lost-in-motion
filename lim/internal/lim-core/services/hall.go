package services

import (
	"time"
	"context"
	"lim/internal/lim-core/models"
)

type HallService interface {
	Create(ctx context.Context, hall *models.Hall) error
	GetByID(ctx context.Context, id uint64) (*models.Hall, error)
	GetByNumber(ctx context.Context, number uint64) (*models.Hall, error)
	GetFreeOnDateTime(ctx context.Context, dateTime time.Time) (map[uint64]models.Hall, error)
}
