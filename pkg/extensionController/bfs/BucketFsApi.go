package bfs

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// BucketFsAPI allows access to BucketFS.
// Users must call the [BucketFsAPI.Close] method to release resources after using the BucketFS API.
type BucketFsAPI interface {
	// ListFiles lists all files in the configured directory recursively.
	ListFiles() ([]BfsFile, error)

	// FindAbsolutePath searches for a file with the given name in BucketFS and returns its absolute path.
	// If multiple files with the same name exist in different folders, this picks an arbitrary file and returns its path.
	// If no file with the given name exists, this will return an error.
	FindAbsolutePath(fileName string) (string, error)

	// Close removes any resources used by the BucketFS API like Exasol UDF SCRIPTS.
	Close() error
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
// Call the [BucketFsAPI.Close] method to release resources after using the BucketFS API.
/* [impl -> dsn~configure-bucketfs-path~1]. */
func CreateBucketFsAPI(bucketFsBasePath string, ctx context.Context, db *sql.DB) (BucketFsAPI, error) {
	if bucketFsBasePath == "" {
		return nil, errors.New("bucketFsBasePath is empty")
	}
	transaction, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create a transaction. Cause: %w", err)
	}
	udfScriptName, err := createUdfScript(transaction)
	if err != nil {
		_ = transaction.Rollback()
		return nil, err
	}
	return &bucketFsAPIImpl{bucketFsBasePath: bucketFsBasePath, udfScriptName: udfScriptName, transaction: transaction}, nil
}

type bucketFsAPIImpl struct {
	bucketFsBasePath string
	udfScriptName    string
	transaction      *sql.Tx
}

/* [impl -> dsn~extension-components~1]. */
func (bfs bucketFsAPIImpl) ListFiles() ([]BfsFile, error) {
	t0 := time.Now()
	statement, err := bfs.transaction.Prepare("SELECT " + bfs.udfScriptName + "(?) ORDER BY FULL_PATH") //nolint:gosec // SQL string concatenation is safe here
	if err != nil {
		return nil, fmt.Errorf("failed to create prepared statement for listing files. Cause: %w", err)
	}
	defer statement.Close()
	result, err := statement.Query(bfs.bucketFsBasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to list files. Cause: %w", err)
	}
	defer result.Close()
	files, err := readQueryResult(result)
	if logrus.IsLevelEnabled(logrus.TraceLevel) {
		for _, file := range files {
			logrus.Tracef("- Found file %q with size %d", file.Path, file.Size)
		}
	}
	logrus.Debugf("Listed %d files under %q in %dms", len(files), bfs.bucketFsBasePath, time.Since(t0).Milliseconds())
	return files, err
}

/* [impl -> dsn~resolving-files-in-bucketfs~1]. */
/* [impl -> dsn~extension-context-bucketfs~1]. */
func (bfs bucketFsAPIImpl) FindAbsolutePath(fileName string) (string, error) {
	t0 := time.Now()
	query := fmt.Sprintf(`SELECT FULL_PATH FROM (SELECT %s(?)) WHERE FILE_NAME = ? ORDER BY FULL_PATH LIMIT 1`, bfs.udfScriptName) //nolint:gosec // SQL string concatenation is safe here
	statement, err := bfs.transaction.Prepare(query)
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
	logrus.Tracef("Found absolute path %q for file %q in %.2fs", absolutePath, fileName, time.Since(t0).Seconds())
	return absolutePath, nil
}

//go:embed udf/list_files_udf.py
var listFilesRecursivelyUdfContent string

func createUdfScript(transaction *sql.Tx) (string, error) {
	t0 := time.Now()
	schemaName := fmt.Sprintf("INTERNAL_%v", t0.Unix())
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
	logrus.Tracef("Created UDF script %s in %.2f", udfScriptName, time.Since(t0).Seconds())
	return udfScriptName, nil
}

func readQueryResult(result *sql.Rows) ([]BfsFile, error) {
	var files []BfsFile
	for result.Next() {
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

func (bfs bucketFsAPIImpl) Close() error {
	if err := bfs.transaction.Rollback(); err != nil {
		return fmt.Errorf("failed to rollback transaction to cleanup resources. Cause: %w", err)
	}
	return nil
}
