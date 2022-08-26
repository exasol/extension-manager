package restAPI

import (
	"database/sql"
	"net/http"

	"github.com/Nightapes/go-rest/pkg/openapi"
)

func InstallExtension(apiContext *ApiContext) *openapi.Put {
	return &openapi.Put{
		Summary:        "Install an extension.",
		Description:    "This installs an extension in a given version, e.g. by creating Adapter Scripts.",
		OperationID:    "InstallExtension",
		Tags:           []string{TagExtension},
		Authentication: authentication,
		RequestBody:    InstallExtensionRequest{},
		Response: map[string]openapi.MethodResponse{
			"204": {Description: "OK"},
		},
		Path:        newPathWithDbQueryParams().Add("installations"),
		HandlerFunc: adaptDbHandler(handleInstallExtension(apiContext)),
	}
}

func handleInstallExtension(apiContext *ApiContext) dbHandler {
	return func(db *sql.DB, writer http.ResponseWriter, request *http.Request) {
		requestBody := InstallExtensionRequest{}
		err := DecodeJSONBody(writer, request, &requestBody)
		if err != nil {
			HandleError(request.Context(), writer, err)
			return
		}
		err = apiContext.Controller.InstallExtension(request.Context(), db, requestBody.ExtensionId, requestBody.ExtensionVersion)

		if err != nil {
			HandleError(request.Context(), writer, err)
			return
		}
		SendNoContent(request.Context(), writer)
	}
}

type InstallExtensionRequest struct {
	ExtensionId      string
	ExtensionVersion string
}
