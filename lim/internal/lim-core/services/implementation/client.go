package servicesImplementation

import (
	"context"
	"net/mail"
	"strconv"

	"lim/internal/lim-core/errors/repositoriesErrors"
	"lim/internal/lim-core/errors/servicesErrors"
	"lim/internal/lim-core/managers"
	"lim/internal/lim-core/repositories"
	"lim/internal/lim-core/services"

	"lim/internal/lim-core/models"

	"github.com/charmbracelet/log"
)

type ClientServiceImplementation struct {
	ClientRepository   repositories.ClientRepository
	TrainingRepository repositories.TrainingRepository
	TransactionManager managers.TransactionManager
	logger             *log.Logger
}

func NewClientServiceImplementation(
	ClientRepository repositories.ClientRepository,
	TrainingRepository repositories.TrainingRepository,
	TransactionManager managers.TransactionManager,
	logger *log.Logger,
) services.ClientService {

	return &ClientServiceImplementation{
		ClientRepository:   ClientRepository,
		TrainingRepository: TrainingRepository,
		TransactionManager: TransactionManager,
		logger:             logger,
	}
}

const TelephoneNumberLen = 10

func (c *ClientServiceImplementation) validate(ctx context.Context, client *models.Client) error {
	_, err := c.ClientRepository.GetByTelephone(ctx, client.Telephone)
	if err != nil && err != repositoriesErrors.EntityDoesNotExists {
		c.logger.Warn("CLIENT! Error in repository GetClientByTelephone", "telephone", client.Telephone, "error", err)
		return err
	} else if err == nil {
		c.logger.Warn("CLIENT! Client already exists", "telephone", client.Telephone)
		return servicesErrors.ClientAlreadyExists
	}

	if len(client.Telephone) != TelephoneNumberLen {
		c.logger.Warn("CLIENT! Client telephone length incorrect", "telephone", client.Telephone)
		return servicesErrors.ClientTelephoneIncorrect
	}

	_, err = strconv.Atoi(client.Telephone)
	if err != nil {
		c.logger.Warn("CLIENT! Client telephone incorrect", "telephone", client.Telephone)
		return servicesErrors.ClientTelephoneIncorrect
	}

	_, err = mail.ParseAddress(client.Mail)
	if err != nil {
		c.logger.Warn("CLIENT! Client mail incorrect", "mail", client.Mail)
		return servicesErrors.ClientMailIncorrect
	}

	return nil
}

func (c *ClientServiceImplementation) GetByTelephone(ctx context.Context, telephone string) (*models.Client, error) {
	client, err := c.ClientRepository.GetByTelephone(ctx, telephone)

	if err != nil && err == repositoriesErrors.EntityDoesNotExists {
		c.logger.Warn("CLIENT! Сlient with this telephone does not exists", "telephone", telephone, "error", err)
		return nil, servicesErrors.ClientDoesNotExists
	} else if err != nil {
		c.logger.Warn("CLIENT! Error in repository GetClientByTelephone", "telephone", telephone, "error", err)
		return nil, err
	}

	c.logger.Debug("CLIENT! Success GetClientByTelephone", "telephone", telephone)
	return client, nil
}

func (c *ClientServiceImplementation) Create(ctx context.Context, client *models.Client) error {
	err := c.validate(ctx, client)
	if err != nil {
		return err
	}

	err = c.ClientRepository.Create(ctx, client)
	if err != nil {
		c.logger.Warn("CLIENT! Error in repository Create", "telephone", client.Telephone, "error", err)
		return err
	}

	c.logger.Info("CLIENT! Success create client", "telephone", client.Telephone, "id", client.ID)
	return nil
}

func (c *ClientServiceImplementation) Login(ctx context.Context, telephone, password string) (*models.Client, error) {
	tempClient, err := c.ClientRepository.GetByTelephone(ctx, telephone)
	if err != nil && err == repositoriesErrors.EntityDoesNotExists {
		c.logger.Warn("CLIENT! Сlient with this telephone does not exists", "telephone", telephone, "error", err)
		return nil, servicesErrors.ClientDoesNotExists
	} else if err != nil {
		c.logger.Warn("CLIENT! Error in repository method GetByTelephone", "telepnone", telephone, "error", err)
		return nil, err
	}

	if password != tempClient.Password {
		c.logger.Warn("CLIENT! Error client password", "telephone", telephone)
		return nil, servicesErrors.InvalidPassword
	}

	c.logger.Info("CLIENT! Success login", "telepnone", telephone, "id", tempClient.ID)
	return tempClient, nil
}

