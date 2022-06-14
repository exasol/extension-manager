package main

import (
	"github.com/exasol/extension-manager/restAPI"

	"github.com/exasol/extension-manager/extensionController"
)

func main() {
	restApi := restAPI.Create(extensionController.Create("../extensionApi/extensionForTesting/"))
	restApi.Serve()
}
