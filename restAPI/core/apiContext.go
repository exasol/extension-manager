package core

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/exasol/exasol-driver-go"
	ctrl "github.com/exasol/extension-manager/extensionController"
)

type ApiContext interface {
	OpenDBConnection(request *http.Request) (*sql.DB, error)
	Controller() ctrl.TransactionController
}

func NewApiContext(controller ctrl.TransactionController) ApiContext {
	return &contextImpl{controller: controller}
}

type contextImpl struct {
	controller ctrl.TransactionController
}

func (c *contextImpl) OpenDBConnection(request *http.Request) (*sql.DB, error) {
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

func createDbConfig(request *http.Request) (*exasol.DSNConfigBuilder, error) {
	query := request.URL.Query()
	config, err := createDbConfigWithAuthentication(request)
	if err != nil {
		return nil, err
	}

	if host := query.Get("dbHost"); host == "" {
		return nil, NewBadRequestErrorF("missing parameter dbHost")
	} else {
		config.Host(host)
	}

	if portString := query.Get("dbPort"); portString == "" {
		return nil, NewBadRequestErrorF("missing parameter dbPort")
	} else {
		if port, err := strconv.Atoi(portString); err != nil {
			return nil, NewBadRequestErrorF("invalid value %q for parameter dbPort", portString)
		} else {
			config.Port(port)
		}
	}
	return config, nil
}

func createDbConfigWithAuthentication(request *http.Request) (*exasol.DSNConfigBuilder, error) {
	query := request.URL.Query()
	accessToken := query.Get("dbAccessToken")
	if accessToken != "" {
		return exasol.NewConfigWithAccessToken(accessToken), nil
	}

	refreshToken := query.Get("dbRefreshToken")
	if refreshToken != "" {
		return exasol.NewConfigWithRefreshToken(refreshToken), nil
	}

	user := query.Get("dbUser")
	if user == "" {
		return nil, NewBadRequestErrorF("missing parameter dbUser")
	}

	password := query.Get("dbPassword")
	if password == "" {
		return nil, NewBadRequestErrorF("missing parameter dbPassword")
	}

	return exasol.NewConfig(user, password), nil
}

func (c *contextImpl) Controller() ctrl.TransactionController {
	return c.controller
}

func CloseDbConnection(db *sql.DB) {
	err := db.Close()
	if err != nil {
		// Strange but not critical. So we just log it and go on.
		fmt.Printf("failed to close db connection. Cause %v", err)
	}
}
