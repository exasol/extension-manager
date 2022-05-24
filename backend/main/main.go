package main

import (
	"backend/extensionController"
	"backend/restAPI"
)

func main() {
	restApi := restAPI.Create(extensionController.Create("../extensionApi/extensionForTesting/"))
	restApi.Serve()
}
