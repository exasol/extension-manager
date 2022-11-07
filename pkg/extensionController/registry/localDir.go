package registry

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/exasol/extension-manager/pkg/apiErrors"
)

func newLocalDirRegistry(dir string) Registry {
	return &localDirRegistry{dir: dir}
}

type localDirRegistry struct {
	dir string
}

// FindExtensions searches for .js files in the local registry directory.
/* [impl -> dsn~extension-definitions-storage~1] */
func (l *localDirRegistry) FindExtensions() ([]string, error) {
	var files []string
	err := filepath.Walk(l.dir, func(path string, info os.FileInfo, err error) error {
		if info != nil && strings.HasSuffix(info.Name(), ".js") {
			files = append(files, info.Name())
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find extensions in %q: %w", l.dir, err)
	}
	return files, nil
}

func (l *localDirRegistry) ReadExtension(id string) (string, error) {
	fileName := path.Join(l.dir, id)
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", apiErrors.NewNotFoundErrorF("extension %q not found", fileName)
		}
		return "", fmt.Errorf("failed to open extension file %q: %w", fileName, err)
	}
	return string(bytes), nil
}
