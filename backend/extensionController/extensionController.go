package extensionController

import (
	"backend"
	"backend/extensionApi"
	"database/sql"
	"fmt"
	"path"
)

type ExtensionController interface {
	GetAllInstallations(dbConnection *sql.DB) ([]*extensionApi.Installation, error)
	GetAllExtensions() ([]*extensionApi.Extension, error)
}

func Create(pathToExtensionFolder string) ExtensionController {
	return &extensionControllerImpl{pathToExtensionFolder: pathToExtensionFolder}
}

type extensionControllerImpl struct {
	pathToExtensionFolder string
}

func (controller *extensionControllerImpl) GetAllExtensions() ([]*extensionApi.Extension, error) {
	extensionPath := path.Join(controller.pathToExtensionFolder, "dist.js")
	extension, err := extensionApi.GetExtensionFromFile(extensionPath)
	if err != nil {
		return nil, err
	}
	return []*extensionApi.Extension{extension}, nil
}

func (controller *extensionControllerImpl) GetAllInstallations(dbConnection *sql.DB) ([]*extensionApi.Installation, error) {
	allScriptTable, err := extensionApi.ReadExaAllScriptTable(dbConnection)
	if err != nil {
		return nil, fmt.Errorf("failed to read EXA_ALL_SCRIPT table. Cause: %w", err)
	}
	sqlClient := backend.ExasolSqlClient{Connection: dbConnection}
	extensions, err := controller.GetAllExtensions()
	if err != nil {
		return nil, err
	}
	var allInstallations []*extensionApi.Installation
	for _, extension := range extensions {
		installations := extension.FindInstallations(sqlClient, allScriptTable)
		allInstallations = append(allInstallations, installations...)
	}
	return allInstallations, nil
}
