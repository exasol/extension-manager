package transaction

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/exasol/extension-manager/pkg/extensionController/bfs"
)

type TransactionStarterMock struct {
	transactionStarter TransactionStarter
	dbMock             *sql.DB
	bfsMock            *bfs.BucketFsMock
}

func CreateTransactionStarterMock(dbMock *sql.DB, bfsMock *bfs.BucketFsMock) *TransactionStarterMock {
	mock := &TransactionStarterMock{
		dbMock:             dbMock,
		bfsMock:            bfsMock,
		transactionStarter: nil,
	}
	mock.SimulateMockTransaction()
	return mock
}

func (m *TransactionStarterMock) GetTransactionStarter() TransactionStarter {
	return m.transactionStarter
}

func (m *TransactionStarterMock) SimulateMockTransaction() {
	m.transactionStarter = func(ctx context.Context, db *sql.DB) (*TransactionContext, error) {
		tx, err := m.dbMock.Begin()
		if err != nil {
			return nil, fmt.Errorf("failed to start mock transaction: %w", err)
		}
		return &TransactionContext{
			context:     context.Background(),
			transaction: tx,
			db:          m.dbMock,
			createBfsClient: func(bucketFsBasePath string, ctx context.Context, db *sql.DB) (bfs.BucketFsAPI, error) {
				return m.bfsMock, nil
			},
			bfsClient: nil,
		}, nil
	}
}
