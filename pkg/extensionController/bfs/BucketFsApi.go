package bfs

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"time"
)

// BucketFsAPI allows access to BucketFS.
type BucketFsAPI interface {
	// ListFiles lists all files in the configured directory recursively.
	ListFiles(ctx context.Context, db *sql.DB) ([]BfsFile, error)
}

// CreateBucketFsAPI creates an instance of BucketFsAPI.
//
// The current implementation uses a Python UDF for accessing BucketFS.
// In the future that implementation might be replaced by direct access.
func CreateBucketFsAPI(bucketFsBasePath string) BucketFsAPI {
	return &bucketFsAPIImpl{bucketFsBasePath: bucketFsBasePath}
}

type bucketFsAPIImpl struct {
	bucketFsBasePath string
}

/* [impl -> dsn~extension-components~1]. */
func (bfs bucketFsAPIImpl) ListFiles(ctx context.Context, db *sql.DB) ([]BfsFile, error) {
	return bfs.listDirInUDF(ctx, db)
}

// BfsFile represents a file in BucketFS.
type BfsFile struct {
	Path string // Absolute path in BucketFS, starting with the base path, e.g. "/buckets/bfsdefault/default/"
	Name string // File name
	Size int    // File size in bytes
}

func (bfs bucketFsAPIImpl) listDirInUDF(ctx context.Context, db *sql.DB) (files []BfsFile, retErr error) {
	transaction, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create a transaction. Cause: %w", err)
	}
	defer func() {
		err = transaction.Rollback()
		if err != nil {
			retErr = fmt.Errorf("failed to rollback transaction. Cause: %w", err)
		}
	}()
	scriptName, err := bfs.createUdfScript(transaction)
	if err != nil {
		return nil, err
	}
	return bfs.queryUdfScript(transaction, scriptName)
}

//go:embed list_files_recursively_udf.py
var listFilesRecursivelyUdfContent string

func (bfs bucketFsAPIImpl) createUdfScript(transaction *sql.Tx) (string, error) {
	schemaName := fmt.Sprintf("INTERNAL_%v", time.Now().Unix())
	_, err := transaction.Exec("CREATE SCHEMA " + schemaName)
	if err != nil {
		return "", fmt.Errorf("failed to create a schema for BucketFS list script. Cause: %w", err)
	}
	scriptName := fmt.Sprintf(`"%s"."LIST_RECURSIVELY"`, schemaName)
	script := fmt.Sprintf(`CREATE OR REPLACE PYTHON3 SCALAR SCRIPT %s ("path" VARCHAR(100))
	EMITS ("FILE_NAME" VARCHAR(250), "FULL_PATH" VARCHAR(500), "SIZE" DECIMAL(18,0)) AS
%s
/`, scriptName, listFilesRecursivelyUdfContent)
	_, err = transaction.Exec(script)
	if err != nil {
		return "", fmt.Errorf("failed to create script for listing bucket. Cause: %w", err)
	}
	return scriptName, nil
}

func (bfs bucketFsAPIImpl) queryUdfScript(transaction *sql.Tx, scriptName string) ([]BfsFile, error) {
	statement, err := transaction.Prepare("SELECT " + scriptName + "(?) ORDER BY FULL_PATH") //nolint:gosec // SQL string concatenation is safe here
	if err != nil {
		return nil, fmt.Errorf("failed to create prepared statement for running list files UDF. Cause: %w", err)
	}
	defer statement.Close()
	result, err := statement.Query(bfs.bucketFsBasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to list files in BucketFS using UDF. Cause: %w", err)
	}
	defer result.Close()
	return readQueryResult(result, err)
}

func readQueryResult(result *sql.Rows, err error) ([]BfsFile, error) {
	var files []BfsFile
	for result.Next() {
		if result.Err() != nil {
			return nil, fmt.Errorf("failed iterating BucketFS list UDF. Cause: %w", err)
		}
		var file BfsFile
		var fileSize float64
		err = result.Scan(&file.Name, &file.Path, &fileSize)
		if err != nil {
			return nil, fmt.Errorf("failed reading result of BucketFS list UDF. Cause: %w", err)
		}
		file.Size = int(fileSize)
		files = append(files, file)
	}
	return files, nil
}
