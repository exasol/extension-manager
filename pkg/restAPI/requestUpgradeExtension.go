package restAPI

import (
	"database/sql"
	"net/http"

	"github.com/Nightapes/go-rest/pkg/openapi"
	"github.com/exasol/extension-manager/pkg/apiErrors"
	"github.com/go-chi/chi/v5"
)

/* [impl -> dsn~upgrade-extension~1]. */
func UpgradeExtension(apiContext *ApiContext) *openapi.Post {
	return &openapi.Post{
		Summary:        "Upgrade an extension.",
		Description:    "This upgrades all instances of an extension to the latest version, e.g. by updating the JAR used in adapter scripts to the latest version.",
		OperationID:    "UpgradeExtension",
		Tags:           []string{TagInstallation},
		Authentication: authentication,
		Response: map[string]openapi.MethodResponse{
			"200": {
				Description: "OK",
				Value:       UpgradeExtensionResponse{PreviousVersion: "1.2.3", NewVersion: "1.3.0"}},
			"404": {
				Description: "Extension not found",
				Value:       apiErrors.NewNotFoundErrorF("Extension not found")},
		},
		Path: newPathWithDbQueryParams().
			Add("installations").
			AddParameter("extensionId", openapi.STRING, "The ID of the installed extension to upgrade").
			Add("upgrade"),
		HandlerFunc: adaptDbHandler(handleUpgradeExtension(apiContext)),
	}
}

func handleUpgradeExtension(apiContext *ApiContext) dbHandler {
	return func(db *sql.DB, writer http.ResponseWriter, request *http.Request) {
		extensionId := chi.URLParam(request, "extensionId")
		result, err := apiContext.Controller.UpgradeExtension(request.Context(), db, extensionId)
		if err != nil {
			HandleError(request.Context(), writer, err)
			return
		}
		SendJSON(request.Context(), writer, UpgradeExtensionResponse{
			PreviousVersion: result.PreviousVersion,
			NewVersion:      result.NewVersion})
	}
}

// Response data for upgrading an extension.
type UpgradeExtensionResponse struct {
	PreviousVersion string `json:"previousVersion"` // Version that was installed before the upgrade.
	NewVersion      string `json:"newVersion"`      // New version that is installed after the upgrade.
}
