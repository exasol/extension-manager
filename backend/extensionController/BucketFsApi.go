package extensionController

import (
	"database/sql"
	"fmt"
	"strings"
)

// BucketFsAPI allows access to BucketFS. Currently, it's implemented by running UDFs via SQL that read the data. In the future that implementation might be replaced by direct access.
type BucketFsAPI interface {
	// ListBuckets returns a list of public buckets
	ListBuckets() ([]string, error)
	// ListFiles lists the files in a given bucket
	ListFiles(bucket string) ([]BfsFile, error)
}

// CreateBucketFsAPI create an instance of BucketFsAPI
func CreateBucketFsAPI(db *sql.DB) BucketFsAPI {
	bucketFsAPI := bucketFsAPIImpl{db: db}
	return &bucketFsAPI
}

type bucketFsAPIImpl struct {
	db *sql.DB
}

func (bfsApi bucketFsAPIImpl) ListBuckets() ([]string, error) {
	files, err := bfsApi.listDirInUDF("/buckets/bfsdefault/")
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(files))
	for _, file := range files {
		names = append(names, file.Name)
	}
	return names, nil
}

func (bfsApi bucketFsAPIImpl) ListFiles(bucket string) ([]BfsFile, error) {
	if strings.Contains(bucket, "/") {
		return nil, fmt.Errorf("invalid bucket name. Bucket name must not contain slashes")
	}
	return bfsApi.listDirInUDF("/buckets/bfsdefault/" + bucket)
}

// BfsFile represents a file in BucketFS
type BfsFile struct {
	Name string
	Size int
}

func (bfsApi bucketFsAPIImpl) listDirInUDF(directory string) ([]BfsFile, error) {
	transaction, err := bfsApi.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to create a transaction. Cause: %w", err)
	}
	_, err = transaction.Exec("CREATE SCHEMA INTERNAL")
	if err != nil {
		return nil, fmt.Errorf("failed to create a schema for bucket_fs_list script. Cause: %w", err)
	}
	_, err = transaction.Exec(`CREATE OR REPLACE PYTHON3 SCALAR SCRIPT INTERNAL."LS" ("my_path" VARCHAR(100)) EMITS ("FILES" VARCHAR(250), "SIZE" DECIMAL(18,0)) AS
import os
def run(ctx):
    for line in os.listdir(ctx.my_path):
        size = os.path.getsize(ctx.my_path + "/" + line)
        ctx.emit(line, size)
/`)
	if err != nil {
		return nil, fmt.Errorf("failed to create script for listing bucket. Cause: %w", err)
	}
	statement, err := transaction.Prepare("SELECT INTERNAL.LS(?)")
	if err != nil {
		return nil, fmt.Errorf("failed to create prepard statement for running list files UDF. Cause: %w", err)
	}
	result, err := statement.Query(directory)
	if err != nil {
		return nil, fmt.Errorf("failed to list files in BucketFS using UDF. Cause: %w", err)
	}
	var files []BfsFile
	for result.Next() {
		var file BfsFile
		err = result.Scan(&file.Name, &file.Size)
		if err != nil {
			return nil, fmt.Errorf("failed reading result of BucketFS list UDF. Cause: %w", err)
		}
		files = append(files, file)
	}
	err = transaction.Rollback()
	if err != nil {
		return nil, fmt.Errorf("failed to rollback transaction. Cause: %w", err)
	}
	return files, nil
}
