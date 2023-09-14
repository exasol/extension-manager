package restAPI

import (
	"fmt"
	"testing"

	"github.com/exasol/extension-manager/pkg/extensionController"
	"github.com/exasol/extension-manager/pkg/integrationTesting"
	"github.com/kinbiko/jsonassert"

	"github.com/stretchr/testify/suite"
)

const (
	EXTENSION_SCHEMA = "test"
	EXTENSION_ID     = "ext-id"
)

type RestAPIIntegrationTestSuite struct {
	suite.Suite
	restApi        *baseRestAPITest
	exasol         *integrationTesting.DbTestSetup
	registryServer *integrationTesting.MockRegistryServer
	assertJSON     *jsonassert.Asserter
}

/* [itest -> dsn~rest-interface~1] */
func TestRestAPIIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(RestAPIIntegrationTestSuite))
}

func (suite *RestAPIIntegrationTestSuite) SetupSuite() {
	suite.registryServer = integrationTesting.NewMockRegistryServer(&suite.Suite)
	suite.registryServer.Start()
	suite.exasol = integrationTesting.StartDbSetup(&suite.Suite)
	suite.assertJSON = jsonassert.New(suite.T())
}

func (suite *RestAPIIntegrationTestSuite) TearDownSuite() {
	suite.registryServer.Close()
	suite.exasol.StopDb()
}

func (suite *RestAPIIntegrationTestSuite) SetupTest() {
	// [itest -> dsn~extension-registry~1]
	ctrl := extensionController.Create(suite.registryServer.IndexUrl(), EXTENSION_SCHEMA)
	suite.restApi = startRestApi(&suite.Suite, false, ctrl)
	suite.registryServer.Reset()
	suite.registryServer.SetRegistryContent(`{"extensions":[]}`)
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

func (suite *RestAPIIntegrationTestSuite) upgradeExtension(extensionId string) string {
	return fmt.Sprintf("%s/installations/%s/upgrade?%s", BASE_URL, extensionId, suite.getValidDbArgs())
}

func (suite *RestAPIIntegrationTestSuite) TestGetAllExtensionsSuccessfully() {
	response := suite.makeGetRequest(suite.listAvailableExtensions())
	suite.assertJSON.Assertf(response, `{"extensions":[]}`)
}

// List installed extensions

/* [itest -> dsn~list-extensions~1] */
func (suite *RestAPIIntegrationTestSuite) TestGetInstallationsSuccessfully() {
	response := suite.makeGetRequest(suite.listInstalledExtensions())
	suite.assertJSON.Assertf(response, `{"installations":[]}`)
}

func (suite *RestAPIIntegrationTestSuite) TestGetInstallationsFailsInvalidUsernamePassword() {
	response := suite.restApi.makeRequestWithAuthHeader("GET", suite.listInstalledExtensions(), createBasicAuthHeader("wrong", "user"), "", 401)
	suite.Contains(response, `{"code":401,"message":"invalid database credentials"`)
}

func (suite *RestAPIIntegrationTestSuite) TestGetInstallationsFailsInvalidBearerToken() {
	response := suite.restApi.makeRequestWithAuthHeader("GET", suite.listInstalledExtensions(), createBearerAuthHeader("invalid-token"), "", 401)
	suite.Contains(response, `{"code":401,"message":"invalid database credentials"`)
}

// Get extension details

func (suite *RestAPIIntegrationTestSuite) TestGetExtensionDetailsFailsForUnknownExtension() {
	response := suite.makeRequest("GET", suite.getExtensionDetails("unknown-ext-id", "version"), 404)
	suite.Contains(response, `{"code":404,"message":"extension \"unknown-ext-id\" not found",`)
}

func (suite *RestAPIIntegrationTestSuite) TestGetExtensionDetailsSucceeds() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithGetInstanceParameterDefinitionFunc(`return [{id: "param1", name: "My param:"+version, type: "string"}]`).
		Build().Publish(suite.registryServer, EXTENSION_ID)
	response := suite.makeRequest("GET", suite.getExtensionDetails(EXTENSION_ID, "ext-version"), 200)
	suite.assertJSON.Assertf(response, `{"id": "ext-id", "version":"ext-version", "parameterDefinitions": [
		{"id":"param1","name":"My param:ext-version","definition":{"id": "param1", "name": "My param:ext-version", "type": "string"}}
	]}`)
}

// List instances

func (suite *RestAPIIntegrationTestSuite) TestListInstancesSuccessfully() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithFindInstancesFunc("context.sqlClient.execute('select 1'); return [{id: 'instId', name: 'instName_ver_'+version}]").
		Build().Publish(suite.registryServer, EXTENSION_ID)
	response := suite.makeGetRequest(suite.listInstances(EXTENSION_ID, "ext-version"))
	suite.assertJSON.Assertf(response, `{"instances":[{"id":"instId","name":"instName_ver_ext-version"}]}`)
}

func (suite *RestAPIIntegrationTestSuite) TestListInstancesQueryFails() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithFindInstancesFunc("context.sqlClient.execute('invalid query'); return [{id: 'instId', name: 'instName_ver'+version}]").
		Build().Publish(suite.registryServer, EXTENSION_ID)
	response := suite.makeRequest("GET", suite.listInstances(EXTENSION_ID, "ext-version"), 500)
	suite.Contains(response, `{"code":500,"message":"Internal server error"`)
}

