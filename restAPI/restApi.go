package restAPI

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/exasol/extension-manager/extensionController"
	log "github.com/sirupsen/logrus"
)

// RestAPI is the interface that provides the REST API server of the extension-manager.
type RestAPI interface {
	// Serve starts the server. This method blocks until the server is stopped or fails.
	Serve()
	// Stop stops the server
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
		Addr:    api.serverAddress,
		Handler: handler,
	}
	log.Printf("Starting server on %s...\n", api.serverAddress)
	err = api.server.ListenAndServe() // blocking
	if err != nil && !api.isStopped() {
		log.Fatalf("failed to start server: %v", err)
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