func (c *ClientServiceImplementation) GetByID(ctx context.Context, id uint64) (*models.Client, error) {
	client, err := c.ClientRepository.GetByID(ctx, id)
	if err != nil {
		c.logger.Warn("CLIENT! Error in repository method GetByID", "id", id, "error", err)
		return nil, err
	}

	c.logger.Debug("CLIENT! Success GetByID", "id", id)
	return client, nil
}

func (c *ClientServiceImplementation) trainingInSameDateTimeCheck(ctx context.Context, client *models.Client, training *models.Training) error {
	curStart := training.DateTime

	clientTrainings, err := c.TrainingRepository.GetAllByClient(ctx, client.ID)
	if err != nil {
		c.logger.Warn("TRAINING! Error in repository method GetAllByClient", "id", client.ID, "error", err)
		return err
	}

	for _, t := range clientTrainings {
		if t.DateTime == curStart {
			c.logger.Warn("CLIENT! There is already an assignment for this time", "id", client.ID, "error", err)
			return servicesErrors.AssignmentOnThisTimeAlreadyExists
		}
	}

	return nil
}

func (c *ClientServiceImplementation) CreateAssignmentChecks(ctx context.Context, client *models.Client, training *models.Training) error {
	if training.PlacesNum == 0 {
		c.logger.Warn("TRAINING! There is no available places num", "id", training.ID)
		return servicesErrors.NoAvailablePlacesNum
	}

	err := c.trainingInSameDateTimeCheck(ctx, client, training)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientServiceImplementation) createAssignment(ctx context.Context, client *models.Client, training *models.Training) error {
	return c.TransactionManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		err := c.ClientRepository.CreateAssignment(ctx, client.ID, training.ID)
		if err != nil {
			c.logger.Warn("CLIENT! Error in repository CreateAssignment", "clientID", client.ID, "trainingID", training.ID, "error", err)
			return err
		}

		err = c.TrainingRepository.ReduceAvailablePlacesNum(ctx, training.ID)
		if err != nil {
			c.logger.Warn("TRAINING! Error in repository ReduceAvailablePlacesNum", "trainingID", training.ID, "error", err)
			return err
		}

		return nil
	})
}

func (c *ClientServiceImplementation) CreateAssignment(ctx context.Context, clientID, trainingID uint64) error {
	client, err := c.ClientRepository.GetByID(ctx, clientID)
	if err != nil && err == repositoriesErrors.EntityDoesNotExists {
		c.logger.Warn("CLIENT! Сlient does not exists", "error", err)
		return servicesErrors.ClientDoesNotExists
	} else if err != nil {
		c.logger.Warn("CLIENT! Error in repository method GetByID", "id", clientID, "error", err)
		return err
	}

	training, err := c.TrainingRepository.GetByID(ctx, trainingID)
	if err != nil && err == repositoriesErrors.EntityDoesNotExists {
		c.logger.Warn("CLIENT! Training does not exists", "error", err)
		return servicesErrors.TrainingDoesNotExists
	} else if err != nil {
		c.logger.Warn("TRAINING! Error in repository method GetByID", "id", trainingID, "error", err)
		return err
	}

	err = c.CreateAssignmentChecks(ctx, client, training)
	if err != nil {
		return err
	}

	err = c.createAssignment(ctx, client, training)
	if err != nil {
		return err
	}

	c.logger.Info("CLIENT! Success create assignment", "clientID", clientID, "trainingID", clientID)
	return nil
}

func (c *ClientServiceImplementation) DeleteAssignment(ctx context.Context, clientID, trainingID uint64) error {
	return c.TransactionManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		err := c.ClientRepository.DeleteAssignment(ctx, clientID, trainingID)
		if err != nil && err == repositoriesErrors.EntityDoesNotExists {
			c.logger.Warn("CLIENT! Assignment does not exists", "error", err)
			return servicesErrors.AssignmentDoesNotExists
		} else if err != nil {
			c.logger.Warn("CLIENT! Error in repository DeleteAssignment", "clientID", clientID, "trainingID", clientID, "error", err)
			return err
		}

		err = c.TrainingRepository.IncreaseAvailablePlacesNum(ctx, trainingID)
		if err != nil {
			c.logger.Warn("TRAINING! Error in repository IncreaseAvailablePlacesNum", "trainingID", trainingID, "error", err)
			return err
		}

		c.logger.Info("CLIENT! Success delete assignment", "clientID", clientID, "trainingID", clientID)
		return nil
	})
}

func (c *ClientServiceImplementation) ChangePassword(ctx context.Context, clientID uint64, password string) error {
	err := c.ClientRepository.ChangePassword(ctx, clientID, password)
	if err != nil {
		c.logger.Warn("CLIENT! Error in repository ChangePassword", "clientID", clientID, "error", err)
		return err
	}

	c.logger.Info("CLIENT! Success change password", "clientID", clientID)
	return nil
}
