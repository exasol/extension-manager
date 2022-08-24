package extensions

import (
	"net/http"

	"github.com/Nightapes/go-rest/pkg/openapi"
	"github.com/exasol/extension-manager/extensionAPI"
	"github.com/exasol/extension-manager/restAPI/core"
	"github.com/exasol/extension-manager/restAPI/models"
	log "github.com/sirupsen/logrus"
)

func ListInstalledExtensions(apiContext core.ApiContext) *openapi.Get {
	return &openapi.Get{
		Summary:        "List installed extensions",
		Description:    "Get a list of all installed extensions.",
		OperationID:    "ListInstalledExtensions",
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
		Path:        core.NewPublicPath().Add("installations"),
		HandlerFunc: handleListInstalledExtensions(apiContext),
	}
}

func handleListInstalledExtensions(apiContext core.ApiContext) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		db, err := apiContext.OpenDBConnection(request)
		if err != nil {
			core.HandleError(request.Context(), writer, err)
			return
		}
		defer core.CloseDbConnection(db)
		installations, err := apiContext.Controller().GetAllInstallations(request.Context(), db)
		if err != nil {
			core.HandleError(request.Context(), writer, err)
			return
		}
		response := createResponse(installations)
		log.Debugf("Installed extensions: %d", len(response.Installations))
		core.SendJSON(request.Context(), writer, response)
	}
}

func createResponse(installations []*extensionAPI.JsExtInstallation) models.InstallationsResponse {
	convertedInstallations := make([]models.InstallationsResponseInstallation, 0, len(installations))
	for _, installation := range installations {
		convertedInstallations = append(convertedInstallations, models.InstallationsResponseInstallation{
			Name: installation.Name, Version: installation.Version, InstanceParameters: installation.InstanceParameters,
		})
	}
	return models.InstallationsResponse{
		Installations: convertedInstallations,
	}
}
