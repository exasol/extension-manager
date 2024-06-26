package restAPI

import (
	"github.com/Nightapes/go-rest/pkg/openapi"
	"github.com/exasol/extension-manager/pkg/extensionController"
)

const (
	/* [impl -> const~use-reserved-schema~1]. */
	EXTENSION_SCHEMA_NAME = "EXA_EXTENSIONS"

	TagExtension    = "Extension"
	TagInstallation = "Installation"
	TagInstance     = "Instance"

	BearerAuth = "DbAccessToken"
	BasicAuth  = "DbUsernamePassword"
)

// AddPublicEndpoints adds the extension manager endpoints to the API.
// The config struct contains configuration options for the extension manager.
/* [impl -> dsn~go-library~1]. */
func AddPublicEndpoints(api *openapi.API, config extensionController.ExtensionManagerConfig) error {
	controller, err := extensionController.CreateWithValidatedConfig(config)
	if err != nil {
		return err
	}
	return addPublicEndpointsWithController(api, false, controller)
}

/* [impl -> dsn~rest-interface~1] */
/* [impl -> dsn~openapi-spec~1]. */
func addPublicEndpointsWithController(api *openapi.API, addCauseToInternalServerError bool, controller extensionController.TransactionController) error {
	api.AddTag(TagExtension, "List and install extensions")
	api.AddTag(TagInstallation, "List and uninstall installed extensions")
	api.AddTag(TagInstance, "Calls to list, create and remove instances of an extension")

	apiContext := NewApiContext(controller, addCauseToInternalServerError)

	if err := api.Get(ListAvailableExtensions(apiContext)); err != nil {
		return err
	}
	if err := api.Get(ListInstalledExtensions(apiContext)); err != nil {
		return err
	}
	if err := api.Get(GetExtensionDetails(apiContext)); err != nil {
		return err
	}
	if err := api.Put(InstallExtension(apiContext)); err != nil {
		return err
	}
	if err := api.Delete(UninstallExtension(apiContext)); err != nil {
		return err
	}
	if err := api.Post(UpgradeExtension(apiContext)); err != nil {
		return err
	}
	if err := api.Post(CreateInstance(apiContext)); err != nil {
		return err
	}
	if err := api.Get(ListInstances(apiContext)); err != nil {
		return err
	}
	if err := api.Delete(DeleteInstance(apiContext)); err != nil {
		return err
	}
	return nil
}
