package context

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/exasol/extension-manager/pkg/backend"
	"github.com/exasol/extension-manager/pkg/extensionController/bfs"
	"github.com/exasol/extension-manager/pkg/extensionController/transaction"
)

func CreateContext(txCtx *transaction.TransactionContext, extensionSchemaName string, bucketFsBasePath string) *ExtensionContext {
	var sqlClient backend.SimpleSQLClient = backend.NewSqlClient(txCtx.GetContext(), txCtx.GetTransaction())
	var bucketFsClient bfs.BucketFsAPI = bfs.CreateBucketFsAPI(bucketFsBasePath)
	return CreateContextWithClient(extensionSchemaName, txCtx, sqlClient, bucketFsClient)
}

func CreateContextWithClient(extensionSchemaName string, txCtx *transaction.TransactionContext, client backend.SimpleSQLClient, bucketFsClient bfs.BucketFsAPI) *ExtensionContext {
	return &ExtensionContext{
		ExtensionSchemaName: extensionSchemaName,
		SqlClient:           &contextSqlClient{client},
		BucketFs: &bucketFsContextImpl{
			bucketFsClient: bucketFsClient,
			context:        txCtx.GetContext(),
			db:             txCtx.GetDBConnection(),
		},
	}
}

// Instances of type ExtensionContext are passed to an extension so that extension can
//   - retrieve context information like the extension schema name (field ExtensionSchemaName)
//   - execute SQL queries against the database using a [SqlClient]
//   - or resolve files in BucketFS using [BucketFs]
type ExtensionContext struct {
	ExtensionSchemaName string           `json:"extensionSchemaName"` // Name of the schema where EM creates all database objects (e.g. scripts or virtual schemas)
	SqlClient           ContextSqlClient `json:"sqlClient"`           // Allows extensions to execute SQL queries and statements
	BucketFs            BucketFsContext  `json:"bucketFs"`            // Allows extensions to interact with BucketFS
}

type ContextSqlClient interface {
	// Execute runs a query that does not return rows, e.g. INSERT or UPDATE.
	Execute(query string, args ...any)

	// Query runs a query that returns rows, typically a SELECT.
	Query(query string, args ...any) backend.QueryResult
}

type contextSqlClient struct {
	client backend.SimpleSQLClient
}

func (c *contextSqlClient) Execute(query string, args ...any) {
	_, err := c.client.Execute(query, args...)
	if err != nil {
		reportError(err)
	}
}

func (c *contextSqlClient) Query(query string, args ...any) backend.QueryResult {
	result, err := c.client.Query(query, args...)
	if err != nil {
		reportError(err)
	}
	return *result
}

// BucketFsContext allows extensions to interact with BucketFS.
/* [impl -> dsn~extension-context-bucketfs~1]. */
type BucketFsContext interface {
	// ResolvePath returns an absolute path for the given filename in BucketFS.
	ResolvePath(fileName string) string
}

type bucketFsContextImpl struct {
	bucketFsClient bfs.BucketFsAPI
	context        context.Context
	db             *sql.DB
}

/* [impl -> dsn~resolving-files-in-bucketfs~1]. */
func (b *bucketFsContextImpl) ResolvePath(fileName string) string {
	path, err := b.bucketFsClient.FindAbsolutePath(b.context, b.db, fileName)
	if err != nil {
		reportError(fmt.Errorf("failed to find absolute path for file %q: %w", fileName, err))
	}
	return path
}

func reportError(err error) {
	// Context functions are called by JavaScript code. The only way to report a failure is to panic.
	// The JS runtime will convert this panic into a thrown JS error.
	panic(err)
}
