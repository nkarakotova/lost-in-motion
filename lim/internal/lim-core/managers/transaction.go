package managers

import (
	"context"
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/jmoiron/sqlx"
)

//go:generate mockgen  -source=transaction.go -destination=mocks/transaction.go
type TransactionManager interface {
	WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error
	Rollback(ctx context.Context, txRes *trmsqlx.CtxGetter, databasex *sqlx.DB)
}