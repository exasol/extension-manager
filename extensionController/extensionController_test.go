package extensionController

import (
	"os"
	"path"
	"testing"

	"github.com/exasol/extension-manager/integrationTesting"

	"github.com/stretchr/testify/suite"
)

const (
	EXTENSION_SCHEMA   = "test"
	EXTENSION_FILENAME = "my-extension.js"
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
	controller := Create(suite.tempExtensionRepo, EXTENSION_SCHEMA)
	extensions, err := controller.GetAllExtensions(suite.Connection)
	suite.NoError(err)
	suite.Assert().Equal(1, len(extensions))
	suite.Assert().Equal("MyDemoExtension", extensions[0].Name, "name")
	suite.Assert().Equal(EXTENSION_FILENAME, extensions[0].Id, "id")
}

func (suite *ExtensionControllerSuite) writeDefaultExtension() {
	integrationTesting.CreateTestExtensionBuilder().
		WithBucketFsUpload(integrationTesting.BucketFsUploadParams{Name: "extension jar", BucketFsFilename: "my-extension.1.2.3.jar", FileSize: 3}).
		WithFindInstallationsFunc(`
		return metadata.allScripts.rows.map(row => {
			return {name: row.schema + "." + row.name, version: "0.1.0", instanceParameters: []}
		});`).
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_FILENAME))
}

func (suite *ExtensionControllerSuite) TestGetAllExtensionsWithMissingJar() {
	integrationTesting.CreateTestExtensionBuilder().
		WithBucketFsUpload(integrationTesting.BucketFsUploadParams{Name: "extension jar", BucketFsFilename: "missing-jar.jar", FileSize: 3}).
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_FILENAME))
	controller := Create(suite.tempExtensionRepo, EXTENSION_SCHEMA)
	dbConnectionWithNoAutocommit, err := suite.Exasol.CreateConnectionWithConfig(false)
	suite.NoError(err)
	defer func() { suite.NoError(dbConnectionWithNoAutocommit.Close()) }()
	extensions, err := controller.GetAllExtensions(dbConnectionWithNoAutocommit)
	suite.NoError(err)
	suite.Assert().Empty(extensions)
}

func (suite *ExtensionControllerSuite) TestGetAllExtensionsThrowingJSError() {
	const jarName = "my-failing-extension-1.2.3.jar"
	integrationTesting.CreateTestExtensionBuilder().
		WithBucketFsUpload(integrationTesting.BucketFsUploadParams{Name: "extension jar", BucketFsFilename: jarName, FileSize: 3}).
		WithFindInstallationsFunc("throw Error(`mock error from js`)").
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_FILENAME))
	suite.NoError(suite.Exasol.UploadStringContent("123", jarName)) // create file with 3B size
	defer func() { suite.NoError(suite.Exasol.DeleteFile(jarName)) }()
	controller := Create(suite.tempExtensionRepo, EXTENSION_SCHEMA)
	extensions, err := controller.GetAllInstallations(suite.Connection)
	suite.ErrorContains(err, `failed to find installations from extension "my-extension.js": failed to find installations for extension "my-extension.js": Error: mock error from js at`)
	suite.Nil(extensions)
}

func (suite *ExtensionControllerSuite) TestGetAllInstallations() {
	suite.writeDefaultExtension()
	fixture := integrationTesting.CreateLuaScriptFixture(suite.Connection)
	controller := Create(suite.tempExtensionRepo, fixture.GetSchemaName())
	defer fixture.Close()
	installations, err := controller.GetAllInstallations(suite.Connection)
	suite.NoError(err)
	suite.Assert().Equal(1, len(installations))
	suite.Assert().Equal(fixture.GetSchemaName()+".MY_SCRIPT", installations[0].Name)
}
