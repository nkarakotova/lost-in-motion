package integrationCoreTests

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"lim/internal/lim-core/errors/servicesErrors"
	transaction_manager "lim/internal/lim-core/managers"
	transaction_manager_implementation "lim/internal/lim-core/managers/implementation"
	"lim/internal/lim-core/services"
	service_implementation "lim/internal/lim-core/services/implementation"
	servicesDataBuilder "lim/internal/lim-core/services/implementation/data_builders"
	postgreSQL "lim/internal/lim-repo/postgreSQL"

	"lim/internal/lim-core/models"
	"lim/internal/lim-repo/flags"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	at_manager "github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/charmbracelet/log"
	"github.com/jmoiron/sqlx"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type ClientIntegrationCoreSuite struct {
	suite.Suite
	db            *sql.DB
	dbx           *sqlx.DB
	txResolver    *trmsqlx.CtxGetter
	trm           transaction_manager.TransactionManager
	clientService services.ClientService
	ctx           context.Context
}

func (s *ClientIntegrationCoreSuite) BeforeEach(t provider.T) {
	var err error
	pf := flags.PostgresFlags{
		Host:     "postgres",
		User:     "natali",
		Password: "12345",
		Port:     "5432",
		DBName:   "postgres",
	}

	s.db, err = pf.InitDB(log.Default())

	if err != nil {
		fmt.Println(err)
	}
	if s.db == nil {
		return
	}

	fields := postgreSQL.PostgresRepositoryFields{DB: s.db}

	s.ctx = context.Background()
	s.dbx = sqlx.NewDb(s.db, "pgx")
	clientRepository := postgreSQL.CreateClientPostgreSQLRepository(&fields)

	s.txResolver = trmsqlx.DefaultCtxGetter
	manager, _ := at_manager.New(trmsqlx.NewDefaultFactory(s.dbx))
	s.trm = transaction_manager_implementation.NewTransactionManagerImplementation(manager)

	logger := log.New(os.Stderr)
	trainingRepository := postgreSQL.CreateTrainingPostgreSQLRepository(&fields)

	s.clientService = service_implementation.NewClientServiceImplementation(clientRepository, trainingRepository, s.trm, logger)
}

func (s *ClientIntegrationCoreSuite) AfterEach(t provider.T) {
	s.db.Close()
}

func (s *ClientIntegrationCoreSuite) TestClientIntegrationCoreCreateSuccess(t provider.T) {
	t.Title("CreateClient: Success")
	t.Tags("Client")
	if IsUnitTestsFailed() {
		t.Skip()
	}
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		_ = s.trm.WithinTransaction(s.ctx, func(txCtx context.Context) error {
			err := s.clientService.Create(txCtx, servicesDataBuilder.NewClientBuilder().Build())
			sCtx.Assert().NoError(err)

			s.trm.Rollback(txCtx, s.txResolver, s.dbx)
			return nil
		})
	})
}

func (s *ClientIntegrationCoreSuite) TestClientIntegrationCoreCreateFailure(t provider.T) {
	t.Title("CreateClient: Failure")
	t.Tags("Client")
	if IsUnitTestsFailed() {
		t.Skip()
	}
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		_ = s.trm.WithinTransaction(s.ctx, func(txCtx context.Context) error {
			err := s.clientService.Create(txCtx, &models.Client{
				Name:      "Diana",
				Telephone: "1111111111",
				Mail:      "d@mail.ru",
				Password:  "111",
			})
			sCtx.Assert().ErrorIs(err, servicesErrors.ClientAlreadyExists)

			s.trm.Rollback(txCtx, s.txResolver, s.dbx)
			return nil
		})
	})
}

func (s *ClientIntegrationCoreSuite) TestClientIntegrationCoreLoginSuccess(t provider.T) {
	t.Title("LoginClient: Success")
	t.Tags("Client")
	if IsUnitTestsFailed() {
		t.Skip()
	}
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		_ = s.trm.WithinTransaction(s.ctx, func(txCtx context.Context) error {
			client, err := s.clientService.Login(txCtx, "1111111111", "111")
			sCtx.Assert().NoError(err)
			sCtx.Assert().Equal(client, &models.Client{
				ID:        1,
				Name:      "Diana",
				Telephone: "1111111111",
				Mail:      "d@mail.ru",
				Password:  "111",
			})

			s.trm.Rollback(txCtx, s.txResolver, s.dbx)
			return nil
		})
	})
}

