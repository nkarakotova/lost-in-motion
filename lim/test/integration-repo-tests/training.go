package integrationRepoTests

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"lim/internal/lim-core/repositories"
	"lim/internal/lim-repo/flags"
	postgreSQL "lim/internal/lim-repo/postgreSQL"
	postgreSQLObjectMother "lim/internal/lim-repo/postgreSQL/object_mothers"

	transaction_manager "lim/internal/lim-core/managers"
	transaction_manager_implementation "lim/internal/lim-core/managers/implementation"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	at_manager "github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/charmbracelet/log"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type TrainingIntegrationRepoSuite struct {
	suite.Suite
	db                 *sql.DB
	trainingRepository repositories.TrainingRepository
	trm                transaction_manager.TransactionManager
	ctx                context.Context
	dbx                *sqlx.DB
	txResolver         *trmsqlx.CtxGetter
}

func (s *TrainingIntegrationRepoSuite) BeforeEach(t provider.T) {
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
	s.trainingRepository = postgreSQL.CreateTrainingPostgreSQLRepository(&fields)

	s.txResolver = trmsqlx.DefaultCtxGetter
	manager, _ := at_manager.New(trmsqlx.NewDefaultFactory(s.dbx))
	s.trm = transaction_manager_implementation.NewTransactionManagerImplementation(manager)
}

func (s *TrainingIntegrationRepoSuite) AfterEach(t provider.T) {
	s.db.Close()
}

func (s *TrainingIntegrationRepoSuite) TestTrainingIntegrationRepoCreateSuccess(t provider.T) {
	t.Title("TrainingCreate: Success")
	t.Tags("Training")
	if IsUnitTestsFailed() {
		t.Skip()
	}
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		_ = s.trm.WithinTransaction(s.ctx, func(txCtx context.Context) error {

			training := postgreSQLObjectMother.CreateTestTraining()
			err := s.trainingRepository.Create(txCtx, training)
			sCtx.Assert().NoError(err)

			s.trm.Rollback(txCtx, s.txResolver, s.dbx)
			return nil
		})
	})
}
