package repositories

import (
	"context"

	"lim/internal/lim-core/models"
)

//go:generate mockgen  -source=client.go -destination=mocks/client.go
type ClientRepository interface {
	Create(ctx context.Context, client *models.Client) error
	GetByID(ctx context.Context, id uint64) (*models.Client, error)
	GetByTelephone(ctx context.Context, telephone string) (*models.Client, error)
	GetByTraining(ctx context.Context, id uint64) ([]models.Client, error)
	CreateAssignment(ctx context.Context, clientID, trainingID uint64) error
	DeleteAssignment(ctx context.Context, clientID, trainingID uint64) error
	ChangePassword(ctx context.Context, clientID uint64, password string) error
}
