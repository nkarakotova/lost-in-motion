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

	"lim/internal/lim-repo/flags"
	"lim/internal/lim-core/models"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	at_manager "github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/charmbracelet/log"
	"github.com/jmoiron/sqlx"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type CoachIntegrationCoreSuite struct {
	suite.Suite
	db            *sql.DB
	dbx           *sqlx.DB
	txResolver    *trmsqlx.CtxGetter
	trm           transaction_manager.TransactionManager
	coachService services.CoachService
	ctx           context.Context
}

func (s *CoachIntegrationCoreSuite) BeforeEach(t provider.T) {
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
	coachRepository := postgreSQL.CreateCoachPostgreSQLRepository(&fields)

	s.txResolver = trmsqlx.DefaultCtxGetter
	manager, _ := at_manager.New(trmsqlx.NewDefaultFactory(s.dbx))
	s.trm = transaction_manager_implementation.NewTransactionManagerImplementation(manager)

	logger := log.New(os.Stderr)
	trainingRepository := postgreSQL.CreateTrainingPostgreSQLRepository(&fields)

	s.coachService = service_implementation.NewCoachServiceImplementation(coachRepository, trainingRepository, logger)
}

func (s *CoachIntegrationCoreSuite) AfterEach(t provider.T) {
	s.db.Close()
}

func (s *CoachIntegrationCoreSuite) TestCoachIntegrationCoreCreateSuccess(t provider.T) {
	t.Title("CreateCoach: Success")
	t.Tags("Coach")
	if IsUnitTestsFailed() {
		t.Skip()
	}
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		_ = s.trm.WithinTransaction(s.ctx, func(txCtx context.Context) error {
			err := s.coachService.Create(txCtx, servicesDataBuilder.NewCoachBuilder().Build())
			sCtx.Assert().NoError(err)

			s.trm.Rollback(txCtx, s.txResolver, s.dbx)
			return nil
		})
	})
}

func (s *CoachIntegrationCoreSuite) TestCoachIntegrationCoreCreateFailure(t provider.T) {
	t.Title("CreateCoach: Failure")
	t.Tags("Coach")
	if IsUnitTestsFailed() {
		t.Skip()
	}
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		_ = s.trm.WithinTransaction(s.ctx, func(txCtx context.Context) error {
			err := s.coachService.Create(txCtx, &models.Coach{
				Name: "Kate",
			})
			sCtx.Assert().ErrorIs(err, servicesErrors.CoachAlreadyExists)

			s.trm.Rollback(txCtx, s.txResolver, s.dbx)
			return nil
		})
	})
}
