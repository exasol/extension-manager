package requests

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/Nightapes/go-rest/pkg/openapi"

	"github.com/exasol/extension-manager/extensionController"
	"github.com/exasol/extension-manager/restAPI/core"
)

func ListAvailableExtensions(apiContext core.ApiContext) *openapi.Get {
	return &openapi.Get{
		Summary:        "List available extensions",
		Description:    "Get a list of all available extensions, i.e. extensions that can be installed.",
		OperationID:    "ListAvailableExtensions",
		Tags:           []string{core.TagExtension},
		Authentication: authentication,
		Response: map[string]openapi.MethodResponse{
			"200": {Description: "List of extensions", Value: ExtensionsResponse{
				Extensions: []ExtensionsResponseExtension{{
					Id:                  "s3-vs",
					Name:                "S3 Virtual Schema",
					Description:         "...",
					InstallableVersions: []string{"1.0.0", "1.2.0"},
				}},
			}},
		},
		Path:        core.NewPathWithDbQueryParams().Add("extensions"),
		HandlerFunc: handleListAvailableExtensions(apiContext),
	}
}

func handleListAvailableExtensions(apiContext core.ApiContext) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		db, err := apiContext.OpenDBConnection(request)
		if err != nil {
			core.HandleError(request.Context(), writer, err)
			return
		}
		defer core.CloseDbConnection(db)
		extensions, err := apiContext.Controller().GetAllExtensions(request.Context(), db)
		if err != nil {
			core.HandleError(request.Context(), writer, err)
			return
		}
		response := convertResponse(extensions)
		log.Debugf("Got %d available extensions", len(response.Extensions))
		core.SendJSON(request.Context(), writer, response)
	}
}

func convertResponse(extensions []*extensionController.Extension) ExtensionsResponse {
	convertedExtensions := make([]ExtensionsResponseExtension, 0, len(extensions))
	for _, extension := range extensions {
		ext := ExtensionsResponseExtension{Id: extension.Id, Name: extension.Name, Description: extension.Description, InstallableVersions: extension.InstallableVersions}
		convertedExtensions = append(convertedExtensions, ext)
	}
	return ExtensionsResponse{
		Extensions: convertedExtensions,
	}
}

// ExtensionsResponse contains all available extensions.
type ExtensionsResponse struct {
	Extensions []ExtensionsResponseExtension `json:"extensions"` // All available extensions.
}

// ExtensionsResponseExtension contains information about an available extension that can be installed.
type ExtensionsResponseExtension struct {
	Id                  string   `json:"id"`                  // ID of the extension. Don't store this as it may change in the future.
	Name                string   `json:"name"`                // The name of the extension to be displayed to the user.
	Description         string   `json:"description"`         // The description of the extension to be displayed to the user.
	InstallableVersions []string `json:"installableVersions"` // A list of versions of this extension available for installation.
}
