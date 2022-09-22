package main

import (
	"flag"

	log "github.com/sirupsen/logrus"

	"github.com/exasol/extension-manager/restAPI"

	"github.com/exasol/extension-manager/extensionController"
)

func main() {
	var extensionRegistryURL = flag.String("extensionRegistryURL", "", "URL of the extension registry index used to find available extensions or the path of a local directory")
	var serverAddress = flag.String("serverAddress", ":8080", `Server address, e.g. ":8080" (all network interfaces) or "localhost:8080" (only local interface)`)
	flag.Parse()
	log.SetLevel(log.DebugLevel)
	startServer(*extensionRegistryURL, *serverAddress)
}

func startServer(pathToExtensionFolder string, serverAddress string) {
	log.Printf("Starting extension manager with extension folder %q", pathToExtensionFolder)
	restApi := restAPI.Create(extensionController.Create(pathToExtensionFolder, restAPI.EXTENSION_SCHEMA_NAME), serverAddress)
	restApi.Serve()
}
