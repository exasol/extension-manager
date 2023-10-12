package transaction

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/exasol/extension-manager/pkg/apiErrors"
	"github.com/exasol/extension-manager/pkg/extensionController/bfs"
)

// TransactionStarter starts a database transaction and returns a new [TransactionContext].
// It allows injecting a mock transaction in unit tests.
type (
	TransactionStarter func(ctx context.Context, db *sql.DB, bucketFsBasePath string) (*TransactionContext, error)
)

// BucketFsClientCreator creates a new [bfs.BucketFsAPI].
// It allows injecting a mock BucketFS client in unit tests.
type (
	BucketFsClientCreator func() (bfs.BucketFsAPI, error)
)

// BeginTransaction starts a new database transaction.
func BeginTransaction(ctx context.Context, db *sql.DB, bucketFsBasePath string) (*TransactionContext, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		if strings.Contains(err.Error(), "Connection exception - authentication failed") {
			return nil, apiErrors.NewUnauthorizedErrorF("invalid database credentials")
		}
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return &TransactionContext{
		context:     ctx,
		db:          db,
		transaction: tx,
		bfsClient:   nil,
		createBfsClient: func() (bfs.BucketFsAPI, error) {
			return bfs.CreateBucketFsAPI(bucketFsBasePath, ctx, db)
		},
	}, nil
}

// TransactionContext contains the state of a running database transaction.
type TransactionContext struct {
	context         context.Context
	db              *sql.DB
	transaction     *sql.Tx
	createBfsClient BucketFsClientCreator
	bfsClient       bfs.BucketFsAPI
}

// GetTransaction returns the current database transaction.
func (ctx *TransactionContext) GetTransaction() *sql.Tx {
	return ctx.transaction
}

// GetDBConnection returns the current database connection.
func (ctx *TransactionContext) GetDBConnection() *sql.DB {
	return ctx.db
}

// GetContext returns the current [context.Context].
func (ctx *TransactionContext) GetContext() context.Context {
	return ctx.context
}

// GetBucketFsClient returns a [bfs.BucketFsAPI].
// This creates a new client if none exists yet or returns the existing one.
func (ctx *TransactionContext) GetBucketFsClient() (bfs.BucketFsAPI, error) {
	if ctx.bfsClient == nil {
		client, err := ctx.createBfsClient()
		if err != nil {
			return nil, err
		}
		ctx.bfsClient = client
	}
	return ctx.bfsClient, nil
}

// Rollback rolls back the transaction and cleans up any resources like the [bfs.BucketFsAPI] if one was created.
func (ctx *TransactionContext) Rollback() {
	_ = ctx.cleanup()
	// Even if Tx.Rollback fails, the transaction will no longer be valid, nor will it have been committed to the database.
	// See https://go.dev/doc/database/execute-transactions
	_ = ctx.transaction.Rollback()
}

// Commit commits the transaction and cleans up any resources like the [bfs.BucketFsAPI] if one was created.
func (ctx *TransactionContext) Commit() error {
	err := ctx.cleanup()
	if err != nil {
		return err
	}
	return ctx.transaction.Commit()
}

func (ctx *TransactionContext) cleanup() error {
	if ctx.bfsClient != nil {
		return ctx.bfsClient.Close()
	}
	return nil
}
