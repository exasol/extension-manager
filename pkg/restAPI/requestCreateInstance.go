package restAPI

import (
	"database/sql"
	"net/http"

	"github.com/Nightapes/go-rest/pkg/openapi"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"

	"github.com/exasol/extension-manager/pkg/extensionController"
)

func CreateInstance(apiContext *ApiContext) *openapi.Post {
	return &openapi.Post{
		Summary:        "Create an instance of an extension.",
		Description:    "This creates a new instance of an extension, e.g. a virtual schema.",
		OperationID:    "CreateInstance",
		Tags:           []string{TagInstance},
		Authentication: authentication,
		RequestBody:    CreateInstanceRequest{ParameterValues: []ParameterValue{{Name: "param1", Value: "value1"}}},
		Response: map[string]openapi.MethodResponse{
			"200": {Description: "OK", Value: CreateInstanceResponse{InstanceId: "id", InstanceName: "new-instance-name"}},
		},
		Path: newPathWithDbQueryParams().Add("installations").
			AddParameter("extensionId", openapi.STRING, "ID of the installed extension for which to create an instance").
			AddParameter("extensionVersion", openapi.STRING, "Version of the installed extension for which to create an instance").
			Add("instances"),
		HandlerFunc: adaptDbHandler(handleCreateInstance(apiContext)),
	}
}

func handleCreateInstance(apiContext *ApiContext) dbHandler {
	return func(db *sql.DB, writer http.ResponseWriter, request *http.Request) {
		requestBody := CreateInstanceRequest{}
		err := DecodeJSONBody(writer, request, &requestBody)
		if err != nil {
			HandleError(request.Context(), writer, err)
			return
		}
		var parameters []extensionController.ParameterValue
		for _, p := range requestBody.ParameterValues {
			parameters = append(parameters, extensionController.ParameterValue{Name: p.Name, Value: p.Value})
		}
		extensionId := chi.URLParam(request, "extensionId")
		extensionVersion := chi.URLParam(request, "extensionVersion")
		instance, err := apiContext.Controller.CreateInstance(request.Context(), db, extensionId, extensionVersion, parameters)
		if err != nil {
			HandleError(request.Context(), writer, err)
			return
		}
		logrus.Debugf("Created instance %q", instance)
		SendJSON(request.Context(), writer, CreateInstanceResponse{InstanceId: instance.Id, InstanceName: instance.Name})
	}
}

// Request data for creating a new instance of an extension.
type CreateInstanceRequest struct {
	ParameterValues []ParameterValue `json:"parameterValues"` // The parameters for the new instance
}

// Parameter values for creating a new instance.
type ParameterValue struct {
	Name  string `json:"name"`  // The name of the parameter
	Value string `json:"value"` // The value of the parameter
}

// Response data for creating a new instance of an extension.
type CreateInstanceResponse struct {
	InstanceId   string `json:"instanceId"`   // The ID of the newly created instance
	InstanceName string `json:"instanceName"` // The name of the newly created instance
}
