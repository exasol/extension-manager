package main

import (
	"github.com/exasol/extension-manager/restAPI"

	"github.com/exasol/extension-manager/extensionController"
)

//go:generate swag init -g restAPI/restApi.go -o generatedApiDocs

func main() {
	restApi := restAPI.Create(extensionController.Create("../extensionApi/extensionForTesting/"))
	restApi.Serve()
}
