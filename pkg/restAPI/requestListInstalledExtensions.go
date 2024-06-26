package restAPI

import (
	"database/sql"
	"net/http"

	"github.com/Nightapes/go-rest/pkg/openapi"
	"github.com/exasol/extension-manager/pkg/extensionAPI"
)

func ListInstalledExtensions(apiContext *ApiContext) *openapi.Get {
	//nolint:exhaustruct // Default values for request are OK
	return &openapi.Get{
		Summary:        "List installed extensions",
		Description:    "Get a list of all installed extensions.",
		OperationID:    "ListInstalledExtensions",
		Tags:           []string{TagInstallation},
		Authentication: authentication,
		Response: map[string]openapi.MethodResponse{
			"200": {Description: "List of extensions", Value: InstallationsResponse{
				Installations: []InstallationsResponseInstallation{
					{ID: "s3-vs", Name: "S3 Virtual Schema", Version: "1.0.0"},
					{ID: "cloud-storage", Name: "Cloud Storage Extension", Version: "1.1.0"}},
			}},
		},
		Path:        newPathWithDbQueryParams().Add("installations"),
		HandlerFunc: adaptDbHandler(apiContext, handleListInstalledExtensions(apiContext)),
	}
}

func handleListInstalledExtensions(apiContext *ApiContext) dbHandler {
	return func(db *sql.DB, writer http.ResponseWriter, request *http.Request) error {
		installations, err := apiContext.Controller.GetInstalledExtensions(request.Context(), db)
		if err != nil {
			return err
		}
		response := createResponse(installations)
		return SendJSON(request.Context(), writer, response)
	}
}

func createResponse(installations []*extensionAPI.JsExtInstallation) InstallationsResponse {
	convertedInstallations := make([]InstallationsResponseInstallation, 0, len(installations))
	for _, installation := range installations {
		convertedInstallations = append(convertedInstallations, InstallationsResponseInstallation{
			ID: installation.ID, Name: installation.Name, Version: installation.Version,
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
	ID      string `json:"id"`
	Name    string `json:"name"`
	Version string `json:"version"`
}
