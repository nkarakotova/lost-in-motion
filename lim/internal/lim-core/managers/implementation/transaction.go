package transactionManager

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"lim/internal/lim-core/managers"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/jmoiron/sqlx"
)

type TransactionManagerImplementation struct {
	manager *manager.Manager
}

func NewTransactionManagerImplementation(manager *manager.Manager) managers.TransactionManager {
	return &TransactionManagerImplementation{manager: manager}
}

func (transactor *TransactionManagerImplementation) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return transactor.manager.Do(ctx, fn)
}

func (transactor *TransactionManagerImplementation) Rollback(ctx context.Context, txRes *trmsqlx.CtxGetter, databasex *sqlx.DB) {
	txRes.DefaultTrOrDB(ctx, databasex).QueryContext(ctx, "rollback error")
}