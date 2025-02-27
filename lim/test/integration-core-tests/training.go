package integrationCoreTests

import (
	"context"
	"database/sql"
	"time"
	"fmt"
	"os"

	transaction_manager "lim/internal/lim-core/managers"
	transaction_manager_implementation "lim/internal/lim-core/managers/implementation"
	"lim/internal/lim-core/services"
	service_implementation "lim/internal/lim-core/services/implementation"
	servicesDataBuilder "lim/internal/lim-core/services/implementation/data_builders"
	postgreSQL "lim/internal/lim-repo/postgreSQL"

	"lim/internal/lim-repo/flags"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	at_manager "github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/charmbracelet/log"
	"github.com/jmoiron/sqlx"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type TrainingIntegrationCoreSuite struct {
	suite.Suite
	db            *sql.DB
	dbx           *sqlx.DB
	txResolver    *trmsqlx.CtxGetter
	trm           transaction_manager.TransactionManager
	trainingService services.TrainingService
	ctx           context.Context
}

func (s *TrainingIntegrationCoreSuite) BeforeEach(t provider.T) {
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
	trainingRepository := postgreSQL.CreateTrainingPostgreSQLRepository(&fields)

	s.txResolver = trmsqlx.DefaultCtxGetter
	manager, _ := at_manager.New(trmsqlx.NewDefaultFactory(s.dbx))
	s.trm = transaction_manager_implementation.NewTransactionManagerImplementation(manager)

	logger := log.New(os.Stderr)
	clientRepository := postgreSQL.CreateClientPostgreSQLRepository(&fields)
	coachRepository := postgreSQL.CreateCoachPostgreSQLRepository(&fields)
	hallRepository := postgreSQL.CreateHallPostgreSQLRepository(&fields)

	s.trainingService = service_implementation.NewTrainingServiceImplementation(trainingRepository, clientRepository, coachRepository, hallRepository, s.trm, logger)
}

func (s *TrainingIntegrationCoreSuite) AfterEach(t provider.T) {
	s.db.Close()
}

func (s *TrainingIntegrationCoreSuite) TestTrainingIntegrationCoreCreateSuccess(t provider.T) {
	t.Title("CreateTraining: Success")
	t.Tags("Training")
	if IsUnitTestsFailed() {
		t.Skip()
	}
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		_ = s.trm.WithinTransaction(s.ctx, func(txCtx context.Context) error {
			err := s.trainingService.Create(txCtx, servicesDataBuilder.NewTrainingBuilder().Build())
			sCtx.Assert().NoError(err)

			s.trm.Rollback(txCtx, s.txResolver, s.dbx)
			return nil
		})
	})
}

func (s *TrainingIntegrationCoreSuite) TestTrainingIntegrationCoreCreateFailure(t provider.T) {
	t.Title("CreateTraining: Failure")
	t.Tags("Training")
	if IsUnitTestsFailed() {
		t.Skip()
	}
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		_ = s.trm.WithinTransaction(s.ctx, func(txCtx context.Context) error {
			err := s.trainingService.Create(txCtx, servicesDataBuilder.NewTrainingBuilder().WithCoachID(4).Build())
			sCtx.Assert().Error(err)

			s.trm.Rollback(txCtx, s.txResolver, s.dbx)
			return nil
		})
	})
}

func (s *TrainingIntegrationCoreSuite) TestTrainingIntegrationCoreGetAllByClientSuccess(t provider.T) {
	t.Title("GetAllByClientTraining: Success")
	t.Tags("Training")
	if IsUnitTestsFailed() {
		t.Skip()
	}
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		_ = s.trm.WithinTransaction(s.ctx, func(txCtx context.Context) error {
			trainings, err := s.trainingService.GetAllByClient(txCtx, uint64(2))
			sCtx.Assert().NoError(err)
			sCtx.Assert().Equal(len(trainings), 1)

			s.trm.Rollback(txCtx, s.txResolver, s.dbx)
			return nil
		})
	})
}

func (s *TrainingIntegrationCoreSuite) TestTrainingIntegrationCoreGetAllByCoachOnDateSuccess(t provider.T) {
	t.Title("GetAllByCoachOnDateTraining: Success")
	t.Tags("Training")
	if IsUnitTestsFailed() {
		t.Skip()
	}
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		_ = s.trm.WithinTransaction(s.ctx, func(txCtx context.Context) error {
			trainings, err := s.trainingService.GetAllByCoachOnDate(txCtx, uint64(1), time.Date(2024, 5, 5, 17, 0, 0, 0, time.UTC))
			sCtx.Assert().NoError(err)
			sCtx.Assert().Equal(2, len(trainings))

			s.trm.Rollback(txCtx, s.txResolver, s.dbx)
			return nil
		})
	})
}
