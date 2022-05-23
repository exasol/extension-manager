package extensionController

import (
	"backend/integrationTesting"
	"github.com/stretchr/testify/suite"
	"io"
	"os"
	"path"
	"testing"
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
	controller := Create(suite.tempExtensionRepo)
	extensions, err := controller.GetAllExtensions()
	suite.NoError(err)
	suite.Assert().Equal("MyDemoExtension", extensions[0].Name)
}

func (suite *ExtensionControllerSuite) TestGetAllInstallations() {
	controller := Create(suite.tempExtensionRepo)
	luaScriptFixture := integrationTesting.CreateLuaScriptFixture(suite.Connection)
	defer luaScriptFixture.Close()
	installations, err := controller.GetAllInstallations(suite.Connection)
	suite.NoError(err)
	suite.Assert().Equal("TEST.MY_SCRIPT", installations[0].Name)
}
