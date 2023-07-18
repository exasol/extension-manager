package restAPI

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/exasol/extension-manager/pkg/extensionController"
	log "github.com/sirupsen/logrus"
)

// RestAPI is the interface that provides the REST API server of the extension-manager.
type RestAPI interface {
	// Serve starts the server. This method blocks until the server is stopped or fails.
	Serve()
	// StartInBackground starts the server in the background and blocks until it is ready, i.e. reacts to HTTP requests.
	StartInBackground()
	// Stop stops the server.
	Stop()
}

// Create creates a new RestAPI.
func Create(controller extensionController.TransactionController, serverAddress string) RestAPI {
	return &restAPIImpl{controller: controller, serverAddress: serverAddress}
}

type restAPIImpl struct {
	controller    extensionController.TransactionController
	serverAddress string
	server        *http.Server
	stopped       *bool
	stoppedMutex  *sync.Mutex
}

func (api *restAPIImpl) StartInBackground() {
	go api.Serve()
	api.waitUntilServerReplies()
}

func (api *restAPIImpl) Serve() {
	if api.server != nil {
		panic("server already running")
	}
	api.setStopped(false)

	handler, _, err := setupStandaloneAPI(api.controller)
	if err != nil {
		log.Fatalf("failed to setup api: %v", err)
	}
	api.server = &http.Server{
		Addr:              api.serverAddress,
		Handler:           handler,
		ReadHeaderTimeout: 3 * time.Second,
	}
	api.startServer()
}

// startServer starts the server on the given serverAddress. This method blocks until the server is stopped or fails.
func (api *restAPIImpl) startServer() {
	log.Printf("Starting server on %s...\n", api.serverAddress)
	err := api.server.ListenAndServe() // blocking
	if err != nil && !api.isStopped() {
		log.Fatalf("failed to start server: %v", err)
	}
}

func (api *restAPIImpl) waitUntilServerReplies() {
	request, err := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("http://%s", api.serverAddress), strings.NewReader(""))
	if err != nil {
		log.Fatalf("failed to create request: %v", err)
	}
	timeout := time.Now().Add(1 * time.Second)
	for {
		response, err := http.DefaultClient.Do(request)
		if err == nil {
			response.Body.Close()
			if response.StatusCode == 404 {
				return
			}
		}
		if time.Now().After(timeout) {
			log.Fatalf("Server did not reply within 1s, error: %v, response: %d %s", err, response.StatusCode, response.Status)
		}
	}
}

func (api *restAPIImpl) setStopped(stopped bool) {
	if api.stopped == nil {
		stopped := false
		api.stopped = &stopped
		api.stoppedMutex = &sync.Mutex{}
	}
	api.stoppedMutex.Lock()
	defer api.stoppedMutex.Unlock()
	*api.stopped = stopped
}

func (api *restAPIImpl) isStopped() bool {
	api.stoppedMutex.Lock()
	defer api.stoppedMutex.Unlock()
	return *api.stopped
}

func (api *restAPIImpl) Stop() {
	if api.server == nil {
		panic("cant stop server since it's not running")
	}
	api.setStopped(true)
	err := api.server.Shutdown(context.Background())
	if err != nil {
		panic(fmt.Sprintf("failed to shutdown rest API server. Cause: %v", err))
	}
	api.server = nil
}
