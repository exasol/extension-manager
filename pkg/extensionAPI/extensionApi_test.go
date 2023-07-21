package extensionAPI

import (
	"fmt"
	"strings"
	"testing"

	"github.com/exasol/extension-manager/pkg/backend"
	"github.com/exasol/extension-manager/pkg/extensionAPI/context"
	"github.com/exasol/extension-manager/pkg/extensionAPI/exaMetadata"
	"github.com/exasol/extension-manager/pkg/extensionController/bfs"
	"github.com/exasol/extension-manager/pkg/integrationTesting"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ExtensionApiSuite struct {
	suite.Suite
	mockSQLClient      sqlClientMock
	mockBucketFsClient bfs.BucketFsMock
}

func TestExtensionApiSuite(t *testing.T) {
	suite.Run(t, new(ExtensionApiSuite))
}

func (suite *ExtensionApiSuite) SetupSuite() {
	// Nothing to do
}

func (suite *ExtensionApiSuite) SetupTest() {
	suite.mockSQLClient = sqlClientMock{}
	suite.mockBucketFsClient = bfs.BucketFsMock{}
}

func (suite *ExtensionApiSuite) TearDownTest() {
	suite.mockSQLClient.AssertExpectations(suite.T())
	suite.mockBucketFsClient.AssertExpectations(suite.T())
}

/* [utest -> dsn~extension-definition~1] */
func (suite *ExtensionApiSuite) TestLoadExtension() {
	extensionContent := integrationTesting.CreateTestExtensionBuilder(suite.T()).Build().AsString()
	extension := suite.loadExtension(extensionContent)
	suite.Equal("MyDemoExtension", extension.Name)
	suite.Equal("Demo category", extension.Category)
}

type sqlClientMock struct {
	mock.Mock
}

func (mock *sqlClientMock) Execute(query string, args ...any) {
	mock.Called(query, args)
}

func (mock *sqlClientMock) Query(query string, args ...any) backend.QueryResult {
	mockArgs := mock.Called(query, args)
	return mockArgs.Get(0).(backend.QueryResult)
}

func (suite *ExtensionApiSuite) TestGetParameterDefinitionsEmptyResult() {
	extensionContent := integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithGetInstanceParameterDefinitionFunc(`return []`).
		Build().AsString()
	extension := suite.loadExtension(extensionContent)
	definitions := extension.extension.GetParameterDefinitions(suite.mockContext(), "extVersion")
	suite.Equal([]interface{}{}, definitions)
}

/* [utest -> dsn~configuration-parameters~1] */
func (suite *ExtensionApiSuite) TestGetParameterDefinitions() {
	extensionContent := integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithGetInstanceParameterDefinitionFunc(`return [{id: "param1", name: "My param", type: "string"}]`).
		Build().AsString()
	extension := suite.loadExtension(extensionContent)
	definitions := extension.extension.GetParameterDefinitions(suite.mockContext(), "extVersion")
	suite.Equal([]interface{}{map[string]interface{}{"id": "param1", "name": "My param", "type": "string"}}, definitions)
}

func (suite *ExtensionApiSuite) TestInstall() {
	extensionContent := integrationTesting.CreateTestExtensionBuilder(suite.T()).Build().AsString()
	extension := suite.loadExtension(extensionContent)
	suite.mockSQLClient.On("Execute", "select 1", []any{}).Return()
	err := extension.Install(suite.mockContext(), "extVersion")
	suite.NoError(err)
}

func (suite *ExtensionApiSuite) TestInstallResolveBucketFsPath() {
	version := "extVersion"
	absolutePath := "/absolute/file/path"
	suite.mockBucketFsClient.SimulateAbsolutePath("my-adapter-"+version+".jar", absolutePath)
	extensionContent := integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithInstallFunc("context.sqlClient.execute(`create script path ${context.bucketFs.resolvePath('my-adapter-'+version+'.jar')}`)").
		Build().AsString()
	extension := suite.loadExtension(extensionContent)
	suite.mockSQLClient.On("Execute", "create script path "+absolutePath, []any{}).Return()
	err := extension.Install(suite.mockContext(), version)
	suite.NoError(err)
}

func (suite *ExtensionApiSuite) TestInstallResolveBucketFsPathFails() {
	version := "extVersion"
	suite.mockBucketFsClient.SimulateAbsolutePathError("my-adapter-"+version+".jar", fmt.Errorf("mock error"))
	extensionContent := integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithInstallFunc("context.sqlClient.execute(`create script path ${context.bucketFs.resolvePath('my-adapter-'+version+'.jar')}`)").
		Build().AsString()
	extension := suite.loadExtension(extensionContent)
	err := extension.Install(suite.mockContext(), version)
	suite.EqualError(err, `failed to install extension "ext-id": failed to find absolute path for file "my-adapter-extVersion.jar": mock error`)
}

func (suite *ExtensionApiSuite) TestJavaScriptConsoleLog() {
	extensionContent := integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithInstallFunc("console.log('test log message')").
		Build().AsString()
	extension := suite.loadExtension(extensionContent)
	err := extension.Install(suite.mockContext(), "extensionVersion")
	suite.NoError(err)
}

func (suite *ExtensionApiSuite) TestJavaScriptConsoleWarn() {
	extensionContent := integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithInstallFunc("console.warn('test warning message')").
		Build().AsString()
	extension := suite.loadExtension(extensionContent)
	err := extension.Install(suite.mockContext(), "extensionVersion")
	suite.NoError(err)
}

func (suite *ExtensionApiSuite) TestJavaScriptConsoleError() {
	extensionContent := integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithInstallFunc("console.error('test error message')").
		Build().AsString()
	extension := suite.loadExtension(extensionContent)
	err := extension.Install(suite.mockContext(), "extensionVersion")
	suite.NoError(err)
}

func (suite *ExtensionApiSuite) TestUninstall() {
	extensionContent := integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithUninstallFunc("context.sqlClient.execute(`uninstall version ${version}`)").
		Build().AsString()
	extension := suite.loadExtension(extensionContent)
	suite.mockSQLClient.On("Execute", "uninstall version extVersion", []any{}).Return()
	err := extension.Uninstall(suite.mockContext(), "extVersion")
	suite.NoError(err)
}

func (suite *ExtensionApiSuite) TestAddInstanceValidParameters() {
	extensionContent := integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithAddInstanceFunc("context.sqlClient.execute('create vs');\n" +
			"return {id: 'instId', name: `instance_${version}_${params.values[0].name}_${params.values[0].value}`};").
		Build().AsString()
	extension := suite.loadExtension(extensionContent)
	suite.mockSQLClient.On("Execute", "create vs", []any{}).Return()
	instance, err := extension.AddInstance(suite.mockContext(), "extensionVersion", &ParameterValues{Values: []ParameterValue{{Name: "p1", Value: "v1"}}})
	suite.NoError(err)
	suite.Equal(&JsExtInstance{Id: "instId", Name: "instance_extensionVersion_p1_v1"}, instance)
}

func (suite *ExtensionApiSuite) TestListInstancesEmptyResult() {
	extensionContent := integrationTesting.CreateTestExtensionBuilder(suite.T()).
		Build().AsString()
	extension := suite.loadExtension(extensionContent)
	instances, err := extension.ListInstances(suite.mockContext(), "ver")
	suite.NoError(err)
	suite.Empty(instances)
}

func (suite *ExtensionApiSuite) TestListInstancesNonEmptyResult() {
	extensionContent := integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithFindInstancesFunc(`return [{id: "instId", name: "instName"}]`).
		Build().AsString()
	extension := suite.loadExtension(extensionContent)
	instances, err := extension.ListInstances(suite.mockContext(), "ver")
	suite.NoError(err)
	suite.Equal([]*JsExtInstance{{Id: "instId", Name: "instName"}}, instances)
}

func (suite *ExtensionApiSuite) TestDeleteInstance() {
	extensionContent := integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithDeleteInstanceFunc("context.sqlClient.execute(`drop instance ${instanceId}`)").
		Build().AsString()
	extension := suite.loadExtension(extensionContent)
	suite.mockSQLClient.On("Execute", "drop instance instId", []any{}).Return()
	err := extension.DeleteInstance(suite.mockContext(), "extVersion", "instId")
	suite.NoError(err)
}

func createMockMetadata() *exaMetadata.ExaMetadata {
	exaAllScripts := exaMetadata.ExaScriptTable{Rows: []exaMetadata.ExaScriptRow{{Name: "test"}}}
	return &exaMetadata.ExaMetadata{AllScripts: exaAllScripts}
}

func (suite *ExtensionApiSuite) TestFindInstallationsCanReadAllScriptsTable() {
	extensionContent := integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithFindInstallationsFunc(`
		return metadata.allScripts.rows.map(row => {
			return {name: row.name, version: "0.1.0"}
		});`).Build().AsString()
	extension := suite.loadExtension(extensionContent)
	result, err := extension.FindInstallations(createMockContext(), createMockMetadata())
	suite.Equal([]*JsExtInstallation{{Name: "test", Version: "0.1.0"}}, result)
	suite.NoError(err)
}

func (suite *ExtensionApiSuite) TestFindInstallationsReturningParameters() {
	extensionContent := integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithFindInstallationsFunc(integrationTesting.
			MockFindInstallationsFunction("test", "0.1.0")).Build().AsString()
	extension := suite.loadExtension(extensionContent)
	result, err := extension.FindInstallations(createMockContext(), createMockMetadata())
	suite.Equal([]*JsExtInstallation{{Name: "test", Version: "0.1.0"}}, result)
	suite.NoError(err)
}

/* [itest -> dsn~extension-compatibility~1] */
func (suite *ExtensionApiSuite) TestLoadExtensionWithCompatibleApiVersion() {
	extensionContent := minimalExtension("0.1.15")
	extension, err := LoadExtension("ext-id", extensionContent)
	suite.NoError(err)
	suite.NotNil(extension)
}

/* [itest -> dsn~extension-compatibility~1] */
func (suite *ExtensionApiSuite) TestLoadExtensionWithIncompatibleApiVersion() {
	extensionContent := minimalExtension("99.0.0")
	extension, err := LoadExtension("ext-id", extensionContent)
	suite.ErrorContains(err, `extension "ext-id" uses incompatible API version "99.0.0". Please update the extension to use supported version "`+supportedApiVersion+`"`)
	suite.Nil(extension)
}

func (suite *ExtensionApiSuite) TestLoadExtensionWithInvalidApiVersion() {
	extensionContent := minimalExtension("invalid")
	extension, err := LoadExtension("ext-id", extensionContent)
	suite.ErrorContains(err, `extension "ext-id" uses invalid API version number "invalid"`)
	suite.Nil(extension)
}

func (suite *ExtensionApiSuite) TestUsingExtensionWithMissingPropertiesWorks() {
	extensionContent := minimalExtension("0.1.15")
	extension, err := LoadExtension("ext-id", extensionContent)
	suite.NoError(err)
	suite.Equal("ext-id", extension.Id)
	suite.Equal("", extension.Category)
	suite.Equal("", extension.Description)
	suite.Empty(extension.BucketFsUploads)
	suite.Empty(extension.InstallableVersions)
}

func (suite *ExtensionApiSuite) TestUsingExtensionWithMissingFunctionFailsGracefully() {
	extensionContent := minimalExtension("0.1.15")
	extension, err := LoadExtension("ext-id", extensionContent)
	suite.NoError(err)
	err = extension.Install(nil, "version")
	suite.EqualError(err, `extension "ext-id" does not support operation "install"`)
}

func minimalExtension(version string) string {
	content := `
	(function(){
		global.installedExtension = {
			extension: {},
			apiVersion: "$VERSION$"
		}
	})()`
	return strings.Replace(content, "$VERSION$", version, 1)
}

func (suite *ExtensionApiSuite) TestLoadExtensionWithoutSettingGlobalVariable() {
	extension, err := LoadExtension("ext-id", `(function(){ })()`)
	suite.EqualError(err, `extension "ext-id" did not set global.installedExtension`)
	suite.Nil(extension)
}

func (suite *ExtensionApiSuite) TestLoadExtensionInvalidJavaScript() {
	extension, err := LoadExtension("ext-id", `invalid javascript`)
	suite.ErrorContains(err, `failed to run extension "ext-id"`)
	suite.ErrorContains(err, "SyntaxError")
	suite.ErrorContains(err, "Unexpected identifier")
	suite.Nil(extension)
}

func (suite *ExtensionApiSuite) mockContext() *context.ExtensionContext {
	return createMockContextWithClients(&suite.mockSQLClient, &suite.mockBucketFsClient)
}

func (suite *ExtensionApiSuite) loadExtension(content string) *JsExtension {
	extension, err := LoadExtension("ext-id", content)
	if err != nil {
		suite.T().Fatalf("loading extension failed: %v", err)
	}
	return extension
}
