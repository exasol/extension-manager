package extensionController

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/exasol/extension-manager/extensionAPI"
	"github.com/exasol/extension-manager/parameterValidator"
)

// ExtensionController is the core part of the extension-manager that provides the extension handling functionality.
type ExtensionController interface {
	// GetAllInstallations searches for installations of any extensions.
	// db is a connection to the Exasol DB with autocommit turned off
	GetAllInstallations(ctx context.Context, db *sql.DB) ([]*extensionAPI.JsExtInstallation, error)

	// InstallExtension installs an extension.
	// db is a connection to the Exasol DB with autocommit turned off
	// extensionId is the ID of the extension to install
	// extensionVersion is the version of the extension to install
	InstallExtension(ctx context.Context, db *sql.DB, extensionId string, extensionVersion string) error

	// GetAllExtensions reports all extension definitions.
	// db is a connection to the Exasol DB with autocommit turned off
	GetAllExtensions(ctx context.Context, db *sql.DB) ([]*Extension, error)

	// CreateInstance creates a new instance of an extension, e.g. a virtual schema and returns it's name.
	// db is a connection to the Exasol DB with autocommit turned off
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
	return &extensionControllerImpl{pathToExtensionFolder: pathToExtensionFolder, extensionSchemaName: extensionSchemaName}
}

type extensionControllerImpl struct {
	pathToExtensionFolder string
	extensionSchemaName   string
}

func (c *extensionControllerImpl) GetAllExtensions(ctx context.Context, db *sql.DB) ([]*Extension, error) {
	jsExtensions, err := c.getAllExtensions()
	if err != nil {
		return nil, err
	}
	bfsFiles, err := listBfsFiles(ctx, db)
	if err != nil {
		return nil, err
	}
	var extensions []*Extension
	for _, jsExtension := range jsExtensions {
		if c.requiredFilesAvailable(jsExtension, bfsFiles) {
			extension := Extension{Id: jsExtension.Id, Name: jsExtension.Name, Description: jsExtension.Description, InstallableVersions: jsExtension.InstallableVersions}
			extensions = append(extensions, &extension)
		}
	}
	return extensions, nil
}

func listBfsFiles(ctx context.Context, db *sql.DB) ([]BfsFile, error) {
	bfsAPI := CreateBucketFsAPI(ctx, db)
	bfsFiles, err := bfsAPI.ListFiles("default")
	if err != nil {
		return nil, fmt.Errorf("failed to search for required files in BucketFS. Cause: %w", err)
	}
	return bfsFiles, nil
}

func (c *extensionControllerImpl) requiredFilesAvailable(extension *extensionAPI.JsExtension, bfsFiles []BfsFile) bool {
	for _, requiredFile := range extension.BucketFsUploads {
		if !existsFileInBfs(bfsFiles, requiredFile) {
			log.Printf("Ignoring extension %q since the required file %q does not exist or has a wrong file size.\n", extension.Name, requiredFile.Name)
			return false
		}
	}
	log.Printf("Required files found for extension %q\n", extension.Name)
	return true
}

func existsFileInBfs(bfsFiles []BfsFile, requiredFile extensionAPI.BucketFsUpload) bool {
	for _, existingFile := range bfsFiles {
		if requiredFile.BucketFsFilename == existingFile.Name && requiredFile.FileSize == existingFile.Size {
			return true
		}
	}
	return false
}

func (c *extensionControllerImpl) getAllExtensions() ([]*extensionAPI.JsExtension, error) {
	var extensions []*extensionAPI.JsExtension
	extensionPaths := FindJSFilesInDir(c.pathToExtensionFolder)
	for _, path := range extensionPaths {
		extension, err := c.loadExtensionFromFile(path)
		if err == nil {
			extensions = append(extensions, extension)
		} else {
			log.Printf("error: Failed to load extension. This extension will be ignored. Cause: %v\n", err)
		}
	}
	return extensions, nil
}

func (c *extensionControllerImpl) loadExtensionById(id string) (*extensionAPI.JsExtension, error) {
	extensionPath := path.Join(c.pathToExtensionFolder, id)
	return c.loadExtensionFromFile(extensionPath)
}

func (c *extensionControllerImpl) loadExtensionFromFile(extensionPath string) (*extensionAPI.JsExtension, error) {
	extension, err := extensionAPI.GetExtensionFromFile(extensionPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load extension from file %q: %w", extensionPath, err)
	}
	return extension, nil
}

