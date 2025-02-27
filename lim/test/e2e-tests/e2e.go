package e2eTests

import (
	"context"
	"database/sql"
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

type E2ESuite struct {
	suite.Suite
	db              *sql.DB
	dbx             *sqlx.DB
	txResolver      *trmsqlx.CtxGetter
	trm             transaction_manager.TransactionManager
	clientService   services.ClientService
	trainingService services.TrainingService
	ctx             context.Context
}

func (s *E2ESuite) BeforeEach(t provider.T) {
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
	coachRepository := postgreSQL.CreateCoachPostgreSQLRepository(&fields)
	hallRepository := postgreSQL.CreateHallPostgreSQLRepository(&fields)

	s.clientService = service_implementation.NewClientServiceImplementation(clientRepository, trainingRepository, s.trm, logger)
	s.trainingService = service_implementation.NewTrainingServiceImplementation(trainingRepository, clientRepository, coachRepository, hallRepository, s.trm, logger)
}

func (s *E2ESuite) AfterEach(t provider.T) {
	s.db.Close()
}

func (s *E2ESuite) TestCreateClientSuccess(t provider.T) {
	t.Title("E2E")
	if IsTestsFailed() {
		t.Skip()
	}
	t.WithNewStep("E2E", func(sCtx provider.StepCtx) {
		_ = s.trm.WithinTransaction(s.ctx, func(txCtx context.Context) error {
			err := s.clientService.Create(txCtx, servicesDataBuilder.NewClientBuilder().Build())
			sCtx.Assert().NoError(err)

			client, err := s.clientService.Login(txCtx, "1234567890", "123")
			sCtx.Assert().NoError(err)
			sCtx.Assert().Equal(client, servicesDataBuilder.NewClientBuilder().WithID(4).Build())

			training, err := s.trainingService.GetByID(txCtx, uint64(1))
			sCtx.Assert().NoError(err)

			cur_places := training.PlacesNum

			err = s.clientService.CreateAssignment(txCtx, client.ID, uint64(1))
			sCtx.Assert().NoError(err)

			training, err = s.trainingService.GetByID(txCtx, uint64(1))
			sCtx.Assert().NoError(err)

			places_after_create_assignment := training.PlacesNum

			sCtx.Assert().Equal(places_after_create_assignment, cur_places-1)

			err = s.clientService.DeleteAssignment(txCtx, client.ID, uint64(1))
			sCtx.Assert().NoError(err)

			training, err = s.trainingService.GetByID(txCtx, uint64(1))
			sCtx.Assert().NoError(err)

			places_after_delete_assignment := training.PlacesNum

			sCtx.Assert().Equal(places_after_delete_assignment, cur_places)

			s.trm.Rollback(txCtx, s.txResolver, s.dbx)
			return nil
		})
	})
}
