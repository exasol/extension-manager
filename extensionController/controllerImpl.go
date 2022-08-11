package extensionController

import (
	"database/sql"
	"fmt"
	"log"
	"path"

	"github.com/exasol/extension-manager/extensionAPI"
)

// ExtensionController is the core part of the extension-manager that provides the extension handling functionality.
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
	CreateInstance(tx *sql.Tx, extensionId string, extensionVersion string, parameterValues []ParameterValue) (string, error)
}

type controllerImpl struct {
	pathToExtensionFolder string
	extensionSchemaName   string
}

func createImpl(pathToExtensionFolder string, extensionSchemaName string) controller {
	return &controllerImpl{pathToExtensionFolder: pathToExtensionFolder, extensionSchemaName: extensionSchemaName}
}

func (c *controllerImpl) GetAllExtensions(bfsFiles []BfsFile) ([]*Extension, error) {
	jsExtensions, err := c.getAllExtensions()
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

func (c *controllerImpl) requiredFilesAvailable(extension *extensionAPI.JsExtension, bfsFiles []BfsFile) bool {
	for _, requiredFile := range extension.BucketFsUploads {
		if !existsFileInBfs(bfsFiles, requiredFile) {
			log.Printf("Ignoring extension %q since the required file %q does not exist or has a wrong file size.\n", extension.Name, requiredFile.Name)
			return false
		}
	}
	return true
}

func (c *controllerImpl) getAllExtensions() ([]*extensionAPI.JsExtension, error) {
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

func (c *controllerImpl) loadExtensionById(id string) (*extensionAPI.JsExtension, error) {
	extensionPath := path.Join(c.pathToExtensionFolder, id)
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

func (c *controllerImpl) CreateInstance(tx *sql.Tx, extensionId string, extensionVersion string, parameterValues []ParameterValue) (string, error) {
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

func (c *controllerImpl) findInstallationByVersion(tx *sql.Tx, context *extensionAPI.ExtensionContext, extension *extensionAPI.JsExtension, version string) (*extensionAPI.JsExtInstallation, error) {
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

func (c *controllerImpl) createExtensionContext(tx *sql.Tx) *extensionAPI.ExtensionContext {
	return extensionAPI.CreateContext(c.extensionSchemaName, tx)
}

func (c *controllerImpl) ensureSchemaExists(tx *sql.Tx) error {
	_, err := tx.Exec(fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS "%s"`, c.extensionSchemaName))
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}
	return nil
}
