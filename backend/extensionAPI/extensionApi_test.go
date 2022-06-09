package extensionAPI

import (
	"backend/integrationTesting"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

type ExtensionApiSuite struct {
	suite.Suite
	validExtensionFile string
}

func TestExtensionApiSuite(t *testing.T) {
	suite.Run(t, new(ExtensionApiSuite))
}

func (suite *ExtensionApiSuite) SetupSuite() {
	suite.validExtensionFile = integrationTesting.CreateTestExtensionBuilder().Build().WriteToTmpFile()
}

func (suite *ExtensionApiSuite) TearDownSuite() {
	err := os.Remove(suite.validExtensionFile)
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

func (suite *ExtensionApiSuite) Test_GetExtensionFromFile_FindInstallations() {
	mockSqlClient := MockSimpleSQLClient{}
	extension, err := GetExtensionFromFile(suite.validExtensionFile)
	suite.NoError(err)
	exaAllScripts := ExaAllScriptTable{Rows: []ExaAllScriptRow{{Name: "test"}}}
	result := extension.FindInstallations(&mockSqlClient, &exaAllScripts)
	suite.Assert().Equal("test", result[0].Name)
	suite.NoError(err)
}

func (suite *ExtensionApiSuite) Test_GetExtensionFromFile_withOutdatedApiVersion() {
	extensionFile := suite.writeExtension(`(function(){
	installedExtension = {
		extension: {},
		apiVersion: "0.0.0"
	}
	})()`)
	_, err := GetExtensionFromFile(extensionFile)
	suite.Error(err)
	suite.Assert().Contains(err.Error(), "incompatible extension API version 0.0.0. Please update the extension to use a supported version of the extension API")
}

func (suite *ExtensionApiSuite) writeExtension(extensionJs string) string {
	extensionFile := path.Join(suite.T().TempDir(), "extension.js")
	suite.NoError(ioutil.WriteFile(extensionFile, []byte(extensionJs), 0600))
	return extensionFile
}
