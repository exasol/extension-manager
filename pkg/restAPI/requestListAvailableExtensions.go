package restAPI

import (
	"database/sql"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/Nightapes/go-rest/pkg/openapi"

	"github.com/exasol/extension-manager/pkg/extensionAPI"
	"github.com/exasol/extension-manager/pkg/extensionController"
)

func ListAvailableExtensions(apiContext *ApiContext) *openapi.Get {
	return &openapi.Get{
		Summary:        "List available extensions",
		Description:    "Get a list of all available extensions, i.e. extensions that can be installed.",
		OperationID:    "ListAvailableExtensions",
		Tags:           []string{TagExtension},
		Authentication: authentication,
		Response: map[string]openapi.MethodResponse{
			"200": {Description: "List of extensions", Value: ExtensionsResponse{
				Extensions: []ExtensionsResponseExtension{{
					Id:                  "s3-vs",
					Name:                "S3 Virtual Schema",
					Category:            "virtual-schema",
					Description:         "...",
					InstallableVersions: []ExtensionVersion{{Name: "1.2.3", Deprecated: true, Latest: false}, {Name: "1.3.0", Latest: true, Deprecated: false}},
				}},
			}},
		},
		Path:        newPathWithDbQueryParams().Add("extensions"),
		HandlerFunc: adaptDbHandler(handleListAvailableExtensions(apiContext)),
	}
}

func handleListAvailableExtensions(apiContext *ApiContext) dbHandler {
	return func(db *sql.DB, writer http.ResponseWriter, request *http.Request) {
		extensions, err := apiContext.Controller.GetAllExtensions(request.Context(), db)
		if err != nil {
			HandleError(request.Context(), writer, err)
			return
		}
		response := convertResponse(extensions)
		log.Debugf("Got %d available extensions", len(response.Extensions))
		SendJSON(request.Context(), writer, response)
	}
}

func convertResponse(extensions []*extensionController.Extension) ExtensionsResponse {
	convertedExtensions := make([]ExtensionsResponseExtension, 0, len(extensions))
	for _, extension := range extensions {
		convertedExtensions = append(convertedExtensions, convertExtension(extension))
	}
	return ExtensionsResponse{Extensions: convertedExtensions}
}

func convertExtension(extension *extensionController.Extension) ExtensionsResponseExtension {
	return ExtensionsResponseExtension{
		Id:                  extension.Id,
		Name:                extension.Name,
		Category:            extension.Category,
		Description:         extension.Description,
		InstallableVersions: convertVersions(extension.InstallableVersions)}
}

func convertVersions(versions []extensionAPI.JsExtensionVersion) []ExtensionVersion {
	result := make([]ExtensionVersion, 0, len(versions))
	for _, v := range versions {
		result = append(result, ExtensionVersion{Name: v.Name, Latest: v.Latest, Deprecated: v.Deprecated})
	}
	return result
}

// ExtensionsResponse contains all available extensions.
type ExtensionsResponse struct {
	Extensions []ExtensionsResponseExtension `json:"extensions"` // All available extensions.
}

// ExtensionsResponseExtension contains information about an available extension that can be installed.
type ExtensionsResponseExtension struct {
	Id                  string             `json:"id"`                  // ID of the extension. Don't store this as it may change when restarting the server.
	Name                string             `json:"name"`                // The name of the extension to be displayed to the user.
	Category            string             `json:"category"`            // The category of the extension, e.g. "driver" or "virtual-schema".
	Description         string             `json:"description"`         // The description of the extension to be displayed to the user.
	InstallableVersions []ExtensionVersion `json:"installableVersions"` // A list of versions of this extension available for installation.
}

type ExtensionVersion struct {
	Name       string `json:"name"`
	Latest     bool   `json:"latest"`
	Deprecated bool   `json:"deprecated"`
}
