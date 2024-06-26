package restAPI

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/exasol/exasol-driver-go"
	"github.com/exasol/exasol-driver-go/pkg/dsn"
	"github.com/exasol/extension-manager/pkg/apiErrors"
)

type generalHandlerFunc = func(writer http.ResponseWriter, request *http.Request)
type dbHandler = func(db *sql.DB, writer http.ResponseWriter, request *http.Request) error

func adaptDbHandler(apiContext *ApiContext, handler dbHandler) generalHandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		db, err := openDBRequest(request)
		if err != nil {
			handleError(request.Context(), apiContext, writer, err)
			return
		}
		defer closeDBRequest(db)
		err = handler(db, writer, request)
		if err != nil {
			handleError(request.Context(), apiContext, writer, err)
		}
	}
}

func openDBRequest(request *http.Request) (*sql.DB, error) {
	config, err := createDbConfig(request)
	if err != nil {
		return nil, err
	}
	config.ValidateServerCertificate(false)
	config.Autocommit(false)
	database, err := sql.Open("exasol", config.String())
	if err != nil {
		return nil, fmt.Errorf("failed to open a database connection. Cause: %w", err)
	}
	return database, nil
}

func createDbConfig(request *http.Request) (*dsn.DSNConfigBuilder, error) {
	query := request.URL.Query()
	config, err := createDbConfigWithAuthentication(request)
	if err != nil {
		return nil, err
	}

	if host := query.Get("dbHost"); host == "" {
		return nil, apiErrors.NewBadRequestErrorF("missing parameter dbHost")
	} else {
		config.Host(host)
	}

	if portString := query.Get("dbPort"); portString == "" {
		return nil, apiErrors.NewBadRequestErrorF("missing parameter dbPort")
	} else {
		if port, err := strconv.Atoi(portString); err != nil {
			return nil, apiErrors.NewBadRequestErrorF("invalid value '%s' for parameter dbPort", portString)
		} else {
			config.Port(port)
		}
	}
	return config, nil
}

func createDbConfigWithAuthentication(request *http.Request) (*dsn.DSNConfigBuilder, error) {
	auth := request.Header.Get("Authorization")
	if auth == "" {
		return nil, apiErrors.NewUnauthorizedErrorF("missing Authorization header")
	}
	parts := strings.Split(auth, " ")
	if len(parts) < 2 {
		return nil, apiErrors.NewUnauthorizedErrorF("invalid Authorization header %q", auth)
	}
	scheme := parts[0]
	switch scheme {
	case "Basic":
		return newUserPasswordConfig(parts[1])
	case "Bearer":
		return exasol.NewConfigWithAccessToken(parts[1]), nil
	default:
		return nil, apiErrors.NewUnauthorizedErrorF("invalid Authorization scheme %q", parts[0])
	}
}

func newUserPasswordConfig(basicAuthCredentials string) (*dsn.DSNConfigBuilder, error) {
	user, password, err := extractUserPassword(basicAuthCredentials)
	if err != nil {
		return nil, err
	}
	return exasol.NewConfig(user, password), nil
}

func extractUserPassword(basicAuthCredentials string) (string, string, error) {
	data, err := base64.StdEncoding.DecodeString(basicAuthCredentials)
	if err != nil {
		return "", "", apiErrors.NewUnauthorizedErrorF("invalid basic auth header %q: %v", basicAuthCredentials, err)
	}
	userPassword := string(data)
	colon := strings.Index(userPassword, ":")
	if colon < 0 {
		return "", "", apiErrors.NewUnauthorizedErrorF("colon missing in basic auth header")
	}
	user := userPassword[:colon]
	password := userPassword[colon+1:]
	return user, password, nil
}

func closeDBRequest(db *sql.DB) {
	err := db.Close()
	if err != nil {
		// Strange but not critical. So we just log it and go on.
		log.Printf("failed to close db connection. Cause %v", err)
	}
}
