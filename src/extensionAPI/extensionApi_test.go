package extensionAPI

import (
	"testing"

	"github.com/exasol/extension-manager/src/backend"
	"github.com/exasol/extension-manager/src/integrationTesting"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ExtensionApiSuite struct {
	suite.Suite
	mockSQLClient sqlClientMock
}

func TestExtensionApiSuite(t *testing.T) {
	suite.Run(t, new(ExtensionApiSuite))
}

func (suite *ExtensionApiSuite) SetupSuite() {
}

func (suite *ExtensionApiSuite) SetupTest() {
	suite.mockSQLClient = sqlClientMock{}
}

func (suite *ExtensionApiSuite) TearDownTest() {
	suite.mockSQLClient.AssertExpectations(suite.T())
}

func (suite *ExtensionApiSuite) Test_LoadExtension() {
	extensionContent := integrationTesting.CreateTestExtensionBuilder(suite.T()).Build().AsString()
	extension := suite.loadExtension(extensionContent)
	suite.Equal("MyDemoExtension", extension.Name)
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

func (suite *ExtensionApiSuite) Test_GetParameterDefinitions_EmptyResult() {
	extensionContent := integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithGetInstanceParameterDefinitionFunc(`return []`).
		Build().AsString()
	extension := suite.loadExtension(extensionContent)
	definitions := extension.extension.GetParameterDefinitions(suite.mockContext(), "extVersion")
	suite.Equal([]interface{}{}, definitions)
}

func (suite *ExtensionApiSuite) Test_GetParameterDefinitions() {
	extensionContent := integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithGetInstanceParameterDefinitionFunc(`return [{id: "param1", name: "My param", type: "string"}]`).
		Build().AsString()
	extension := suite.loadExtension(extensionContent)
	definitions := extension.extension.GetParameterDefinitions(suite.mockContext(), "extVersion")
	suite.Equal([]interface{}{map[string]interface{}{"id": "param1", "name": "My param", "type": "string"}}, definitions)
}

func (suite *ExtensionApiSuite) Test_Install() {
	extensionContent := integrationTesting.CreateTestExtensionBuilder(suite.T()).Build().AsString()
	extension := suite.loadExtension(extensionContent)
	suite.mockSQLClient.On("Execute", "select 1", []any{}).Return()
	err := extension.Install(suite.mockContext(), "extVersion")
	suite.NoError(err)
}

func (suite *ExtensionApiSuite) Test_Install_ResolveBucketFsPath() {
	extensionContent := integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithInstallFunc("context.sqlClient.execute(`create script path ${context.bucketFs.resolvePath('my-adapter-'+version+'.jar')}`)").
		Build().AsString()
	extension := suite.loadExtension(extensionContent)
	suite.mockSQLClient.On("Execute", "create script path /buckets/bfsdefault/default/my-adapter-extensionVersion.jar", []any{}).Return()
	err := extension.Install(suite.mockContext(), "extensionVersion")
	suite.NoError(err)
}

func (suite *ExtensionApiSuite) Test_Install_ConsoleLog() {
	extensionContent := integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithInstallFunc("console.log('test log message')").
		Build().AsString()
	extension := suite.loadExtension(extensionContent)
	err := extension.Install(suite.mockContext(), "extensionVersion")
	suite.NoError(err)
}

func (suite *ExtensionApiSuite) Test_Uninstall() {
	extensionContent := integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithUninstallFunc("context.sqlClient.execute(`uninstall version ${version}`)").
		Build().AsString()
	extension := suite.loadExtension(extensionContent)
	suite.mockSQLClient.On("Execute", "uninstall version extVersion", []any{}).Return()
	err := extension.Uninstall(suite.mockContext(), "extVersion")
	suite.NoError(err)
}

func (suite *ExtensionApiSuite) Test_AddInstance_validParameters() {
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

func (suite *ExtensionApiSuite) Test_ListInstances_EmptyResult() {
	extensionContent := integrationTesting.CreateTestExtensionBuilder(suite.T()).
		Build().AsString()
	extension := suite.loadExtension(extensionContent)
	instances, err := extension.ListInstances(suite.mockContext(), "ver")
	suite.NoError(err)
	suite.Empty(instances)
}

func (suite *ExtensionApiSuite) Test_ListInstances_NonEmptyResult() {
	extensionContent := integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithFindInstancesFunc(`return [{id: "instId", name: "instName"}]`).
		Build().AsString()
	extension := suite.loadExtension(extensionContent)
	instances, err := extension.ListInstances(suite.mockContext(), "ver")
	suite.NoError(err)
	suite.Equal([]*JsExtInstance{{Id: "instId", Name: "instName"}}, instances)
}

func (suite *ExtensionApiSuite) Test_DeleteInstance() {
	extensionContent := integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithDeleteInstanceFunc("context.sqlClient.execute(`drop instance ${instanceId}`)").
		Build().AsString()
	extension := suite.loadExtension(extensionContent)
	suite.mockSQLClient.On("Execute", "drop instance instId", []any{}).Return()
	err := extension.DeleteInstance(suite.mockContext(), "extVersion", "instId")
	suite.NoError(err)
}

func createMockMetadata() *ExaMetadata {
	exaAllScripts := ExaScriptTable{Rows: []ExaScriptRow{{Name: "test"}}}
	return &ExaMetadata{AllScripts: exaAllScripts}
}

func (suite *ExtensionApiSuite) Test_FindInstallationsCanReadAllScriptsTable() {
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

func (suite *ExtensionApiSuite) Test_FindInstallationsReturningParameters() {
	extensionContent := integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithFindInstallationsFunc(integrationTesting.
			MockFindInstallationsFunction("test", "0.1.0")).Build().AsString()
	extension := suite.loadExtension(extensionContent)
	result, err := extension.FindInstallations(createMockContext(), createMockMetadata())
	suite.Equal([]*JsExtInstallation{{Name: "test", Version: "0.1.0"}}, result)
	suite.NoError(err)
}

func (suite *ExtensionApiSuite) Test_LoadExtension_withOutdatedApiVersion() {
	extensionContent := `
	(function(){
		global.installedExtension = {
			extension: {},
			apiVersion: "0.0.0"
		}
	})()`
	extension, err := LoadExtension("ext-id", extensionContent)
	suite.ErrorContains(err, `incompatible extension API version "0.0.0". Please update the extension to use supported version "`+SupportedApiVersion+`"`)
	suite.Nil(extension)
}

func (suite *ExtensionApiSuite) Test_LoadExtension_withoutSettingGlobalVariable() {
	extension, err := LoadExtension("ext-id", `(function(){ })()`)
	suite.EqualError(err, `extension "ext-id" did not set global.installedExtension`)
	suite.Nil(extension)
}

func (suite *ExtensionApiSuite) Test_LoadExtension_invalidJavaScript() {
	extension, err := LoadExtension("ext-id", `invalid javascript`)
	suite.ErrorContains(err, `failed to run extension "ext-id"`)
	suite.ErrorContains(err, "SyntaxError")
	suite.ErrorContains(err, "Unexpected identifier")
	suite.Nil(extension)
}

func (suite *ExtensionApiSuite) mockContext() *ExtensionContext {
	return createMockContextWithSqlClient(&suite.mockSQLClient)
}

func (suite *ExtensionApiSuite) loadExtension(content string) *JsExtension {
	extension, err := LoadExtension("ext-id", content)
	if err != nil {
		suite.T().Fatalf("loading extension failed: %v", err)
	}
	return extension
}