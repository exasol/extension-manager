package extensionAPI

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/exasol/extension-manager/integrationTesting"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ExtensionApiSuite struct {
	suite.Suite
	validExtensionFile string
}

func TestExtensionApiSuite(t *testing.T) {
	suite.Run(t, new(ExtensionApiSuite))
}

func (suite *ExtensionApiSuite) SetupSuite() {
	suite.validExtensionFile = integrationTesting.CreateTestExtensionBuilder().WithFindInstallationsFunc(`
		return exaAllScripts.rows.map(row => {
			return {name: row.name, version: "0.1.0", instanceParameters: []}
		});`).Build().WriteToTmpFile()
}

func (suite *ExtensionApiSuite) TearDownSuite() {
	deleteFileAndCheckError(suite.validExtensionFile)
}

func deleteFileAndCheckError(file string) {
	err := os.Remove(file)
	if err != nil {
		panic(err)
	}
}

func (suite *ExtensionApiSuite) Test_GetExtensionFromFile() {
	extension, err := GetExtensionFromFile(suite.validExtensionFile)
	suite.NoError(err)
	suite.Assert().Equal("MyDemoExtension", extension.Name)
}

type MockSimpleSQLClient struct {
	mock.Mock
}

func (mock *MockSimpleSQLClient) RunQuery(query string) {
	mock.Called(query)
}

func (suite *ExtensionApiSuite) Test_GetExtensionFromFile_Install() {
	mockSQLClient := MockSimpleSQLClient{}
	mockSQLClient.On("RunQuery", "CREATE ADAPTER SCRIPT ...").Return()
	extension, err := GetExtensionFromFile(suite.validExtensionFile)
	suite.NoError(err)
	extension.Install(&mockSQLClient)
	suite.NoError(err)
	mockSQLClient.AssertCalled(suite.T(), "RunQuery", "CREATE ADAPTER SCRIPT ...")
}

func (suite *ExtensionApiSuite) Test_FindInstallationsCanReadAllScriptsTable() {
	extensionFile := integrationTesting.CreateTestExtensionBuilder().WithFindInstallationsFunc(`
		return exaAllScripts.rows.map(row => {
			return {name: row.name, version: "0.1.0", instanceParameters: []}
		});`).Build().WriteToTmpFile()
	defer deleteFileAndCheckError(extensionFile)
	mockSqlClient := MockSimpleSQLClient{}
	extension, err := GetExtensionFromFile(extensionFile)
	suite.NoError(err)
	exaAllScripts := ExaAllScriptTable{Rows: []ExaAllScriptRow{{Name: "test"}}}
	result := extension.FindInstallations(&mockSqlClient, &exaAllScripts)
	suite.Assert().Equal("test", result[0].Name)
	suite.NoError(err)
}

func (suite *ExtensionApiSuite) Test_FindInstallationsReturningParameters() {
	extensionFile := integrationTesting.CreateTestExtensionBuilder().WithFindInstallationsFunc(integrationTesting.
		MockFindInstallationsFunction("test", "0.1.0", `[{
		id: "param1",
		name: "My param",
		type: "string"
	}]`)).Build().WriteToTmpFile()
	defer deleteFileAndCheckError(extensionFile)
	mockSqlClient := MockSimpleSQLClient{}
	extension, err := GetExtensionFromFile(extensionFile)
	suite.NoError(err)
	exaAllScripts := ExaAllScriptTable{Rows: []ExaAllScriptRow{{Name: "test"}}}
	result := extension.FindInstallations(&mockSqlClient, &exaAllScripts)
	suite.Assert().Equal("test", result[0].Name)
	suite.Assert().Equal("0.1.0", result[0].Version)
	suite.Assert().Equal([]interface{}{map[string]interface{}{"id": "param1", "name": "My param", "type": "string"}}, result[0].InstanceParameters)
	suite.NoError(err)
}

func (suite *ExtensionApiSuite) Test_GetExtensionFromFile_withOutdatedApiVersion() {
	extensionFile := suite.writeExtension(`(function(){
	global.installedExtension = {
		extension: {},
		apiVersion: "0.0.0"
	}
	})()`)
	_, err := GetExtensionFromFile(extensionFile)
	suite.Error(err)
	suite.Assert().Contains(err.Error(), `incompatible extension API version "0.0.0". Please update the extension to use supported version "`+SupportedApiVersion+`"`)
}

func (suite *ExtensionApiSuite) Test_GetExtensionFromFile_withoutSettingGlobalVariable() {
	extensionFile := suite.writeExtension(`(function(){ })()`)
	_, err := GetExtensionFromFile(extensionFile)
	suite.Error(err)
	suite.Assert().Contains(err.Error(), "extension did not set global.installedExtension")
}

func (suite *ExtensionApiSuite) Test_GetExtensionFromFile_invalidJavaScript() {
	extensionFile := suite.writeExtension(`invalid javascript`)
	_, err := GetExtensionFromFile(extensionFile)
	suite.Error(err)
	suite.Assert().Contains(err.Error(), "failed to run extension file")
	suite.Assert().Contains(err.Error(), "SyntaxError")
	suite.Assert().Contains(err.Error(), "Unexpected identifier")
}

func (suite *ExtensionApiSuite) writeExtension(extensionJs string) string {
	extensionFile := path.Join(suite.T().TempDir(), "extension.js")
	suite.NoError(ioutil.WriteFile(extensionFile, []byte(extensionJs), 0600))
	return extensionFile
}
