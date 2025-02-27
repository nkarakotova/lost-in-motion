package servicesImplementation

import (
	"context"
	"time"

	"lim/internal/lim-core/errors/repositoriesErrors"
	"lim/internal/lim-core/errors/servicesErrors"
	"lim/internal/lim-core/models"
	"lim/internal/lim-core/repositories"
	"lim/internal/lim-core/services"

	"github.com/charmbracelet/log"
)

type HallServiceImplementation struct {
	HallRepository     repositories.HallRepository
	TrainingRepository repositories.TrainingRepository
	logger             *log.Logger
}

func NewHallServiceImplementation(
	HallRepository repositories.HallRepository,
	TrainingRepository repositories.TrainingRepository,
	logger *log.Logger,
) services.HallService {

	return &HallServiceImplementation{
		HallRepository:     HallRepository,
		TrainingRepository: TrainingRepository,
		logger:             logger,
	}
}

func (h *HallServiceImplementation) validate(ctx context.Context, hall *models.Hall) error {
	_, err := h.HallRepository.GetByNumber(ctx, hall.Number)
	if err != nil && err != repositoriesErrors.EntityDoesNotExists {
		h.logger.Warn("HALL! Error in repository GetByNumber", "number", hall.Number, "error", err)
		return err
	} else if err == nil {
		h.logger.Warn("HALL! Hall already exists", "number", hall.Number)
		return servicesErrors.HallAlreadyExists
	}

	return nil
}

func (h *HallServiceImplementation) GetByNumber(ctx context.Context, number uint64) (*models.Hall, error) {
	hall, err := h.HallRepository.GetByNumber(ctx, number)
	if err != nil {
		h.logger.Warn("HALL! Error in repository GetByNumber", "number", number, "error", err)
		return nil, err
	}

	h.logger.Debug("HALL! Success GetByNumber", "number", number)
	return hall, nil
}

func (h *HallServiceImplementation) Create(ctx context.Context, hall *models.Hall) error {
	err := h.validate(ctx, hall)
	if err != nil {
		return err
	}

	err = h.HallRepository.Create(ctx, hall)
	if err != nil {
		h.logger.Warn("HALL! Error in repository Create", "number", hall.Number, "error", err)
		return err
	}

	h.logger.Info("HALL! Success create hall", "number", hall.Number)
	return nil
}

func (h *HallServiceImplementation) GetByID(ctx context.Context, id uint64) (*models.Hall, error) {
	hall, err := h.HallRepository.GetByID(ctx, id)
	if err != nil {
		h.logger.Warn("HALL! Error in repository method GetByID", "id", id, "error", err)
		return nil, err
	}

	h.logger.Debug("HALL! Success GetByID", "id", id)
	return hall, nil
}

func (h *HallServiceImplementation) GetFreeOnDateTime(ctx context.Context, slot time.Time) (map[uint64]models.Hall, error) {
	trainings, err := h.TrainingRepository.GetAllByDateTime(ctx, slot)
	if err != nil {
		h.logger.Warn("TRAINING! Error in repository method GetAllByDateTime", "error", err)
		return nil, err
	}

	freeHalls, err := h.HallRepository.GetAll(ctx)
	if err != nil {
		h.logger.Warn("HALL! Error in repository method GetAll", "error", err)
		return nil, err
	}

	for _, t := range trainings {
		delete(freeHalls, t.HallID)
	}

	h.logger.Debug("HALL! Success GetFreeOnDateTime")
	return freeHalls, nil
}
