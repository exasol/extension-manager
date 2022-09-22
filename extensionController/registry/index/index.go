package index

import (
	"encoding/json"
	"fmt"
	"io"
)

type RegistryIndex struct {
	Extensions []Extension `json:"extensions"`
}

type Extension struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// Decode parses the content of the given reader and returns a RegistryIndex.
func Decode(reader io.Reader) (RegistryIndex, error) {
	decoder := json.NewDecoder(reader)
	content := RegistryIndex{}
	err := decoder.Decode(&content)
	if err != nil {
		return content, fmt.Errorf("failed to decode registry content: %w", err)
	}
	return content, nil
}

// GetExtensionIDs returns the IDs of all extensions contained in this index.
func (i RegistryIndex) GetExtensionIDs() []string {
	ids := make([]string, 0, len(i.Extensions))
	for _, e := range i.Extensions {
		ids = append(ids, e.ID)
	}
	return ids
}

// GetExtension returns the extension with the given ID and true or and empty extension and false.
func (i RegistryIndex) GetExtension(id string) (extension Extension, ok bool) {
	for _, ext := range i.Extensions {
		if ext.ID == id {
			return ext, true
		}
	}
	return Extension{}, false
}
