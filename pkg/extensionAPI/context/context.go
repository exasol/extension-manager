package context

import (
	"github.com/exasol/extension-manager/pkg/backend"
	"github.com/exasol/extension-manager/pkg/extensionAPI/exaMetadata"
	"github.com/exasol/extension-manager/pkg/extensionController/bfs"
	"github.com/exasol/extension-manager/pkg/extensionController/transaction"
)

func CreateContext(txCtx *transaction.TransactionContext, extensionSchemaName string, bucketFsBasePath string) *ExtensionContext {
	var sqlClient backend.SimpleSQLClient = backend.NewSqlClient(txCtx.GetContext(), txCtx.GetTransaction())
	var bucketFsClient bfs.BucketFsAPI = bfs.CreateBucketFsAPI(bucketFsBasePath)
	var metadataReader exaMetadata.ExaMetadataReader = exaMetadata.CreateExaMetaDataReader()
	return CreateContextWithClient(extensionSchemaName, txCtx, sqlClient, bucketFsClient, metadataReader)
}

func CreateContextWithClient(extensionSchemaName string, txCtx *transaction.TransactionContext,
	client backend.SimpleSQLClient, bucketFsClient bfs.BucketFsAPI, metadataReader exaMetadata.ExaMetadataReader) *ExtensionContext {
	return &ExtensionContext{
		ExtensionSchemaName: extensionSchemaName,
		SqlClient:           &contextSqlClient{client},
		BucketFs: &bucketFsContextImpl{
			bucketFsClient: bucketFsClient,
			context:        txCtx.GetContext(),
			db:             txCtx.GetDBConnection(),
		},
		Metadata: &metadataContextImpl{
			transaction:    txCtx.GetTransaction(),
			schemaName:     extensionSchemaName,
			metadataReader: metadataReader,
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
	Metadata            MetadataContext  `json:"metadata"`            // Allows extensions to read Exasol metadata tables
}

// reportError panics with the given error.
//
// Context functions are called by JavaScript code. The only way to report a failure is to panic.
// The JS runtime will convert this panic into a thrown JS error.
func reportError(err error) {
	panic(err)
}
