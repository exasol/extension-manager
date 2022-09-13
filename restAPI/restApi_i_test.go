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

func (suite *RestAPIIntegrationTestSuite) TestGetAllExtensionsSuccessfully() {
	response := suite.makeGetRequest("/api/v1/extensions?" + suite.getValidDbArgs())
	suite.assertJSON.Assertf(response, `{"extensions":[]}`)
}

func (suite *RestAPIIntegrationTestSuite) TestGetInstallationsSuccessfully() {
	response := suite.makeGetRequest("/api/v1/installations?" + suite.getValidDbArgs())
	suite.assertJSON.Assertf(response, `{"installations":[]}`)
}

func (suite *RestAPIIntegrationTestSuite) TestGetInstallationsFails_InvalidUsernamePassword() {
	response := suite.restApi.makeRequestWithAuthHeader("GET", "/api/v1/installations?"+suite.getValidDbArgs(), createBasicAuthHeader("wrong", "user"), "", 401)
	suite.Regexp(`{"code":401,"message":"invalid database credentials".*`, response)
}

func (suite *RestAPIIntegrationTestSuite) TestGetInstallationsFails_InvalidBearerToken() {
	response := suite.restApi.makeRequestWithAuthHeader("GET", "/api/v1/installations?"+suite.getValidDbArgs(), "Bearer invalid", "", 401)
	suite.Regexp(`{"code":401,"message":"invalid database credentials".*`, response)
}

func (suite *RestAPIIntegrationTestSuite) TestListInstancesSuccessfully() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithFindInstancesFunc("context.sqlClient.execute('select 1'); return [{id: 'instId', name: 'instName_ver'+version}]").
		Build().WriteToFile(path.Join(suite.tempExtensionRepo, "ext-id"))
	response := suite.makeGetRequest("/api/v1/extension/ext-id/ver/instances?" + suite.getValidDbArgs())
	suite.assertJSON.Assertf(response, `{"instances":[{"id":"instId","name":"instName_verver"}]}`)
}

func (suite *RestAPIIntegrationTestSuite) TestListInstancesQueryFails() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithFindInstancesFunc("context.sqlClient.execute('invalid query'); return [{id: 'instId', name: 'instName_ver'+version}]").
		Build().WriteToFile(path.Join(suite.tempExtensionRepo, "ext-id"))
	response := suite.makeRequest("GET", "/api/v1/extension/ext-id/ver/instances?"+suite.getValidDbArgs(), "", 500)
	suite.Contains(response, `{"code":500,"message":"Internal server error"`)
}

func (suite *RestAPIIntegrationTestSuite) TestDeleteInstanceSuccessfully() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithDeleteInstanceFunc("context.sqlClient.execute('select 1')").
		Build().WriteToFile(path.Join(suite.tempExtensionRepo, "ext-id"))
	response := suite.makeRequest("DELETE", "/api/v1/extension/ext-id/instance/inst-id?"+suite.getValidDbArgs(), "", 204)
	suite.Equal("", response)
}

func (suite *RestAPIIntegrationTestSuite) TestDeleteInstanceFails() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithDeleteInstanceFunc("context.sqlClient.execute('invalid query')").
		Build().WriteToFile(path.Join(suite.tempExtensionRepo, "ext-id"))
	response := suite.makeRequest("DELETE", "/api/v1/extension/ext-id/instance/inst-id?"+suite.getValidDbArgs(), "", 500)
	suite.Contains(response, `{"code":500,"message":"Internal server error"`)
}

func (suite *RestAPIIntegrationTestSuite) TestDeleteExtensionSuccessfully() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithUninstallFunc("context.sqlClient.execute('select 1')").
		Build().WriteToFile(path.Join(suite.tempExtensionRepo, "ext-id"))
	response := suite.makeRequest("DELETE", "/api/v1/extension/ext-id/version/version?"+suite.getValidDbArgs(), "", 204)
	suite.Equal("", response)
}

func (suite *RestAPIIntegrationTestSuite) TestExtensionInstanceFails() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithUninstallFunc("context.sqlClient.execute('invalid query')").
		Build().WriteToFile(path.Join(suite.tempExtensionRepo, "ext-id"))
	response := suite.makeRequest("DELETE", "/api/v1/extension/ext-id/version/version?"+suite.getValidDbArgs(), "", 500)
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
	info := suite.exasol.ConnectionInfo
	return suite.restApi.makeRequestWithAuthHeader(method, path, createBasicAuthHeader(info.User, info.Password), body, expectedStatusCode)
}
