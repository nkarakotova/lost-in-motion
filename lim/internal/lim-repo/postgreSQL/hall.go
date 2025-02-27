package postgreSQL

import (
	"context"
	"database/sql"

	"lim/internal/lim-core/repositories"

	"lim/internal/lim-core/errors/repositoriesErrors"

	"lim/internal/lim-core/models"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/jinzhu/copier"
	"github.com/jmoiron/sqlx"
)

type HallPostgreSQL struct {
	ID     uint64 `db:"hall_id"`
	Number uint64 `db:"number"`
}

type HallPostgreSQLRepository struct {
	db         *sqlx.DB
	txResolver *trmsqlx.CtxGetter
}

func NewHallPostgreSQLRepository(db *sqlx.DB) repositories.HallRepository {
	return &HallPostgreSQLRepository{db: db, txResolver: trmsqlx.DefaultCtxGetter}
}

func (h *HallPostgreSQLRepository) Create(ctx context.Context, hall *models.Hall) error {
	query := `insert into halls(number) values($1) returning hall_id;`

	err := h.txResolver.DefaultTrOrDB(ctx, h.db).QueryRowxContext(ctx, query, hall.Number).Scan(&hall.ID)
	if err != nil {
		return err
	}

	return nil
}

func (h *HallPostgreSQLRepository) GetByID(ctx context.Context, id uint64) (*models.Hall, error) {
	query := `select * from halls where hall_id=$1;`

	hallDB := &HallPostgreSQL{}
	err := h.txResolver.DefaultTrOrDB(ctx, h.db).GetContext(ctx, hallDB, query, id)
	if err == sql.ErrNoRows {
		return nil, repositoriesErrors.EntityDoesNotExists
	} else if err != nil {
		return nil, err
	}

	hallModels := &models.Hall{}
	err = copier.Copy(hallModels, hallDB)
	if err != nil {
		return nil, err
	}

	return hallModels, nil
}

func (h *HallPostgreSQLRepository) GetByNumber(ctx context.Context, number uint64) (*models.Hall, error) {
	query := `select * from halls where number=$1;`

	hallDB := &HallPostgreSQL{}
	err := h.txResolver.DefaultTrOrDB(ctx, h.db).GetContext(ctx, hallDB, query, number)
	if err == sql.ErrNoRows {
		return nil, repositoriesErrors.EntityDoesNotExists
	} else if err != nil {
		return nil, err
	}

	hallModels := &models.Hall{}
	err = copier.Copy(hallModels, hallDB)
	if err != nil {
		return nil, err
	}

	return hallModels, nil
}

func (h *HallPostgreSQLRepository) GetAll(ctx context.Context) (map[uint64]models.Hall, error) {
	query := `select * from halls;`

	hallDB := []HallPostgreSQL{}
	err := h.txResolver.DefaultTrOrDB(ctx, h.db).SelectContext(ctx, &hallDB, query)
	if err == sql.ErrNoRows {
		return nil, repositoriesErrors.EntityDoesNotExists
	} else if err != nil {
		return nil, err
	}

	hallModels := make(map[uint64]models.Hall)
	for i := range hallDB {
		hall := models.Hall{}
		err = copier.Copy(&hall, &hallDB[i])
		if err != nil {
			return nil, err
		}

		hallModels[hall.ID] = hall
	}

	return hallModels, nil
}
