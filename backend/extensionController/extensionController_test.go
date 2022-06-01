package extensionController

import (
	"backend/integrationTesting"
	"io"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ExtensionControllerSuite struct {
	integrationTesting.IntegrationTestSuite
	tempExtensionRepo string
}

func TestExtensionControllerSuite(t *testing.T) {
	suite.Run(t, new(ExtensionControllerSuite))
}

func (suite *ExtensionControllerSuite) SetupSuite() {
	suite.IntegrationTestSuite.SetupSuite()
	tempExtensionRepo, err := os.MkdirTemp(os.TempDir(), "ExtensionControllerSuite")
	if err != nil {
		panic(err)
	}
	extensionPath := integrationTesting.GetExtensionForTesting("../")
	suite.copyToExtensionRepo(extensionPath, tempExtensionRepo)
	suite.tempExtensionRepo = tempExtensionRepo
}

func (suite *ExtensionControllerSuite) TearDownSuite() {
	err := os.RemoveAll(suite.tempExtensionRepo)
	if err != nil {
		panic(err)
	}
	suite.IntegrationTestSuite.TearDownSuite()
}

func (suite *ExtensionControllerSuite) copyToExtensionRepo(extensionPath string, tempExtensionRepo string) {
	extensionFile, err := os.Open(extensionPath)
	if err != nil {
		panic(err)
	}
	targetPath := path.Join(tempExtensionRepo, "myExtension.js")
	targetFile, err := os.Create(targetPath)
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(targetFile, extensionFile)
	if err != nil {
		panic(err)
	}
}

func (suite *ExtensionControllerSuite) TestGetAllExtensions() {
	suite.NoError(suite.Exasol.UploadStringContent("123", "my-extension.1.2.3.jar")) // create file with 3B size
	defer func() { suite.NoError(suite.Exasol.DeleteFile("my-extension.1.2.3.jar")) }()
	controller := Create(suite.tempExtensionRepo)
	dbConnectionWithNoAutocommit, err := suite.Exasol.CreateConnectionWithConfig(false)
	suite.NoError(err)
	defer func() { suite.NoError(dbConnectionWithNoAutocommit.Close()) }()
	extensions, err := controller.GetAllExtensions(dbConnectionWithNoAutocommit)
	suite.NoError(err)
	suite.Assert().Equal("MyDemoExtension", extensions[0].Name)
}

func (suite *ExtensionControllerSuite) TestGetAllExtensionsWithMissingJar() {
	controller := Create(suite.tempExtensionRepo)
	dbConnectionWithNoAutocommit, err := suite.Exasol.CreateConnectionWithConfig(false)
	suite.NoError(err)
	defer func() { suite.NoError(dbConnectionWithNoAutocommit.Close()) }()
	extensions, err := controller.GetAllExtensions(dbConnectionWithNoAutocommit)
	suite.NoError(err)
	suite.Assert().Empty(extensions)
}

func (suite *ExtensionControllerSuite) TestGetAllInstallations() {
	controller := Create(suite.tempExtensionRepo)
	luaScriptFixture := integrationTesting.CreateLuaScriptFixture(suite.Connection)
	defer luaScriptFixture.Close()
	installations, err := controller.GetAllInstallations(suite.Connection)
	suite.NoError(err)
	suite.Assert().Equal("TEST.MY_SCRIPT", installations[0].Name)
}
