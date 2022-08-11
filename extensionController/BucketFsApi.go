package extensionController

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// BucketFsAPI allows access to BucketFS. Currently, it's implemented by running UDFs via SQL that read the data.
// In the future that implementation might be replaced by direct access.
type BucketFsAPI interface {
	// ListBuckets returns a list of public buckets
	ListBuckets() ([]string, error)
	// ListFiles lists the files in a given bucket
	ListFiles(bucket string) ([]BfsFile, error)
}

// CreateBucketFsAPI creates an instance of BucketFsAPI
func CreateBucketFsAPI(ctx context.Context, db *sql.DB) BucketFsAPI {
	return &bucketFsAPIImpl{db: db, ctx: ctx}
}

type bucketFsAPIImpl struct {
	db  *sql.DB
	ctx context.Context
}

func (bfs bucketFsAPIImpl) ListBuckets() ([]string, error) {
	files, err := bfs.listDirInUDF("/buckets/bfsdefault/")
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(files))
	for _, file := range files {
		names = append(names, file.Name)
	}
	return names, nil
}

func (bfs bucketFsAPIImpl) ListFiles(bucket string) ([]BfsFile, error) {
	if strings.Contains(bucket, "/") {
		return nil, fmt.Errorf("invalid bucket name. Bucket name must not contain slashes")
	}
	return bfs.listDirInUDF("/buckets/bfsdefault/" + bucket)
}

// BfsFile represents a file in BucketFS
type BfsFile struct {
	Name string
	Size int
}

func (bfs bucketFsAPIImpl) listDirInUDF(directory string) (files []BfsFile, retErr error) {
	transaction, err := bfs.db.BeginTx(bfs.ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create a transaction. Cause: %w", err)
	}
	defer func() {
		err = transaction.Rollback()
		if err != nil {
			retErr = fmt.Errorf("failed to rollback transaction. Cause: %w", err)
		}
	}()
	schemaName := fmt.Sprintf("INTERNAL_%v", time.Now().Unix())
	_, err = transaction.Exec("CREATE SCHEMA " + schemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to create a schema for bucket_fs_list script. Cause: %w", err)
	}
	_, err = transaction.Exec(`CREATE OR REPLACE PYTHON3 SCALAR SCRIPT ` + schemaName + `."LS" ("my_path" VARCHAR(100)) EMITS ("FILES" VARCHAR(250), "SIZE" DECIMAL(18,0)) AS
import os
def run(ctx):
    for line in os.listdir(ctx.my_path):
        size = os.path.getsize(ctx.my_path + "/" + line)
        ctx.emit(line, size)
/`)
	if err != nil {
		return nil, fmt.Errorf("failed to create script for listing bucket. Cause: %w", err)
	}
	statement, err := transaction.Prepare("SELECT " + schemaName + ".LS(?)")
	if err != nil {
		return nil, fmt.Errorf("failed to create prepared statement for running list files UDF. Cause: %w", err)
	}
	result, err := statement.Query(directory)
	if err != nil {
		return nil, fmt.Errorf("failed to list files in BucketFS using UDF. Cause: %w", err)
	}
	for result.Next() {
		var file BfsFile
		var fileSize float64
		err = result.Scan(&file.Name, &fileSize)
		if err != nil {
			return nil, fmt.Errorf("failed reading result of BucketFS list UDF. Cause: %w", err)
		}
		file.Size = int(fileSize)
		files = append(files, file)
	}
	return files, nil
}
