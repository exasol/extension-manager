package core

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"io"
	"net/http"
	"strings"

	"github.com/exasol/extension-manager/apiErrors"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
)

// SendJSON to writer
func SendJSON(ctx context.Context, writer http.ResponseWriter, data interface{}) {
	SendJSONWithStatus(ctx, 200, writer, data)
}

func SendNoContent(ctx context.Context, writer http.ResponseWriter) {
	SendJSONWithStatus(ctx, http.StatusNoContent, writer, nil)
}

func SendJSONWithStatus(ctx context.Context, status int, writer http.ResponseWriter, data interface{}) {
	logger := GetLogger(ctx)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	if log.IsLevelEnabled(log.TraceLevel) {
		jsonData, _ := json.MarshalIndent(data, "", "    ")
		logger.Debugf("Send json %s", jsonData)
	}
	if data != nil {
		encoder := json.NewEncoder(writer)
		encoder.SetEscapeHTML(false)
		encodeErr := encoder.Encode(data)
		if encodeErr != nil {
			logger.Warnf("Could not send json: %s", encodeErr.Error())
		}
	} else {
		logger.Warnf("No data")
	}
}

func HandleError(context context.Context, writer http.ResponseWriter, err error) {
	errorToSend := convertToApiError(err)
	sendError(errorToSend, context, writer)
}

func convertToApiError(err error) *apiErrors.APIError {
	switch apiError := err.(type) {
	default:
		log.Errorf("Internal error: %s", err.Error())
		return apiErrors.NewInternalServerError(err).(*apiErrors.APIError)

	case *apiErrors.APIError:
		if apiError.OriginalError != nil {
			log.Errorf("Error: %s (original: %s)", err.Error(), apiError.OriginalError.Error())
		} else {
			log.Errorf("Error: %s", err.Error())
		}
		return err.(*apiErrors.APIError)
	}
}

func sendError(a *apiErrors.APIError, context context.Context, writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(a.Status)
	if context != nil && a.Status != http.StatusUnauthorized {
		a.RequestID = middleware.GetReqID(context)
		if a.Timestamp == "" {
			a.Timestamp = time.Now().Format(time.RFC3339)
		}
		a.APIID = getContextValue(context, APIIDKey)
	}

	err := json.NewEncoder(writer).Encode(a)
	if err != nil {
		logger := GetLogger(context)
		logger.Errorf("Could not send simple error to client %s", err.Error())
	}
}

type ContextKeyAPIID int

const APIIDKey ContextKeyAPIID = 1

func GetLogger(context context.Context) *log.Entry {
	fields := log.Fields{}
	if id := getContextValue(context, APIIDKey); id != "" {
		fields["api"] = id
	}
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
	if value := request.Header.Get("Content-Type"); value != "application/json" {
		return apiErrors.NewAPIError(http.StatusBadRequest, "Content-Type header is not application/json")
	}

	request.Body = http.MaxBytesReader(writer, request.Body, 1048576)

	dec := json.NewDecoder(request.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
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

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return apiErrors.NewBadRequestErrorF("Request body must only contain a single JSON object")
	}

	return nil
}
