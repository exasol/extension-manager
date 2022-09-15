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

const (
	BASE_URL                  = "/api/v1/extensionmanager"
	LIST_AVAILABLE_EXTENSIONS = BASE_URL + "/extensions"
	LIST_INSTALLED_EXTENSIONS = BASE_URL + "/installations"
	INSTALL_EXT_URL           = BASE_URL + "/extensions/ext-id/ext-version/install"
	UNINSTALL_EXT_URL         = BASE_URL + "/installations/ext-id/ext-version"
	DELETE_INSTANCE_URL       = BASE_URL + "/installations/ext-id/ext-version/instances/inst-id"
	LIST_INSTANCES_URL        = BASE_URL + "/installations/ext-id/ext-version/instances"
	CREATE_INSTANCE_URL       = BASE_URL + "/installations/ext-id/ext-version/instances"
)

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

// GetInstalledExtensions

func (suite *RestAPISuite) TestGetInstallationsSuccessfully() {
	suite.controller.On("GetInstalledExtensions", mock.Anything, mock.Anything).Return([]*extensionAPI.JsExtInstallation{{Name: "test", Version: "0.1.0", InstanceParameters: []interface{}{map[string]interface{}{"id": "param1", "name": "My param", "type": "string"}}}}, nil)
	for _, test := range authSuccessTests {
		suite.Run(test.authHeader, func() {
			responseString := suite.restApi.makeRequestWithAuthHeader("GET", LIST_INSTALLED_EXTENSIONS+"?dbHost=host&dbPort=8563&", test.authHeader, "", 200)
			suite.assertJSON.Assertf(responseString, `{"installations":[{"name":"test","version":"0.1.0","instanceParameters":[{"id":"param1","name":"My param","type":"string"}]}]}`)
		})
	}
}

func (suite *RestAPISuite) TestGetInstallationsFailed() {
	suite.controller.On("GetInstalledExtensions", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("mock error"))
	responseString := suite.makeRequest("GET", LIST_INSTALLED_EXTENSIONS+"?dbHost=host&dbPort=8563", "", 500)
	suite.Regexp(`{"code":500,"message":"Internal server error",`, responseString)
}

// GetAllExtensions

func (suite *RestAPISuite) TestGetAllExtensionsSuccessfully() {
	suite.controller.On("GetAllExtensions", mock.Anything, mock.Anything).Return([]*extensionController.Extension{{Id: "ext-id", Name: "my-extension", Description: "a cool extension", InstallableVersions: []string{"0.1.0"}}}, nil)
	for _, test := range authSuccessTests {
		suite.Run(test.authHeader, func() {
			responseString := suite.restApi.makeRequestWithAuthHeader("GET", LIST_AVAILABLE_EXTENSIONS+"?dbHost=host&dbPort=8563&", test.authHeader, "", 200)
			suite.assertJSON.Assertf(responseString, `{"extensions":[{"id": "ext-id", "name":"my-extension","description":"a cool extension","installableVersions":["0.1.0"]}]}`)
		})
	}
}

func (suite *RestAPISuite) TestGetAllExtensionsFails() {
	suite.controller.On("GetAllExtensions", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("mock error"))
	responseString := suite.makeRequest("GET", LIST_AVAILABLE_EXTENSIONS+"?dbHost=host&dbPort=8563", "", 500)
	suite.Regexp(`{"code":500,"message":"Internal server error",.*`, responseString)
}

// Install extension

func (suite *RestAPISuite) TestInstallExtensionsSuccessfully() {
	suite.controller.On("InstallExtension", mock.Anything, mock.Anything, "ext-id", "ext-version").Return(nil)
	for _, test := range authSuccessTests {
		suite.Run(test.authHeader, func() {
			responseString := suite.restApi.makeRequestWithAuthHeader("PUT", INSTALL_EXT_URL+"?dbHost=host&dbPort=8563&", test.authHeader, `{}`, 204)
			suite.Equal("", responseString)
		})
	}
}

