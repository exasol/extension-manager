package restAPI

import (
	"database/sql"
	"net/http"

	"github.com/Nightapes/go-rest/pkg/openapi"
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
		},
		Path: newPathWithDbQueryParams().Add("extension").
			AddParameter("extensionId", openapi.STRING, "The ID of the extension for which to delete an instance").
			Add("instance").
			AddParameter("instanceId", openapi.STRING, "The ID of the instance to delete"),
		HandlerFunc: adaptDbHandler(handleDeleteInstance(apiContext)),
	}
}

func handleDeleteInstance(apiContext *ApiContext) dbHandler {
	return func(db *sql.DB, writer http.ResponseWriter, request *http.Request) {
		extensionId := chi.URLParam(request, "extensionId")
		instanceId := chi.URLParam(request, "instanceId")
		err := apiContext.Controller.DeleteInstance(request.Context(), db, extensionId, instanceId)
		if err != nil {
			HandleError(request.Context(), writer, err)
			return
		}
		SendNoContent(request.Context(), writer)
	}
}
