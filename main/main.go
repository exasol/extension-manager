package main

import (
	"extension-manager/extensionController"
	"extension-manager/restAPI"
)

func main() {
	restApi := restAPI.Create(extensionController.Create("../extensionApi/extensionForTesting/"))
	restApi.Serve()
}
