package restAPI

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/exasol/extension-manager/extensionController"
	"github.com/exasol/extension-manager/integrationTesting"
	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/suite"
)

const (
	EXTENSION_SCHEMA     = "test"
	DEFAULT_EXTENSION_ID = "testing-extension.js"
)

type RestAPIIntegrationTestSuite struct {
	integrationTesting.IntegrationTestSuite
	tempExtensionRepo string
	assertJSON        *jsonassert.Asserter
	restAPI           RestAPI
	baseUrl           string
}

func TestRestAPIIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(RestAPIIntegrationTestSuite))
}

func (suite *RestAPIIntegrationTestSuite) SetupSuite() {
	suite.IntegrationTestSuite.SetupSuite()
	suite.assertJSON = jsonassert.New(suite.T())
}

func (suite *RestAPIIntegrationTestSuite) SetupTest() {
	ctrl := extensionController.Create(suite.tempExtensionRepo, EXTENSION_SCHEMA)
	suite.restAPI = Create(ctrl, "localhost:8081")
	suite.baseUrl = "http://localhost:8081"
	go suite.restAPI.Serve()
	time.Sleep(10 * time.Millisecond) // give the server some time to become ready
}

func (suite *RestAPIIntegrationTestSuite) TearDownTest() {
	suite.restAPI.Stop()
}

func (suite *RestAPIIntegrationTestSuite) TestGetAllExtensionsSuccessfully() {
	responseString := suite.makeGetRequest("/extensions?" + suite.getDbArgs())
	suite.assertJSON.Assertf(responseString, `{"extensions":[]}`)
}

func (suite *RestAPIIntegrationTestSuite) TestGetInstallationsSuccessfully() {
	responseString := suite.makeGetRequest("/installations?" + suite.getDbArgs())
	suite.assertJSON.Assertf(responseString, `{"installations":[]}`)
}

func (suite *RestAPIIntegrationTestSuite) getDbArgs() string {
	info, err := suite.Exasol.GetConnectionInfo()
	if err != nil {
		suite.FailNowf("error getting connection info: %v", err.Error())
	}
	return fmt.Sprintf("dbHost=%s&dbPort=%d&dbUser=%s&dbPass=%s", info.Host, info.Port, info.User, info.Password)
}

func (suite *RestAPIIntegrationTestSuite) makeGetRequest(path string) string {
	return suite.makeRequest("GET", path, "", 200)
}

func (suite *RestAPIIntegrationTestSuite) makeRequest(method string, path string, body string, expectedStatusCode int) string {
	request, err := http.NewRequest(method, suite.baseUrl+path, strings.NewReader(body))
	suite.NoError(err)
	response, err := http.DefaultClient.Do(request)
	suite.NoError(err)
	suite.Equal(expectedStatusCode, response.StatusCode)
	defer func() { suite.NoError(response.Body.Close()) }()
	bytes, err := io.ReadAll(response.Body)
	suite.NoError(err)
	return string(bytes)
}
