package restAPI

import (
	"database/sql"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/exasol/extension-manager/extensionAPI"
	"github.com/exasol/extension-manager/extensionController"
	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RestAPISuite struct {
	suite.Suite
	assertJSON *jsonassert.Asserter
	controller *MockExtensionController
	restAPI    RestAPI
}

func TestRestApiSuite(t *testing.T) {
	suite.Run(t, new(RestAPISuite))
}

type MockExtensionController struct {
	mock.Mock
}

func (suite *RestAPISuite) SetupSuite() {
	suite.assertJSON = jsonassert.New(suite.T())
}

func (suite *RestAPISuite) SetupTest() {
	suite.controller = &MockExtensionController{}
	suite.restAPI = Create(suite.controller)
	go suite.restAPI.Serve()
	time.Sleep(10 * time.Millisecond) // give the server some time to become ready
}

func (suite *RestAPISuite) TearDownTest() {
	suite.restAPI.Stop()
}

func (mock *MockExtensionController) GetAllInstallations(dbConnection *sql.DB) ([]*extensionAPI.JsExtInstallation, error) {
	args := mock.Called(dbConnection)
	return args.Get(0).([]*extensionAPI.JsExtInstallation), args.Error(1)
}

func (mock *MockExtensionController) GetAllExtensions(dbConnectionWithNoAutocommit *sql.DB) ([]*extensionController.Extension, error) {
	args := mock.Called(dbConnectionWithNoAutocommit)
	return args.Get(0).([]*extensionController.Extension), args.Error(1)
}

func (suite *RestAPISuite) TestGetInstallations() {
	suite.controller.On("GetAllInstallations", mock.Anything).Return([]*extensionAPI.JsExtInstallation{{Name: "test", Version: "0.1.0", InstanceParameters: []interface{}{map[string]interface{}{"id": "param1", "name": "My param", "type": "string"}}}}, nil)
	responseString := suite.makeGetRequest("/installations?dbHost=host&dbPort=8563&dbUser=user&dbPass=password")
	suite.assertJSON.Assertf(responseString, `{"installations":[{"name":"test","version":"0.1.0","instanceParameters":[{"id":"param1","name":"My param","type":"string"}]}]}`)
}

func (suite *RestAPISuite) TestGetExtensions() {
	suite.controller.On("GetAllExtensions", mock.Anything).Return([]*extensionController.Extension{{Name: "my-extension", Description: "a cool extension", InstallableVersions: []string{"0.1.0"}}}, nil)
	responseString := suite.makeGetRequest("/extensions?dbHost=host&dbPort=8563&dbUser=user&dbPass=password")
	suite.assertJSON.Assertf(responseString, `{"extensions":[{"name":"my-extension","description":"a cool extension","installableVersions":["0.1.0"]}]}`)
}

func (suite *RestAPISuite) makeGetRequest(path string) string {
	const apiHost = "http://localhost:8080"
	response, err := http.Get(apiHost + path)
	suite.NoError(err)
	defer func() { suite.NoError(response.Body.Close()) }()
	bytes, err := ioutil.ReadAll(response.Body)
	suite.NoError(err)
	return string(bytes)
}
