package registry

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/exasol/extension-manager/pkg/apiErrors"
	"github.com/exasol/extension-manager/pkg/extensionController/registry/index"
	log "github.com/sirupsen/logrus"
)

func newHttpRegistry(url string) Registry {
	log.Debugf("Creating HTTP registry for %q", url)
	return &httpRegistry{url: url, index: nil}
}

type httpRegistry struct {
	url   string
	index *index.RegistryIndex
}

/* [impl -> dsn~extension-registry~1] */
/* [impl -> dsn~extension-definitions-storage~1]. */
func (h *httpRegistry) FindExtensions() ([]string, error) {
	index, err := h.getIndex()
	if err != nil {
		return nil, err
	}
	return index.GetExtensionIDs(), nil
}

/* [impl -> dsn~extension-registry.cache~1]. */
func (h *httpRegistry) getIndex() (*index.RegistryIndex, error) {
	if h.index == nil {
		index, err := loadIndex(h.url)
		if err != nil {
			return nil, err
		}
		h.index = index
	}
	return h.index, nil
}

func loadIndex(url string) (*index.RegistryIndex, error) {
	t0 := time.Now()
	response, err := getResponse(url)
	if err != nil {
		return nil, fmt.Errorf("failed to load index from %q: %w", url, err)
	}
	defer response.Body.Close()
	index, err := index.Decode(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode index from %q: %w", url, err)
	}
	log.Debugf("Loaded registry index with %d extensions from %q in %dms", len(index.Extensions), url, time.Since(t0).Milliseconds())
	return &index, nil
}

func getResponse(url string) (*http.Response, error) {
	request, err := http.NewRequestWithContext(context.Background(), "GET", url, strings.NewReader(""))
	if err != nil {
		return nil, err
	}
	response, err := http.DefaultClient.Do(request)
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
	index, err := h.getIndex()
	if err != nil {
		return "", err
	}
	ext, ok := index.GetExtension(id)
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
	defer response.Body.Close()
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}
	return string(bytes), nil
}
