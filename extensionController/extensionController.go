package extensionController

import (
	"database/sql"
	"fmt"
	"log"
	"path"

	"github.com/exasol/extension-manager/backend"
	"github.com/exasol/extension-manager/extensionAPI"
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
}

type Extension struct {
	Id                  string
	Name                string
	Description         string
	InstallableVersions []string
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
			log.Printf("error: Failed to load extension. This extension will be ignored. Cause: %v\n", err.Error())
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
		return nil, fmt.Errorf("failed to load extension from file %q: %v", extensionPath, err.Error())
	}
	_, fileName := path.Split(extensionPath)
	extension.Id = fileName
	return extension, nil
}

func (controller *extensionControllerImpl) InstallExtension(dbConnection *sql.DB, extensionId string, extensionVersion string) error {
	extension, err := controller.getJsExtensionById(extensionId)
	if err != nil {
		return fmt.Errorf("failed to load extension with id %q: %v", extensionId, err)
	}
	sqlClient := backend.ExasolSqlClient{Connection: dbConnection}
	return extension.Install(sqlClient, extensionVersion)
}

func (controller *extensionControllerImpl) GetAllInstallations(dbConnection *sql.DB) ([]*extensionAPI.JsExtInstallation, error) {
	metadata, err := extensionAPI.ReadMetadataTables(dbConnection, controller.extensionSchemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata tables. Cause: %w", err)
	}
	sqlClient := backend.ExasolSqlClient{Connection: dbConnection}
	extensions, err := controller.getAllJsExtensions()
	if err != nil {
		return nil, err
	}
	var allInstallations []*extensionAPI.JsExtInstallation
	for _, extension := range extensions {
		installations, err := extension.FindInstallations(sqlClient, metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to find installations from extension %q: %v", extension.Id, err)
		} else {
			allInstallations = append(allInstallations, installations...)
		}
	}
	return allInstallations, nil
}
