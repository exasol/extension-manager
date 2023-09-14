package restAPI

import (
	"database/sql"
	"net/http"

	"github.com/Nightapes/go-rest/pkg/openapi"
	"github.com/exasol/extension-manager/pkg/apiErrors"
	"github.com/exasol/extension-manager/pkg/parameterValidator"
	"github.com/go-chi/chi/v5"
)

/* [impl -> dsn~parameter-versioning~1]. */
func GetExtensionDetails(apiContext *ApiContext) *openapi.Get {
	//nolint:exhaustruct // Default values for request are OK
	return &openapi.Get{
		Summary:        "Get details about an extension version.",
		Description:    "This returns details about an extension version, e.g. the parameter definitions required for creating an instance.",
		OperationID:    "GetExtensionDetails",
		Tags:           []string{TagExtension},
		Authentication: authentication,
		Response: map[string]openapi.MethodResponse{
			"200": {
				Description: "OK",
				Value: ExtensionDetailsResponse{Id: "s3-vs", Version: "1.2.3", ParamDefinitions: []ParamDefinition{
					{Id: "s3Bucket", Name: "S3 Bucket Name",
						RawDefinition: map[string]interface{}{"id": "s3Bucket", "name": "S3 Bucket Name", "type": "string", "required": true}}}}},
			"404": {
				Description: "Extension not found or creating instances not supported for this extension",
				Value:       apiErrors.NewNotFoundErrorF("Creating instances not supported")},
		},
		Path: newPathWithDbQueryParams().
			Add("extensions").
			AddParameter("extensionId", openapi.STRING, "ID of the extension").
			AddParameter("extensionVersion", openapi.STRING, "Version of the extension"),
		HandlerFunc: adaptDbHandler(apiContext, handleGetParameterDefinitions(apiContext)),
	}
}

func handleGetParameterDefinitions(apiContext *ApiContext) dbHandler {
	return func(db *sql.DB, writer http.ResponseWriter, request *http.Request) error {
		extensionId := chi.URLParam(request, "extensionId")
		extensionVersion := chi.URLParam(request, "extensionVersion")
		definitions, err := apiContext.Controller.GetParameterDefinitions(request.Context(), db, extensionId, extensionVersion)
		if err != nil {
			return err
		}
		response := ExtensionDetailsResponse{Id: extensionId, Version: extensionVersion, ParamDefinitions: convertParamDefinitions(definitions)}
		return SendJSON(request.Context(), writer, response)
	}
}

func convertParamDefinitions(definitions []parameterValidator.ParameterDefinition) []ParamDefinition {
	result := make([]ParamDefinition, 0, len(definitions))
	for _, d := range definitions {
		result = append(result, convertParamDefinition(d))
	}
	return result
}

func convertParamDefinition(d parameterValidator.ParameterDefinition) ParamDefinition {
	return ParamDefinition{
		Id: d.Id, Name: d.Name, RawDefinition: d.RawDefinition,
	}
}

// ExtensionDetailsResponse is the response for the GetExtensionDetails request.
type ExtensionDetailsResponse struct {
	Id               string            `json:"id"`                   // ID of this extension
	Version          string            `json:"version"`              // Version of this extension
	ParamDefinitions []ParamDefinition `json:"parameterDefinitions"` // Parameters required for creating an instance of this extension.
}

// This represents a parameter required for creating a new instance of an extension.
type ParamDefinition struct {
	Id            string      `json:"id"`         // ID of this parameter
	Name          string      `json:"name"`       // Name of this parameter
	RawDefinition interface{} `json:"definition"` // Raw parameter definition to be used as input for the Parameter Validator (https://github.com/exasol/extension-parameter-validator)
}
