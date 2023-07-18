package extensionController

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/exasol/extension-manager/pkg/apiErrors"
	"github.com/exasol/extension-manager/pkg/extensionAPI"
	"github.com/exasol/extension-manager/pkg/extensionController/registry"
	"github.com/exasol/extension-manager/pkg/parameterValidator"
)

// controller is the core part of the extension-manager that provides the extension handling functionality.
type controller interface {
	// GetAllExtensions reports all extension definitions.
	GetAllExtensions(bfsFiles []BfsFile) ([]*Extension, error)

	// GetAllInstallations searches for installations of any extensions.
	GetAllInstallations(tx *sql.Tx) ([]*extensionAPI.JsExtInstallation, error)

	// GetParameterDefinitions returns the parameter definitions required for installing a given extension version.
	GetParameterDefinitions(tx *sql.Tx, extensionId string, extensionVersion string) ([]parameterValidator.ParameterDefinition, error)

	// InstallExtension installs an extension.
	// extensionId is the ID of the extension to install
	// extensionVersion is the version of the extension to install
	InstallExtension(tx *sql.Tx, extensionId string, extensionVersion string) error

	// UninstallExtension removes an extension.
	// extensionId is the ID of the extension to uninstall
	// extensionVersion is the version of the extension to uninstall
	UninstallExtension(tx *sql.Tx, extensionId string, extensionVersion string) error

	// CreateInstance creates a new instance of an extension, e.g. a virtual schema and returns it's name.
	CreateInstance(tx *sql.Tx, extensionId string, extensionVersion string, parameterValues []ParameterValue) (*extensionAPI.JsExtInstance, error)

	// FindInstances returns a list of all instances for the given version.
	FindInstances(tx *sql.Tx, extensionId string, extensionVersion string) ([]*extensionAPI.JsExtInstance, error)

	// DeleteInstance deletes instance with the given ID.
	DeleteInstance(tx *sql.Tx, extensionId, extensionVersion, instanceId string) error
}

type controllerImpl struct {
	registry       registry.Registry
	schema         string
	metaDataReader extensionAPI.ExaMetadataReader
}

func createImpl(config ExtensionManagerConfig) controller {
	return &controllerImpl{
		registry:       registry.NewRegistry(config.ExtensionRegistryURL),
		metaDataReader: extensionAPI.CreateExaMetaDataReader(),
		schema:         config.ExtensionSchema,
	}
}

/* [impl -> dsn~list-extensions~1]. */
func (c *controllerImpl) GetAllExtensions(bfsFiles []BfsFile) ([]*Extension, error) {
	jsExtensions, err := c.getAllExtensions()
	if err != nil {
		return nil, err
	}
	var extensions []*Extension
	for _, jsExtension := range jsExtensions {
		if c.requiredFilesAvailable(jsExtension, bfsFiles) {
			extensions = append(extensions, convertExtension(jsExtension))
		}
	}
	return extensions, nil
}

func convertExtension(jsExtension *extensionAPI.JsExtension) *Extension {
	return &Extension{
		Id:                  jsExtension.Id,
		Name:                jsExtension.Name,
		Category:            jsExtension.Category,
		Description:         jsExtension.Description,
		InstallableVersions: jsExtension.InstallableVersions}
}

func (c *controllerImpl) requiredFilesAvailable(extension *extensionAPI.JsExtension, bfsFiles []BfsFile) bool {
	for _, requiredFile := range extension.BucketFsUploads {
		if !existsFileInBfs(bfsFiles, requiredFile) {
			log.Debugf("Ignoring extension %q since the required file %q (%q) does not exist or has a wrong file size.\n", extension.Name, requiredFile.Name, requiredFile.BucketFsFilename)
			return false
		}
	}
	return true
}

func existsFileInBfs(bfsFiles []BfsFile, requiredFile extensionAPI.BucketFsUpload) bool {
	for _, existingFile := range bfsFiles {
		if requiredFile.BucketFsFilename == existingFile.Name && requiredFile.FileSize == existingFile.Size {
			return true
		}
		if requiredFile.BucketFsFilename == existingFile.Name {
			log.Debugf("File %q exists but has wrong size %d, expected %d bytes", existingFile.Name, existingFile.Size, requiredFile.FileSize)
		}
	}
	return false
}

func (c *controllerImpl) getAllExtensions() ([]*extensionAPI.JsExtension, error) {
	extensionIds, err := c.registry.FindExtensions()
	if err != nil {
		return nil, err
	}
	extensions := make([]*extensionAPI.JsExtension, 0, len(extensionIds))
	for _, id := range extensionIds {
		extension, err := c.loadExtensionById(id)
		if err != nil {
			return nil, fmt.Errorf("failed to load extension %q: %w", id, err)
		}
		extensions = append(extensions, extension)
	}
	return extensions, nil
}

