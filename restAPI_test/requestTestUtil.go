package restAPI_test

import (
	"encoding/base64"
	"io"
	"net/http"
	"strings"

	"github.com/exasol/extension-manager/restAPI"
	"github.com/stretchr/testify/suite"
)

type baseRestAPITest struct {
	suite.Suite
	baseUrl string
	restAPI restAPI.RestAPI
}

func (suite *baseRestAPITest) TearDownTest() {
	suite.restAPI.Stop()
}

func (suite *baseRestAPITest) makeRequestWithAuthHeader(method string, path string, authHeader string, body string, expectedStatusCode int) string {
	request, err := http.NewRequest(method, suite.baseUrl+path, strings.NewReader(body))
	request.Header.Add("Authorization", authHeader)
	if body != "" {
		request.Header.Add("Content-Type", "application/json")
	}
	suite.NoError(err)
	response, err := http.DefaultClient.Do(request)
	suite.NoError(err)
	suite.Equal(expectedStatusCode, response.StatusCode)
	defer func() { suite.NoError(response.Body.Close()) }()
	bytes, err := io.ReadAll(response.Body)
	suite.NoError(err)
	responseBody := string(bytes)
	return responseBody
}

func createBasicAuthHeader(user, password string) string {
	basicAuth := base64.StdEncoding.EncodeToString([]byte(user + ":" + password))
	return "Basic " + basicAuth
}