func (suite *RestAPIIntegrationTestSuite) TestListInstancesQueryFailsForUnknownExtension() {
	response := suite.makeRequest("GET", suite.listInstances("unknown-ext-id", "ext-version"), 404)
	suite.Contains(response, `{"code":404,"message":"extension \"unknown-ext-id\" not found"`)
}

// Delete instance

func (suite *RestAPIIntegrationTestSuite) TestDeleteInstanceSuccessfully() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithDeleteInstanceFunc("context.sqlClient.execute('select 1')").
		Build().Publish(suite.registryServer, EXTENSION_ID)
	response := suite.makeRequest("DELETE", suite.deleteInstance(EXTENSION_ID, "ext-version", "inst-id"), 204)
	suite.Equal("", response)
}

func (suite *RestAPIIntegrationTestSuite) TestDeleteInstanceFails() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithDeleteInstanceFunc("context.sqlClient.execute('invalid query')").
		Build().Publish(suite.registryServer, EXTENSION_ID)
	response := suite.makeRequest("DELETE", suite.deleteInstance(EXTENSION_ID, "ext-version", "inst-id"), 500)
	suite.Contains(response, `{"code":500,"message":"Internal server error"`)
}

func (suite *RestAPIIntegrationTestSuite) TestDeleteInstanceFailsForUnknownExtension() {
	response := suite.makeRequest("DELETE", suite.deleteInstance("unknown-ext-id", "ext-version", "inst-id"), 404)
	suite.Contains(response, `{"code":404,"message":"extension \"unknown-ext-id\" not found"`)
}

// Uninstall extension

func (suite *RestAPIIntegrationTestSuite) TestUninstallExtensionSuccessfully() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithUninstallFunc("context.sqlClient.execute('select 1')").
		Build().Publish(suite.registryServer, EXTENSION_ID)
	response := suite.makeRequest("DELETE", suite.uninstallExtension(EXTENSION_ID, "ext-version"), 204)
	suite.Equal("", response)
}

func (suite *RestAPIIntegrationTestSuite) TestUninstallExtensionFails() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithUninstallFunc("context.sqlClient.execute('invalid query')").
		Build().Publish(suite.registryServer, EXTENSION_ID)
	response := suite.makeRequest("DELETE", suite.uninstallExtension(EXTENSION_ID, "ext-version"), 500)
	suite.Contains(response, `{"code":500,"message":"Internal server error"`)
}

func (suite *RestAPIIntegrationTestSuite) TestUninstallExtensionFailsForUnknownExtension() {
	response := suite.makeRequest("DELETE", suite.uninstallExtension("unknown-ext-id", "ext-version"), 404)
	suite.Contains(response, `{"code":404,"message":"extension \"unknown-ext-id\" not found"`)
}

// Upgrade

/* [itest -> dsn~upgrade-extension~1] */
func (suite *RestAPIIntegrationTestSuite) TestUpgradeExtensionSuccessfully() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithUpgradeFunc("context.sqlClient.execute('select 1'); return { previousVersion: 'old', newVersion: 'new' };").
		Build().Publish(suite.registryServer, EXTENSION_ID)
	response := suite.makeRequest("POST", suite.upgradeExtension(EXTENSION_ID), 200)
	suite.assertJSON.Assertf(response, `{"previousVersion":"old","newVersion":"new"}`)
}

func (suite *RestAPIIntegrationTestSuite) TestUpgradeFails() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithUpgradeFunc("context.sqlClient.execute('invalid query'); return { previousVersion: 'old', newVersion: 'new' };").
		Build().Publish(suite.registryServer, EXTENSION_ID)
	response := suite.makeRequest("POST", suite.upgradeExtension(EXTENSION_ID), 500)
	suite.Contains(response, `{"code":500,"message":"Internal server error"`)
}

func (suite *RestAPIIntegrationTestSuite) TestUpgradeFailsForUnknownExtension() {
	response := suite.makeRequest("POST", suite.upgradeExtension("unknown-ext-id"), 404)
	suite.Contains(response, `{"code":404,"message":"extension \"unknown-ext-id\" not found"`)
}

/* [itest -> dsn~openapi-spec~1] */
func (suite *RestAPIIntegrationTestSuite) TestGetOpenApiHtml() {
	response := suite.makeGetRequest("/openapi/index.html")
	suite.Contains(response, "\n<!DOCTYPE html>")
}

/* [itest -> dsn~openapi-spec~1] */
func (suite *RestAPIIntegrationTestSuite) TestGetOpenApiJson() {
	response := suite.makeGetRequest("/openapi.json")
	suite.Contains(response, `"openapi": "3.0.0",`)
}

func (suite *RestAPIIntegrationTestSuite) getValidDbArgs() string {
	return suite.getDbArgsWithUserPassword()
}

func (suite *RestAPIIntegrationTestSuite) getDbArgsWithUserPassword() string {
	info := suite.exasol.ConnectionInfo
	return fmt.Sprintf("dbHost=%s&dbPort=%d", info.Host, info.Port)
}

func (suite *RestAPIIntegrationTestSuite) makeGetRequest(path string) string {
	return suite.makeRequest("GET", path, 200)
}

func (suite *RestAPIIntegrationTestSuite) makeRequest(method string, path string, expectedStatusCode int) string {
	suite.T().Helper()
	info := suite.exasol.ConnectionInfo
	body := ""
	return suite.restApi.makeRequestWithAuthHeader(method, path, createBasicAuthHeader(info.User, info.Password), body, expectedStatusCode)
}
