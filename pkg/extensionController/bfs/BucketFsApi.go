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

	// FindAbsolutePath searches for a file with the given name in BucketFS and returns its absolute path.
	// If multiple files with the same name exist in different folders, this picks an arbitrary file and returns its path.
	// If no file with the given name exists, this will return an error.
	FindAbsolutePath(ctx context.Context, db *sql.DB, fileName string) (string, error)
}

// BfsFile represents a file in BucketFS.
type BfsFile struct {
	Path string // Absolute path in BucketFS, starting with the base path, e.g. "/buckets/bfsdefault/default/"
	Name string // File name
	Size int    // File size in bytes
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
func (bfs bucketFsAPIImpl) ListFiles(ctx context.Context, db *sql.DB) (files []BfsFile, retErr error) {
	transaction, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create a transaction. Cause: %w", err)
	}
	defer func() {
		if err = transaction.Rollback(); err != nil {
			retErr = fmt.Errorf("failed to rollback transaction. Cause: %w", err)
		}
	}()
	udfScriptName, err := bfs.createUdfScript(transaction)
	if err != nil {
		return nil, err
	}
	return bfs.queryBucketFsContent(transaction, udfScriptName)
}

func (bfs bucketFsAPIImpl) FindAbsolutePath(ctx context.Context, db *sql.DB, fileName string) (absolutePath string, retErr error) {
	transaction, err := db.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create a transaction. Cause: %w", err)
	}
	defer func() {
		if err = transaction.Rollback(); err != nil {
			retErr = fmt.Errorf("failed to rollback transaction. Cause: %w", err)
		}
	}()
	udfScriptName, err := bfs.createUdfScript(transaction)
	if err != nil {
		return "", err
	}
	return bfs.queryAbsoluteFilePath(transaction, udfScriptName, fileName)
}

//go:embed list_files_recursively_udf.py
var listFilesRecursivelyUdfContent string

func (bfs bucketFsAPIImpl) createUdfScript(transaction *sql.Tx) (string, error) {
	schemaName := fmt.Sprintf("INTERNAL_%v", time.Now().Unix())
	_, err := transaction.Exec("CREATE SCHEMA " + schemaName)
	if err != nil {
		return "", fmt.Errorf("failed to create a schema for BucketFS list script. Cause: %w", err)
	}
	udfScriptName := fmt.Sprintf(`"%s"."LIST_RECURSIVELY"`, schemaName)
	script := fmt.Sprintf(`CREATE OR REPLACE PYTHON3 SCALAR SCRIPT %s ("path" VARCHAR(100))
	EMITS ("FILE_NAME" VARCHAR(250), "FULL_PATH" VARCHAR(500), "SIZE" DECIMAL(18,0)) AS
%s
/`, udfScriptName, listFilesRecursivelyUdfContent)
	_, err = transaction.Exec(script)
	if err != nil {
		return "", fmt.Errorf("failed to create UDF script for listing bucket. Cause: %w", err)
	}
	return udfScriptName, nil
}

func (bfs bucketFsAPIImpl) queryBucketFsContent(transaction *sql.Tx, udfScriptName string) ([]BfsFile, error) {
	statement, err := transaction.Prepare("SELECT " + udfScriptName + "(?) ORDER BY FULL_PATH") //nolint:gosec // SQL string concatenation is safe here
	if err != nil {
		return nil, fmt.Errorf("failed to create prepared statement for running list files UDF. Cause: %w", err)
	}
	defer statement.Close()
	result, err := statement.Query(bfs.bucketFsBasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to list files in BucketFS using UDF. Cause: %w", err)
	}
	defer result.Close()
	return readQueryResult(result)
}

func readQueryResult(result *sql.Rows) ([]BfsFile, error) {
	var files []BfsFile
	for result.Next() {
		if result.Err() != nil {
			return nil, fmt.Errorf("failed iterating BucketFS list UDF. Cause: %w", result.Err())
		}
		var file BfsFile
		var fileSize float64
		err := result.Scan(&file.Name, &file.Path, &fileSize)
		if err != nil {
			return nil, fmt.Errorf("failed reading result of BucketFS list UDF. Cause: %w", err)
		}
		file.Size = int(fileSize)
		files = append(files, file)
	}
	return files, nil
}

func (bfs bucketFsAPIImpl) queryAbsoluteFilePath(transaction *sql.Tx, udfScriptName string, fileName string) (string, error) {
	statement, err := transaction.Prepare(`SELECT FULL_PATH FROM (SELECT ` + udfScriptName + `(?)) ORDER BY FULL_PATH LIMIT 1`) //nolint:gosec // SQL string concatenation is safe here
	if err != nil {
		return "", fmt.Errorf("failed to create prepared statement for running list files UDF. Cause: %w", err)
	}
	defer statement.Close()
	result, err := statement.Query(bfs.bucketFsBasePath, fileName)
	if err != nil {
		return "", fmt.Errorf("failed to find absolute path in BucketFS using UDF. Cause: %w", err)
	}
	defer result.Close()
	if !result.Next() {
		if result.Err() != nil {
			return "", fmt.Errorf("failed iterating absolute path results. Cause: %w", result.Err())
		}
		return "", fmt.Errorf("file %q not found in BucketFS", fileName)
	}
	var absolutePath string
	err = result.Scan(&absolutePath)
	if err != nil {
		return "", fmt.Errorf("failed reading absolute path. Cause: %w", err)
	}
	return absolutePath, nil
}
