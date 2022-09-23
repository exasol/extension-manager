package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	log "github.com/sirupsen/logrus"

	"github.com/exasol/extension-manager/restAPI"

	"github.com/exasol/extension-manager/extensionController"
)

func main() {
	var extensionRegistryURL = flag.String("extensionRegistryURL", "", "URL of the extension registry index used to find available extensions or the path of a local directory")
	var serverAddress = flag.String("serverAddress", ":8080", `Server address, e.g. ":8080" (all network interfaces) or "localhost:8080" (only local interface)`)
	var openAPIOutputPath = flag.String("openAPIOutputPath", "", "Generate the OpenAPI spec at the given path instead of starting the server")
	flag.Parse()
	log.SetLevel(log.DebugLevel)
	if openAPIOutputPath != nil {
		err := generateOpenAPISpec(*openAPIOutputPath)
		if err != nil {
			panic(fmt.Sprintf("failed to generate OpenAPI: %v", err))
		}
	} else {
		startServer(*extensionRegistryURL, *serverAddress)
	}
}

func startServer(pathToExtensionFolder string, serverAddress string) {
	log.Printf("Starting extension manager with extension folder %q", pathToExtensionFolder)
	restApi := restAPI.Create(extensionController.Create(pathToExtensionFolder, restAPI.EXTENSION_SCHEMA_NAME), serverAddress)
	restApi.Serve()
}

func generateOpenAPISpec(filename string) error {
	json, err := generateOpenAPIJson()
	if err != nil {
		return err
	}
	return writeFile(filename, json)
}

func writeFile(filename string, content []byte) error {
	dir := path.Dir(filename)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, content, 0600)
	if err != nil {
		return err
	}
	fmt.Printf("Wrote OpenAPI spec to %s", filename)
	return nil
}

func generateOpenAPIJson() ([]byte, error) {
	api, err := restAPI.CreateOpenApi()
	if err != nil {
		return nil, err
	}
	err = restAPI.AddPublicEndpoints(api, restAPI.ExtensionManagerConfig{})
	if err != nil {
		return nil, err
	}
	json, err := api.ToJSON()
	if err != nil {
		return nil, err
	}
	return json, err
}
