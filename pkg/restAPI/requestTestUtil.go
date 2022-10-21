package restAPI

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/exasol/extension-manager/src/extensionController"
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
	t.suite.T().Helper()
	url := t.baseUrl + path
	request, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		t.suite.FailNowf("Creating request %s %s failed: %v", method, url, err)
	}
	request.Header.Add("Authorization", authHeader)
	if body != "" {
		request.Header.Add("Content-Type", "application/json")
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.suite.FailNowf("Request %s %s failed: %v", method, url, err)
	}
	t.suite.Equal(expectedStatusCode, response.StatusCode, "Got response status %q", response.Status)
	defer func() { t.suite.NoError(response.Body.Close()) }()
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		t.suite.FailNowf("Reading body failed: %v", err.Error())
	}
	responseBody := string(bytes)
	return responseBody
}

func createBasicAuthHeader(user, password string) string {
	basicAuth := base64.StdEncoding.EncodeToString([]byte(user + ":" + password))
	return "Basic " + basicAuth
}

func createBearerAuthHeader(token string) string {
	return "Bearer " + token
}
