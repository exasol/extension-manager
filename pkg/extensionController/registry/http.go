package registry

import (
	"fmt"
	"io"
	"net/http"

	"github.com/exasol/extension-manager/pkg/apiErrors"
	"github.com/exasol/extension-manager/pkg/extensionController/registry/index"
)

func newHttpRegistry(url string) Registry {
	return &httpRegistry{url: url}
}

type httpRegistry struct {
	url   string
	index *index.RegistryIndex
}

func (h *httpRegistry) FindExtensions() ([]string, error) {
	err := h.loadIndex()
	if err != nil {
		return nil, err
	}
	return h.index.GetExtensionIDs(), nil
}

func (h *httpRegistry) loadIndex() error {
	response, err := getResponse(h.url)
	if err != nil {
		return err
	}
	index, err := index.Decode(response.Body)
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

func (h *httpRegistry) ReadExtension(id string) (string, error) {
	err := h.loadIndex()
	if err != nil {
		return "", err
	}
	ext, ok := h.index.GetExtension(id)
	if !ok {
		return "", apiErrors.NewNotFoundErrorF("extension %q not found", id)
	}
	extContent, err := getUrlContent(ext.URL)
	if err != nil {
		return "", fmt.Errorf("failed to load extension %q: %w", id, err)
	}
	return extContent, nil
}

func getUrlContent(url string) (string, error) {
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
