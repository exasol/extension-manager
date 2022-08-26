package restAPI_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/exasol/extension-manager/extensionAPI"
	"github.com/exasol/extension-manager/extensionController"
	"github.com/exasol/extension-manager/restAPI"
	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RestAPISuite struct {
	baseRestAPITest
	assertJSON *jsonassert.Asserter
	controller *MockExtensionController
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
	suite.restAPI = restAPI.Create(suite.controller, "localhost:8080")
	suite.baseUrl = "http://localhost:8080/api/v1"
	go suite.restAPI.Serve()
	time.Sleep(10 * time.Millisecond) // give the server some time to become ready
}

func (mock *MockExtensionController) InstallExtension(ctx context.Context, db *sql.DB, extensionId string, extensionVersion string) error {
	args := mock.Called(ctx, db, extensionId, extensionVersion)
	return args.Error(0)
}

func (mock *MockExtensionController) GetInstalledExtensions(ctx context.Context, db *sql.DB) ([]*extensionAPI.JsExtInstallation, error) {
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
	restAPI := restAPI.Create(controller, "localhost:8080")
	suite.Panics(restAPI.Stop)
}

var authSuccessTests = []struct{ authHeader string }{
	{authHeader: "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ=="},
	{authHeader: "Bearer token"}}

func (suite *RestAPISuite) TestGetInstallationsSuccessfully() {
	suite.controller.On("GetInstalledExtensions", mock.Anything, mock.Anything).Return([]*extensionAPI.JsExtInstallation{{Name: "test", Version: "0.1.0", InstanceParameters: []interface{}{map[string]interface{}{"id": "param1", "name": "My param", "type": "string"}}}}, nil)
	for _, test := range authSuccessTests {
		suite.Run(test.authHeader, func() {
			responseString := suite.makeRequestWithAuthHeader("GET", "/installations?dbHost=host&dbPort=8563&", test.authHeader, "", 200)
			suite.assertJSON.Assertf(responseString, `{"installations":[{"name":"test","version":"0.1.0","instanceParameters":[{"id":"param1","name":"My param","type":"string"}]}]}`)
		})
	}
}

func (suite *RestAPISuite) TestGetInstallationsFailed() {
	suite.controller.On("GetInstalledExtensions", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("mock error"))
	responseString := suite.makeRequest("GET", "/installations?dbHost=host&dbPort=8563", "", 500)
	suite.Regexp(`{"code":500,"message":"Internal server error",`, responseString)
}

func (suite *RestAPISuite) TestGetAllExtensionsSuccessfully() {
	suite.controller.On("GetAllExtensions", mock.Anything, mock.Anything).Return([]*extensionController.Extension{{Id: "ext-id", Name: "my-extension", Description: "a cool extension", InstallableVersions: []string{"0.1.0"}}}, nil)
	for _, test := range authSuccessTests {
		suite.Run(test.authHeader, func() {
			responseString := suite.makeRequestWithAuthHeader("GET", "/extensions?dbHost=host&dbPort=8563&", test.authHeader, "", 200)
			suite.assertJSON.Assertf(responseString, `{"extensions":[{"id": "ext-id", "name":"my-extension","description":"a cool extension","installableVersions":["0.1.0"]}]}`)
		})
	}
}

func (suite *RestAPISuite) TestGetAllExtensionsFails() {
	suite.controller.On("GetAllExtensions", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("mock error"))
	responseString := suite.makeRequest("GET", "/extensions?dbHost=host&dbPort=8563", "", 500)
	suite.Regexp(`{"code":500,"message":"Internal server error",.*`, responseString)
}

func (suite *RestAPISuite) TestInstallExtensionsSuccessfully() {
	suite.controller.On("InstallExtension", mock.Anything, mock.Anything, "ext-id", "ver").Return(nil)
	for _, test := range authSuccessTests {
		suite.Run(test.authHeader, func() {
			responseString := suite.makeRequestWithAuthHeader("PUT", "/installations?dbHost=host&dbPort=8563&", test.authHeader, `{"extensionId": "ext-id", "extensionVersion": "ver"}`, 204)
			suite.Equal("", responseString)
		})
	}
}

func (suite *RestAPISuite) TestInstallExtensionsFailed() {
	suite.controller.On("InstallExtension", mock.Anything, mock.Anything, "ext-id", "ver").Return(fmt.Errorf("mock error"))
	responseString := suite.makeRequest("PUT", "/installations?extensionId=ext-id&extensionVersion=ver&dbHost=host&dbPort=8563", `{"extensionId": "ext-id", "extensionVersion": "ver"}`, 500)
	suite.Regexp(`{"code":500,"message":"Internal server error",.*`, responseString)
}

