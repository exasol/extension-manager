package core

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

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