func (c *extensionControllerImpl) InstallExtension(ctx context.Context, db *sql.DB, extensionId string, extensionVersion string) error {
	tx, err := beginTransaction(ctx, db)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	err = c.installExtensionWithTx(tx, extensionId, extensionVersion)
	if err != nil {
		tx.Commit()
	}
	return err
}

func (c *extensionControllerImpl) installExtensionWithTx(tx *sql.Tx, extensionId string, extensionVersion string) error {
	extension, err := c.loadExtensionById(extensionId)
	if err != nil {
		return fmt.Errorf("failed to load extension with id %q: %w", extensionId, err)
	}
	err = c.ensureSchemaExists(tx)
	if err != nil {
		return err
	}
	return extension.Install(c.createExtensionContext(tx), extensionVersion)
}

func (c *extensionControllerImpl) GetAllInstallations(ctx context.Context, db *sql.DB) ([]*extensionAPI.JsExtInstallation, error) {
	tx, err := beginTransaction(ctx, db)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	return c.getAllInstallationsWithTx(tx)
}

func (c *extensionControllerImpl) getAllInstallationsWithTx(tx *sql.Tx) ([]*extensionAPI.JsExtInstallation, error) {
	metadata, err := extensionAPI.ReadMetadataTables(tx, c.extensionSchemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata tables. Cause: %w", err)
	}
	extensions, err := c.getAllExtensions()
	if err != nil {
		return nil, err
	}
	extensionContext := c.createExtensionContext(tx)
	var allInstallations []*extensionAPI.JsExtInstallation
	for _, extension := range extensions {
		installations, err := extension.FindInstallations(extensionContext, metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to find installations: %v", err)
		} else {
			allInstallations = append(allInstallations, installations...)
		}
	}
	return allInstallations, nil
}

func (c *extensionControllerImpl) CreateInstance(ctx context.Context, db *sql.DB, extensionId string, extensionVersion string, parameterValues []ParameterValue) (string, error) {
	tx, err := beginTransaction(ctx, db)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()
	instanceName, err := c.createInstanceWithTx(tx, extensionId, extensionVersion, parameterValues)
	if err != nil {
		tx.Commit()
	}
	return instanceName, err
}

func (c *extensionControllerImpl) createInstanceWithTx(tx *sql.Tx, extensionId string, extensionVersion string, parameterValues []ParameterValue) (string, error) {
	extension, err := c.loadExtensionById(extensionId)
	if err != nil {
		return "", fmt.Errorf("failed to load extension with id %q: %w", extensionId, err)
	}
	err = c.ensureSchemaExists(tx)
	if err != nil {
		return "", err
	}
	params := extensionAPI.ParameterValues{}
	for _, p := range parameterValues {
		params.Values = append(params.Values, extensionAPI.ParameterValue{Name: p.Name, Value: p.Value})
	}

	extensionContext := c.createExtensionContext(tx)
	installation, err := c.findInstallationByVersion(tx, extensionContext, extension, extensionVersion)
	if err != nil {
		return "", fmt.Errorf("failed to find installations: %w", err)
	}

	err = validateParameters(installation.InstanceParameters, params)
	if err != nil {
		return "", err
	}

	instance, err := extension.AddInstance(extensionContext, extensionVersion, &params)
	if err != nil {
		return "", err
	}
	if instance == nil {
		return "", fmt.Errorf("extension did not return an instance")
	}
	return instance.Name, nil
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

func (c *extensionControllerImpl) findInstallationByVersion(tx *sql.Tx, context *extensionAPI.ExtensionContext, extension *extensionAPI.JsExtension, version string) (*extensionAPI.JsExtInstallation, error) {
	metadata, err := extensionAPI.ReadMetadataTables(tx, c.extensionSchemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata tables. Cause: %w", err)
	}

	installations, err := extension.FindInstallations(context, metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to find installations. Cause: %w", err)
	}
	var availableVersions []string
	for _, i := range installations {
		if i.Version == version {
			return i, nil
		}
		availableVersions = append(availableVersions, i.Version)
	}
	return nil, fmt.Errorf("version %q not found for extension %q, available versions: %q", version, extension.Id, availableVersions)
}

func beginTransaction(ctx context.Context, db *sql.DB) (*sql.Tx, error) {
	return db.BeginTx(ctx, nil)
}

func (c *extensionControllerImpl) createExtensionContext(tx *sql.Tx) *extensionAPI.ExtensionContext {
	return extensionAPI.CreateContext(c.extensionSchemaName, tx)
}

func (c *extensionControllerImpl) ensureSchemaExists(tx *sql.Tx) error {
	_, err := tx.Exec(fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS "%s"`, c.extensionSchemaName))
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}
	return nil
}
