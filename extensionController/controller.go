package extensionController

import (
	"context"
	"database/sql"
	"fmt"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/exasol/extension-manager/apiErrors"
	"github.com/exasol/extension-manager/extensionAPI"
	"github.com/exasol/extension-manager/parameterValidator"
)

// controller is the core part of the extension-manager that provides the extension handling functionality.
type controller interface {
	// GetAllExtensions reports all extension definitions.
	GetAllExtensions(bfsFiles []BfsFile) ([]*Extension, error)

	// GetAllInstallations searches for installations of any extensions.
	GetAllInstallations(tx *sql.Tx) ([]*extensionAPI.JsExtInstallation, error)

	// InstallExtension installs an extension.
	// extensionId is the ID of the extension to install
	// extensionVersion is the version of the extension to install
	InstallExtension(tx *sql.Tx, extensionId string, extensionVersion string) error

	// CreateInstance creates a new instance of an extension, e.g. a virtual schema and returns it's name.
	CreateInstance(tx *sql.Tx, extensionId string, extensionVersion string, parameterValues []ParameterValue) (*extensionAPI.JsExtInstance, error)

	// FindInstances returns a list of all instances for the given version.
	FindInstances(tx *sql.Tx, extensionId string, extensionVersion string) ([]*extensionAPI.JsExtInstance, error)

	// DeleteInstance deletes instance with the given ID.
	DeleteInstance(tx *sql.Tx, extensionId string, instanceId string) error
}

type controllerImpl struct {
	extensionFolder string
	schema          string
	metaDataReader  extensionAPI.ExaMetadataReader
}

func createImpl(extensionFolder string, schema string) controller {
	return &controllerImpl{extensionFolder: extensionFolder, schema: schema, metaDataReader: extensionAPI.CreateExaMetaDataReader()}
}

func (c *controllerImpl) GetAllExtensions(bfsFiles []BfsFile) ([]*Extension, error) {
	jsExtensions := c.getAllExtensions()
	var extensions []*Extension
	for _, jsExtension := range jsExtensions {
		if c.requiredFilesAvailable(jsExtension, bfsFiles) {
			extension := Extension{Id: jsExtension.Id, Name: jsExtension.Name, Description: jsExtension.Description, InstallableVersions: jsExtension.InstallableVersions}
			extensions = append(extensions, &extension)
		}
	}
	return extensions, nil
}

func (c *controllerImpl) requiredFilesAvailable(extension *extensionAPI.JsExtension, bfsFiles []BfsFile) bool {
	for _, requiredFile := range extension.BucketFsUploads {
		if !existsFileInBfs(bfsFiles, requiredFile) {
			log.Printf("Ignoring extension %q since the required file %q (%q) does not exist or has a wrong file size.\n", extension.Name, requiredFile.Name, requiredFile.BucketFsFilename)
			return false
		}
	}
	return true
}

func (c *controllerImpl) getAllExtensions() []*extensionAPI.JsExtension {
	var extensions []*extensionAPI.JsExtension
	extensionPaths := FindJSFilesInDir(c.extensionFolder)
	for _, path := range extensionPaths {
		extension, err := c.loadExtensionFromFile(path)
		if err == nil {
			extensions = append(extensions, extension)
		} else {
			log.Printf("error: Failed to load extension. This extension will be ignored. Cause: %v\n", err)
		}
	}
	return extensions
}

func (c *controllerImpl) loadExtensionById(id string) (*extensionAPI.JsExtension, error) {
	extensionPath := path.Join(c.extensionFolder, id)
	return c.loadExtensionFromFile(extensionPath)
}

func (c *controllerImpl) loadExtensionFromFile(extensionPath string) (*extensionAPI.JsExtension, error) {
	extension, err := extensionAPI.GetExtensionFromFile(extensionPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load extension from file %q: %w", extensionPath, err)
	}
	return extension, nil
}

