package restAPI

import (
	"database/sql"
	"net/http"

	"github.com/Nightapes/go-rest/pkg/openapi"
	"github.com/exasol/extension-manager/pkg/apiErrors"
	"github.com/go-chi/chi/v5"
)

func ListInstances(apiContext *ApiContext) *openapi.Get {
	return &openapi.Get{
		Summary:        "List all instances of an extension.",
		Description:    "This lists all instances of an extension, e.g. virtual schema.",
		OperationID:    "ListInstances",
		Tags:           []string{TagInstance},
		Authentication: authentication,
		Response: map[string]openapi.MethodResponse{
			"200": {Description: "OK", Value: ListInstancesResponse{Instances: []Instance{{Id: "s3-vs-1", Name: "SALES_S3_VS"}}}},
			"404": {
				Description: "Extension not found",
				Value:       apiErrors.NewNotFoundErrorF("Extension not found")},
		},
		Path: newPathWithDbQueryParams().Add("installations").
			AddParameter("extensionId", openapi.STRING, "The ID of the installed extension for which to get the instances").
			AddParameter("extensionVersion", openapi.STRING, "The version of the installed extension for which to get the instances").
			Add("instances"),
		HandlerFunc: adaptDbHandler(handleListInstances(apiContext)),
	}
}

func handleListInstances(apiContext *ApiContext) dbHandler {
	return func(db *sql.DB, writer http.ResponseWriter, request *http.Request) {
		id := chi.URLParam(request, "extensionId")
		version := chi.URLParam(request, "extensionVersion")
		rawInstances, err := apiContext.Controller.FindInstances(request.Context(), db, id, version)
		if err != nil {
			HandleError(request.Context(), writer, err)
			return
		}
		instances := make([]Instance, 0, len(rawInstances))
		for _, i := range rawInstances {
			instances = append(instances, Instance{Id: i.Id, Name: i.Name})
		}
		SendJSON(request.Context(), writer, ListInstancesResponse{Instances: instances})
	}
}

// Response data for listing all instances of an extension.
type ListInstancesResponse struct {
	Instances []Instance `json:"instances"` // Instances of the extension.
}

// Instance represents an instance of an extension, e.g. a virtual schema.
type Instance struct {
	Id   string `json:"id"`   // The ID of the instance
	Name string `json:"name"` // The name of the instance
}
