package extensionController

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/exasol/extension-manager/pkg/extensionAPI"
	"github.com/exasol/extension-manager/pkg/extensionController/bfs"
	"github.com/exasol/extension-manager/pkg/extensionController/transaction"
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

	// UpgradeExtension upgrades an installed extension to the latest version.
	// db is a connection to the Exasol DB
	// extensionId is the ID of the extension to uninstall
	UpgradeExtension(ctx context.Context, db *sql.DB, extensionId string) (*extensionAPI.JsUpgradeResult, error)

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
	/* [impl -> dsn~configure-bucketfs-path~1] */
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
//
// Deprecated: Use function [CreateWithValidatedConfig] which additionally validates the given configuration.
func CreateWithConfig(config ExtensionManagerConfig) TransactionController {
	controller, _ := CreateWithValidatedConfig(config)
	return controller
}

// CreateWithValidatedConfig validates the configuration and creates a new instance of [TransactionController].
func CreateWithValidatedConfig(config ExtensionManagerConfig) (TransactionController, error) {
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}
	controller := createImpl(config)
	transactionController := &transactionControllerImpl{
		controller:         controller,
		transactionStarter: transaction.BeginTransaction,
		config:             config,
	}
	return transactionController, nil
}

func validateConfig(config ExtensionManagerConfig) error {
	if config.BucketFSBasePath == "" {
		return errors.New("missing BucketFSBasePath")
	}
	if config.ExtensionRegistryURL == "" {
		return errors.New("missing ExtensionRegistryURL")
	}
	if config.ExtensionSchema == "" {
		return errors.New("missing ExtensionSchema")
	}
	return nil
}

type transactionControllerImpl struct {
	controller         controller
	transactionStarter transaction.TransactionStarter
	config             ExtensionManagerConfig
}

func (c *transactionControllerImpl) GetAllExtensions(ctx context.Context, db *sql.DB) ([]*Extension, error) {
	bfsFiles, err := c.listBfsFiles(ctx, db)
	if err != nil {
		return nil, err
	}
	return c.controller.GetAllExtensions(bfsFiles)
}

func (c *transactionControllerImpl) listBfsFiles(ctx context.Context, db *sql.DB) ([]bfs.BfsFile, error) {
	txCtx, err := c.beginTransaction(ctx, db)
	if err != nil {
		return nil, err
	}
	defer txCtx.Rollback()
	bfsClient, err := txCtx.GetBucketFsClient()
	if err != nil {
		return nil, fmt.Errorf("failed to search for required files in BucketFS. Cause: %w", err)
	}
	bfsFiles, err := bfsClient.ListFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to search for required files in BucketFS. Cause: %w", err)
	}
	return bfsFiles, nil
}

func (c *transactionControllerImpl) InstallExtension(ctx context.Context, db *sql.DB, extensionId string, extensionVersion string) (returnErr error) {
	txCtx, err := c.beginTransaction(ctx, db)
	if err != nil {
		return err
	}
	defer txCtx.Rollback()
	err = c.controller.InstallExtension(txCtx, extensionId, extensionVersion)
	if err == nil {
		err = txCtx.Commit()
		if err != nil {
			return err
		}
	}
	return err
}

func (c *transactionControllerImpl) UninstallExtension(ctx context.Context, db *sql.DB, extensionId string, extensionVersion string) (returnErr error) {
	tx, err := c.beginTransaction(ctx, db)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	err = c.controller.UninstallExtension(tx, extensionId, extensionVersion)
	if err == nil {
		err = tx.Commit()
		if err != nil {
			return err
		}
	}
	return err
}

/* [impl -> dsn~upgrade-extension~1]. */
func (c *transactionControllerImpl) UpgradeExtension(ctx context.Context, db *sql.DB, extensionId string) (result *extensionAPI.JsUpgradeResult, returnErr error) {
	tx, err := c.beginTransaction(ctx, db)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	result, err = c.controller.UpgradeExtension(tx, extensionId)
	if err == nil {
		err = tx.Commit()
		if err != nil {
			return nil, err
		}
	}
	return result, err
}

func (c *transactionControllerImpl) GetInstalledExtensions(ctx context.Context, db *sql.DB) ([]*extensionAPI.JsExtInstallation, error) {
	tx, err := c.beginTransaction(ctx, db)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	return c.controller.GetAllInstallations(tx)
}

func (c *transactionControllerImpl) GetParameterDefinitions(ctx context.Context, db *sql.DB, extensionId string, extensionVersion string) ([]parameterValidator.ParameterDefinition, error) {
	tx, err := c.beginTransaction(ctx, db)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	return c.controller.GetParameterDefinitions(tx, extensionId, extensionVersion)
}

func (c *transactionControllerImpl) CreateInstance(ctx context.Context, db *sql.DB, extensionId string, extensionVersion string, parameterValues []ParameterValue) (*extensionAPI.JsExtInstance, error) {
	tx, err := c.beginTransaction(ctx, db)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	instance, err := c.controller.CreateInstance(tx, extensionId, extensionVersion, parameterValues)
	if err == nil {
		err = tx.Commit()
		if err != nil {
			return nil, err
		}
	}
	return instance, err
}

func (c *transactionControllerImpl) FindInstances(ctx context.Context, db *sql.DB, extensionId string, extensionVersion string) ([]*extensionAPI.JsExtInstance, error) {
	tx, err := c.beginTransaction(ctx, db)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	return c.controller.FindInstances(tx, extensionId, extensionVersion)
}

func (c *transactionControllerImpl) DeleteInstance(ctx context.Context, db *sql.DB, extensionId, extensionVersion, instanceId string) error {
	tx, err := c.beginTransaction(ctx, db)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	err = c.controller.DeleteInstance(tx, extensionId, extensionVersion, instanceId)
	if err == nil {
		err = tx.Commit()
		if err != nil {
			return err
		}
	}
	return err
}

func (c *transactionControllerImpl) beginTransaction(ctx context.Context, db *sql.DB) (*transaction.TransactionContext, error) {
	tx, err := c.transactionStarter(ctx, db, c.config.BucketFSBasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	return tx, nil
}
