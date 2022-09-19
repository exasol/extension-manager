package restAPI

import (
	"fmt"
	"path"
	"testing"

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
	suite.Suite
	restApi           *baseRestAPITest
	exasol            *integrationTesting.DbTestSetup
	tempExtensionRepo string
	assertJSON        *jsonassert.Asserter
}

func TestRestAPIIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(RestAPIIntegrationTestSuite))
}

func (suite *RestAPIIntegrationTestSuite) SetupSuite() {
	suite.exasol = integrationTesting.StartDbSetup(&suite.Suite)
	suite.assertJSON = jsonassert.New(suite.T())
}

func (suite *RestAPIIntegrationTestSuite) TearDownSuite() {
	suite.exasol.StopDb()
}

func (suite *RestAPIIntegrationTestSuite) SetupTest() {
	ctrl := extensionController.Create(suite.tempExtensionRepo, EXTENSION_SCHEMA)
	suite.restApi = startRestApi(&suite.Suite, ctrl)
}

func (suite *RestAPIIntegrationTestSuite) TearDownTest() {
	suite.restApi.restAPI.Stop()
}

func (suite *RestAPIIntegrationTestSuite) listAvailableExtensions() string {
	return LIST_AVAILABLE_EXTENSIONS + "?" + suite.getValidDbArgs()
}

func (suite *RestAPIIntegrationTestSuite) listInstalledExtensions() string {
	return LIST_INSTALLED_EXTENSIONS + "?" + suite.getValidDbArgs()
}

func (suite *RestAPIIntegrationTestSuite) getExtensionDetails(extensionId, extensionVersion string) string {
	return fmt.Sprintf("%s/extensions/%s/%s?%s", BASE_URL, extensionId, extensionVersion, suite.getValidDbArgs())
}

func (suite *RestAPIIntegrationTestSuite) listInstances(extensionId, extensionVersion string) string {
	return fmt.Sprintf("%s/installations/%s/%s/instances?%s", BASE_URL, extensionId, extensionVersion, suite.getValidDbArgs())
}

func (suite *RestAPIIntegrationTestSuite) deleteInstance(extensionId, extensionVersion, instanceId string) string {
	return fmt.Sprintf("%s/installations/%s/%s/instances/%s?%s", BASE_URL, extensionId, extensionVersion, instanceId, suite.getValidDbArgs())
}

func (suite *RestAPIIntegrationTestSuite) uninstallExtension(extensionId, extensionVersion string) string {
	return fmt.Sprintf("%s/installations/%s/%s?%s", BASE_URL, extensionId, extensionVersion, suite.getValidDbArgs())
}

func (suite *RestAPIIntegrationTestSuite) TestGetAllExtensionsSuccessfully() {
	response := suite.makeGetRequest(suite.listAvailableExtensions())
	suite.assertJSON.Assertf(response, `{"extensions":[]}`)
}

// List installed extensions

func (suite *RestAPIIntegrationTestSuite) TestGetInstallationsSuccessfully() {
	response := suite.makeGetRequest(suite.listInstalledExtensions())
	suite.assertJSON.Assertf(response, `{"installations":[]}`)
}

func (suite *RestAPIIntegrationTestSuite) TestGetInstallationsFails_InvalidUsernamePassword() {
	response := suite.restApi.makeRequestWithAuthHeader("GET", suite.listInstalledExtensions(), createBasicAuthHeader("wrong", "user"), "", 401)
	suite.Regexp(`{"code":401,"message":"invalid database credentials".*`, response)
}

func (suite *RestAPIIntegrationTestSuite) TestGetInstallationsFails_InvalidBearerToken() {
	response := suite.restApi.makeRequestWithAuthHeader("GET", suite.listInstalledExtensions(), createBearerAuthHeader("invalid-token"), "", 401)
	suite.Regexp(`{"code":401,"message":"invalid database credentials".*`, response)
}

// Get extension details

func (suite *RestAPIIntegrationTestSuite) TestGetExtensionDetailsFailsForUnknownExtension() {
	response := suite.makeRequest("GET", suite.getExtensionDetails("unknown-extension", "version"), "", 500)
	suite.Regexp(`{"code":500,"message":"Internal server error".*`, response)
}

