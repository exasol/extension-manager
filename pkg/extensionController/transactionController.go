package extensionController

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/exasol/extension-manager/pkg/apiErrors"
	"github.com/exasol/extension-manager/pkg/extensionAPI"
	"github.com/exasol/extension-manager/pkg/parameterValidator"
)

// TransactionController is the core part of the extension-manager that provides the extension handling functionality.
// All of it's methods expect a [context.Context] and [*sql.DB] as arguments.
// The controller will take care of transaction handling,
// i.e. it will create a new transaction and commit or rollback if necessary.
type TransactionController interface {
	// GetAllExtensions reports all extension definitions.
	// db is a connection to the Exasol DB
	GetAllExtensions(ctx context.Context, db *sql.DB) ([]*Extension, error)

	// GetInstalledExtensions searches for installations of any extensions.
	// db is a connection to the Exasol DB
	GetInstalledExtensions(ctx context.Context, db *sql.DB) ([]*extensionAPI.JsExtInstallation, error)

	// GetParameterDefinitions returns the parameter definitions required for installing a given extension version.
	GetParameterDefinitions(ctx context.Context, db *sql.DB, extensionId string, extensionVersion string) ([]parameterValidator.ParameterDefinition, error)

	// InstallExtension installs an extension.
	// db is a connection to the Exasol DB
	// extensionId is the ID of the extension to install
	// extensionVersion is the version of the extension to install
	InstallExtension(ctx context.Context, db *sql.DB, extensionId string, extensionVersion string) error

	// UninstallExtension uninstalls an extension.
	// db is a connection to the Exasol DB
	// extensionId is the ID of the extension to uninstall
	// extensionVersion is the version of the extension to uninstall
	UninstallExtension(ctx context.Context, db *sql.DB, extensionId string, extensionVersion string) error

	// CreateInstance creates a new instance of an extension, e.g. a virtual schema and returns it's name.
	// db is a connection to the Exasol DB
	CreateInstance(ctx context.Context, db *sql.DB, extensionId string, extensionVersion string, parameterValues []ParameterValue) (*extensionAPI.JsExtInstance, error)

	// FindInstances returns a list of all instances for the given version.
	FindInstances(ctx context.Context, db *sql.DB, extensionId string, extensionVersion string) ([]*extensionAPI.JsExtInstance, error)

	// DeleteInstance deletes instance with the given ID.
	DeleteInstance(ctx context.Context, db *sql.DB, extensionId, extensionVersion, instanceId string) error
}

type Extension struct {
	Id                  string
	Name                string
	Category            string
	Description         string
	InstallableVersions []extensionAPI.JsExtensionVersion
}

type ParameterValue struct {
	Name  string
	Value string
}

// ExtInstallation represents the installation of an Extension.
type ExtInstallation struct {
}

// Configuration options for the extension manager.
type ExtensionManagerConfig struct {
	// URL of the extension registry index used to find available extensions.
	// This can also be the path of a local directory for local testing.
	ExtensionRegistryURL string
	// BucketFS base path where to search for extension files, e.g. "/buckets/bfsdefault/default/".
	BucketFSBasePath string
	// Schema where extensions are searched for and new extensions are created, e.g. "EXA_EXTENSIONS".
	ExtensionSchema string
}

// Create creates a new instance of [TransactionController].
//
// Deprecated: Use function [CreateWithConfig] which allows specifying additional configuration options.
func Create(extensionRegistryURL string, schema string) TransactionController {
	return CreateWithConfig(ExtensionManagerConfig{
		ExtensionRegistryURL: extensionRegistryURL,
		BucketFSBasePath:     "/buckets/bfsdefault/default/",
		ExtensionSchema:      schema,
	})
}

// CreateWithConfig creates a new instance of [TransactionController] with more configuration options.
func CreateWithConfig(config ExtensionManagerConfig) TransactionController {
	controller := createImpl(config)
	return &transactionControllerImpl{
		controller: controller,
		bucketFs:   CreateBucketFsAPI(config.BucketFSBasePath)}
}

type transactionControllerImpl struct {
	controller controller
	bucketFs   BucketFsAPI
}

func (c *transactionControllerImpl) GetAllExtensions(ctx context.Context, db *sql.DB) ([]*Extension, error) {
	bfsFiles, err := c.listBfsFiles(ctx, db)
	if err != nil {
		return nil, err
	}
	return c.controller.GetAllExtensions(bfsFiles)
}

