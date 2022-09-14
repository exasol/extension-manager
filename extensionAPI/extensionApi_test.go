package extensionAPI

import (
	"os"
	"path"
	"testing"

	"github.com/exasol/extension-manager/backend"
	"github.com/exasol/extension-manager/integrationTesting"

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

func (suite *ExtensionApiSuite) Test_GetExtensionFromFile() {
	extensionFile := integrationTesting.CreateTestExtensionBuilder(suite.T()).Build().WriteToTmpFile()
	extension := suite.loadExtension(extensionFile)
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

func (suite *ExtensionApiSuite) Test_Install() {
	extensionFile := integrationTesting.CreateTestExtensionBuilder(suite.T()).Build().WriteToTmpFile()
	extension := suite.loadExtension(extensionFile)
	suite.mockSQLClient.On("Execute", "select 1", []any{}).Return()
	err := extension.Install(suite.mockContext(), "extVersion")
	suite.NoError(err)
}

func (suite *ExtensionApiSuite) Test_Install_ResolveBucketFsPath() {
	extensionFile := integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithInstallFunc("context.sqlClient.execute(`create script path ${context.bucketFs.resolvePath('my-adapter-'+version+'.jar')}`)").
		Build().WriteToTmpFile()
	extension := suite.loadExtension(extensionFile)
	suite.mockSQLClient.On("Execute", "create script path /buckets/bfsdefault/default/my-adapter-extensionVersion.jar", []any{}).Return()
	err := extension.Install(suite.mockContext(), "extensionVersion")
	suite.NoError(err)
}

func (suite *ExtensionApiSuite) Test_Uninstall() {
	extensionFile := integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithUninstallFunc("context.sqlClient.execute(`uninstall version ${version}`)").
		Build().WriteToTmpFile()
	extension := suite.loadExtension(extensionFile)
	suite.mockSQLClient.On("Execute", "uninstall version extVersion", []any{}).Return()
	err := extension.Uninstall(suite.mockContext(), "extVersion")
	suite.NoError(err)
}

func (suite *ExtensionApiSuite) Test_AddInstance_validParameters() {
	extensionFile := integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithAddInstanceFunc("context.sqlClient.execute('create vs');\n" +
			"return {id: 'instId', name: `instance_${version}_${params.values[0].name}_${params.values[0].value}`};").
		Build().WriteToTmpFile()
	extension := suite.loadExtension(extensionFile)
	suite.mockSQLClient.On("Execute", "create vs", []any{}).Return()
	instance, err := extension.AddInstance(suite.mockContext(), "extensionVersion", &ParameterValues{Values: []ParameterValue{{Name: "p1", Value: "v1"}}})
	suite.NoError(err)
	suite.Equal(&JsExtInstance{Id: "instId", Name: "instance_extensionVersion_p1_v1"}, instance)
}

func (suite *ExtensionApiSuite) Test_ListInstances_EmptyResult() {
	extensionFile := integrationTesting.CreateTestExtensionBuilder(suite.T()).
		Build().WriteToTmpFile()
	extension := suite.loadExtension(extensionFile)
	instances, err := extension.ListInstances(suite.mockContext(), "ver")
	suite.NoError(err)
	suite.Empty(instances)
}

func (suite *ExtensionApiSuite) Test_ListInstances_NonEmptyResult() {
	extensionFile := integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithFindInstancesFunc(`return [{id: "instId", name: "instName"}]`).
		Build().WriteToTmpFile()
	extension := suite.loadExtension(extensionFile)
	instances, err := extension.ListInstances(suite.mockContext(), "ver")
	suite.NoError(err)
	suite.Equal([]*JsExtInstance{{Id: "instId", Name: "instName"}}, instances)
}

func (suite *ExtensionApiSuite) Test_DeleteInstance() {
	extensionFile := integrationTesting.CreateTestExtensionBuilder(suite.T()).
		Build().WriteToTmpFile()
	extension := suite.loadExtension(extensionFile)
	suite.mockSQLClient.On("Execute", "drop instance instId", []any{}).Return()
	err := extension.DeleteInstance(suite.mockContext(), "instId")
	suite.NoError(err)
}

func createMockMetadata() *ExaMetadata {
	exaAllScripts := ExaScriptTable{Rows: []ExaScriptRow{{Name: "test"}}}
	return &ExaMetadata{AllScripts: exaAllScripts}
}

func (suite *ExtensionApiSuite) Test_FindInstallationsCanReadAllScriptsTable() {
	extensionFile := integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithFindInstallationsFunc(`
		return metadata.allScripts.rows.map(row => {
			return {name: row.name, version: "0.1.0", instanceParameters: []}
		});`).Build().WriteToTmpFile()
	extension := suite.loadExtension(extensionFile)
	result, err := extension.FindInstallations(createMockContext(), createMockMetadata())
	suite.Equal([]*JsExtInstallation{{Name: "test", Version: "0.1.0", InstanceParameters: []interface{}{}}}, result)
	suite.NoError(err)
}

func (suite *ExtensionApiSuite) Test_FindInstallationsReturningParameters() {
	extensionFile := integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithFindInstallationsFunc(integrationTesting.
			MockFindInstallationsFunction("test", "0.1.0", `[{
		id: "param1",
		name: "My param",
		type: "string"
	}]`)).Build().WriteToTmpFile()
	extension := suite.loadExtension(extensionFile)
	result, err := extension.FindInstallations(createMockContext(), createMockMetadata())
	suite.Equal([]*JsExtInstallation{{Name: "test", Version: "0.1.0", InstanceParameters: []interface{}{map[string]interface{}{"id": "param1", "name": "My param", "type": "string"}}}}, result)
	suite.NoError(err)
}

func (suite *ExtensionApiSuite) Test_GetExtensionFromFile_withOutdatedApiVersion() {
	extensionFile := suite.writeExtension(`
	(function(){
		global.installedExtension = {
			extension: {},
			apiVersion: "0.0.0"
		}
	})()`)
	extension, err := GetExtensionFromFile(extensionFile)
	suite.ErrorContains(err, `incompatible extension API version "0.0.0". Please update the extension to use supported version "`+SupportedApiVersion+`"`)
	suite.Nil(extension)
}

func (suite *ExtensionApiSuite) Test_GetExtensionFromFile_withoutSettingGlobalVariable() {
	extensionFile := suite.writeExtension(`(function(){ })()`)
	extension, err := GetExtensionFromFile(extensionFile)
	suite.EqualError(err, "extension did not set global.installedExtension")
	suite.Nil(extension)
}

func (suite *ExtensionApiSuite) Test_GetExtensionFromFile_invalidJavaScript() {
	extensionFile := suite.writeExtension(`invalid javascript`)
	extension, err := GetExtensionFromFile(extensionFile)
	suite.ErrorContains(err, "failed to run extension file")
	suite.ErrorContains(err, "SyntaxError")
	suite.ErrorContains(err, "Unexpected identifier")
	suite.Nil(extension)
}

func (suite *ExtensionApiSuite) Test_GetExtensionFromFile_invalidFileName() {
	extension, err := GetExtensionFromFile("no-such-file")
	suite.ErrorContains(err, "failed to open extension file no-such-file")
	suite.Nil(extension)
}

func (suite *ExtensionApiSuite) writeExtension(extensionJs string) string {
	extensionFile := path.Join(suite.T().TempDir(), "extension.js")
	suite.NoError(os.WriteFile(extensionFile, []byte(extensionJs), 0600))
	return extensionFile
}

func (suite *ExtensionApiSuite) mockContext() *ExtensionContext {
	return createMockContextWithSqlClient(&suite.mockSQLClient)
}

func (suite *ExtensionApiSuite) loadExtension(extensionFile string) *JsExtension {
	extension, err := GetExtensionFromFile(extensionFile)
	if err != nil {
		suite.T().Fatalf("loading extension failed: %v", err)
	}
	return extension
}
