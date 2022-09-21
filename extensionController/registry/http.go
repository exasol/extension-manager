package registry

import "fmt"

func newHttpRegistry(url string) Registry {
	return &httpRegistry{url: url}
}

type httpRegistry struct {
	url string
}

func (h *httpRegistry) FindExtensions() ([]string, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (h *httpRegistry) ReadExtension(id string) (string, error) {
	return "", fmt.Errorf("unimplemented")
}