func (suite *RestAPISuite) TestInstallExtensionsFailed() {
	suite.controller.On("InstallExtension", mock.Anything, mock.Anything, "ext-id", "ext-version").Return(fmt.Errorf("mock error"))
	responseString := suite.makeRequest("PUT", INSTALL_EXT_URL+"?dbHost=host&dbPort=8563", `{}`, 500)
	suite.Regexp(`{"code":500,"message":"Internal server error",.*`, responseString)
}

// Uninstall extension

func (suite *RestAPISuite) TestUninstallExtensionsSuccessfully() {
	suite.controller.On("UninstallExtension", mock.Anything, mock.Anything, "ext-id", "ext-version").Return(nil)
	for _, test := range authSuccessTests {
		suite.Run(test.authHeader, func() {
			responseString := suite.restApi.makeRequestWithAuthHeader("DELETE", UNINSTALL_EXT_URL+"?dbHost=host&dbPort=8563&", test.authHeader, "", 204)
			suite.Equal("", responseString)
		})
	}
}

func (suite *RestAPISuite) TestUninstallExtensionsFailed() {
	suite.controller.On("UninstallExtension", mock.Anything, mock.Anything, "ext-id", "ext-version").Return(fmt.Errorf("mock error"))
	responseString := suite.makeRequest("DELETE", UNINSTALL_EXT_URL+"?extensionId=ext-id&extensionVersion=ver&dbHost=host&dbPort=8563", "", 500)
	suite.Regexp(`{"code":500,"message":"Internal server error",.*`, responseString)
}

// Create instance

func (suite *RestAPISuite) TestCreateInstanceSuccessfully() {
	suite.controller.On("CreateInstance", mock.Anything, mock.Anything, "ext-id", "ext-version", []extensionController.ParameterValue{{Name: "p1", Value: "v1"}}).
		Return(&extensionAPI.JsExtInstance{Id: "instId", Name: "instName"}, nil)
	for _, test := range authSuccessTests {
		suite.Run(test.authHeader, func() {
			responseString := suite.restApi.makeRequestWithAuthHeader("POST", CREATE_INSTANCE_URL+"?dbHost=host&dbPort=8563&", test.authHeader,
				`{"parameterValues": [{"name":"p1", "value":"v1"}]}`, 200)
			suite.Equal(`{"instanceId":"instId","instanceName":"instName"}`+"\n", responseString)
		})
	}
}

func (suite *RestAPISuite) TestCreateInstanceFailed_invalidPayload() {
	suite.controller.On("CreateInstance", mock.Anything, mock.Anything, "ext-id", "ext-version", []extensionController.ParameterValue{{Name: "p1", Value: "v1"}}).Return(&extensionAPI.JsExtInstance{Id: "instId", Name: "instName"}, nil)
	responseString := suite.makeRequest("POST", CREATE_INSTANCE_URL+"?dbHost=host&dbPort=8563",
		`invalid payload`, 400)
	suite.Regexp("{\"code\":400,\"message\":\"Request body contains badly-formed JSON \\(at position 1\\)\".*", responseString)
}

func (suite *RestAPISuite) TestCreateInstanceFailed() {
	suite.controller.On("CreateInstance", mock.Anything, mock.Anything, "ext-id", "ext-version", []extensionController.ParameterValue{{Name: "p1", Value: "v1"}}).Return(nil, fmt.Errorf("mock error"))
	responseString := suite.makeRequest("POST", CREATE_INSTANCE_URL+"?dbHost=host&dbPort=8563",
		`{"parameterValues": [{"name":"p1", "value":"v1"}]}`, 500)
	suite.Regexp(`{"code":500,"message":"Internal server error",.*`, responseString)
}

// List instances

