package restAPI

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"strings"
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
	baseUrl    string
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
	suite.restAPI = Create(suite.controller, "localhost:8080")
	suite.baseUrl = "http://localhost:8080"
	go suite.restAPI.Serve()
	time.Sleep(10 * time.Millisecond) // give the server some time to become ready
}

func (suite *RestAPISuite) TearDownTest() {
	suite.restAPI.Stop()
}

func (mock *MockExtensionController) InstallExtension(ctx context.Context, db *sql.DB, extensionId string, extensionVersion string) error {
	args := mock.Called(ctx, db, extensionId, extensionVersion)
	return args.Error(0)
}

func (mock *MockExtensionController) GetAllInstallations(ctx context.Context, db *sql.DB) ([]*extensionAPI.JsExtInstallation, error) {
	args := mock.Called(ctx, db)
	if installations, ok := args.Get(0).([]*extensionAPI.JsExtInstallation); ok {
		return installations, args.Error(1)
	}
	return nil, args.Error(1)
}

func (mock *MockExtensionController) GetAllExtensions(ctx context.Context, db *sql.DB) ([]*extensionController.Extension, error) {
	args := mock.Called(ctx, db)
	if extensions, ok := args.Get(0).([]*extensionController.Extension); ok {
		return extensions, args.Error(1)
	} else {
		return nil, args.Error(1)
	}
}

func (mock *MockExtensionController) CreateInstance(ctx context.Context, db *sql.DB, extensionId string, extensionVersion string, parameterValues []extensionController.ParameterValue) (string, error) {
	args := mock.Called(ctx, db, extensionId, extensionVersion, parameterValues)
	return args.String(0), args.Error(1)
}

func (suite *RestAPISuite) TestStopWithoutStartFails() {
	controller := &MockExtensionController{}
	restAPI := Create(controller, "localhost:8080")
	suite.Panics(restAPI.Stop)
}

func (suite *RestAPISuite) TestInvalidParameterFormat() {
	responseString := suite.makeRequest("GET", "/extensions?dbHost=host;invalid&dbPort=8563&dbUser=user&dbPassword=password", "", 500)
	suite.Equal("Request failed: missing parameter dbHost", responseString)
}

var authSuccessTests = []struct{ parameters string }{
	{parameters: "dbUser=user&dbPassword=password"},
	{parameters: "dbAccessToken=token"},
	{parameters: "dbRefreshToken=token"}}

func (suite *RestAPISuite) TestGetInstallationsSuccessfully() {
	suite.controller.On("GetAllInstallations", mock.Anything, mock.Anything).Return([]*extensionAPI.JsExtInstallation{{Name: "test", Version: "0.1.0", InstanceParameters: []interface{}{map[string]interface{}{"id": "param1", "name": "My param", "type": "string"}}}}, nil)
	for _, test := range authSuccessTests {
		suite.Run(test.parameters, func() {
			responseString := suite.makeGetRequest("/installations?dbHost=host&dbPort=8563&" + test.parameters)
			suite.assertJSON.Assertf(responseString, `{"installations":[{"name":"test","version":"0.1.0","instanceParameters":[{"id":"param1","name":"My param","type":"string"}]}]}`)
		})
	}
}

func (suite *RestAPISuite) TestGetInstallationsFailed() {
	suite.controller.On("GetAllInstallations", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("mock error"))
	responseString := suite.makeRequest("GET", "/installations?dbHost=host&dbPort=8563&dbUser=user&dbPassword=password", "", 500)
	suite.Equal("Request failed: mock error", responseString)
}

func (suite *RestAPISuite) TestGetAllExtensionsSuccessfully() {
	suite.controller.On("GetAllExtensions", mock.Anything, mock.Anything).Return([]*extensionController.Extension{{Id: "ext-id", Name: "my-extension", Description: "a cool extension", InstallableVersions: []string{"0.1.0"}}}, nil)
	for _, test := range authSuccessTests {
		suite.Run(test.parameters, func() {
			responseString := suite.makeGetRequest("/extensions?dbHost=host&dbPort=8563&" + test.parameters)
			suite.assertJSON.Assertf(responseString, `{"extensions":[{"id": "ext-id", "name":"my-extension","description":"a cool extension","installableVersions":["0.1.0"]}]}`)
		})
	}
}

