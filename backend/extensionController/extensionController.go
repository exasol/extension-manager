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
	GetAllExtensions() ([]*extensionAPI.Extension, error)
	// GetAllInstallations searches for installations of any extensions
	GetAllInstallations(dbConnection *sql.DB) ([]*extensionAPI.Installation, error)
}

// Create am instance of ExtensionController
func Create(pathToExtensionFolder string) ExtensionController {
	return &extensionControllerImpl{pathToExtensionFolder: pathToExtensionFolder}
}

type extensionControllerImpl struct {
	pathToExtensionFolder string
}

func (controller *extensionControllerImpl) GetAllExtensions() ([]*extensionAPI.Extension, error) {
	var extensions []*extensionAPI.Extension
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

func (controller *extensionControllerImpl) GetAllInstallations(dbConnection *sql.DB) ([]*extensionAPI.Installation, error) {
	allScriptTable, err := extensionAPI.ReadExaAllScriptTable(dbConnection)
	if err != nil {
		return nil, fmt.Errorf("failed to read EXA_ALL_SCRIPT table. Cause: %w", err)
	}
	sqlClient := backend.ExasolSqlClient{Connection: dbConnection}
	extensions, err := controller.GetAllExtensions()
	if err != nil {
		return nil, err
	}
	var allInstallations []*extensionAPI.Installation
	for _, extension := range extensions {
		installations := extension.FindInstallations(sqlClient, allScriptTable)
		allInstallations = append(allInstallations, installations...)
	}
	return allInstallations, nil
}