func (suite *RestAPISuite) TestListInstancesSuccessfully() {
	suite.controller.On("FindInstances", mock.Anything, mock.Anything, "ext-id", "ext-version").Return([]*extensionAPI.JsExtInstance{{Id: "instId", Name: "instName"}}, nil)
	for _, test := range authSuccessTests {
		suite.Run(test.authHeader, func() {
			responseString := suite.restApi.makeRequestWithAuthHeader("GET", LIST_INSTANCES_URL+"?dbHost=host&dbPort=8563&", test.authHeader, "", 200)
			suite.Equal(`{"instances":[{"id":"instId","name":"instName"}]}`+"\n", responseString)
		})
	}
}

func (suite *RestAPISuite) TestListInstancesFailed_genericError() {
	suite.controller.On("FindInstances", mock.Anything, mock.Anything, "ext-id", "ext-version").Return(nil, fmt.Errorf("mock"))
	responseString := suite.restApi.makeRequestWithAuthHeader("GET", LIST_INSTANCES_URL+"?dbHost=host&dbPort=8563&", "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==", "", 500)
	suite.Contains(responseString, "{\"code\":500,\"message\":\"Internal server error\"")
}

func (suite *RestAPISuite) TestListInstancesFailed_apiError() {
	suite.controller.On("FindInstances", mock.Anything, mock.Anything, "ext-id", "ext-version").Return(nil, apiErrors.NewAPIError(432, "mock"))
	responseString := suite.restApi.makeRequestWithAuthHeader("GET", LIST_INSTANCES_URL+"?dbHost=host&dbPort=8563&", "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==", "", 432)
	suite.Contains(responseString, "{\"code\":432,\"message\":\"mock\",")
}

// Delete instance

func (suite *RestAPISuite) TestDeleteInstanceSuccessfully() {
	suite.controller.On("DeleteInstance", mock.Anything, mock.Anything, "ext-id", "ext-version", "inst-id").Return(nil)
	for _, test := range authSuccessTests {
		suite.Run(test.authHeader, func() {
			responseString := suite.restApi.makeRequestWithAuthHeader("DELETE", DELETE_INSTANCE_URL+"?dbHost=host&dbPort=8563&", test.authHeader, "", 204)
			suite.Equal("", responseString)
		})
	}
}

func (suite *RestAPISuite) TestDeleteInstanceFailed_genericError() {
	suite.controller.On("DeleteInstance", mock.Anything, mock.Anything, "ext-id", "ext-version", "inst-id").Return(fmt.Errorf("mock"))
	responseString := suite.restApi.makeRequestWithAuthHeader("DELETE", DELETE_INSTANCE_URL+"?dbHost=host&dbPort=8563&", "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==", "", 500)
	suite.Contains(responseString, "{\"code\":500,\"message\":\"Internal server error\"")
}

func (suite *RestAPISuite) TestDeleteInstanceFailed_apiError() {
	suite.controller.On("DeleteInstance", mock.Anything, mock.Anything, "ext-id", "ext-version", "inst-id").Return(apiErrors.NewAPIError(432, "mock"))
	responseString := suite.restApi.makeRequestWithAuthHeader("DELETE", DELETE_INSTANCE_URL+"?dbHost=host&dbPort=8563&", "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==", "", 432)
	suite.Contains(responseString, "{\"code\":432,\"message\":\"mock\",")
}