func (suite *RestAPISuite) TestGetAllExtensionsFails() {
	suite.controller.On("GetAllExtensions", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("mock error"))
	responseString := suite.makeRequest("GET", "/extensions?dbHost=host&dbPort=8563&dbUser=user&dbPassword=password", "", 500)
	suite.Equal("Request failed: mock error", responseString)
}

func (suite *RestAPISuite) TestInstallExtensionsSuccessfully() {
	suite.controller.On("InstallExtension", mock.Anything, mock.Anything, "ext-id", "ver").Return(nil)
	for _, test := range authSuccessTests {
		suite.Run(test.parameters, func() {
			responseString := suite.makePutRequest("/installations?extensionId=ext-id&extensionVersion=ver&dbHost=host&dbPort=8563&" + test.parameters)
			suite.Equal("", responseString)
		})
	}
}

func (suite *RestAPISuite) TestInstallExtensionsFailed() {
	suite.controller.On("InstallExtension", mock.Anything, mock.Anything, "ext-id", "ver").Return(fmt.Errorf("mock error"))
	responseString := suite.makeRequest("PUT", "/installations?extensionId=ext-id&extensionVersion=ver&dbHost=host&dbPort=8563&dbUser=user&dbPassword=password", "", 500)
	suite.Equal("Request failed: error installing extension: mock error", responseString)
}

func (suite *RestAPISuite) TestCreateInstanceSuccessfully() {
	suite.controller.On("CreateInstance", mock.Anything, mock.Anything, "ext-id", "ver", []extensionController.ParameterValue{{Name: "p1", Value: "v1"}}).Return("instanceName", nil)
	for _, test := range authSuccessTests {
		suite.Run(test.parameters, func() {
			responseString := suite.makeRequest("PUT", "/instances?dbHost=host&dbPort=8563&"+test.parameters,
				`{"extensionId": "ext-id", "extensionVersion": "ver", "parameterValues": [{"name":"p1", "value":"v1"}]}`, 200)
			suite.Equal(`{"instanceName":"instanceName"}`, responseString)
		})
	}
}

func (suite *RestAPISuite) TestCreateInstanceFailed_invalidPayload() {
	suite.controller.On("CreateInstance", mock.Anything, mock.Anything, "ext-id", "ver", []extensionController.ParameterValue{{Name: "p1", Value: "v1"}}).Return("instanceName", nil)
	responseString := suite.makeRequest("PUT", "/instances?dbHost=host&dbPort=8563&dbUser=user&dbPassword=password",
		`invalid payload`, 400)
	suite.Equal("Request failed: invalid request: invalid character 'i' looking for beginning of value", responseString)
}

func (suite *RestAPISuite) TestCreateInstanceFailed() {
	suite.controller.On("CreateInstance", mock.Anything, mock.Anything, "ext-id", "ver", []extensionController.ParameterValue{{Name: "p1", Value: "v1"}}).Return("", fmt.Errorf("mock error"))
	responseString := suite.makeRequest("PUT", "/instances?dbHost=host&dbPort=8563&dbUser=user&dbPassword=password",
		`{"extensionId": "ext-id", "extensionVersion": "ver", "parameterValues": [{"name":"p1", "value":"v1"}]}`, 500)
	suite.Equal("Request failed: error installing extension: mock error", responseString)
}

