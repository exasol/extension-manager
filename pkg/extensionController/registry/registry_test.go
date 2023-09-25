package registry

import (
	"fmt"
	"testing"
)

/* [utest -> dsn~extension-registry~1]. */
func TestNewRegistry(t *testing.T) {
	const typeLocal = "*registry.localDirRegistry"
	const typeHttp = "*registry.httpRegistry"
	tests := []struct {
		input        string
		expectedType string
	}{
		{"http://url", typeHttp},
		{"HTTP://URL", typeHttp},
		{"https://url", typeHttp},
		{"HTTPS://URL", typeHttp},
		{"http/url", typeLocal},
		{"https:/url", typeLocal},
		{"https:/url", typeLocal},
		{"https/url", typeLocal},
		{"https/url", typeLocal},
		{"relative/path", typeLocal},
		{"/absolute/path", typeLocal},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := NewRegistry(test.input)
			actualType := fmt.Sprintf("%T", result)
			if actualType != test.expectedType {
				t.Errorf("expected type %s but got %s", test.expectedType, actualType)
			}
		})
	}
}
