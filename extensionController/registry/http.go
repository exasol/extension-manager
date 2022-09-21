package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func newHttpRegistry(url string) Registry {
	return &httpRegistry{url: url}
}

type httpRegistry struct {
	url   string
	index *RegistryIndex
}

func (h *httpRegistry) FindExtensions() ([]string, error) {
	err := h.loadIndex()
	if err != nil {
		return nil, err
	}
	return h.index.getExtensionIDs(), nil
}

func (h *httpRegistry) loadIndex() error {
	response, err := getResponse(h.url)
	if err != nil {
		return err
	}
	index, err := decodeRegistryIndex(response.Body)
	if err != nil {
		return err
	}
	h.index = &index
	return nil
}

func getResponse(url string) (*http.Response, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		bytes, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("registry at %s returned status %q and response %q", url, response.Status, bytes)
	}
	return response, nil
}

func decodeRegistryIndex(reader io.Reader) (RegistryIndex, error) {
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()
	content := RegistryIndex{}
	err := decoder.Decode(&content)
	if err != nil {
		return content, fmt.Errorf("failed to decode registry content: %w", err)
	}
	return content, nil
}

func (h *httpRegistry) ReadExtension(id string) (string, error) {
	return "", fmt.Errorf("unimplemented")
}

func getUrl(url string) (string, error) {
	response, err := getResponse(url)
	if err != nil {
		return "", err
	}
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}
	return string(bytes), nil
}
