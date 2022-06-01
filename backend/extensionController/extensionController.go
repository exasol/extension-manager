package extensionController

import (
	"backend"
	"backend/extensionAPI"
	"database/sql"
	"fmt"
)

// ExtensionController is the core part of the extension-manager that provides the extension handling functionality.
type ExtensionController interface {
	// GetAllExtensions reports all extension definitions
	GetAllExtensions(dbConnectionWithNoAutocommit *sql.DB) ([]*Extension, error)
	// GetAllInstallations searches for installations of any extensions
	GetAllInstallations(dbConnection *sql.DB) ([]*extensionAPI.JsExtInstallation, error)
}

type Extension struct {
	Name                string
	Description         string
	InstallableVersions []string
}

// ExtInstallation represents the installation of an Extension
type ExtInstallation struct {
}

// Create am instance of ExtensionController
func Create(pathToExtensionFolder string) ExtensionController {
	return &extensionControllerImpl{pathToExtensionFolder: pathToExtensionFolder}
}

type extensionControllerImpl struct {
	pathToExtensionFolder string
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
		if controller.checkRequiredFiles(jsExtension, bfsFiles) {
			extension := Extension{Name: jsExtension.Name, Description: jsExtension.Description, InstallableVersions: jsExtension.InstallableVersions}
			extensions = append(extensions, &extension)
		}
	}
	return extensions, nil
}

func (controller *extensionControllerImpl) checkRequiredFiles(jsExtension *extensionAPI.JsExtension, bfsFiles []BfsFile) bool {
	for _, requiredFile := range jsExtension.BucketFsUploads {
		if !controller.existsFileInBfs(bfsFiles, requiredFile) {
			fmt.Printf("ignoring extension %v since the required file %v does not exist or has a wrong file size.", jsExtension.Name, requiredFile.Name)
			return false
		}
	}
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
		extension, err := extensionAPI.GetExtensionFromFile(extensionPath)
		if err == nil {
			extensions = append(extensions, extension)
		} else {
			fmt.Printf("error: Failed to load extension form file %v. This extension will be ignored. Cause: %v", extensionPath, err.Error())
		}
	}
	return extensions, nil
}

func (controller *extensionControllerImpl) GetAllInstallations(dbConnection *sql.DB) ([]*extensionAPI.JsExtInstallation, error) {
	allScriptTable, err := extensionAPI.ReadExaAllScriptTable(dbConnection)
	if err != nil {
		return nil, fmt.Errorf("failed to read EXA_ALL_SCRIPT table. Cause: %w", err)
	}
	sqlClient := backend.ExasolSqlClient{Connection: dbConnection}
	extensions, err := controller.getAllJsExtensions()
	if err != nil {
		return nil, err
	}
	var allInstallations []*extensionAPI.JsExtInstallation
	for _, extension := range extensions {
		installations := extension.FindInstallations(sqlClient, allScriptTable)
		allInstallations = append(allInstallations, installations...)
	}
	return allInstallations, nil
}
