package servicesImplementation

import (
	"context"
	"time"

	"lim/internal/lim-core/errors/servicesErrors"
	"lim/internal/lim-core/errors/repositoriesErrors"
	"lim/internal/lim-core/managers"
	"lim/internal/lim-core/repositories"
	"lim/internal/lim-core/services"

	"lim/internal/lim-core/models"

	"github.com/charmbracelet/log"
)

type TrainingServiceImplementation struct {
	TrainingRepository repositories.TrainingRepository
	ClientRepository   repositories.ClientRepository
	CoachRepository    repositories.CoachRepository
	HallRepository     repositories.HallRepository
	TransactionManager managers.TransactionManager
	logger             *log.Logger
}

func NewTrainingServiceImplementation(
	TrainingRepository repositories.TrainingRepository,
	ClientRepository repositories.ClientRepository,
	CoachRepository repositories.CoachRepository,
	HallRepository repositories.HallRepository,
	TransactionManager managers.TransactionManager,
	logger *log.Logger,
) services.TrainingService {

	return &TrainingServiceImplementation{
		TrainingRepository: TrainingRepository,
		ClientRepository:   ClientRepository,
		CoachRepository:    CoachRepository,
		HallRepository:     HallRepository,
		TransactionManager: TransactionManager,
		logger:             logger,
	}
}

func (t *TrainingServiceImplementation) validate(ctx context.Context, training *models.Training) error {
	h, m, s := training.DateTime.Clock()
	if h < services.FirstTrainingTime || h > services.LastTrainingTime || m != 0 || s != 0 {
		t.logger.Warn("TRAINING! Incorrect start time", "id", training.ID)
		return servicesErrors.IncorrectTrainingTime
	}

	_, err := t.CoachRepository.GetByID(ctx, training.CoachID)
	if err != nil && err == repositoriesErrors.EntityDoesNotExists {
		t.logger.Warn("TRAINING! Coach does not exists", "id", training.CoachID)
		return servicesErrors.CoachDoesNotExists
	} else if err != nil {
		t.logger.Warn("COACH! Error in repository GetByID", "id", training.CoachID)
		return err
	}

	_, err = t.HallRepository.GetByID(ctx, training.HallID)
	if err != nil && err == repositoriesErrors.EntityDoesNotExists {
		t.logger.Warn("TRAINING! Hall does not exists", "id", training.HallID)
		return servicesErrors.HallDoesNotExists
	} else if err != nil {
		t.logger.Warn("HALL! Error in repository GetByID", "id", training.HallID)
		return err
	}

	trainings, err := t.GetAllByDateTime(ctx, training.DateTime)
	if err != nil {
		t.logger.Warn("TRAINING! Error in repository GetAllByDateTime", "id", training.ID, "error", err)
		return err
	}
	for _, tr := range trainings {
		if tr.CoachID == training.CoachID || tr.HallID == training.HallID {
			t.logger.Warn("TRAINING! Busy date time", "id", training.ID)
			return servicesErrors.BusyDateTime
		}
	}

	return nil
}

func (t *TrainingServiceImplementation) Create(ctx context.Context, training *models.Training) error {
	err := t.validate(ctx, training)
	if err != nil {
		return err
	}

	err = t.TrainingRepository.Create(ctx, training)
	if err != nil {
		t.logger.Warn("TRAINING! Error in repository Create", "id", training.ID, "error", err)
		return err
	}

	t.logger.Info("TRAINING! Success create training", "id", training.ID)
	return nil
}

func (t *TrainingServiceImplementation) Delete(ctx context.Context, id uint64) error {
	err := t.TrainingRepository.Delete(ctx, id)

	if err != nil && err == repositoriesErrors.EntityDoesNotExists {
		t.logger.Warn("TRAINING! Training does not exists", "id", id)
		return servicesErrors.TrainingDoesNotExists
	} else if err != nil {
		t.logger.Warn("TRAINING! Error in repository Delete", "id", id, "error", err)
		return err
	}

	t.logger.Info("TRAINING! Success delete training", "id", id)
	return nil
}

func (t *TrainingServiceImplementation) GetByID(ctx context.Context, id uint64) (*models.Training, error) {
	training, err := t.TrainingRepository.GetByID(ctx, id)

	if err != nil {
		t.logger.Warn("TRAINING! Error in repository method GetByID", "id", id, "error", err)
		return nil, err
	}

	t.logger.Debug("TRAINING! Success GetByID", "id", id)
	return training, nil
}

func (t *TrainingServiceImplementation) GetAllByClient(ctx context.Context, id uint64) ([]models.Training, error) {
	trainings, err := t.TrainingRepository.GetAllByClient(ctx, id)

	if err != nil {
		t.logger.Warn("TRAINING! Error in repository method GetAllByClient", "id", id, "err", err)
		return nil, err
	}

	t.logger.Debug("TRAINING! Success GetAllByClient", "id", id)
	return trainings, nil
}

func (t *TrainingServiceImplementation) GetAllByCoachOnDate(ctx context.Context, id uint64, date time.Time) ([]models.Training, error) {
	trainings, err := t.TrainingRepository.GetAllByCoachOnDate(ctx, id, date)
	if err != nil {
		t.logger.Warn("TRAINING! Error in service method GetAllByCoachOnDate", "id", id, "err", err)
		return nil, err
	}

	t.logger.Debug("TRAINING! Successfully service method GetAllByCoachOnDate", "id", id)
	return trainings, nil
}

func (t *TrainingServiceImplementation) GetAllByDateTime(ctx context.Context, dateTime time.Time) ([]models.Training, error) {
	trainings, err := t.TrainingRepository.GetAllByDateTime(ctx, dateTime)

	if err != nil {
		t.logger.Warn("TRAINING! Error in repository method GetAllByDateTime", "dateTime", dateTime, "err", err)
		return nil, err
	}

	t.logger.Debug("TRAINING! Success GetAllByDateTime", "dateTime", dateTime)
	return trainings, nil
}

func (t *TrainingServiceImplementation) GetAllBetweenDateTime(ctx context.Context, start time.Time, end time.Time) ([]models.Training, error) {
	if start.After(end) {
		t.logger.Warn("TRAINING! Start after end", "start", start, "end", end)
		return nil, servicesErrors.StartDateAfterEndDate
	}

	trainings, err := t.TrainingRepository.GetAllBetweenDateTime(ctx, start, end)
	if err != nil {
		t.logger.Warn("TRAINING! Error in repository method GetAllBetweenDateTime", "start", start, "end", end, "err", err)
		return nil, err
	}

	t.logger.Debug("TRAINING! Success GetAllBetweenDateTime", "start", start, "end", end)
	return trainings, nil
}

func (t *TrainingServiceImplementation) GetAllByClientBetweenDateTime(ctx context.Context, id uint64, start time.Time, end time.Time) ([]models.Training, error) {
	if start.After(end) {
		t.logger.Warn("TRAINING! Start after end", "start", start, "end", end)
		return nil, servicesErrors.StartDateAfterEndDate
	}

	trainings, err := t.TrainingRepository.GetAllByClientBetweenDateTime(ctx, id, start, end)
	if err != nil {
		t.logger.Warn("TRAINING! Error in repository method GetAllByClientBetweenDateTime", "start", start, "end", end, "err", err)
		return nil, err
	}

	t.logger.Debug("TRAINING! Success GetAllByClientBetweenDateTime", "start", start, "end", end)
	return trainings, nil
}
