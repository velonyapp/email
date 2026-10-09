package mysql

import (
	"context"
	"database/sql"

	"github.com/velonyapp/email/internal/application/port"
)

var _ port.UnitOfWork = (*unitOfWork)(nil)

type unitOfWork struct {
	db *sql.DB
}

func NewUnitOfWork(db *sql.DB) port.UnitOfWork {
	return &unitOfWork{db: db}
}

func (uow *unitOfWork) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	if _, ok := txFromContext(ctx); ok {
		return fn(ctx)
	}

	tx, err := uow.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	txCtx := withTx(ctx, tx)

	if err := fn(txCtx); err != nil {
		return err
	}

	return tx.Commit()
}
