// +build unit

package postgreSQL

import (
	"context"
	"database/sql"
	"testing"

	"lim/internal/lim-core/models"
	"lim/internal/lim-core/repositories"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"

	"lim/internal/lim-core/errors/repositoriesErrors"
	postgreSQLObjectMother "lim/internal/lim-repo/postgreSQL/object_mothers"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type HallRepoSuite struct {
	suite.Suite
	db         *sql.DB
	mock       sqlmock.Sqlmock
	repository repositories.HallRepository
	ctx        context.Context
}

func (s *HallRepoSuite) BeforeEach(t provider.T) {
	var err error
	s.db, s.mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	dbx := sqlx.NewDb(s.db, "pgx")
	s.repository = NewHallPostgreSQLRepository(dbx)
	s.ctx = context.Background()
}

func (s *HallRepoSuite) AfterEach(t provider.T) {
	s.db.Close()
}

func (s *HallRepoSuite) TestHallMockGetByNumberSuccess(t provider.T) {
	t.Title("GetHallByNumber: Success")
	t.Tags("Hall")
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`select * from halls where number=$1;`).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"hall_id", "number"}).
				AddRow(1, 1))

		new_hall := postgreSQLObjectMother.CreateTestHall()
		hall, err := s.repository.GetByNumber(s.ctx, 1)

		sCtx.Assert().NoError(err)
		sCtx.Assert().NotNil(hall)
		sCtx.Assert().Equal(new_hall, hall)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *HallRepoSuite) TestHallMockGetByNumberFailure(t provider.T) {
	t.Title("GetHallByNumber: Failure")
	t.Tags("Hall")
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`select * from halls where number=$1;`).WithArgs(1).WillReturnError(sql.ErrNoRows)

		_, err := s.repository.GetByNumber(s.ctx, 1)

		sCtx.Assert().ErrorIs(err, repositoriesErrors.EntityDoesNotExists)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *HallRepoSuite) TestHallMockCreateSuccess(t provider.T) {
	t.Title("CreateHall: Success")
	t.Tags("Hall")
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`insert into halls(number) values($1) returning hall_id;`).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"hall_id"}).AddRow(1))

		hall := postgreSQLObjectMother.CreateTestHall()
		err := s.repository.Create(s.ctx, hall)

		sCtx.Assert().NoError(err)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *HallRepoSuite) TestHallMockCreateFailure(t provider.T) {
	t.Title("CreateHall: Failure")
	t.Tags("Hall")
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`insert into halls(number) values($1) returning hall_id;`).
			WithArgs(1)

		hall := postgreSQLObjectMother.CreateTestHall()
		err := s.repository.Create(s.ctx, hall)

		sCtx.Assert().Error(err)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *HallRepoSuite) TestHallMockGetByIDSuccess(t provider.T) {
	t.Title("GetHallByID: Success")
	t.Tags("Hall")
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`select * from halls where hall_id=$1;`).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"hall_id", "number"}).
				AddRow(1, 1))

		new_hall := postgreSQLObjectMother.CreateTestHall()
		hall, err := s.repository.GetByID(s.ctx, 1)

		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(new_hall, hall)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *HallRepoSuite) TestHallMockGetByIDFailure(t provider.T) {
	t.Title("GetHallByID: Failure")
	t.Tags("Hall")
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`select * from halls where hall_id=$1;`).WithArgs(1).WillReturnError(sql.ErrNoRows)

		_, err := s.repository.GetByID(s.ctx, 1)

		sCtx.Assert().ErrorIs(err, repositoriesErrors.EntityDoesNotExists)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *HallRepoSuite) TestHallMockGetAllSuccess(t provider.T) {
	t.Title("HallMockGetAll: Success")
	t.Tags("Hall")
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`select * from halls;`).
			WillReturnRows(sqlmock.NewRows([]string{"hall_id", "number"}).
				AddRow(1, 1))

		new_halls := map[uint64]models.Hall{1: *postgreSQLObjectMother.CreateTestHall()}
		hall, err := s.repository.GetAll(s.ctx)

		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(new_halls, hall)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *HallRepoSuite) TestHallMockGetAllFailure(t provider.T) {
	t.Title("HallMockGetAll: Failure")
	t.Tags("Hall")
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`select * from halls;`).WillReturnError(sql.ErrNoRows)

		_, err := s.repository.GetAll(s.ctx)

		sCtx.Assert().ErrorIs(err, repositoriesErrors.EntityDoesNotExists)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestHallSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(HallRepoSuite))
}
