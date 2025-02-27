package postgreSQL

import (
	"context"
	"database/sql"

	"lim/internal/lim-core/repositories"

	"lim/internal/lim-core/errors/repositoriesErrors"

	"github.com/jinzhu/copier"
	"github.com/jmoiron/sqlx"
	"lim/internal/lim-core/models"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
)

type CoachPostgreSQL struct {
	ID   uint64 `db:"coach_id"`
	Name string `db:"name"`
}

type CoahcPostgreSQLRepository struct {
	db         *sqlx.DB
	txResolver *trmsqlx.CtxGetter
}

func NewCoachPostgreSQLRepository(db *sqlx.DB) repositories.CoachRepository {
	return &CoahcPostgreSQLRepository{db: db, txResolver: trmsqlx.DefaultCtxGetter}
}

func (c *CoahcPostgreSQLRepository) Create(ctx context.Context, coach *models.Coach) error {
	query := `insert into coaches(name) values($1) returning coach_id;`

	err := c.txResolver.DefaultTrOrDB(ctx, c.db).QueryRowxContext(ctx, query, coach.Name).Scan(&coach.ID)
	if err != nil {
		return err
	}

	return nil
}

func (c *CoahcPostgreSQLRepository) GetByID(ctx context.Context, id uint64) (*models.Coach, error) {
	query := `select * from coaches where coach_id = $1;`

	coachDB := &CoachPostgreSQL{}
	err := c.txResolver.DefaultTrOrDB(ctx, c.db).GetContext(ctx, coachDB, query, id)
	if err == sql.ErrNoRows {
		return nil, repositoriesErrors.EntityDoesNotExists
	} else if err != nil {
		return nil, err
	}

	coachModels := &models.Coach{}
	err = copier.Copy(coachModels, coachDB)
	if err != nil {
		return nil, err
	}

	return coachModels, nil
}

func (c *CoahcPostgreSQLRepository) GetByName(ctx context.Context, name string) (*models.Coach, error) {
	query := `select * from coaches where name = $1;`

	coachDB := &CoachPostgreSQL{}
	err := c.txResolver.DefaultTrOrDB(ctx, c.db).GetContext(ctx, coachDB, query, name)
	if err == sql.ErrNoRows {
		return nil, repositoriesErrors.EntityDoesNotExists
	} else if err != nil {
		return nil, err
	}

	coachModels := &models.Coach{}
	err = copier.Copy(coachModels, coachDB)
	if err != nil {
		return nil, err
	}

	return coachModels, nil
}

func (c *CoahcPostgreSQLRepository) GetAll(ctx context.Context) ([]models.Coach, error) {
	query := `select * from coaches;`

	coachDB := []CoachPostgreSQL{}
	err := c.txResolver.DefaultTrOrDB(ctx, c.db).SelectContext(ctx, &coachDB, query)
	if err == sql.ErrNoRows {
		return nil, repositoriesErrors.EntityDoesNotExists
	} else if err != nil {
		return nil, err
	}

	coachModels := []models.Coach{}
	for i := range coachDB {
		coach := models.Coach{}
		err = copier.Copy(&coach, &coachDB[i])
		if err != nil {
			return nil, err
		}

		coachModels = append(coachModels, coach)
	}

	return coachModels, nil
}
