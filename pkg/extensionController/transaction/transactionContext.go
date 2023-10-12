package transaction

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/exasol/extension-manager/pkg/apiErrors"
	"github.com/exasol/extension-manager/pkg/extensionController/bfs"
)

type (
	TransactionStarter func(ctx context.Context, db *sql.DB) (*TransactionContext, error)
)

type (
	BucketFsClientCreator func(bucketFsBasePath string, ctx context.Context, db *sql.DB) (bfs.BucketFsAPI, error)
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
	return &TransactionContext{
		context:         ctx,
		db:              db,
		transaction:     tx,
		createBfsClient: bfs.CreateBucketFsAPI,
		bfsClient:       nil,
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

func (ctx *TransactionContext) GetTransaction() *sql.Tx {
	return ctx.transaction
}

func (ctx *TransactionContext) GetDBConnection() *sql.DB {
	return ctx.db
}

func (ctx *TransactionContext) GetContext() context.Context {
	return ctx.context
}

func (ctx *TransactionContext) GetBucketFsClient() (bfs.BucketFsAPI, error) {
	if ctx.bfsClient == nil {
		client, err := ctx.createBfsClient("", ctx.context, ctx.db)
		fmt.Printf("Created bucket FS Client %v\n", client)
		if err != nil {
			return nil, err
		}
		ctx.bfsClient = client
	}
	return ctx.bfsClient, nil
}

func (ctx *TransactionContext) Rollback() {
	_ = ctx.cleanup()
	// Even if Tx.Rollback fails, the transaction will no longer be valid, nor will it have been committed to the database.
	// See https://go.dev/doc/database/execute-transactions
	_ = ctx.transaction.Rollback()
}

func (ctx *TransactionContext) Commit() error {
	err := ctx.cleanup()
	if err != nil {
		return err
	}
	return ctx.transaction.Commit()
}

func (ctx *TransactionContext) cleanup() error {
	if ctx.bfsClient != nil {
		fmt.Printf("Closing bfs client %v\n", ctx.bfsClient)
		return ctx.bfsClient.Close()
	}
	fmt.Println("bfs client is nil, not closing")
	return nil
}
