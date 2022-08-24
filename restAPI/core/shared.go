package core

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Nightapes/go-rest/pkg/openapi"

	"github.com/go-chi/chi/v5/middleware"
)

const (
	TagExtension = "Extension"
	TagInstance  = "Instance"

	BearerAuth = "BearerAuth"
	BasicAuth  = "BasicAuth"
)

type ExaPath struct {
	*openapi.PathBuilder
}

// NewPublicPath /api/v1
func NewPublicPath() *ExaPath {
	return &ExaPath{getV1PublicBasePath(openapi.NewPathBuilder())}
}

func getV1PublicBasePath(builder *openapi.PathBuilder) *openapi.PathBuilder {
	return builder.
		Add("api").
		Add("v1")
}

func (e *ExaPath) GetInstalledExtensionsBasePath() *ExaPath {
	e.Add("installedExtensions")
	return e
}

type APIError struct {
	Status        int    `json:"code"`                // HTTP status code
	Message       string `json:"message"`             // human-readable message
	RequestID     string `json:"requestID,omitempty"` // ID to identify the request that caused this error
	Timestamp     string `json:"timestamp,omitempty" jsonschema:"format=date-time" `
	APIID         string `json:"apiID,omitempty"` // Corresponding API action
	OriginalError error  `json:"-"`
}

func (a *APIError) Error() string {
	return a.Message
}

func (a APIError) Send(context context.Context, writer http.ResponseWriter) {
	logger := GetLogger(context)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(a.Status)
	if context != nil && a.Status != http.StatusUnauthorized {
		a.RequestID = middleware.GetReqID(context)
		if a.Timestamp == "" {
			a.Timestamp = time.Now().Format(time.RFC3339)
		}
		a.APIID = GetContextValue(context, APIIDKey)
	}

	err := json.NewEncoder(writer).Encode(a)
	if err != nil {
		logger.Errorf("Could not send simple error to client %s", err.Error())
	}
}