func (suite *RestAPISuite) TestRequestsFailForMissingParameters() {
	var tests = []struct {
		method        string
		url           string
		parameters    string
		expectedError string
	}{
		{"GET", "/extensions", "dbPort=8563&dbUser=user&dbPassword=password", "missing parameter dbHost"},
		{"GET", "/extensions", "dbHost=host&dbUser=user&dbPassword=password", "missing parameter dbPort"},
		{"GET", "/extensions", "dbHost=host&dbPort=invalidPort&dbUser=user&dbPassword=password", "invalid value \"invalidPort\" for parameter dbPort"},
		{"GET", "/extensions", "dbHost=host&dbPort=8563&dbPassword=password", "missing parameter dbUser"},
		{"GET", "/extensions", "dbHost=host&dbPort=8563&dbUser=user", "missing parameter dbPassword"},

		{"GET", "/installations", "dbPort=8563&dbUser=user&dbPassword=password", "missing parameter dbHost"},
		{"GET", "/installations", "dbHost=host&dbUser=user&dbPassword=password", "missing parameter dbPort"},
		{"GET", "/installations", "dbHost=host&dbPort=invalidPort&dbUser=user&dbPassword=password", "invalid value \"invalidPort\" for parameter dbPort"},
		{"GET", "/installations", "dbHost=host&dbPort=8563&dbPassword=password", "missing parameter dbUser"},
		{"GET", "/installations", "dbHost=host&dbPort=8563&dbUser=user", "missing parameter dbPassword"},

		{"PUT", "/installations", "extensionId=ext-id&extensionVersion=ver&dbPort=8563&dbUser=user&dbPassword=password", "missing parameter dbHost"},
		{"PUT", "/installations", "extensionId=ext-id&extensionVersion=ver&dbHost=host&dbUser=user&dbPassword=password", "missing parameter dbPort"},
		{"PUT", "/installations", "extensionId=ext-id&extensionVersion=ver&dbHost=host&dbPort=invalidPort&dbUser=user&dbPassword=password", "invalid value \"invalidPort\" for parameter dbPort"},
		{"PUT", "/installations", "extensionId=ext-id&extensionVersion=ver&dbHost=host&dbPort=8563&dbPassword=password", "missing parameter dbUser"},
		{"PUT", "/installations", "extensionId=ext-id&dbHost=host&dbPort=8563&dbUser=user&dbPassword=password", "missing parameter extensionVersion"},
		{"PUT", "/installations", "extensionVersion=ver&dbHost=host&dbPort=8563&dbUser=user&dbPassword=password", "missing parameter extensionId"},

		{"PUT", "/instances", "dbPort=8563&dbUser=user&dbPassword=password", "missing parameter dbHost"},
		{"PUT", "/instances", "dbHost=host&dbUser=user&dbPassword=password", "missing parameter dbPort"},
		{"PUT", "/instances", "dbHost=host&dbPort=invalidPort&dbUser=user&dbPassword=password", "invalid value \"invalidPort\" for parameter dbPort"},
		{"PUT", "/instances", "dbHost=host&dbPort=8563&dbPassword=password", "missing parameter dbUser"},
		{"PUT", "/instances", "dbHost=host&dbPort=8563&dbUser=user", "missing parameter dbPassword"},
	}
	suite.controller.On("GetAllExtensions", mock.Anything, mock.Anything).Return([]*extensionController.Extension{{Name: "my-extension", Description: "a cool extension", InstallableVersions: []string{"0.1.0"}}}, nil)
	suite.controller.On("GetAllInstallations", mock.Anything, mock.Anything).Return([]*extensionAPI.JsExtInstallation{{Name: "test", Version: "0.1.0", InstanceParameters: []interface{}{map[string]interface{}{"id": "param1", "name": "My param", "type": "string"}}}}, nil)
	suite.controller.On("InstallExtension", mock.Anything, mock.Anything, "ext-id", "ver").Return(nil)
	suite.controller.On("CreateInstance", mock.Anything, mock.Anything, "ext-id", "ver", mock.Anything).Return("instanceName", nil)
	for _, test := range tests {
		suite.Run(fmt.Sprintf("Request %s %s?%s results in error message %q", test.method, test.url, test.parameters, test.expectedError), func() {
			completePath := fmt.Sprintf("%s?%s", test.url, test.parameters)
			responseString := suite.makeRequest(test.method, completePath, "", 500)
			suite.Equal(fmt.Sprintf("Request failed: %s", test.expectedError), responseString, fmt.Sprintf("Expected request %s to fail", completePath))
		})
	}
}

func (suite *RestAPISuite) makeGetRequest(path string) string {
	return suite.makeRequest("GET", path, "", 200)
}

func (suite *RestAPISuite) makePutRequest(path string) string {
	return suite.makeRequest("PUT", path, "", 200)
}

func (suite *RestAPISuite) makeRequest(method string, path string, body string, expectedStatusCode int) string {
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
