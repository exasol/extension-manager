package restAPI

import (
	"database/sql"
	"net/http"

	"github.com/Nightapes/go-rest/pkg/openapi"
	"github.com/exasol/extension-manager/pkg/apiErrors"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

/* [impl -> dsn~upgrade-extension~1]. */
func UpgradeExtension(apiContext *ApiContext) *openapi.Post {
	//nolint:exhaustruct // Default values for request are OK
	return &openapi.Post{
		Summary:        "Upgrade an extension.",
		Description:    "This upgrades all instances of an extension to the latest version, e.g. by updating the JAR used in adapter scripts to the latest version.",
		OperationID:    "UpgradeExtension",
		Tags:           []string{TagInstallation},
		Authentication: authentication,
		Response: map[string]openapi.MethodResponse{
			"200": {
				Description: "Extension upgraded successfully",
				Value:       UpgradeExtensionResponse{PreviousVersion: "1.2.3", NewVersion: "1.3.0"}},
			"412": {
				Description: "Extension already installed in the latest version",
				Value:       apiErrors.NewNotFoundErrorF("Latest version 1.3.0 is already installed")},
			"404": {
				Description: "Extension not found or not installed",
				Value:       apiErrors.NewNotFoundErrorF("Extension not found")},
		},
		Path: newPathWithDbQueryParams().
			Add("installations").
			AddParameter("extensionId", openapi.STRING, "The ID of the installed extension to upgrade").
			Add("upgrade"),
		HandlerFunc: adaptDbHandler(apiContext, handleUpgradeExtension(apiContext)),
	}
}

func handleUpgradeExtension(apiContext *ApiContext) dbHandler {
	return func(db *sql.DB, writer http.ResponseWriter, request *http.Request) error {
		extensionId := chi.URLParam(request, "extensionId")
		result, err := apiContext.Controller.UpgradeExtension(request.Context(), db, extensionId)
		if err != nil {
			logrus.Warnf("Upgrading of extension %q failed: %v", extensionId, err)
			return err
		}
		logrus.Infof("Successfully upgraded extension %q from version %s to %s", extensionId, result.PreviousVersion, result.NewVersion)
		return SendJSON(request.Context(), writer, UpgradeExtensionResponse{
			PreviousVersion: result.PreviousVersion,
			NewVersion:      result.NewVersion})
	}
}

// Response data for upgrading an extension.
type UpgradeExtensionResponse struct {
	PreviousVersion string `json:"previousVersion"` // Version that was installed before the upgrade.
	NewVersion      string `json:"newVersion"`      // New version that is installed after the upgrade.
}
