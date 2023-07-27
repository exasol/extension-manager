package restAPI

import (
	"database/sql"
	"net/http"

	"github.com/Nightapes/go-rest/pkg/openapi"
	"github.com/exasol/extension-manager/pkg/extensionAPI"
)

func ListInstalledExtensions(apiContext *ApiContext) *openapi.Get {
	return &openapi.Get{
		Summary:        "List installed extensions",
		Description:    "Get a list of all installed extensions.",
		OperationID:    "ListInstalledExtensions",
		Tags:           []string{TagInstallation},
		Authentication: authentication,
		Response: map[string]openapi.MethodResponse{
			"200": {Description: "List of extensions", Value: InstallationsResponse{
				Installations: []InstallationsResponseInstallation{
					{Name: "s3-vs", Version: "1.0.0"},
					{Name: "s3-vs", Version: "1.1.0"}},
			}},
		},
		Path:        newPathWithDbQueryParams().Add("installations"),
		HandlerFunc: adaptDbHandler(handleListInstalledExtensions(apiContext)),
	}
}

func handleListInstalledExtensions(apiContext *ApiContext) dbHandler {
	return func(db *sql.DB, writer http.ResponseWriter, request *http.Request) {
		installations, err := apiContext.Controller.GetInstalledExtensions(request.Context(), db)
		if err != nil {
			HandleError(request.Context(), writer, err)
			return
		}
		response := createResponse(installations)
		SendJSON(request.Context(), writer, response)
	}
}

func createResponse(installations []*extensionAPI.JsExtInstallation) InstallationsResponse {
	convertedInstallations := make([]InstallationsResponseInstallation, 0, len(installations))
	for _, installation := range installations {
		convertedInstallations = append(convertedInstallations, InstallationsResponseInstallation{
			Name: installation.Name, Version: installation.Version,
		})
	}
	return InstallationsResponse{
		Installations: convertedInstallations,
	}
}

// InstallationsResponse contains all installed extensions.
type InstallationsResponse struct {
	Installations []InstallationsResponseInstallation `json:"installations"`
}

// InstallationsResponseInstallation contains information about installed extensions.
type InstallationsResponseInstallation struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
