package integrationTesting

import (
	"bytes"
	"fmt"
	"os/exec"
	"path"
)

var isBuilt = false

// GetExtensionForTesting builds a dummy extension and returns its path
// The parameter name must point to the root of this go project, so that this method can resolve relative paths.
func GetExtensionForTesting(pathToProjectRoot string) string {
	extensionForTestingDir := path.Join(pathToProjectRoot, "integrationTesting", "extensionForTesting")
	if !isBuilt {
		isBuilt = true
		installCommand := exec.Command("npm", "ci")
		installCommand.Dir = extensionForTestingDir
		err := installCommand.Run()
		if err != nil {
			panic(fmt.Sprintf("Failed to install node modules (run 'npm ci') for extensionForTesting. Cause: %v", err.Error()))
		}
		buildCommand := exec.Command("npm", "run", "build")
		buildCommand.Dir = extensionForTestingDir
		err = buildCommand.Run()
		if err != nil {
			var stderr bytes.Buffer
			buildCommand.Stderr = &stderr
			fmt.Println(stderr.String())
			panic(fmt.Sprintf("Failed to build extensionForTesting. Cause: %v", err.Error()))
		}
	}
	return path.Join(extensionForTestingDir, "dist.js")
}
