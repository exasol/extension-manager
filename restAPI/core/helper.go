package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

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
	switch apiError := err.(type) {
	default:
		log.Errorf("Internal error: %s", err.Error())
		NewInternalServerError(err).(*APIError).Send(context, writer)
	case *APIError:
		if apiError.OriginalError != nil {
			log.Errorf("Error: %s (original: %s)", err.Error(), apiError.OriginalError.Error())
		}
		apiError.Send(context, writer)
	}
}

func NewInternalServerError(originalError error) error {
	return &APIError{
		Status:        http.StatusInternalServerError,
		Message:       "Internal server error",
		OriginalError: originalError,
	}
}

func NewBadRequestErrorF(format string, a ...interface{}) error {
	return NewAPIErrorF(http.StatusBadRequest, format, a...)
}

func NewAPIErrorF(status int, format string, a ...interface{}) error {
	return &APIError{
		Status:  status,
		Message: fmt.Sprintf(format, a...),
	}
}

func NewAPIError(status int, message string) error {
	return &APIError{
		Status:  status,
		Message: message,
	}
}

type ContextKeyAPIID int

const APIIDKey ContextKeyAPIID = 1

func GetLogger(context context.Context) *log.Entry {
	fields := log.Fields{}

	if id := GetContextValue(context, APIIDKey); id != "" {
		fields["api"] = id
	}

	if id := middleware.GetReqID(context); id != "" {
		fields["request"] = id
	}

	return log.WithFields(fields)
}

func GetContextValue(ctx context.Context, id interface{}) string {
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
		return NewAPIError(http.StatusBadRequest, "Content-Type header is not application/json")
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
			return NewBadRequestErrorF("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return NewBadRequestErrorF("Request body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			return NewBadRequestErrorF("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return NewBadRequestErrorF("Request body contains unknown field %q", fieldName)

		case errors.Is(err, io.EOF):
			return NewBadRequestErrorF("Request body must not be empty")

		case err.Error() == "http: request body too large":
			return NewBadRequestErrorF("Request body must not be larger than 1MB")

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return NewBadRequestErrorF("Request body must only contain a single JSON object")
	}

	return nil
}
