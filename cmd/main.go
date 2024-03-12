package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/exasol/extension-manager/pkg/restAPI"

	"github.com/exasol/extension-manager/pkg/extensionController"
)

func main() {
	var extensionRegistryURL = flag.String("extensionRegistryURL", "", "URL of the extension registry index used to find available extensions or the path of a local directory")
	var serverAddress = flag.String("serverAddress", ":8080", `Server address, e.g. ":8080" (all network interfaces) or "localhost:8080" (only local interface)`)
	var openAPIOutputPath = flag.String("openAPIOutputPath", "", "Generate the OpenAPI spec at the given path instead of starting the server")
	var addCauseToInternalServerError = flag.Bool("addCauseToInternalServerError", false, "Add cause of internal server errors (status 500) to the error message. Don't use this in production!")
	flag.Parse()
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&simpleFormatter{})
	if openAPIOutputPath != nil && *openAPIOutputPath != "" {
		err := generateOpenAPISpec(*openAPIOutputPath)
		if err != nil {
			fmt.Printf("failed to generate OpenAPI to %q: %v\n", *openAPIOutputPath, err)
			os.Exit(1)
		}
	} else {
		err := startServer(*extensionRegistryURL, *serverAddress, *addCauseToInternalServerError)
		if err != nil {
			fmt.Printf("failed to start server: %v\n", err)
			os.Exit(1)
		}
	}
}

func startServer(pathToExtensionFolder string, serverAddress string, addCauseToInternalServerError bool) error {
	if pathToExtensionFolder == "" {
		return errors.New("please specify extension registry with parameter '-extensionRegistryURL'")
	}
	log.Printf("Starting extension manager with extension folder %q", pathToExtensionFolder)
	controller, err := extensionController.CreateWithValidatedConfig(extensionController.ExtensionManagerConfig{
		ExtensionRegistryURL: pathToExtensionFolder,
		ExtensionSchema:      restAPI.EXTENSION_SCHEMA_NAME,
		BucketFSBasePath:     "/buckets/bfsdefault/default/"})
	if err != nil {
		return err
	}
	restApi := restAPI.Create(controller, serverAddress, addCauseToInternalServerError)
	restApi.Serve()
	return nil
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
	fmt.Printf("Wrote OpenAPI spec to %s\n", filename)
	return nil
}

func generateOpenAPIJson() ([]byte, error) {
	api, err := restAPI.CreateOpenApi()
	if err != nil {
		return nil, err
	}
	dummyConfiguration := extensionController.ExtensionManagerConfig{ExtensionRegistryURL: "dummy", BucketFSBasePath: "dummy", ExtensionSchema: "dummy"}
	err = restAPI.AddPublicEndpoints(api, dummyConfiguration)
	if err != nil {
		return nil, err
	}
	json, err := api.ToJSON()
	if err != nil {
		return nil, err
	}
	return json, err
}

type simpleFormatter struct {
}

func (f *simpleFormatter) Format(entry *log.Entry) ([]byte, error) {
	b := &bytes.Buffer{}
	b.WriteString(strings.ToUpper(entry.Level.String()))
	b.WriteByte(' ')
	b.WriteString(entry.Message)
	b.WriteByte('\n')
	return b.Bytes(), nil
}
