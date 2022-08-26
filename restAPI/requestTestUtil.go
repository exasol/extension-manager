package restAPI

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/exasol/extension-manager/extensionController"
	"github.com/stretchr/testify/suite"
)

func startRestApi(suite *suite.Suite, controller extensionController.TransactionController) *baseRestAPITest {
	hostAndPort := "localhost:8081"
	api := baseRestAPITest{
		suite:   suite,
		restAPI: Create(controller, hostAndPort),
		baseUrl: fmt.Sprintf("http://%s", hostAndPort)}
	api.restAPI.StartInBackground()
	return &api
}

type baseRestAPITest struct {
	suite   *suite.Suite
	baseUrl string
	restAPI RestAPI
}

func (t *baseRestAPITest) makeRequestWithAuthHeader(method string, path string, authHeader string, body string, expectedStatusCode int) string {
	request, err := http.NewRequest(method, t.baseUrl+path, strings.NewReader(body))
	request.Header.Add("Authorization", authHeader)
	if body != "" {
		request.Header.Add("Content-Type", "application/json")
	}
	t.suite.NoError(err)
	response, err := http.DefaultClient.Do(request)
	t.suite.NoError(err)
	t.suite.Equal(expectedStatusCode, response.StatusCode)
	defer func() { t.suite.NoError(response.Body.Close()) }()
	bytes, err := io.ReadAll(response.Body)
	t.suite.NoError(err)
	responseBody := string(bytes)
	return responseBody
}

func createBasicAuthHeader(user, password string) string {
	basicAuth := base64.StdEncoding.EncodeToString([]byte(user + ":" + password))
	return "Basic " + basicAuth
}
