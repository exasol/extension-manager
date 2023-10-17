package context

import (
	"fmt"

	"github.com/exasol/extension-manager/pkg/extensionController/transaction"
)

// BucketFsContext allows extensions to interact with BucketFS.
/* [impl -> dsn~extension-context-bucketfs~1]. */
type BucketFsContext interface {
	// ResolvePath returns an absolute path for the given filename in BucketFS.
	ResolvePath(fileName string) string
}

type bucketFsContextImpl struct {
	txCtx *transaction.TransactionContext
}

/* [impl -> dsn~resolving-files-in-bucketfs~1]. */
func (b *bucketFsContextImpl) ResolvePath(fileName string) string {
	path, err := b.resolvePath(fileName)
	if err != nil {
		reportError(fmt.Errorf("failed to find absolute path for file %q: %w", fileName, err))
	}
	return path
}

func (b *bucketFsContextImpl) resolvePath(fileName string) (string, error) {
	bfsClient, err := b.txCtx.GetBucketFsClient()
	if err != nil {
		return "", err
	}
	return bfsClient.FindAbsolutePath(fileName)
}
