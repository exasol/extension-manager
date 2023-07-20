package context

import (
	"context"
	"database/sql"

	"github.com/exasol/extension-manager/pkg/backend"
)

func CreateContext(ctx context.Context, extensionSchemaName string, tx *sql.Tx, db *sql.DB) *ExtensionContext {
	var client SimpleSQLClient = backend.NewSqlClient(ctx, tx)
	return CreateContextWithClient(extensionSchemaName, client)
}

func CreateContextWithClient(extensionSchemaName string, client SimpleSQLClient) *ExtensionContext {
	return &ExtensionContext{
		ExtensionSchemaName: extensionSchemaName,
		SqlClient:           client,
		BucketFs:            &bucketFsContextImpl{},
	}
}

// Instances of type ExtensionContext are passed to an extension so that extension can
//   - retrieve context information like the extension schema name (field ExtensionSchemaName)
//   - execute SQL queries against the database using a [SqlClient]
//   - or resolve files in BucketFS using [BucketFs]
type ExtensionContext struct {
	ExtensionSchemaName string          `json:"extensionSchemaName"` // Name of the schema where EM creates all database objects (e.g. scripts or virtual schemas)
	SqlClient           SimpleSQLClient `json:"sqlClient"`           // Allows extensions to execute SQL queries and statements
	BucketFs            BucketFsContext `json:"bucketFs"`            // Allows extensions to interact with BucketFS
}

// BucketFsContext allows extensions to interact with BucketFS.
type BucketFsContext interface {
	// ResolvePath returns an absolute path for the given filename in BucketFS.
	ResolvePath(fileName string) string
}

type bucketFsContextImpl struct{}

func (b *bucketFsContextImpl) ResolvePath(fileName string) string {
	return "/buckets/bfsdefault/default/" + fileName
}

// Extensions use this SQL client to execute queries.
type SimpleSQLClient interface {
	// Execute runs a query that does not return rows, e.g. INSERT or UPDATE.
	Execute(query string, args ...any)

	// Query runs a query that returns rows, typically a SELECT.
	Query(query string, args ...any) backend.QueryResult
}
