package extensions

import (
	"net/http"

	"github.com/Nightapes/go-rest/pkg/openapi"

	"github.com/exasol/extension-manager/extensionController"
	"github.com/exasol/extension-manager/restAPI/core"
)

func CreateInstance(apiContext core.ApiContext) *openapi.Put {
	return &openapi.Put{
		Summary:        "Create an instance of an extension.",
		Description:    "This creates a new instance of an extension, e.g. a virtual schema.",
		OperationID:    "CreateInstance",
		Tags:           []string{core.TagInstance},
		Authentication: map[string][]string{core.BearerAuth: {}},
		RequestBody:    CreateInstanceRequest{ExtensionId: "s3-vs", ExtensionVersion: "1.1.0", ParameterValues: []ParameterValue{{Name: "param1", Value: "value1"}}},
		Response: map[string]openapi.MethodResponse{
			"200": {Description: "OK", Value: CreateInstanceResponse{InstanceName: "new-instance-name"}},
		},
		Path:        core.NewPublicPath().Add("instances"),
		HandlerFunc: handleCreateInstance(apiContext),
	}
}

func handleCreateInstance(apiContext core.ApiContext) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		db, err := apiContext.OpenDBConnection(request)
		if err != nil {
			core.HandleError(request.Context(), writer, err)
			return
		}
		defer core.CloseDbConnection(db)

		requestBody := CreateInstanceRequest{}
		err = core.DecodeJSONBody(writer, request, &requestBody)
		if err != nil {
			core.HandleError(request.Context(), writer, err)
			return
		}
		var parameters []extensionController.ParameterValue
		for _, p := range requestBody.ParameterValues {
			parameters = append(parameters, extensionController.ParameterValue{Name: p.Name, Value: p.Value})
		}
		instanceName, err := apiContext.Controller().CreateInstance(request.Context(), db, requestBody.ExtensionId, requestBody.ExtensionVersion, parameters)
		if err != nil {
			core.HandleError(request.Context(), writer, err)
			return
		}
		core.SendJSON(request.Context(), writer, CreateInstanceResponse{InstanceName: instanceName})
	}
}

// Request data for creating a new instance of an extension.
type CreateInstanceRequest struct {
	ExtensionId      string           `json:"extensionId"`      // The ID of the extension
	ExtensionVersion string           `json:"extensionVersion"` // The version of the extension
	ParameterValues  []ParameterValue `json:"parameterValues"`  // The parameters for the new instance
}

// Parameter values for creating a new instance.
type ParameterValue struct {
	Name  string `json:"name"`  // The name of the parameter
	Value string `json:"value"` // The value of the parameter
}

// Response data for creating a new instance of an extension.
type CreateInstanceResponse struct {
	InstanceName string `json:"instanceName"` // The name of the newly created instance
}
