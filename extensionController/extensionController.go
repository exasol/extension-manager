package extensionController

import (
	"context"
	"database/sql"
	"fmt"

	"strings"

	"github.com/exasol/extension-manager/extensionAPI"
	"github.com/exasol/extension-manager/parameterValidator"
)

// ExtensionController is the core part of the extension-manager that provides the extension handling functionality.
type ExtensionController interface {
	// GetAllInstallations searches for installations of any extensions.
	// db is a connection to the Exasol DB
	GetAllInstallations(ctx context.Context, db *sql.DB) ([]*extensionAPI.JsExtInstallation, error)

	// InstallExtension installs an extension.
	// db is a connection to the Exasol DB
	// extensionId is the ID of the extension to install
	// extensionVersion is the version of the extension to install
	InstallExtension(ctx context.Context, db *sql.DB, extensionId string, extensionVersion string) error

	// GetAllExtensions reports all extension definitions.
	// db is a connection to the Exasol DB
	GetAllExtensions(ctx context.Context, db *sql.DB) ([]*Extension, error)

	// CreateInstance creates a new instance of an extension, e.g. a virtual schema and returns it's name.
	// db is a connection to the Exasol DB
	CreateInstance(ctx context.Context, db *sql.DB, extensionId string, extensionVersion string, parameterValues []ParameterValue) (string, error)
}

type Extension struct {
	Id                  string
	Name                string
	Description         string
	InstallableVersions []string
}

type ParameterValue struct {
	Name  string
	Value string
}

// ExtInstallation represents the installation of an Extension
type ExtInstallation struct {
}

// Create an instance of ExtensionController
func Create(pathToExtensionFolder string, extensionSchemaName string) ExtensionController {
	ctrl := createImpl(pathToExtensionFolder, extensionSchemaName)
	return &extensionControllerImpl{impl: ctrl}
}

type extensionControllerImpl struct {
	impl controller
}

func (c *extensionControllerImpl) GetAllExtensions(ctx context.Context, db *sql.DB) ([]*Extension, error) {
	bfsFiles, err := listBfsFiles(ctx, db)
	if err != nil {
		return nil, err
	}
	return c.impl.GetAllExtensions(bfsFiles)
}

func listBfsFiles(ctx context.Context, db *sql.DB) ([]BfsFile, error) {
	bfsAPI := CreateBucketFsAPI(ctx, db)
	bfsFiles, err := bfsAPI.ListFiles("default")
	if err != nil {
		return nil, fmt.Errorf("failed to search for required files in BucketFS. Cause: %w", err)
	}
	return bfsFiles, nil
}

func existsFileInBfs(bfsFiles []BfsFile, requiredFile extensionAPI.BucketFsUpload) bool {
	for _, existingFile := range bfsFiles {
		if requiredFile.BucketFsFilename == existingFile.Name && requiredFile.FileSize == existingFile.Size {
			return true
		}
	}
	return false
}

func (c *extensionControllerImpl) InstallExtension(ctx context.Context, db *sql.DB, extensionId string, extensionVersion string) (returnErr error) {
	tx, err := beginTransaction(ctx, db)
	if err != nil {
		return err
	}
	defer rollback(tx)
	err = c.impl.InstallExtension(tx, extensionId, extensionVersion)
	if err == nil {
		err = tx.Commit()
		if err != nil {
			return err
		}
	}
	return err
}

func (c *extensionControllerImpl) GetAllInstallations(ctx context.Context, db *sql.DB) ([]*extensionAPI.JsExtInstallation, error) {
	tx, err := beginTransaction(ctx, db)
	if err != nil {
		return nil, err
	}
	defer rollback(tx)
	return c.impl.GetAllInstallations(tx)
}

func (c *extensionControllerImpl) CreateInstance(ctx context.Context, db *sql.DB, extensionId string, extensionVersion string, parameterValues []ParameterValue) (string, error) {
	tx, err := beginTransaction(ctx, db)
	if err != nil {
		return "", err
	}
	defer rollback(tx)
	instanceName, err := c.impl.CreateInstance(tx, extensionId, extensionVersion, parameterValues)
	if err == nil {
		err = tx.Commit()
		if err != nil {
			return "", err
		}
	}
	return instanceName, err
}

func validateParameters(parameterDefinitions []interface{}, params extensionAPI.ParameterValues) error {
	validator, err := parameterValidator.New()
	if err != nil {
		return fmt.Errorf("failed to create parameter validator: %w", err)
	}
	result, err := validator.ValidateParameters(parameterDefinitions, params)
	if err != nil {
		return fmt.Errorf("failed to validate parameters: %w", err)
	}
	message := ""
	for _, r := range result {
		message += r.Message + ", "
	}
	message = strings.TrimSuffix(message, ", ")
	if message != "" {
		return fmt.Errorf("invalid parameters: %s", message)
	}
	return nil
}

func beginTransaction(ctx context.Context, db *sql.DB) (*sql.Tx, error) {
	return db.BeginTx(ctx, nil)
}

func rollback(tx *sql.Tx) {
	// Even if Tx.Rollback fails, the transaction will no longer be valid, nor will it have been committed to the database.
	// See https://go.dev/doc/database/execute-transactions
	_ = tx.Rollback()
}
