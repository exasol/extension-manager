package restAPI

import (
	"database/sql"
	"net/http"

	"github.com/Nightapes/go-rest/pkg/openapi"
	"github.com/exasol/extension-manager/pkg/apiErrors"
	"github.com/go-chi/chi/v5"
)

func DeleteInstance(apiContext *ApiContext) *openapi.Delete {
	return &openapi.Delete{
		Summary:        "Delete an instances of an extension.",
		Description:    "This deletes a single instances of an extension, e.g. a virtual schema.",
		OperationID:    "DeleteInstance",
		Tags:           []string{TagInstance},
		Authentication: authentication,
		Response: map[string]openapi.MethodResponse{
			"204": {Description: "OK"},
			"404": {
				Description: "Extension or instance not found",
				Value:       apiErrors.NewNotFoundErrorF("Extension not found")},
		},
		Path: newPathWithDbQueryParams().Add("installations").
			AddParameter("extensionId", openapi.STRING, "The ID of the extension for which to delete an instance").
			AddParameter("extensionVersion", openapi.STRING, "The version of the installed extension for which to delete an instance").
			Add("instances").
			AddParameter("instanceId", openapi.STRING, "The ID of the instance to delete"),
		HandlerFunc: adaptDbHandler(handleDeleteInstance(apiContext)),
	}
}

func handleDeleteInstance(apiContext *ApiContext) dbHandler {
	return func(db *sql.DB, writer http.ResponseWriter, request *http.Request) error {
		extensionId := chi.URLParam(request, "extensionId")
		extensionVersion := chi.URLParam(request, "extensionVersion")
		instanceId := chi.URLParam(request, "instanceId")
		err := apiContext.Controller.DeleteInstance(request.Context(), db, extensionId, extensionVersion, instanceId)
		if err != nil {
			return err
		}
		return SendNoContent(request.Context(), writer)
	}
}