func (suite *RestAPISuite) TestCreateInstanceSuccessfully() {
	suite.controller.On("CreateInstance", mock.Anything, mock.Anything, "ext-id", "ver", []extensionController.ParameterValue{{Name: "p1", Value: "v1"}}).Return("instanceName", nil)
	for _, test := range authSuccessTests {
		suite.Run(test.authHeader, func() {
			responseString := suite.makeRequestWithAuthHeader("PUT", "/instances?dbHost=host&dbPort=8563&", test.authHeader,
				`{"extensionId": "ext-id", "extensionVersion": "ver", "parameterValues": [{"name":"p1", "value":"v1"}]}`, 200)
			suite.Equal("{\"instanceName\":\"instanceName\"}\n", responseString)
		})
	}
}

func (suite *RestAPISuite) TestCreateInstanceFailed_invalidPayload() {
	suite.controller.On("CreateInstance", mock.Anything, mock.Anything, "ext-id", "ver", []extensionController.ParameterValue{{Name: "p1", Value: "v1"}}).Return("instanceName", nil)
	responseString := suite.makeRequest("PUT", "/instances?dbHost=host&dbPort=8563",
		`invalid payload`, 400)
	suite.Regexp("{\"code\":400,\"message\":\"Request body contains badly-formed JSON \\(at position 1\\)\".*", responseString)
}

func (suite *RestAPISuite) TestCreateInstanceFailed() {
	suite.controller.On("CreateInstance", mock.Anything, mock.Anything, "ext-id", "ver", []extensionController.ParameterValue{{Name: "p1", Value: "v1"}}).Return("", fmt.Errorf("mock error"))
	responseString := suite.makeRequest("PUT", "/instances?dbHost=host&dbPort=8563",
		`{"extensionId": "ext-id", "extensionVersion": "ver", "parameterValues": [{"name":"p1", "value":"v1"}]}`, 500)
	suite.Regexp(`{"code":500,"message":"Internal server error",.*`, responseString)
}

func (suite *RestAPISuite) TestRequestsFailForMissingParameters() {
	var tests = []struct {
		method        string
		url           string
		parameters    string
		expectedError string
	}{
		{"GET", "/extensions", "dbPort=8563", "missing parameter dbHost"},
		{"GET", "/extensions", "dbHost=host", "missing parameter dbPort"},
		{"GET", "/extensions", "dbHost=host&dbPort=invalidPort", "invalid value 'invalidPort' for parameter dbPort"},

		{"GET", "/installations", "dbPort=8563", "missing parameter dbHost"},
		{"GET", "/installations", "dbHost=host", "missing parameter dbPort"},
		{"GET", "/installations", "dbHost=host&dbPort=invalidPort", "invalid value 'invalidPort' for parameter dbPort"},

		{"PUT", "/installations", "extensionId=ext-id&extensionVersion=ver&dbPort=8563", "missing parameter dbHost"},
		{"PUT", "/installations", "extensionId=ext-id&extensionVersion=ver&dbHost=host", "missing parameter dbPort"},
		{"PUT", "/installations", "extensionId=ext-id&extensionVersion=ver&dbHost=host&dbPort=invalidPort", "invalid value 'invalidPort' for parameter dbPort"},

		{"PUT", "/instances", "dbPort=8563", "missing parameter dbHost"},
		{"PUT", "/instances", "dbHost=host", "missing parameter dbPort"},
		{"PUT", "/instances", "dbHost=host&dbPort=invalidPort", "invalid value 'invalidPort' for parameter dbPort"},
	}
	suite.controller.On("GetAllExtensions", mock.Anything, mock.Anything).Return([]*extensionController.Extension{{Name: "my-extension", Description: "a cool extension", InstallableVersions: []string{"0.1.0"}}}, nil)
	suite.controller.On("GetInstalledExtensions", mock.Anything, mock.Anything).Return([]*extensionAPI.JsExtInstallation{{Name: "test", Version: "0.1.0", InstanceParameters: []interface{}{map[string]interface{}{"id": "param1", "name": "My param", "type": "string"}}}}, nil)
	suite.controller.On("InstallExtension", mock.Anything, mock.Anything, "ext-id", "ver").Return(nil)
	suite.controller.On("CreateInstance", mock.Anything, mock.Anything, "ext-id", "ver", mock.Anything).Return("instanceName", nil)
	for _, test := range tests {
		suite.Run(fmt.Sprintf("Request %s %s?%s results in error message %q", test.method, test.url, test.parameters, test.expectedError), func() {
			completePath := fmt.Sprintf("%s?%s", test.url, test.parameters)
			responseString := suite.makeRequest(test.method, completePath, "", 400)
			suite.Regexp(fmt.Sprintf(`{"code":400,"message":"%s"`, test.expectedError), responseString)
		})
	}
}

func (suite *RestAPISuite) makeRequest(method, path, body string, expectedStatus int) string {
	authHeader := createBasicAuthHeader("user", "password")
	return suite.makeRequestWithAuthHeader(method, path, authHeader, body, expectedStatus)
}
