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

type HallIntegrationCoreSuite struct {
	suite.Suite
	db            *sql.DB
	dbx           *sqlx.DB
	txResolver    *trmsqlx.CtxGetter
	trm           transaction_manager.TransactionManager
	hallService   services.HallService
	ctx           context.Context
}

func (s *HallIntegrationCoreSuite) BeforeEach(t provider.T) {
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
	hallRepository := postgreSQL.CreateHallPostgreSQLRepository(&fields)

	s.txResolver = trmsqlx.DefaultCtxGetter
	manager, _ := at_manager.New(trmsqlx.NewDefaultFactory(s.dbx))
	s.trm = transaction_manager_implementation.NewTransactionManagerImplementation(manager)

	logger := log.New(os.Stderr)
	trainingRepository := postgreSQL.CreateTrainingPostgreSQLRepository(&fields)

	s.hallService = service_implementation.NewHallServiceImplementation(hallRepository, trainingRepository, logger)
}

func (s *HallIntegrationCoreSuite) AfterEach(t provider.T) {
	s.db.Close()
}

func (s *HallIntegrationCoreSuite) TestHallIntegrationCoreCreateSuccess(t provider.T) {
	t.Title("CreateHall: Success")
	t.Tags("Hall")
	if IsUnitTestsFailed() {
		t.Skip()
	}
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		_ = s.trm.WithinTransaction(s.ctx, func(txCtx context.Context) error {
			err := s.hallService.Create(txCtx, &models.Hall{
				Number: 7,
			})
			sCtx.Assert().NoError(err)

			s.trm.Rollback(txCtx, s.txResolver, s.dbx)
			return nil
		})
	})
}

func (s *HallIntegrationCoreSuite) TestCoachIntegrationCoreCreateFailure(t provider.T) {
	t.Title("CreateHall: Failure")
	t.Tags("Hall")
	if IsUnitTestsFailed() {
		t.Skip()
	}
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		_ = s.trm.WithinTransaction(s.ctx, func(txCtx context.Context) error {
			err := s.hallService.Create(txCtx, servicesDataBuilder.NewHallBuilder().Build())
			sCtx.Assert().ErrorIs(err, servicesErrors.HallAlreadyExists)

			s.trm.Rollback(txCtx, s.txResolver, s.dbx)
			return nil
		})
	})
}