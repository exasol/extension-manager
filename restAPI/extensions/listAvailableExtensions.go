package extensions

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/Nightapes/go-rest/pkg/openapi"

	"github.com/exasol/extension-manager/extensionController"
	"github.com/exasol/extension-manager/restAPI/core"
	"github.com/exasol/extension-manager/restAPI/models"
)

func ListAvailableExtensions(apiContext core.ApiContext) *openapi.Get {
	return &openapi.Get{
		Summary:        "List available extensions",
		Description:    "Get a list of all available extensions, i.e. extensions that can be installed.",
		OperationID:    "ListAvailableExtensions",
		Tags:           []string{core.TagExtension},
		Authentication: map[string][]string{core.BearerAuth: {}},
		Response: map[string]openapi.MethodResponse{
			"200": {Description: "List of extensions", Value: models.ExtensionsResponse{
				Extensions: []models.ExtensionsResponseExtension{{
					Id:                  "s3-vs",
					Name:                "S3 Virtual Schema",
					Description:         "...",
					InstallableVersions: []string{"1.0.0", "1.2.0"},
				}},
			}},
		},
		Path:        core.NewPublicPath().Add("extensions"),
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

func convertResponse(extensions []*extensionController.Extension) models.ExtensionsResponse {
	convertedExtensions := make([]models.ExtensionsResponseExtension, 0, len(extensions))
	for _, extension := range extensions {
		ext := models.ExtensionsResponseExtension{Id: extension.Id, Name: extension.Name, Description: extension.Description, InstallableVersions: extension.InstallableVersions}
		convertedExtensions = append(convertedExtensions, ext)
	}
	return models.ExtensionsResponse{
		Extensions: convertedExtensions,
	}
}
