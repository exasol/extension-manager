package restAPI

import (
	"github.com/Nightapes/go-rest/pkg/openapi"
	"github.com/exasol/extension-manager/extensionController"
	"github.com/exasol/extension-manager/restAPI/core"
	"github.com/exasol/extension-manager/restAPI/requests"
)

const EXTENSION_SCHEMA_NAME = "EXA_EXTENSIONS"

// Configuration options for the extension manager.
type ExtensionManagerConfig struct {
	ExtensionFolder string // Path to the local folder containing the extensions as .js files.
}

// AddPublicEndpoints adds the extension manager endpoints to the API.
// The config struct contains configuration options for the extension manager.
func AddPublicEndpoints(api *openapi.API, config ExtensionManagerConfig) error {
	controller := extensionController.Create(config.ExtensionFolder, EXTENSION_SCHEMA_NAME)
	return addPublicEndpointsWithController(api, controller)
}

func addPublicEndpointsWithController(api *openapi.API, controller extensionController.TransactionController) error {
	api.AddTag(core.TagExtension, "Calls to list, install and uninstall extensions")
	api.AddTag(core.TagInstance, "Calls to list, create and remove instances of an extension")

	apiContext := core.NewApiContext(controller)

	if err := api.Get(requests.ListAvailableExtensions(apiContext)); err != nil {
		return err
	}
	if err := api.Get(requests.ListInstalledExtensions(apiContext)); err != nil {
		return err
	}
	if err := api.Put(requests.InstallExtension(apiContext)); err != nil {
		return err
	}
	if err := api.Put(requests.CreateInstance(apiContext)); err != nil {
		return err
	}
	return nil
}