func (suite *RestAPISuite) TestRequestsFailForMissingParameters() {
	var tests = []struct {
		method        string
		url           string
		parameters    string
		expectedError string
	}{
		{"GET", LIST_AVAILABLE_EXTENSIONS, "dbPort=8563", "missing parameter dbHost"},
		{"GET", LIST_AVAILABLE_EXTENSIONS, "dbHost=host", "missing parameter dbPort"},
		{"GET", LIST_AVAILABLE_EXTENSIONS, "dbHost=host&dbPort=invalidPort", "invalid value 'invalidPort' for parameter dbPort"},

		{"GET", LIST_INSTALLED_EXTENSIONS, "dbPort=8563", "missing parameter dbHost"},
		{"GET", LIST_INSTALLED_EXTENSIONS, "dbHost=host", "missing parameter dbPort"},
		{"GET", LIST_INSTALLED_EXTENSIONS, "dbHost=host&dbPort=invalidPort", "invalid value 'invalidPort' for parameter dbPort"},

		{"PUT", INSTALL_EXT_URL, "extensionId=ext-id&extensionVersion=ext-version&dbPort=8563", "missing parameter dbHost"},
		{"PUT", INSTALL_EXT_URL, "extensionId=ext-id&extensionVersion=ext-version&dbHost=host", "missing parameter dbPort"},
		{"PUT", INSTALL_EXT_URL, "extensionId=ext-id&extensionVersion=ext-version&dbHost=host&dbPort=invalidPort", "invalid value 'invalidPort' for parameter dbPort"},

		{"POST", CREATE_INSTANCE_URL, "dbPort=8563", "missing parameter dbHost"},
		{"POST", CREATE_INSTANCE_URL, "dbHost=host", "missing parameter dbPort"},
		{"POST", CREATE_INSTANCE_URL, "dbHost=host&dbPort=invalidPort", "invalid value 'invalidPort' for parameter dbPort"},

		{"GET", LIST_INSTANCES_URL, "dbPort=8563", "missing parameter dbHost"},
		{"GET", LIST_INSTANCES_URL, "dbHost=host", "missing parameter dbPort"},
		{"GET", LIST_INSTANCES_URL, "dbHost=host&dbPort=invalidPort", "invalid value 'invalidPort' for parameter dbPort"},

		{"DELETE", DELETE_INSTANCE_URL, "dbPort=8563", "missing parameter dbHost"},
		{"DELETE", DELETE_INSTANCE_URL, "dbHost=host", "missing parameter dbPort"},
		{"DELETE", DELETE_INSTANCE_URL, "dbHost=host&dbPort=invalidPort", "invalid value 'invalidPort' for parameter dbPort"},

		{"DELETE", UNINSTALL_EXT_URL, "dbPort=8563", "missing parameter dbHost"},
		{"DELETE", UNINSTALL_EXT_URL, "dbHost=host", "missing parameter dbPort"},
		{"DELETE", UNINSTALL_EXT_URL, "dbHost=host&dbPort=invalidPort", "invalid value 'invalidPort' for parameter dbPort"},
	}
	suite.controller.On("GetAllExtensions", mock.Anything, mock.Anything).Return([]*extensionController.Extension{{Name: "my-extension", Description: "a cool extension", InstallableVersions: []string{"0.1.0"}}}, nil)
	suite.controller.On("GetInstalledExtensions", mock.Anything, mock.Anything).Return([]*extensionAPI.JsExtInstallation{{Name: "test", Version: "0.1.0", InstanceParameters: []interface{}{map[string]interface{}{"id": "param1", "name": "My param", "type": "string"}}}}, nil)
	suite.controller.On("InstallExtension", mock.Anything, mock.Anything, "ext-id", "ext-version").Return(nil)
	suite.controller.On("CreateInstance", mock.Anything, mock.Anything, "ext-id", "ext-version", mock.Anything).Return(&extensionAPI.JsExtInstance{Id: "instId", Name: "instName"}, nil)
	for _, test := range tests {
		suite.Run(fmt.Sprintf("Request %s %s?%s results in error message %q", test.method, test.url, test.parameters, test.expectedError), func() {
			completePath := fmt.Sprintf("%s?%s", test.url, test.parameters)
			responseString := suite.makeRequest(test.method, completePath, "", 400)
			suite.Regexp(fmt.Sprintf(`{"code":400,"message":"%s"`, test.expectedError), responseString)
		})
	}
}

func (suite *RestAPISuite) makeRequest(method, path, body string, expectedStatus int) string {
	suite.T().Helper()
	authHeader := createBasicAuthHeader("user", "password")
	return suite.restApi.makeRequestWithAuthHeader(method, path, authHeader, body, expectedStatus)
}
