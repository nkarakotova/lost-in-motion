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

type HallIntegrationRepoSuite struct {
	suite.Suite
	db                 *sql.DB
	hallRepository     repositories.HallRepository
	trm                transaction_manager.TransactionManager
	ctx                context.Context
	dbx                *sqlx.DB
	txResolver         *trmsqlx.CtxGetter
}

func (s *HallIntegrationRepoSuite) BeforeEach(t provider.T) {
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
	s.hallRepository = postgreSQL.CreateHallPostgreSQLRepository(&fields)

	s.txResolver = trmsqlx.DefaultCtxGetter
	manager, _ := at_manager.New(trmsqlx.NewDefaultFactory(s.dbx))
	s.trm = transaction_manager_implementation.NewTransactionManagerImplementation(manager)
}

func (s *HallIntegrationRepoSuite) AfterEach(t provider.T) {
	s.db.Close()
}

func (s *HallIntegrationRepoSuite) TestHallIntegrationRepoCreateSuccess(t provider.T) {
	t.Title("HallCreate: Success")
	t.Tags("Hall")
	if IsUnitTestsFailed() {
		t.Skip()
	}
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		_ = s.trm.WithinTransaction(s.ctx, func(txCtx context.Context) error {

			hall := postgreSQLObjectMother.CreateTestHall()
			hall.Number = 7
			err := s.hallRepository.Create(txCtx, hall)
			sCtx.Assert().NoError(err)

			s.trm.Rollback(txCtx, s.txResolver, s.dbx)
			return nil
		})
	})
}