func (s *ClientIntegrationCoreSuite) TestClientIntegrationCoreLoginFailure(t provider.T) {
	t.Title("LoginClient: Failure")
	t.Tags("Client")
	if IsUnitTestsFailed() {
		t.Skip()
	}
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		_ = s.trm.WithinTransaction(s.ctx, func(txCtx context.Context) error {
			_, err := s.clientService.Login(txCtx, "1111111111", "11111")
			sCtx.Assert().Error(err)

			s.trm.Rollback(txCtx, s.txResolver, s.dbx)
			return nil
		})
	})
}

func (s *ClientIntegrationCoreSuite) TestClientIntegrationCoreGetByTelephoneSuccess(t provider.T) {
	t.Title("GetByTelephoneClient: Success")
	t.Tags("Client")
	if IsUnitTestsFailed() {
		t.Skip()
	}
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		_ = s.trm.WithinTransaction(s.ctx, func(txCtx context.Context) error {
			client, err := s.clientService.GetByTelephone(txCtx, "1111111111")
			sCtx.Assert().NoError(err)
			sCtx.Assert().Equal(client, &models.Client{
				ID:        1,
				Name:      "Diana",
				Telephone: "1111111111",
				Mail:      "d@mail.ru",
				Password:  "111",
			})

			s.trm.Rollback(txCtx, s.txResolver, s.dbx)
			return nil
		})
	})
}

func (s *ClientIntegrationCoreSuite) TestClientIntegrationCoreGetByTelephoneFailure(t provider.T) {
	t.Title("GetByTelephoneClient: Failure")
	t.Tags("Client")
	if IsUnitTestsFailed() {
		t.Skip()
	}
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		_ = s.trm.WithinTransaction(s.ctx, func(txCtx context.Context) error {
			_, err := s.clientService.GetByTelephone(txCtx, "1111111110")
			sCtx.Assert().Error(err)

			s.trm.Rollback(txCtx, s.txResolver, s.dbx)
			return nil
		})
	})
}

func (s *ClientIntegrationCoreSuite) TestClientIntegrationCoreCreateAssignmentSuccess(t provider.T) {
	t.Title("ClientCreateAssignment: Success")
	t.Tags("Client")
	if IsUnitTestsFailed() {
		t.Skip()
	}
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		_ = s.trm.WithinTransaction(s.ctx, func(txCtx context.Context) error {

			err := s.clientService.CreateAssignment(txCtx, uint64(1), uint64(3))
			sCtx.Assert().NoError(err)

			s.trm.Rollback(txCtx, s.txResolver, s.dbx)
			return nil
		})
	})
}

func (s *ClientIntegrationCoreSuite) TestClientIntegrationCoreCreateAssignmentFailure(t provider.T) {
	t.Title("ClientCreateAssignment: Failure")
	t.Tags("Client")
	if IsUnitTestsFailed() {
		t.Skip()
	}
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		_ = s.trm.WithinTransaction(s.ctx, func(txCtx context.Context) error {

			err := s.clientService.CreateAssignment(txCtx, uint64(1), uint64(2))
			sCtx.Assert().Error(err)

			s.trm.Rollback(txCtx, s.txResolver, s.dbx)
			return nil
		})
	})
}

func (s *ClientIntegrationCoreSuite) TestClientIntegrationCoreDeleteAssignmentSuccess(t provider.T) {
	t.Title("ClientDeleteAssignment: Success")
	t.Tags("Client")
	if IsUnitTestsFailed() {
		t.Skip()
	}
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		_ = s.trm.WithinTransaction(s.ctx, func(txCtx context.Context) error {

			err := s.clientService.DeleteAssignment(txCtx, uint64(1), uint64(2))
			sCtx.Assert().NoError(err)

			s.trm.Rollback(txCtx, s.txResolver, s.dbx)
			return nil
		})
	})
}

func (s *ClientIntegrationCoreSuite) TestClientIntegrationCoreDeleteAssignmentFailure(t provider.T) {
	t.Title("ClientDeleteAssignment: Failure")
	t.Tags("Client")
	if IsUnitTestsFailed() {
		t.Skip()
	}
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		_ = s.trm.WithinTransaction(s.ctx, func(txCtx context.Context) error {

			err := s.clientService.DeleteAssignment(txCtx, uint64(1), uint64(1))
			sCtx.Assert().Error(err)

			s.trm.Rollback(txCtx, s.txResolver, s.dbx)
			return nil
		})
	})
}
