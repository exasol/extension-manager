package extensionController

import (
	"database/sql"
	"fmt"
	"log"
	"path"

	"github.com/exasol/extension-manager/extensionAPI"
	"github.com/exasol/extension-manager/parameterValidator"
)

// ExtensionController is the core part of the extension-manager that provides the extension handling functionality.
type ExtensionController interface {
	// GetAllInstallations searches for installations of any extensions.
	// dbConnection is a connection to the Exasol DB with autocommit turned off
	GetAllInstallations(dbConnection *sql.DB) ([]*extensionAPI.JsExtInstallation, error)

	// InstallExtension installs an extension.
	// dbConnection is a connection to the Exasol DB with autocommit turned off
	// extensionId is the ID of the extension to install
	// extensionVersion is the version of the extension to install
	InstallExtension(dbConnection *sql.DB, extensionId string, extensionVersion string) error

	// GetAllExtensions reports all extension definitions.
	// dbConnection is a connection to the Exasol DB with autocommit turned off
	GetAllExtensions(dbConnection *sql.DB) ([]*Extension, error)

	// CreateInstance creates a new instance of an extension, e.g. a virtual schema and returns it's name.
	// dbConnection is a connection to the Exasol DB with autocommit turned off
	CreateInstance(db *sql.DB, extensionId string, extensionVersion string, parameterValues []ParameterValue) (string, error)
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

func (controller *extensionControllerImpl) GetAllExtensions(dbConnectionWithNoAutocommit *sql.DB) ([]*Extension, error) {
	jsExtensions, err := controller.getAllJsExtensions()
	if err != nil {
		return nil, err
	}
	var extensions []*Extension
	for _, jsExtension := range jsExtensions {
		bfsAPI := CreateBucketFsAPI(dbConnectionWithNoAutocommit)
		bfsFiles, err := bfsAPI.ListFiles("default")
		if err != nil {
			return nil, fmt.Errorf("failed to search for required files in BucketFS. Cause: %w", err)
		}
		if controller.requiredFilesAvailable(jsExtension, bfsFiles) {
			extension := Extension{Id: jsExtension.Id, Name: jsExtension.Name, Description: jsExtension.Description, InstallableVersions: jsExtension.InstallableVersions}
			extensions = append(extensions, &extension)
		}
	}
	return extensions, nil
}

func (controller *extensionControllerImpl) requiredFilesAvailable(jsExtension *extensionAPI.JsExtension, bfsFiles []BfsFile) bool {
	for _, requiredFile := range jsExtension.BucketFsUploads {
		if !controller.existsFileInBfs(bfsFiles, requiredFile) {
			fmt.Printf("ignoring extension %q since the required file %q does not exist or has a wrong file size.\n", jsExtension.Name, requiredFile.Name)
			return false
		}
	}
	log.Printf("Required files found for extension %q\n", jsExtension.Name)
	return true
}

func (controller *extensionControllerImpl) existsFileInBfs(bfsFiles []BfsFile, requiredFile extensionAPI.BucketFsUpload) bool {
	for _, existingFile := range bfsFiles {
		if requiredFile.BucketFsFilename == existingFile.Name && requiredFile.FileSize == existingFile.Size {
			return true
		}
	}
	return false
}

func (controller *extensionControllerImpl) getAllJsExtensions() ([]*extensionAPI.JsExtension, error) {
	var extensions []*extensionAPI.JsExtension
	extensionPaths := FindJSFilesInDir(controller.pathToExtensionFolder)
	for _, extensionPath := range extensionPaths {
		extension, err := controller.getJsExtension(extensionPath)
		if err == nil {
			extensions = append(extensions, extension)
		} else {
			log.Printf("error: Failed to load extension. This extension will be ignored. Cause: %v\n", err)
		}
	}
	return extensions, nil
}

func (controller *extensionControllerImpl) getJsExtensionById(id string) (*extensionAPI.JsExtension, error) {
	extensionPath := path.Join(controller.pathToExtensionFolder, id)
	return controller.getJsExtension(extensionPath)
}

func (controller *extensionControllerImpl) getJsExtension(extensionPath string) (*extensionAPI.JsExtension, error) {
	extension, err := extensionAPI.GetExtensionFromFile(extensionPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load extension from file %q: %w", extensionPath, err)
	}
	return extension, nil
}

func (controller *extensionControllerImpl) InstallExtension(dbConnection *sql.DB, extensionId string, extensionVersion string) error {
	extension, err := controller.getJsExtensionById(extensionId)
	if err != nil {
		return fmt.Errorf("failed to load extension with id %q: %w", extensionId, err)
	}
	err = controller.ensureSchemaExists(dbConnection)
	if err != nil {
		return err
	}
	return extension.Install(controller.createContext(dbConnection), extensionVersion)
}

func (controller *extensionControllerImpl) GetAllInstallations(dbConnection *sql.DB) ([]*extensionAPI.JsExtInstallation, error) {
	metadata, err := extensionAPI.ReadMetadataTables(dbConnection, controller.extensionSchemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata tables. Cause: %w", err)
	}
	extensions, err := controller.getAllJsExtensions()
	if err != nil {
		return nil, err
	}
	context := controller.createContext(dbConnection)
	var allInstallations []*extensionAPI.JsExtInstallation
	for _, extension := range extensions {
		installations, err := extension.FindInstallations(context, metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to find installations: %v", err)
		} else {
			allInstallations = append(allInstallations, installations...)
		}
	}
	return allInstallations, nil
}

func (controller *extensionControllerImpl) CreateInstance(db *sql.DB, extensionId string, extensionVersion string, parameterValues []ParameterValue) (string, error) {
	extension, err := controller.getJsExtensionById(extensionId)
	if err != nil {
		return "", fmt.Errorf("failed to load extension with id %q: %w", extensionId, err)
	}
	err = controller.ensureSchemaExists(db)
	if err != nil {
		return "", err
	}
	params := extensionAPI.ParameterValues{}
	for _, p := range parameterValues {
		params.Values = append(params.Values, extensionAPI.ParameterValue{Name: p.Name, Value: p.Value})
	}

	context := controller.createContext(db)
	installation, err := controller.findInstallationByVersion(db, context, extension, extensionVersion)
	if err != nil {
		return "", fmt.Errorf("failed to find installations: %w", err)
	}

	err = validateParameters(installation.InstanceParameters, params)
	if err != nil {
		return "", err
	}

	instance, err := extension.AddInstance(controller.createContext(db), extensionVersion, &params)
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
	if len(result) > 0 {
		message := ""
		for i, r := range result {
			if i > 0 {
				message += ", "
			}
			message += r.Message
		}
		return fmt.Errorf("invalid parameters: %s", message)
	}
	return nil
}

func (controller *extensionControllerImpl) findInstallationByVersion(db *sql.DB, context *extensionAPI.ExtensionContext, extension *extensionAPI.JsExtension, version string) (*extensionAPI.JsExtInstallation, error) {
	metadata, err := extensionAPI.ReadMetadataTables(db, controller.extensionSchemaName)
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

func (controller *extensionControllerImpl) createContext(dbConnection *sql.DB) *extensionAPI.ExtensionContext {
	return extensionAPI.CreateContext(controller.extensionSchemaName, dbConnection)
}

func (controller *extensionControllerImpl) ensureSchemaExists(db *sql.DB) error {
	_, err := db.Exec(fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS "%s"`, controller.extensionSchemaName))
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}
	return nil
}
