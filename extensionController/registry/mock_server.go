package registry

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/suite"
)

type mockRegistryServer struct {
	server          *httptest.Server
	suite           *suite.Suite
	registryContent string
	files           map[string]string
}

const REGISTRY_PATH = "/registry.json"

func newMockRegistryServer(suite *suite.Suite) *mockRegistryServer {
	return &mockRegistryServer{suite: suite, registryContent: ""}
}

func (s *mockRegistryServer) start() {
	router := chi.NewRouter()
	router.MethodFunc(http.MethodGet, REGISTRY_PATH, func(w http.ResponseWriter, r *http.Request) {
		if s.registryContent != "" {
			s.sendResponse(w, s.registryContent, 200)
		} else {
			s.sendResponse(w, "no content defined for registry", 404)
		}
	})
	router.MethodFunc(http.MethodGet, "/*", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if content, ok := s.files[path]; ok {
			s.sendResponse(w, content, 200)
		} else {
			s.sendResponse(w, fmt.Sprintf("no content defined for path %q", path), 404)
		}
	})
	s.server = httptest.NewServer(router)
}

func (s *mockRegistryServer) sendResponse(w http.ResponseWriter, content string, status int) {
	w.WriteHeader(status)
	_, err := w.Write([]byte(content))
	s.suite.NoError(err)
}

func (s *mockRegistryServer) setRegistryContent(content string) {
	s.registryContent = content
}

func (s *mockRegistryServer) setPathContent(path, content string) {
	s.files[path] = content
}

func (s *mockRegistryServer) reset() {
	s.registryContent = ""
	s.files = make(map[string]string)
}

func (s *mockRegistryServer) baseUrl() string {
	return s.server.URL
}

func (s *mockRegistryServer) indexUrl() string {
	return s.server.URL + REGISTRY_PATH
}

func (s *mockRegistryServer) close() {
	s.server.Close()
}
