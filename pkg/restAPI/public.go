package restAPI

import (
	"github.com/Nightapes/go-rest/pkg/openapi"
	"github.com/exasol/extension-manager/pkg/extensionController"
)

const (
	/* [impl -> const~use-reserved-schema~1] */
	EXTENSION_SCHEMA_NAME = "EXA_EXTENSIONS"

	TagExtension    = "Extension"
	TagInstallation = "Installation"
	TagInstance     = "Instance"

	BearerAuth = "DbAccessToken"
	BasicAuth  = "DbUsernamePassword"
)

// Configuration options for the extension manager.
type ExtensionManagerConfig struct {
	ExtensionRegistryURL string // URL of the extension registry index used to find available extensions. This can also be the path of a local directory.
}

// AddPublicEndpoints adds the extension manager endpoints to the API.
// The config struct contains configuration options for the extension manager.
/* [impl -> dsn~go-library~1] */
func AddPublicEndpoints(api *openapi.API, config ExtensionManagerConfig) error {
	controller := extensionController.Create(config.ExtensionRegistryURL, EXTENSION_SCHEMA_NAME)
	return addPublicEndpointsWithController(api, controller)
}

/* [impl -> dsn~rest-interface~1] */
/* [impl -> dsn~openapi-spec~1] */
func addPublicEndpointsWithController(api *openapi.API, controller extensionController.TransactionController) error {
	api.AddTag(TagExtension, "List and install extensions")
	api.AddTag(TagInstallation, "List and uninstall installed extensions")
	api.AddTag(TagInstance, "Calls to list, create and remove instances of an extension")

	apiContext := NewApiContext(controller)

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
