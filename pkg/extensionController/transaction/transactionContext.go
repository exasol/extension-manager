package transaction

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/exasol/extension-manager/pkg/apiErrors"
)

// BeginTransaction starts a new database transaction.
func BeginTransaction(ctx context.Context, db *sql.DB) (*TransactionContext, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		if strings.Contains(err.Error(), "Connection exception - authentication failed") {
			return nil, apiErrors.NewUnauthorizedErrorF("invalid database credentials")
		}
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return &TransactionContext{context: ctx, db: db, transaction: tx}, nil
}

// TransactionContext contains the state of a running database transaction.
type TransactionContext struct {
	context     context.Context
	db          *sql.DB
	transaction *sql.Tx
}

func (ctx *TransactionContext) GetTransaction() *sql.Tx {
	return ctx.transaction
}

func (ctx *TransactionContext) GetDBConnection() *sql.DB {
	return ctx.db
}

func (ctx *TransactionContext) GetContext() context.Context {
	return ctx.context
}

func (ctx *TransactionContext) Rollback() {
	// Even if Tx.Rollback fails, the transaction will no longer be valid, nor will it have been committed to the database.
	// See https://go.dev/doc/database/execute-transactions
	_ = ctx.transaction.Rollback()
}

func (ctx *TransactionContext) Commit() error {
	return ctx.transaction.Commit()
}
