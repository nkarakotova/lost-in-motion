package servicesImplementation

import (
	"context"
	"time"

	"lim/internal/lim-core/errors/repositoriesErrors"
	"lim/internal/lim-core/errors/servicesErrors"
	"lim/internal/lim-core/repositories"
	"lim/internal/lim-core/services"

	"lim/internal/lim-core/models"

	"github.com/charmbracelet/log"
)

type CoachServiceImplementation struct {
	CoachRepository    repositories.CoachRepository
	TrainingRepository repositories.TrainingRepository
	logger             *log.Logger
}

func NewCoachServiceImplementation(
	CoachRepository repositories.CoachRepository,
	TrainingRepository repositories.TrainingRepository,
	logger *log.Logger,
) services.CoachService {

	return &CoachServiceImplementation{
		CoachRepository:    CoachRepository,
		TrainingRepository: TrainingRepository,
		logger:             logger,
	}
}

func (c *CoachServiceImplementation) validate(ctx context.Context, coach *models.Coach) error {
	_, err := c.CoachRepository.GetByName(ctx, coach.Name)
	if err != nil && err != repositoriesErrors.EntityDoesNotExists {
		c.logger.Warn("COACH! Error in repository GetByName", "name", coach.Name, "error", err)
		return err
	} else if err == nil {
		c.logger.Warn("COACH! Coach already exists", "name", coach.Name)
		return servicesErrors.CoachAlreadyExists
	}

	return nil
}

func (c *CoachServiceImplementation) GetByName(ctx context.Context, name string) (*models.Coach, error) {
	coach, err := c.CoachRepository.GetByName(ctx, name)
	if err != nil {
		c.logger.Warn("COACH! Error in repository GetByName", "name", name, "error", err)
		return nil, err
	}

	c.logger.Debug("COACH! Success GetByName", "name", name)
	return coach, nil
}

func (c *CoachServiceImplementation) Create(ctx context.Context, coach *models.Coach) error {
	err := c.validate(ctx, coach)
	if err != nil {
		return err
	}

	err = c.CoachRepository.Create(ctx, coach)
	if err != nil {
		c.logger.Warn("COACH! Error in repository Create", "name", coach.Name, "error", err)
		return err
	}

	c.logger.Info("COACH! Success create coach", "name", coach.Name)
	return nil
}

func (c *CoachServiceImplementation) GetByID(ctx context.Context, id uint64) (*models.Coach, error) {
	coach, err := c.CoachRepository.GetByID(ctx, id)
	if err != nil {
		c.logger.Warn("COACH! Error in repository method GetByID", "id", id, "error", err)
		return nil, err
	}

	c.logger.Debug("COACH! Success GetByID", "id", id)
	return coach, nil
}

func (c *CoachServiceImplementation) getAllSlots(date time.Time) []time.Time {
	slots := make([]time.Time, services.LastTrainingTime-services.FirstTrainingTime)
	slots[0] = date.Add(time.Duration(services.FirstTrainingTime) * time.Hour)

	for idx := 1; idx < len(slots); idx++ {
		slots[idx] = slots[idx-1].Add(time.Hour)
	}

	return slots
}

func (c *CoachServiceImplementation) GetFreeTimeOnDate(ctx context.Context, id uint64, date time.Time) ([]time.Time, error) {
	trainings, err := c.TrainingRepository.GetAllByCoachOnDate(ctx, id, date)
	if err != nil {
		c.logger.Warn("TRAINING! Error in repository method GetAllByCoachOnDate", "id", id, "err", err)
		return nil, err
	}

	slots := c.getAllSlots(date)

	for _, t := range trainings {
		time := t.DateTime
		for i, slot := range slots {
			if time.Equal(slot) {
				slots = append(slots[:i], slots[i+1:]...)
				break
			}
		}
		if len(slots) == 0 {
			break
		}
	}

	c.logger.Debug("COACH! Success GetFreeTimeOnDate", "id", id)
	return slots, nil
}

func (c *CoachServiceImplementation) GetAll(ctx context.Context) ([]models.Coach, error) {
	coaches, err := c.CoachRepository.GetAll(ctx)
	if err != nil {
		c.logger.Warn("COACH! Error in repository method GetAll", "err", err)
		return nil, err
	}

	c.logger.Debug("COACH! Success GetAllByDirection")
	return coaches, nil
}
