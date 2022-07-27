package main

import (
	"flag"
	"log"

	"github.com/exasol/extension-manager/restAPI"

	"github.com/exasol/extension-manager/extensionController"
)

//go:generate swag init -g restAPI/restApi.go -o generatedApiDocs
//go:generate npm --prefix parameterValidator ci
//go:generate npm --prefix parameterValidator run build

const EXTENSION_SCHEMA_NAME = "EXA_EXTENSIONS"

func main() {
	var pathToExtensionFolder = flag.String("pathToExtensionFolder", "../extensionApi/extensionForTesting/", "Path to folder containing extensions as .js files")
	var serverAddress = flag.String("serverAddress", ":8080", `Server address, e.g. ":8080" (all network interfaces) or "localhost:8080" (only local interface)`)
	flag.Parse()
	startServer(*pathToExtensionFolder, *serverAddress)
}

func startServer(pathToExtensionFolder string, serverAddress string) {
	log.Printf("Starting extension manager with extension folder %q", pathToExtensionFolder)
	restApi := restAPI.Create(extensionController.Create(pathToExtensionFolder, EXTENSION_SCHEMA_NAME), serverAddress)
	restApi.Serve()
}