func (c *transactionControllerImpl) listBfsFiles(ctx context.Context, db *sql.DB) ([]BfsFile, error) {
	bfsFiles, err := c.bucketFs.ListFiles(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("failed to search for required files in BucketFS. Cause: %w", err)
	}
	return bfsFiles, nil
}

func (c *transactionControllerImpl) InstallExtension(ctx context.Context, db *sql.DB, extensionId string, extensionVersion string) (returnErr error) {
	txCtx, err := beginTransaction(ctx, db)
	if err != nil {
		return err
	}
	defer txCtx.rollback()
	err = c.controller.InstallExtension(txCtx, extensionId, extensionVersion)
	if err == nil {
		err = txCtx.transaction.Commit()
		if err != nil {
			return err
		}
	}
	return err
}

func (c *transactionControllerImpl) UninstallExtension(ctx context.Context, db *sql.DB, extensionId string, extensionVersion string) (returnErr error) {
	tx, err := beginTransaction(ctx, db)
	if err != nil {
		return err
	}
	defer tx.rollback()
	err = c.controller.UninstallExtension(tx, extensionId, extensionVersion)
	if err == nil {
		err = tx.commit()
		if err != nil {
			return err
		}
	}
	return err
}

func (c *transactionControllerImpl) GetInstalledExtensions(ctx context.Context, db *sql.DB) ([]*extensionAPI.JsExtInstallation, error) {
	tx, err := beginTransaction(ctx, db)
	if err != nil {
		return nil, err
	}
	defer tx.rollback()
	return c.controller.GetAllInstallations(tx)
}

func (c *transactionControllerImpl) GetParameterDefinitions(ctx context.Context, db *sql.DB, extensionId string, extensionVersion string) ([]parameterValidator.ParameterDefinition, error) {
	tx, err := beginTransaction(ctx, db)
	if err != nil {
		return nil, err
	}
	defer tx.rollback()
	return c.controller.GetParameterDefinitions(tx, extensionId, extensionVersion)
}

func (c *transactionControllerImpl) CreateInstance(ctx context.Context, db *sql.DB, extensionId string, extensionVersion string, parameterValues []ParameterValue) (*extensionAPI.JsExtInstance, error) {
	tx, err := beginTransaction(ctx, db)
	if err != nil {
		return nil, err
	}
	defer tx.rollback()
	instance, err := c.controller.CreateInstance(tx, extensionId, extensionVersion, parameterValues)
	if err == nil {
		err = tx.commit()
		if err != nil {
			return nil, err
		}
	}
	return instance, err
}

func (c *transactionControllerImpl) FindInstances(ctx context.Context, db *sql.DB, extensionId string, extensionVersion string) ([]*extensionAPI.JsExtInstance, error) {
	tx, err := beginTransaction(ctx, db)
	if err != nil {
		return nil, err
	}
	defer tx.rollback()
	return c.controller.FindInstances(tx, extensionId, extensionVersion)
}

func (c *transactionControllerImpl) DeleteInstance(ctx context.Context, db *sql.DB, extensionId, extensionVersion, instanceId string) error {
	tx, err := beginTransaction(ctx, db)
	if err != nil {
		return err
	}
	defer tx.rollback()
	err = c.controller.DeleteInstance(tx, extensionId, extensionVersion, instanceId)
	if err == nil {
		err = tx.commit()
		if err != nil {
			return err
		}
	}
	return err
}

func beginTransaction(ctx context.Context, db *sql.DB) (*transactionContext, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		if strings.Contains(err.Error(), "Connection exception - authentication failed") {
			return nil, apiErrors.NewUnauthorizedErrorF("invalid database credentials")
		}
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return &transactionContext{context: ctx, db: db, transaction: tx}, nil
}

type transactionContext struct {
	context     context.Context
	db          *sql.DB
	transaction *sql.Tx
}

func (ctx *transactionContext) rollback() {
	// Even if Tx.Rollback fails, the transaction will no longer be valid, nor will it have been committed to the database.
	// See https://go.dev/doc/database/execute-transactions
	_ = ctx.transaction.Rollback()
}

func (ctx *transactionContext) commit() error {
	return ctx.transaction.Commit()
}
