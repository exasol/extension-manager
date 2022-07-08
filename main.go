package main

import (
	"flag"
	"log"

	"github.com/exasol/extension-manager/restAPI"

	"github.com/exasol/extension-manager/extensionController"
)

//go:generate swag init -g restAPI/restApi.go -o generatedApiDocs

const EXTENSION_SCHEMA_NAME = "EXA_EXTENSIONS"

func main() {
	var pathToExtensionFolder = flag.String("pathToExtensionFolder", "../extensionApi/extensionForTesting/", "Path to folder containing extensions as .js files")
	flag.Parse()
	startServer(*pathToExtensionFolder)
}

func startServer(pathToExtensionFolder string) {
	log.Printf("Starting extension manager with extension folder %q", pathToExtensionFolder)
	restApi := restAPI.Create(extensionController.Create(pathToExtensionFolder, EXTENSION_SCHEMA_NAME))
	restApi.Serve()
}