func (c *controllerImpl) GetAllInstallations(tx *sql.Tx) ([]*extensionAPI.JsExtInstallation, error) {
	metadata, err := c.metaDataReader.ReadMetadataTables(tx, c.schema)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata tables. Cause: %w", err)
	}
	extensions := c.getAllExtensions()
	extensionContext := c.createExtensionContext(tx)
	var allInstallations []*extensionAPI.JsExtInstallation
	for _, extension := range extensions {
		installations, err := extension.FindInstallations(extensionContext, metadata)
		if err != nil {
			return nil, apiErrors.NewAPIErrorWithCause(fmt.Sprintf("failed to find installations for extension %q", extension.Name), err)
		} else {
			allInstallations = append(allInstallations, installations...)
		}
	}
	return allInstallations, nil
}

func (c *controllerImpl) InstallExtension(tx *sql.Tx, extensionId string, extensionVersion string) error {
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

func (c *controllerImpl) CreateInstance(tx *sql.Tx, extensionId string, extensionVersion string, parameterValues []ParameterValue) (*extensionAPI.JsExtInstance, error) {
	extension, err := c.loadExtensionById(extensionId)
	if err != nil {
		return nil, fmt.Errorf("failed to load extension with id %q: %w", extensionId, err)
	}
	err = c.ensureSchemaExists(tx)
	if err != nil {
		return nil, err
	}
	params := extensionAPI.ParameterValues{}
	for _, p := range parameterValues {
		params.Values = append(params.Values, extensionAPI.ParameterValue{Name: p.Name, Value: p.Value})
	}

	extensionContext := c.createExtensionContext(tx)
	installation, err := c.findInstallationByVersion(tx, extensionContext, extension, extensionVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to find installations: %w", err)
	}

	err = validateParameters(installation.InstanceParameters, params)
	if err != nil {
		return nil, err
	}

	instance, err := extension.AddInstance(extensionContext, extensionVersion, &params)
	if err != nil {
		return nil, err
	}
	if instance == nil {
		return nil, fmt.Errorf("extension did not return an instance")
	}
	return instance, nil
}

func (c *controllerImpl) DeleteInstance(tx *sql.Tx, extensionId string, instanceId string) error {
	log.Printf("ctrl delete instance %s %s\n", extensionId, instanceId)
	extension, err := c.loadExtensionById(extensionId)
	if err != nil {
		return fmt.Errorf("failed to load extension with id %q: %w", extensionId, err)
	}
	return extension.DeleteInstance(c.createExtensionContext(tx), extensionId)
}

func (c *controllerImpl) FindInstances(tx *sql.Tx, extensionId string, extensionVersion string) ([]*extensionAPI.JsExtInstance, error) {
	extension, err := c.loadExtensionById(extensionId)
	if err != nil {
		return nil, fmt.Errorf("failed to load extension with id %q: %w", extensionId, err)
	}
	return extension.ListInstances(c.createExtensionContext(tx), extensionVersion)
}

func (c *controllerImpl) findInstallationByVersion(tx *sql.Tx, context *extensionAPI.ExtensionContext, extension *extensionAPI.JsExtension, version string) (*extensionAPI.JsExtInstallation, error) {
	metadata, err := c.metaDataReader.ReadMetadataTables(tx, c.schema)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata tables. Cause: %w", err)
	}

	installations, err := extension.FindInstallations(context, metadata)
	if err != nil {
		return nil, apiErrors.NewAPIErrorWithCause(fmt.Sprintf("failed to find installations for extension %q", extension.Name), err)
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

func (c *controllerImpl) createExtensionContext(tx *sql.Tx) *extensionAPI.ExtensionContext {
	return extensionAPI.CreateContext(context.TODO(), c.schema, tx)
}

func (c *controllerImpl) ensureSchemaExists(tx *sql.Tx) error {
	_, err := tx.Exec(fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS "%s"`, c.schema))
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}
	return nil
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
		return apiErrors.NewBadRequestErrorF("invalid parameters: %s", message)
	}
	return nil
}
