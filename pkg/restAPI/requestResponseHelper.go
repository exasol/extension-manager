package restAPI

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/exasol/extension-manager/pkg/apiErrors"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
)

const (
	ContentTypeJson   = "application/json"
	HeaderContentType = "Content-Type"
)

// SendJSON converts the given data to JSON and sends it to the writer.
func SendJSON(ctx context.Context, writer http.ResponseWriter, data interface{}) error {
	return SendJSONWithStatus(ctx, 200, writer, data)
}

func SendNoContent(ctx context.Context, writer http.ResponseWriter) error {
	return SendJSONWithStatus(ctx, http.StatusNoContent, writer, nil)
}

func SendJSONWithStatus(ctx context.Context, status int, writer http.ResponseWriter, data interface{}) error {
	logger := GetLogger(ctx)
	writer.Header().Set(HeaderContentType, ContentTypeJson)
	writer.WriteHeader(status)

	if log.IsLevelEnabled(log.TraceLevel) {
		jsonData, err := json.MarshalIndent(data, "", "    ")
		if err != nil {
			logger.Warnf("Failed to format json data for logging: %q", data)
		} else {
			logger.Debugf("Send json %s", jsonData)
		}
	}
	if data != nil {
		encoder := json.NewEncoder(writer)
		encoder.SetEscapeHTML(false)
		encodeErr := encoder.Encode(data)
		if encodeErr != nil {
			err := fmt.Errorf("Could not send json: %w", encodeErr)
			logger.Warnf(err.Error())
			return err

		}
	} else if status != http.StatusNoContent {
		logger.Warnf("No response data for status %d", status)
	}
	return nil
}

func handleError(context context.Context, apiContext *ApiContext, writer http.ResponseWriter, err error) {
	log.Errorf("Error processing request: %v", err)
	errorToSend := apiErrors.UnwrapAPIError(err)
	sendError(errorToSend, context, apiContext, writer)
}

func sendError(a *apiErrors.APIError, context context.Context, apiContext *ApiContext, writer http.ResponseWriter) {
	writer.Header().Set(HeaderContentType, ContentTypeJson)
	writer.WriteHeader(a.Status)
	if context != nil && a.Status != http.StatusUnauthorized {
		a.RequestID = middleware.GetReqID(context)
	}
	if apiContext.addCauseToInternalServerError && a.Status == http.StatusInternalServerError && a.OriginalError != nil {
		a.Message = a.Message + ": " + a.OriginalError.Error()
	}
	err := json.NewEncoder(writer).Encode(a)
	if err != nil {
		logger := GetLogger(context)
		logger.Errorf("Could not send simple error to client %s", err.Error())
	}
}

func GetLogger(context context.Context) *log.Entry {
	fields := log.Fields{}
	if id := middleware.GetReqID(context); id != "" {
		fields["request"] = id
	}
	return log.WithFields(fields)
}

func getContextValue(ctx context.Context, id interface{}) string {
	if ctx == nil {
		return ""
	}
	if reqID, ok := ctx.Value(id).(string); ok {
		return reqID
	}
	return ""
}

func DecodeJSONBody(writer http.ResponseWriter, request *http.Request, dst interface{}) error {
	if value := request.Header.Get(HeaderContentType); value != ContentTypeJson {
		return apiErrors.NewAPIError(http.StatusBadRequest, "Content-Type header is not application/json")
	}

	request.Body = http.MaxBytesReader(writer, request.Body, 1048576)
	dec := json.NewDecoder(request.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&dst)
	if err != nil {
		return convertError(err)
	}

	return verifyNoMoreJsonContent(dec)
}

func verifyNoMoreJsonContent(dec *json.Decoder) error {
	err := dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return apiErrors.NewBadRequestErrorF("Request body must only contain a single JSON object")
	}
	return nil
}

func convertError(err error) error {
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError
	switch {
	case errors.As(err, &syntaxError):
		return apiErrors.NewBadRequestErrorF("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)

	case errors.Is(err, io.ErrUnexpectedEOF):
		return apiErrors.NewBadRequestErrorF("Request body contains badly-formed JSON")

	case errors.As(err, &unmarshalTypeError):
		return apiErrors.NewBadRequestErrorF("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)

	case strings.HasPrefix(err.Error(), "json: unknown field "):
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		return apiErrors.NewBadRequestErrorF("Request body contains unknown field %q", fieldName)

	case errors.Is(err, io.EOF):
		return apiErrors.NewBadRequestErrorF("Request body must not be empty")

	case err.Error() == "http: request body too large":
		return apiErrors.NewBadRequestErrorF("Request body must not be larger than 1MB")

	default:
		return err
	}
}
