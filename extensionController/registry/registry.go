package registry

import (
	"strings"
)

// Registry allows listing and loading extension files
type Registry interface {
	// FindExtensions finds all available extensions and returns their IDs.
	FindExtensions() ([]string, error)

	// ReadExtension loads and returns the extension content as a string.
	ReadExtension(id string) (string, error)
}

// NewRegistry creates a new extension registry.
// The argument can be an HTTP(S) URL or the path of a local directory.
// This returns a matching registry implementation depending on the argument.
func NewRegistry(urlOrPath string) Registry {
	if isHttpUrl(urlOrPath) {
		return newHttpRegistry(urlOrPath)
	}
	return newLocalDirRegistry(urlOrPath)
}

func isHttpUrl(urlOrPath string) bool {
	lowerCaseUrlOrPath := strings.ToLower(urlOrPath)
	return strings.HasPrefix(lowerCaseUrlOrPath, "http://") || strings.HasPrefix(lowerCaseUrlOrPath, "https://")
}
