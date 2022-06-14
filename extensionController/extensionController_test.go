package extensionController

import (
	"os"
	"path"
	"testing"

	"github.com/exasol/extension-manager/integrationTesting"

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
}

func (suite *ExtensionControllerSuite) TearDownSuite() {
	suite.IntegrationTestSuite.TearDownSuite()
}

func (suite *ExtensionControllerSuite) SetupTest() {
	tempExtensionRepo, err := os.MkdirTemp(os.TempDir(), "ExtensionControllerSuite")
	if err != nil {
		panic(err)
	}
	suite.tempExtensionRepo = tempExtensionRepo
}

func (suite *ExtensionControllerSuite) AfterTest(suiteName, testName string) {
	err := os.RemoveAll(suite.tempExtensionRepo)
	if err != nil {
		panic(err)
	}
}

func (suite *ExtensionControllerSuite) TestGetAllExtensions() {
	suite.writeDefaultExtension()
	suite.NoError(suite.Exasol.UploadStringContent("123", "my-extension.1.2.3.jar")) // create file with 3B size
	defer func() { suite.NoError(suite.Exasol.DeleteFile("my-extension.1.2.3.jar")) }()
	controller := Create(suite.tempExtensionRepo)
	extensions, err := controller.GetAllExtensions(suite.Connection)
	suite.NoError(err)
	suite.Assert().Equal("MyDemoExtension", extensions[0].Name)
}

func (suite *ExtensionControllerSuite) writeDefaultExtension() {
	integrationTesting.CreateTestExtensionBuilder().
		WithBucketFsUpload(integrationTesting.BucketFsUploadParams{Name: "extension jar", BucketFsFilename: "my-extension.1.2.3.jar", FileSize: 3}).
		WithFindInstallationsFunc(`
		return exaAllScripts.rows.map(row => {
			return {name: row.name, version: "0.1.0", instanceParameters: []}
		});`).
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, "myExtension.js"))
}

func (suite *ExtensionControllerSuite) TestGetAllExtensionsWithMissingJar() {
	integrationTesting.CreateTestExtensionBuilder().
		WithBucketFsUpload(integrationTesting.BucketFsUploadParams{Name: "extension jar", BucketFsFilename: "missing-jar.jar", FileSize: 3}).
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, "myExtension.js"))
	controller := Create(suite.tempExtensionRepo)
	dbConnectionWithNoAutocommit, err := suite.Exasol.CreateConnectionWithConfig(false)
	suite.NoError(err)
	defer func() { suite.NoError(dbConnectionWithNoAutocommit.Close()) }()
	extensions, err := controller.GetAllExtensions(dbConnectionWithNoAutocommit)
	suite.NoError(err)
	suite.Assert().Empty(extensions)
}

func (suite *ExtensionControllerSuite) TestGetAllInstallations() {
	suite.writeDefaultExtension()
	controller := Create(suite.tempExtensionRepo)
	luaScriptFixture := integrationTesting.CreateLuaScriptFixture(suite.Connection)
	defer luaScriptFixture.Close()
	installations, err := controller.GetAllInstallations(suite.Connection)
	suite.NoError(err)
	suite.Assert().Equal("TEST.MY_SCRIPT", installations[0].Name)
}
