package services

import (
	"context"
	"lim/internal/lim-core/models"
)

type ClientService interface {
	Create(ctx context.Context, client *models.Client) error
	Login(ctx context.Context, telephone, password string) (*models.Client, error)
	GetByID(ctx context.Context, id uint64) (*models.Client, error)
	GetByTelephone(ctx context.Context, login string) (*models.Client, error)
	CreateAssignment(ctx context.Context, clientID, trainingID uint64) error
	DeleteAssignment(ctx context.Context, clientID, trainingID uint64) error
	ChangePassword(ctx context.Context, clientID uint64, password string) error
}
