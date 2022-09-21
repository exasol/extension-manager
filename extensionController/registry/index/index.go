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

func Decode(reader io.Reader) (RegistryIndex, error) {
	decoder := json.NewDecoder(reader)
	content := RegistryIndex{}
	err := decoder.Decode(&content)
	if err != nil {
		return content, fmt.Errorf("failed to decode registry content: %w", err)
	}
	return content, nil
}

func (i RegistryIndex) GetExtensionIDs() []string {
	ids := make([]string, 0, len(i.Extensions))
	for _, e := range i.Extensions {
		ids = append(ids, e.ID)
	}
	return ids
}