func (suite *RestAPIIntegrationTestSuite) TestGetExtensionDetailsSucceeds() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithGetInstanceParameterDefinitionFunc(`return [{id: "param1", name: "My param:"+version, type: "string"}]`).
		Build().WriteToFile(path.Join(suite.tempExtensionRepo, "ext-id"))
	response := suite.makeRequest("GET", suite.getExtensionDetails("ext-id", "ext-version"), "", 200)
	suite.assertJSON.Assertf(response, `{"id": "ext-id", "version":"ext-version", "parameterDefinitions": [
		{"id":"param1","name":"My param:ext-version","definition":{"id": "param1", "name": "My param:ext-version", "type": "string"}}
	]}`)
}

// List instances

func (suite *RestAPIIntegrationTestSuite) TestListInstancesSuccessfully() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithFindInstancesFunc("context.sqlClient.execute('select 1'); return [{id: 'instId', name: 'instName_ver_'+version}]").
		Build().WriteToFile(path.Join(suite.tempExtensionRepo, "ext-id"))
	response := suite.makeGetRequest(suite.listInstances("ext-id", "ext-version"))
	suite.assertJSON.Assertf(response, `{"instances":[{"id":"instId","name":"instName_ver_ext-version"}]}`)
}

func (suite *RestAPIIntegrationTestSuite) TestListInstancesQueryFails() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithFindInstancesFunc("context.sqlClient.execute('invalid query'); return [{id: 'instId', name: 'instName_ver'+version}]").
		Build().WriteToFile(path.Join(suite.tempExtensionRepo, "ext-id"))
	response := suite.makeRequest("GET", suite.listInstances("ext-id", "ext-version"), "", 500)
	suite.Contains(response, `{"code":500,"message":"Internal server error"`)
}

// Delete instance

func (suite *RestAPIIntegrationTestSuite) TestDeleteInstanceSuccessfully() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithDeleteInstanceFunc("context.sqlClient.execute('select 1')").
		Build().WriteToFile(path.Join(suite.tempExtensionRepo, "ext-id"))
	response := suite.makeRequest("DELETE", suite.deleteInstance("ext-id", "ext-version", "inst-id"), "", 204)
	suite.Equal("", response)
}

func (suite *RestAPIIntegrationTestSuite) TestDeleteInstanceFails() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithDeleteInstanceFunc("context.sqlClient.execute('invalid query')").
		Build().WriteToFile(path.Join(suite.tempExtensionRepo, "ext-id"))
	response := suite.makeRequest("DELETE", suite.deleteInstance("ext-id", "ext-version", "inst-id"), "", 500)
	suite.Contains(response, `{"code":500,"message":"Internal server error"`)
}

func (suite *RestAPIIntegrationTestSuite) TestUninstallExtensionSuccessfully() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithUninstallFunc("context.sqlClient.execute('select 1')").
		Build().WriteToFile(path.Join(suite.tempExtensionRepo, "ext-id"))
	response := suite.makeRequest("DELETE", suite.uninstallExtension("ext-id", "ext-version"), "", 204)
	suite.Equal("", response)
}

func (suite *RestAPIIntegrationTestSuite) TestExtensionInstanceFails() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithUninstallFunc("context.sqlClient.execute('invalid query')").
		Build().WriteToFile(path.Join(suite.tempExtensionRepo, "ext-id"))
	response := suite.makeRequest("DELETE", suite.uninstallExtension("ext-id", "ext-version"), "", 500)
	suite.Contains(response, `{"code":500,"message":"Internal server error"`)
}

func (suite *RestAPIIntegrationTestSuite) TestGetOpenApiHtml() {
	response := suite.makeGetRequest("/openapi/index.html")
	suite.Regexp("\n<!DOCTYPE html>.*", response)
}

func (suite *RestAPIIntegrationTestSuite) TestGetOpenApiJson() {
	response := suite.makeGetRequest("/openapi.json")
	suite.Regexp(".*\"openapi\": \"3\\.0\\.0\",.*", response)
}

func (suite *RestAPIIntegrationTestSuite) getValidDbArgs() string {
	return suite.getDbArgsWithUserPassword()
}

func (suite *RestAPIIntegrationTestSuite) getDbArgsWithUserPassword() string {
	info := suite.exasol.ConnectionInfo
	return fmt.Sprintf("dbHost=%s&dbPort=%d", info.Host, info.Port)
}

func (suite *RestAPIIntegrationTestSuite) makeGetRequest(path string) string {
	return suite.makeRequest("GET", path, "", 200)
}

func (suite *RestAPIIntegrationTestSuite) makeRequest(method string, path string, body string, expectedStatusCode int) string {
	suite.T().Helper()
	info := suite.exasol.ConnectionInfo
	return suite.restApi.makeRequestWithAuthHeader(method, path, createBasicAuthHeader(info.User, info.Password), body, expectedStatusCode)
}
