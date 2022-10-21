package restAPI

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/Nightapes/go-rest/pkg/openapi"

	"github.com/exasol/extension-manager/apiErrors"
	"github.com/exasol/extension-manager/extensionController"

	httpswagger "github.com/swaggo/http-swagger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
)

func setupStandaloneAPI(controller extensionController.TransactionController) (http.Handler, *openapi.API, error) {
	api, err := CreateOpenApi()
	if err != nil {
		return nil, nil, err
	}
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(loggerMiddleware())
	r.Use(middleware.Recoverer)

	err = addPublicEndpointsWithController(api, controller)
	if err != nil {
		return nil, nil, err
	}

	openApiHandler, err := api.OpenAPIHandlerFunc()
	if err != nil {
		return nil, nil, err
	}
	r.Group(func(r chi.Router) {
		r.MethodFunc(http.MethodGet, "/openapi.json", openApiHandler)
		r.Method(http.MethodGet, "/openapi/*", httpswagger.Handler(
			httpswagger.URL("/openapi.json"),
		))
	})

	r.Group(func(r chi.Router) {
		for _, handleConfig := range api.GetHandleFunc() {
			log.Tracef("Add func %s %s", handleConfig.Method, handleConfig.Path)
			r.With(middleware.Timeout(60*time.Second)).Method(handleConfig.Method, handleConfig.Path, handleConfig.HandlerFunc)
		}
	})

	return r, api, nil
}

func CreateOpenApi() (*openapi.API, error) {
	api := openapi.NewOpenAPI()
	api.Title = "Exasol Extension Manager REST-API"
	api.Description = "Managed extensions and instances of extensions"
	api.Version = "1.0"
	if err := api.WithBasicAuth(BasicAuth); err != nil {
		return nil, err
	}
	if err := api.WithBearerAuth(BearerAuth, "bearer", "JWT"); err != nil {
		return nil, err
	}
	api.DefaultResponse(&openapi.MethodResponse{
		Description: "Default error",
		Value: &apiErrors.APIError{
			Status:        500,
			Message:       "Something went wrong.",
			RequestID:     "Rn3x8gcEInnHt205B4c7QZ",
			Timestamp:     "2021-01-02T15:04:05Z07:00",
			APIID:         "",
			OriginalError: nil,
		},
	})
	return api, nil
}

func loggerMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			buf := &bytes.Buffer{}
			logger := GetLogger(r.Context())
			t1 := time.Now()
			defer func() {

				scheme := "http"
				if r.TLS != nil {
					scheme = "https"
				}

				fmt.Fprintf(buf, "\"%s %s://%s%s %s\" ", r.Method, scheme, r.Host, r.RequestURI, r.Proto)
				buf.WriteString("from ")
				buf.WriteString(r.RemoteAddr)
				buf.WriteString(" - ")
				fmt.Fprintf(buf, "%d %dB in %s", ww.Status(), ww.BytesWritten(), time.Since(t1))
				logger.Infof(buf.String())
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
