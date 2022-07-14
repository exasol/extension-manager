package extensionController

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FindJSFilesInDir searches for .js files in the given directory pathToExtensionFolder
func FindJSFilesInDir(pathToExtensionFolder string) []string {
	var files []string
	err := filepath.Walk(pathToExtensionFolder, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".js") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(fmt.Sprintf("failed to load extensions from extension folder %v. Cause %v", pathToExtensionFolder, err))
	}
	return files
}
