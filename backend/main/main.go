package main

import (
	"backend/extensionController"
	"backend/respApi"
)

func main() {
	restApi := respApi.Create(extensionController.Create("../extensionApi/extensionForTesting/"))
	restApi.Serve()
}
