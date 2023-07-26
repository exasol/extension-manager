package context

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/exasol/extension-manager/pkg/extensionController/bfs"
)

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