func (c *controllerImpl) loadExtensionById(id string) (*extensionAPI.JsExtension, error) {
	content, err := c.registry.ReadExtension(id)
	if err != nil {
		return nil, err
	}
	extension, err := extensionAPI.LoadExtension(id, content)
	if err != nil {
		return nil, err
	}
	return extension, nil
}

func (c *controllerImpl) GetAllInstallations(tx *sql.Tx) ([]*extensionAPI.JsExtInstallation, error) {
	metadata, err := c.metaDataReader.ReadMetadataTables(tx, c.schema)
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
			return nil, apiErrors.NewAPIErrorWithCause(fmt.Sprintf("failed to find installations for extension %q", extension.Name), err)
		} else {
			allInstallations = append(allInstallations, installations...)
		}
	}
	return allInstallations, nil
}

func (c *controllerImpl) GetParameterDefinitions(tx *sql.Tx, extensionId string, extensionVersion string) ([]parameterValidator.ParameterDefinition, error) {
	extension, err := c.loadExtensionById(extensionId)
	if err != nil {
		return nil, extensionLoadingFailed(extensionId, err)
	}
	rawDefinitions, err := extension.GetParameterDefinitions(c.createExtensionContext(tx), extensionVersion)
	if err != nil {
		return nil, err
	}
	definitions, err := parameterValidator.ConvertDefinitions(rawDefinitions)
	if err != nil {
		return nil, err
	}
	return definitions, nil
}

func (c *controllerImpl) InstallExtension(tx *sql.Tx, extensionId string, extensionVersion string) error {
	extension, err := c.loadExtensionById(extensionId)
	if err != nil {
		return extensionLoadingFailed(extensionId, err)
	}
	err = c.ensureSchemaExists(tx)
	if err != nil {
		return err
	}
	return extension.Install(c.createExtensionContext(tx), extensionVersion)
}

func (c *controllerImpl) UninstallExtension(tx *sql.Tx, extensionId string, extensionVersion string) error {
	extension, err := c.loadExtensionById(extensionId)
	if err != nil {
		return extensionLoadingFailed(extensionId, err)
	}
	return extension.Uninstall(c.createExtensionContext(tx), extensionVersion)
}

func (c *controllerImpl) CreateInstance(tx *sql.Tx, extensionId string, extensionVersion string, parameterValues []ParameterValue) (*extensionAPI.JsExtInstance, error) {
	extension, err := c.loadExtensionById(extensionId)
	if err != nil {
		return nil, extensionLoadingFailed(extensionId, err)
	}
	err = c.ensureSchemaExists(tx)
	if err != nil {
		return nil, err
	}
	params := extensionAPI.ParameterValues{}
	for _, p := range parameterValues {
		params.Values = append(params.Values, extensionAPI.ParameterValue{Name: p.Name, Value: p.Value})
	}

	paramDefinitions, err := c.GetParameterDefinitions(tx, extensionId, extensionVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get parameter definitions: %w", err)
	}

	err = validateParameters(paramDefinitions, params)
	if err != nil {
		return nil, err
	}

	extensionContext := c.createExtensionContext(tx)
	instance, err := extension.AddInstance(extensionContext, extensionVersion, &params)
	if err != nil {
		return nil, err
	}
	if instance == nil {
		return nil, fmt.Errorf("extension %q did not return an instance", extensionId)
	}
	return instance, nil
}

func (c *controllerImpl) DeleteInstance(tx *sql.Tx, extensionId, extensionVersion, instanceId string) error {
	extension, err := c.loadExtensionById(extensionId)
	if err != nil {
		return extensionLoadingFailed(extensionId, err)
	}
	return extension.DeleteInstance(c.createExtensionContext(tx), extensionVersion, instanceId)
}

func (c *controllerImpl) FindInstances(tx *sql.Tx, extensionId string, extensionVersion string) ([]*extensionAPI.JsExtInstance, error) {
	extension, err := c.loadExtensionById(extensionId)
	if err != nil {
		return nil, extensionLoadingFailed(extensionId, err)
	}
	return extension.ListInstances(c.createExtensionContext(tx), extensionVersion)
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

/* [impl -> dsn~parameter-types~1]. */
func validateParameters(parameterDefinitions []parameterValidator.ParameterDefinition, params extensionAPI.ParameterValues) error {
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

func extensionLoadingFailed(extensionId string, err error) error {
	return fmt.Errorf("failed to load extension %q: %w", extensionId, err)
}
