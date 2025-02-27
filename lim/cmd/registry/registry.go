package registry

import (
	"fmt"
	"lim/cmd/config"
	transaction_manager "lim/internal/lim-core/managers"
	transaction_manager_implementation "lim/internal/lim-core/managers/implementation"
	"lim/internal/lim-core/repositories"
	"lim/internal/lim-core/services"
	servicesImplementation "lim/internal/lim-core/services/implementation"
	postgres_config "lim/internal/lim-repo/config"
	postgres_repo "lim/internal/lim-repo/postgreSQL"
	"os"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	at_manager "github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/charmbracelet/log"
	"github.com/jmoiron/sqlx"
)

const (
	DefaultLogFile  = "fatal.log"
	DefaultLogLevel = log.FatalLevel
)

type AppServiceFields struct {
	ClientService   services.ClientService
	CoachService    services.CoachService
	HallService     services.HallService
	TrainingService services.TrainingService
}

type AppRepositoryFields struct {
	ClientRepository   repositories.ClientRepository
	CoachRepository    repositories.CoachRepository
	HallRepository     repositories.HallRepository
	TrainingRepository repositories.TrainingRepository
}

type AppManagerFields struct {
	TransactionManager transaction_manager.TransactionManager
}

type App struct {
	Config       config.Config
	Repositories *AppRepositoryFields
	Managers     *AppManagerFields
	Services     *AppServiceFields
	Logger       *log.Logger
}

func (a *App) initRepositories(fields *postgres_repo.PostgresRepositoryFields) *AppRepositoryFields {
	f := &AppRepositoryFields{
		ClientRepository:   postgres_repo.CreateClientPostgreSQLRepository(fields),
		CoachRepository:    postgres_repo.CreateCoachPostgreSQLRepository(fields),
		HallRepository:     postgres_repo.CreateHallPostgreSQLRepository(fields),
		TrainingRepository: postgres_repo.CreateTrainingPostgreSQLRepository(fields),
	}

	a.Logger.Info("Success initialization of repositories")
	return f
}

func (a *App) initManagers(manager *at_manager.Manager) *AppManagerFields {
	f := &AppManagerFields{
		TransactionManager: transaction_manager_implementation.NewTransactionManagerImplementation(manager),
	}

	a.Logger.Info("Success initialization of repositories")
	return f
}

func (a *App) initServices(r *AppRepositoryFields, m *AppManagerFields) *AppServiceFields {
	f := &AppServiceFields{
		ClientService:   servicesImplementation.NewClientServiceImplementation(r.ClientRepository, r.TrainingRepository, m.TransactionManager, a.Logger),
		CoachService:    servicesImplementation.NewCoachServiceImplementation(r.CoachRepository, r.TrainingRepository, a.Logger),
		HallService:     servicesImplementation.NewHallServiceImplementation(r.HallRepository, r.TrainingRepository, a.Logger),
		TrainingService: servicesImplementation.NewTrainingServiceImplementation(r.TrainingRepository, r.ClientRepository, r.CoachRepository, r.HallRepository, m.TransactionManager, a.Logger),
	}

	a.Logger.Info("Success initialization of services")
	return f
}

func (a *App) initLogger() {
	f, err_not_file := os.OpenFile(a.Config.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err_not_file != nil {
		var err error
		f, err = os.OpenFile(DefaultLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(fmt.Errorf("%s", err.Error()))
		}
	}

	Logger := log.New(f)

	log.SetFormatter(log.LogfmtFormatter)
	Logger.SetReportTimestamp(true)
	Logger.SetReportCaller(true)

	flag_not_level := false
	if a.Config.LogLevel == "debug" {
		Logger.SetLevel(log.DebugLevel)
	} else if a.Config.LogLevel == "info" {
		Logger.SetLevel(log.InfoLevel)
	} else if a.Config.LogLevel == "warn" {
		Logger.SetLevel(log.WarnLevel)
	} else if a.Config.LogLevel == "fatal" {
		Logger.SetLevel(log.FatalLevel)
	} else {
		flag_not_level = true
		Logger.SetLevel(DefaultLogLevel)
	}

	Logger.Print("\n")
	Logger.Info("Success initialization of new Logger!")

	a.Logger = Logger

	if err_not_file != nil {
		a.Logger.Warn(fmt.Errorf("can't set log file from config with err: %s, setting default log file %s", err_not_file.Error(), DefaultLogFile))
	}
	if flag_not_level {
		a.Logger.Warn("Error log level, set default: DEBU")
	}
}

func (a *App) Init() error {
	a.initLogger()

	pgConfig := postgres_config.Config(a.Config)

	fields, err := postgres_repo.CreatePostgresRepositoryFields(pgConfig.Postgres, a.Logger)
	if err != nil {
		log.Print(fmt.Errorf("%s", err.Error()))
		a.Logger.Fatal("Error create postgres repository fields", "err", fmt.Errorf("%s", err.Error()))
		return err
	}

	dbx := sqlx.NewDb(fields.DB, "pgx")
	manager, err := at_manager.New(trmsqlx.NewDefaultFactory(dbx))
	if err != nil {
		return err
	}

	a.Repositories = a.initRepositories(fields)
	a.Managers = a.initManagers(manager)
	a.Services = a.initServices(a.Repositories, a.Managers)

	services.FirstTrainingTime = a.Config.FirstTrainingTime
	services.LastTrainingTime = a.Config.LastTrainingTime

	return nil
}

func (a *App) Run() error {
	err := a.Init()

	if err != nil {
		a.Logger.Error("Error init app", "err", fmt.Errorf("%s", err.Error()))
		return err
	}

	return nil
}
