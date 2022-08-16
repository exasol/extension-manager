package extensionAPI

import (
	"database/sql"

	"github.com/exasol/extension-manager/backend"
)

const BUCKETFS_PATH_PREFIX = "/buckets/bfsdefault/default/"

func CreateContextWithClient(extensionSchemaName string, client SimpleSQLClient) *ExtensionContext {
	return &ExtensionContext{
		ExtensionSchemaName: extensionSchemaName,
		SqlClient:           client,
		BucketFs:            &bucketFsContextImpl{},
	}
}

func CreateContext(extensionSchemaName string, tx *sql.Tx) *ExtensionContext {
	var client SimpleSQLClient = backend.NewSqlClient(tx)
	return CreateContextWithClient(extensionSchemaName, client)
}

type ExtensionContext struct {
	ExtensionSchemaName string          `json:"extensionSchemaName"`
	BucketFs            BucketFsContext `json:"bucketFs"`
	SqlClient           SimpleSQLClient `json:"sqlClient"`
}

type BucketFsContext interface {
	ResolvePath(fileName string) string
}

type bucketFsContextImpl struct{}

func (b *bucketFsContextImpl) ResolvePath(fileName string) string {
	return BUCKETFS_PATH_PREFIX + fileName
}
