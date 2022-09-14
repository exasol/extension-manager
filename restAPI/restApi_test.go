package restAPI

import (
	"fmt"
	"testing"

	"github.com/exasol/extension-manager/apiErrors"
	"github.com/exasol/extension-manager/extensionAPI"
	"github.com/exasol/extension-manager/extensionController"
	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RestAPISuite struct {
	suite.Suite
	restApi    *baseRestAPITest
	assertJSON *jsonassert.Asserter
	controller *mockExtensionController
}

func TestRestApiSuite(t *testing.T) {
	suite.Run(t, new(RestAPISuite))
}

func (suite *RestAPISuite) SetupSuite() {
	suite.assertJSON = jsonassert.New(suite.T())
}

func (suite *RestAPISuite) SetupTest() {
	suite.controller = &mockExtensionController{}
	suite.restApi = startRestApi(&suite.Suite, suite.controller)
}

func (suite *RestAPISuite) TearDownTest() {
	suite.restApi.restAPI.Stop()
}

func (suite *RestAPISuite) TestStopWithoutStartFails() {
	controller := &mockExtensionController{}
	restAPI := Create(controller, "localhost:8082")
	suite.Panics(restAPI.Stop)
}

var authSuccessTests = []struct{ authHeader string }{
	{authHeader: "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ=="},
	{authHeader: "Bearer token"}}

func (suite *RestAPISuite) TestGetInstallationsSuccessfully() {
	suite.controller.On("GetInstalledExtensions", mock.Anything, mock.Anything).Return([]*extensionAPI.JsExtInstallation{{Name: "test", Version: "0.1.0", InstanceParameters: []interface{}{map[string]interface{}{"id": "param1", "name": "My param", "type": "string"}}}}, nil)
	for _, test := range authSuccessTests {
		suite.Run(test.authHeader, func() {
			responseString := suite.restApi.makeRequestWithAuthHeader("GET", "/api/v1/installations?dbHost=host&dbPort=8563&", test.authHeader, "", 200)
			suite.assertJSON.Assertf(responseString, `{"installations":[{"name":"test","version":"0.1.0","instanceParameters":[{"id":"param1","name":"My param","type":"string"}]}]}`)
		})
	}
}

func (suite *RestAPISuite) TestGetInstallationsFailed() {
	suite.controller.On("GetInstalledExtensions", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("mock error"))
	responseString := suite.makeRequest("GET", "/api/v1/installations?dbHost=host&dbPort=8563", "", 500)
	suite.Regexp(`{"code":500,"message":"Internal server error",`, responseString)
}

func (suite *RestAPISuite) TestGetAllExtensionsSuccessfully() {
	suite.controller.On("GetAllExtensions", mock.Anything, mock.Anything).Return([]*extensionController.Extension{{Id: "ext-id", Name: "my-extension", Description: "a cool extension", InstallableVersions: []string{"0.1.0"}}}, nil)
	for _, test := range authSuccessTests {
		suite.Run(test.authHeader, func() {
			responseString := suite.restApi.makeRequestWithAuthHeader("GET", "/api/v1/extensions?dbHost=host&dbPort=8563&", test.authHeader, "", 200)
			suite.assertJSON.Assertf(responseString, `{"extensions":[{"id": "ext-id", "name":"my-extension","description":"a cool extension","installableVersions":["0.1.0"]}]}`)
		})
	}
}

func (suite *RestAPISuite) TestGetAllExtensionsFails() {
	suite.controller.On("GetAllExtensions", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("mock error"))
	responseString := suite.makeRequest("GET", "/api/v1/extensions?dbHost=host&dbPort=8563", "", 500)
	suite.Regexp(`{"code":500,"message":"Internal server error",.*`, responseString)
}

// Install extension

func (suite *RestAPISuite) TestInstallExtensionsSuccessfully() {
	suite.controller.On("InstallExtension", mock.Anything, mock.Anything, "ext-id", "ver").Return(nil)
	for _, test := range authSuccessTests {
		suite.Run(test.authHeader, func() {
			responseString := suite.restApi.makeRequestWithAuthHeader("PUT", "/api/v1/installations?dbHost=host&dbPort=8563&", test.authHeader, `{"extensionId": "ext-id", "extensionVersion": "ver"}`, 204)
			suite.Equal("", responseString)
		})
	}
}

func (suite *RestAPISuite) TestInstallExtensionsFailed() {
	suite.controller.On("InstallExtension", mock.Anything, mock.Anything, "ext-id", "ver").Return(fmt.Errorf("mock error"))
	responseString := suite.makeRequest("PUT", "/api/v1/installations?extensionId=ext-id&extensionVersion=ver&dbHost=host&dbPort=8563", `{"extensionId": "ext-id", "extensionVersion": "ver"}`, 500)
	suite.Regexp(`{"code":500,"message":"Internal server error",.*`, responseString)
}

// Uninstall extension

func (suite *RestAPISuite) TestUninstallExtensionsSuccessfully() {
	suite.controller.On("UninstallExtension", mock.Anything, mock.Anything, "ext-id", "ver").Return(nil)
	for _, test := range authSuccessTests {
		suite.Run(test.authHeader, func() {
			responseString := suite.restApi.makeRequestWithAuthHeader("DELETE", "/api/v1/extension/ext-id/version/ver?dbHost=host&dbPort=8563&", test.authHeader, "", 204)
			suite.Equal("", responseString)
		})
	}
}

func (suite *RestAPISuite) TestUninstallExtensionsFailed() {
	suite.controller.On("UninstallExtension", mock.Anything, mock.Anything, "ext-id", "ver").Return(fmt.Errorf("mock error"))
	responseString := suite.makeRequest("DELETE", "/api/v1/extension/ext-id/version/ver?extensionId=ext-id&extensionVersion=ver&dbHost=host&dbPort=8563", "", 500)
	suite.Regexp(`{"code":500,"message":"Internal server error",.*`, responseString)
}

// Create instance

func (suite *RestAPISuite) TestCreateInstanceSuccessfully() {
	suite.controller.On("CreateInstance", mock.Anything, mock.Anything, "ext-id", "ver", []extensionController.ParameterValue{{Name: "p1", Value: "v1"}}).Return(&extensionAPI.JsExtInstance{Id: "instId", Name: "instName"}, nil)
	for _, test := range authSuccessTests {
		suite.Run(test.authHeader, func() {
			responseString := suite.restApi.makeRequestWithAuthHeader("PUT", "/api/v1/instances?dbHost=host&dbPort=8563&", test.authHeader,
				`{"extensionId": "ext-id", "extensionVersion": "ver", "parameterValues": [{"name":"p1", "value":"v1"}]}`, 200)
			suite.Equal(`{"instanceId":"instId","instanceName":"instName"}`+"\n", responseString)
		})
	}
}

func (suite *RestAPISuite) TestCreateInstanceFailed_invalidPayload() {
	suite.controller.On("CreateInstance", mock.Anything, mock.Anything, "ext-id", "ver", []extensionController.ParameterValue{{Name: "p1", Value: "v1"}}).Return(&extensionAPI.JsExtInstance{Id: "instId", Name: "instName"}, nil)
	responseString := suite.makeRequest("PUT", "/api/v1/instances?dbHost=host&dbPort=8563",
		`invalid payload`, 400)
	suite.Regexp("{\"code\":400,\"message\":\"Request body contains badly-formed JSON \\(at position 1\\)\".*", responseString)
}

func (suite *RestAPISuite) TestCreateInstanceFailed() {
	suite.controller.On("CreateInstance", mock.Anything, mock.Anything, "ext-id", "ver", []extensionController.ParameterValue{{Name: "p1", Value: "v1"}}).Return(nil, fmt.Errorf("mock error"))
	responseString := suite.makeRequest("PUT", "/api/v1/instances?dbHost=host&dbPort=8563",
		`{"extensionId": "ext-id", "extensionVersion": "ver", "parameterValues": [{"name":"p1", "value":"v1"}]}`, 500)
	suite.Regexp(`{"code":500,"message":"Internal server error",.*`, responseString)
}

// List instances

func (suite *RestAPISuite) TestListInstancesSuccessfully() {
	suite.controller.On("FindInstances", mock.Anything, mock.Anything, "ext-id", "ver").Return([]*extensionAPI.JsExtInstance{{Id: "instId", Name: "instName"}}, nil)
	for _, test := range authSuccessTests {
		suite.Run(test.authHeader, func() {
			responseString := suite.restApi.makeRequestWithAuthHeader("GET", "/api/v1/extension/ext-id/ver/instances?dbHost=host&dbPort=8563&", test.authHeader, "", 200)
			suite.Equal(`{"instances":[{"id":"instId","name":"instName"}]}`+"\n", responseString)
		})
	}
}

func (suite *RestAPISuite) TestListInstancesFailed_genericError() {
	suite.controller.On("FindInstances", mock.Anything, mock.Anything, "ext-id", "ver").Return(nil, fmt.Errorf("mock"))
	responseString := suite.restApi.makeRequestWithAuthHeader("GET", "/api/v1/extension/ext-id/ver/instances?dbHost=host&dbPort=8563&", "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==", "", 500)
	suite.Contains(responseString, "{\"code\":500,\"message\":\"Internal server error\"")
}

func (suite *RestAPISuite) TestListInstancesFailed_apiError() {
	suite.controller.On("FindInstances", mock.Anything, mock.Anything, "ext-id", "ver").Return(nil, apiErrors.NewAPIError(432, "mock"))
	responseString := suite.restApi.makeRequestWithAuthHeader("GET", "/api/v1/extension/ext-id/ver/instances?dbHost=host&dbPort=8563&", "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==", "", 432)
	suite.Contains(responseString, "{\"code\":432,\"message\":\"mock\",")
}

// Delete instance

func (suite *RestAPISuite) TestDeleteInstanceSuccessfully() {
	suite.controller.On("DeleteInstance", mock.Anything, mock.Anything, "ext-id", "inst-id").Return(nil)
	for _, test := range authSuccessTests {
		suite.Run(test.authHeader, func() {
			responseString := suite.restApi.makeRequestWithAuthHeader("DELETE", "/api/v1/extension/ext-id/instance/inst-id?dbHost=host&dbPort=8563&", test.authHeader, "", 204)
			suite.Equal("", responseString)
		})
	}
}

func (suite *RestAPISuite) TestDeleteInstanceFailed_genericError() {
	suite.controller.On("DeleteInstance", mock.Anything, mock.Anything, "ext-id", "inst-id").Return(fmt.Errorf("mock"))
	responseString := suite.restApi.makeRequestWithAuthHeader("DELETE", "/api/v1/extension/ext-id/instance/inst-id?dbHost=host&dbPort=8563&", "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==", "", 500)
	suite.Contains(responseString, "{\"code\":500,\"message\":\"Internal server error\"")
}

func (suite *RestAPISuite) TestDeleteInstanceFailed_apiError() {
	suite.controller.On("DeleteInstance", mock.Anything, mock.Anything, "ext-id", "inst-id").Return(apiErrors.NewAPIError(432, "mock"))
	responseString := suite.restApi.makeRequestWithAuthHeader("DELETE", "/api/v1/extension/ext-id/instance/inst-id?dbHost=host&dbPort=8563&", "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==", "", 432)
	suite.Contains(responseString, "{\"code\":432,\"message\":\"mock\",")
}

func (suite *RestAPISuite) TestRequestsFailForMissingParameters() {
	var tests = []struct {
		method        string
		url           string
		parameters    string
		expectedError string
	}{
		{"GET", "/api/v1/extensions", "dbPort=8563", "missing parameter dbHost"},
		{"GET", "/api/v1/extensions", "dbHost=host", "missing parameter dbPort"},
		{"GET", "/api/v1/extensions", "dbHost=host&dbPort=invalidPort", "invalid value 'invalidPort' for parameter dbPort"},

		{"GET", "/api/v1/installations", "dbPort=8563", "missing parameter dbHost"},
		{"GET", "/api/v1/installations", "dbHost=host", "missing parameter dbPort"},
		{"GET", "/api/v1/installations", "dbHost=host&dbPort=invalidPort", "invalid value 'invalidPort' for parameter dbPort"},

		{"PUT", "/api/v1/installations", "extensionId=ext-id&extensionVersion=ver&dbPort=8563", "missing parameter dbHost"},
		{"PUT", "/api/v1/installations", "extensionId=ext-id&extensionVersion=ver&dbHost=host", "missing parameter dbPort"},
		{"PUT", "/api/v1/installations", "extensionId=ext-id&extensionVersion=ver&dbHost=host&dbPort=invalidPort", "invalid value 'invalidPort' for parameter dbPort"},

		{"PUT", "/api/v1/instances", "dbPort=8563", "missing parameter dbHost"},
		{"PUT", "/api/v1/instances", "dbHost=host", "missing parameter dbPort"},
		{"PUT", "/api/v1/instances", "dbHost=host&dbPort=invalidPort", "invalid value 'invalidPort' for parameter dbPort"},

		{"GET", "/api/v1/extension/extId/extVersion/instances", "dbPort=8563", "missing parameter dbHost"},
		{"GET", "/api/v1/extension/extId/extVersion/instances", "dbHost=host", "missing parameter dbPort"},
		{"GET", "/api/v1/extension/extId/extVersion/instances", "dbHost=host&dbPort=invalidPort", "invalid value 'invalidPort' for parameter dbPort"},

		{"DELETE", "/api/v1/extension/extId/instance/inst-id", "dbPort=8563", "missing parameter dbHost"},
		{"DELETE", "/api/v1/extension/extId/instance/inst-id", "dbHost=host", "missing parameter dbPort"},
		{"DELETE", "/api/v1/extension/extId/instance/inst-id", "dbHost=host&dbPort=invalidPort", "invalid value 'invalidPort' for parameter dbPort"},

		{"DELETE", "/api/v1/extension/extId/version/ver", "dbPort=8563", "missing parameter dbHost"},
		{"DELETE", "/api/v1/extension/extId/version/ver", "dbHost=host", "missing parameter dbPort"},
		{"DELETE", "/api/v1/extension/extId/version/ver", "dbHost=host&dbPort=invalidPort", "invalid value 'invalidPort' for parameter dbPort"},
	}
	suite.controller.On("GetAllExtensions", mock.Anything, mock.Anything).Return([]*extensionController.Extension{{Name: "my-extension", Description: "a cool extension", InstallableVersions: []string{"0.1.0"}}}, nil)
	suite.controller.On("GetInstalledExtensions", mock.Anything, mock.Anything).Return([]*extensionAPI.JsExtInstallation{{Name: "test", Version: "0.1.0", InstanceParameters: []interface{}{map[string]interface{}{"id": "param1", "name": "My param", "type": "string"}}}}, nil)
	suite.controller.On("InstallExtension", mock.Anything, mock.Anything, "ext-id", "ver").Return(nil)
	suite.controller.On("CreateInstance", mock.Anything, mock.Anything, "ext-id", "ver", mock.Anything).Return(&extensionAPI.JsExtInstance{Id: "instId", Name: "instName"}, nil)
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
	return suite.restApi.makeRequestWithAuthHeader(method, path, authHeader, body, expectedStatus)
}
