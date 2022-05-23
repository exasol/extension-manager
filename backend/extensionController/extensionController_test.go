package extensionController

import (
	"backend/integrationTesting"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ExtensionControllerSuite struct {
	integrationTesting.IntegrationTestSuite
}

func TestExtensionControllerSuite(t *testing.T) {
	suite.Run(t, new(ExtensionControllerSuite))
}

func (suite *ExtensionControllerSuite) TestGetAllExtensions() {
	controller := Create("../extensionApi/extensionForTesting/")
	extensions, err := controller.GetAllExtensions()
	suite.NoError(err)
	suite.Assert().Equal("MyDemoExtension", extensions[0].Name)
}

func (suite *ExtensionControllerSuite) TestGetAllInstallations() {
	controller := Create("../extensionApi/extensionForTesting/")
	luaScriptFixture := integrationTesting.CreateLuaScriptFixture(suite.Connection)
	defer luaScriptFixture.Close()
	installations, err := controller.GetAllInstallations(suite.Connection)
	suite.NoError(err)
	suite.Assert().Equal("TEST.MY_SCRIPT", installations[0].Name)
}
