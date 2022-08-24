package restAPI

import (
	"github.com/Nightapes/go-rest/pkg/openapi"
	"github.com/exasol/extension-manager/extensionController"
	"github.com/exasol/extension-manager/restAPI/core"
	"github.com/exasol/extension-manager/restAPI/extensions"
)

func AddPublicEndpoints(api *openapi.API, config ExtensionManagerConfig) error {
	controller := extensionController.Create(config.ExtensionFolder, config.Schema)
	return addPublicEndpointsWithController(api, controller)
}

func addPublicEndpointsWithController(api *openapi.API, controller extensionController.TransactionController) error {
	api.AddTag(core.TagExtension, "Calls to list, install and uninstall extensions")
	api.AddTag(core.TagInstance, "Calls to list, create and remove instances of an extension")

	apiContext := core.NewApiContext(controller)

	if err := api.Get(extensions.ListAvailableExtensions(apiContext)); err != nil {
		return err
	}
	if err := api.Get(extensions.ListInstalledExtensions(apiContext)); err != nil {
		return err
	}
	if err := api.Put(extensions.InstallExtension(apiContext)); err != nil {
		return err
	}
	if err := api.Put(extensions.CreateInstance(apiContext)); err != nil {
		return err
	}
	return nil
}

type ExtensionManagerConfig struct {
	ExtensionFolder string
	Schema          string
}
