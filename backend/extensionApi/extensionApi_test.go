package extensionApi

import (
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"os/exec"
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
	suite.validExtensionFile = buildExtensionForTesting(suite.T())
}

func (suite *ExtensionApiSuite) Test_GetExtensionFromFile() {
	extension, err := GetExtensionFromFile(suite.validExtensionFile, LoggingSimpleSqlClient{})
	suite.NoError(err)
	suite.Assert().Equal("MyDemoExtension", extension.GetName())
}

type MockSimpleSqlClient struct {
	mock.Mock
}

func (mock *MockSimpleSqlClient) RunSqlQuery(query string) {
	mock.Called(query)
}

func (suite *ExtensionApiSuite) Test_GetExtensionFromFile_install() {
	mockSqlClient := MockSimpleSqlClient{}
	mockSqlClient.On("RunSqlQuery", "CREATE ADAPTER SCRIPT ...").Return()
	extension, err := GetExtensionFromFile(suite.validExtensionFile, &mockSqlClient)
	suite.NoError(err)
	err = extension.Install()
	suite.NoError(err)
	mockSqlClient.AssertCalled(suite.T(), "RunSqlQuery", "CREATE ADAPTER SCRIPT ...")
}

func (suite *ExtensionApiSuite) Test_GetExtensionFromFile_withOutdatedApiVersion() {
	extensionFile := suite.writeExtension(`(function(){
	installedExtension = {
		extension: {},
		apiVersion: "0.0.0"
	}
	})()`)
	_, err := GetExtensionFromFile(extensionFile, LoggingSimpleSqlClient{})
	suite.Error(err)
	suite.Assert().Contains(err.Error(), "incompatible extension API version 0.0.0. Please update the extension to use a supported version of the extension API")
}

func (suite *ExtensionApiSuite) Test_GetExtensionFromFile_withNonCallToInstallExtension() {
	extensionFile := suite.writeExtension(`(function(){
	})()`)
	_, err := GetExtensionFromFile(extensionFile, LoggingSimpleSqlClient{})
	suite.Error(err)
	suite.Assert().Contains(err.Error(), "invalid installedExtension. The provided JS file did not set the installedExtension variable or it is not an object. Make sure that the extension fill calls installExtension")
}

func (suite *ExtensionApiSuite) writeExtension(extensionJs string) string {
	extensionFile := path.Join(suite.T().TempDir(), "extension.js")
	suite.NoError(ioutil.WriteFile(extensionFile, []byte(extensionJs), 0600))
	return extensionFile
}

func buildExtensionForTesting(t *testing.T) string {
	const extensionForTestingDir = "extensionForTesting"
	installCommand := exec.Command("npm", "install")
	installCommand.Dir = extensionForTestingDir
	err := installCommand.Run()
	if err != nil {
		t.Errorf("Failed to install node modules (run 'npm install') for extensionForTesting. Cause: %v", err.Error())
	}
	buildCommand := exec.Command("npm", "run", "build")
	buildCommand.Dir = extensionForTestingDir
	err = buildCommand.Run()
	if err != nil {
		t.Errorf("Failed to build extensionForTesting. Cause: %v", err.Error())
	}
	return path.Join(extensionForTestingDir, "dist.js")
}
