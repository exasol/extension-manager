package extensionAPI

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ErrorHandlingExtensionSuite struct {
	suite.Suite
	rawExtension *rawJsExtension
	extension    *JsExtension
}

func TestJsExtensionSuite(t *testing.T) {
	suite.Run(t, new(ErrorHandlingExtensionSuite))
}

func (suite *ErrorHandlingExtensionSuite) SetupSuite() {
	suite.rawExtension = &rawJsExtension{Name: "name", Description: "desc", InstallableVersions: []string{"v1", "v2"}, BucketFsUploads: []BucketFsUpload{{Name: "uploadName"}}}
	suite.extension = wrapExtension(suite.rawExtension, "id", nil)
}

func (suite *ErrorHandlingExtensionSuite) TestProperties() {
	suite.Assert().Equal(&JsExtension{
		Id:                  "id",
		Name:                "name",
		Description:         "desc",
		InstallableVersions: []string{"v1", "v2"},
		BucketFsUploads:     []BucketFsUpload{{Name: "uploadName"}},
		extension:           suite.rawExtension},
		suite.extension)
}

func createMockContextWithSqlClient(sqlClient SimpleSQLClient) *ExtensionContext {
	return CreateContextWithClient("extension_schema", sqlClient)
}

func createMockContext() *ExtensionContext {
	var client SimpleSQLClient = &MockSimpleSQLClient{}
	return CreateContextWithClient("extension_schema", client)
}

func (suite *ErrorHandlingExtensionSuite) TestFindInstallationsSuccessful() {
	expectedInstallations := []*JsExtInstallation{{Name: "instName"}}
	suite.rawExtension.FindInstallations = func(context *ExtensionContext, metadata *ExaMetadata) []*JsExtInstallation {
		return expectedInstallations
	}
	installations, err := suite.extension.FindInstallations(createMockContext(), &ExaMetadata{})
	suite.NoError(err)
	suite.Equal(expectedInstallations, installations)
}

func (suite *ErrorHandlingExtensionSuite) TestFindInstallationsFailure() {
	suite.rawExtension.FindInstallations = func(context *ExtensionContext, metadata *ExaMetadata) []*JsExtInstallation {
		panic("mock error")
	}
	installations, err := suite.extension.FindInstallations(createMockContext(), &ExaMetadata{})
	suite.EqualError(err, "failed to find installations for extension \"id\": mock error")
	suite.Nil(installations)
}

func (suite *ErrorHandlingExtensionSuite) TestInstallSuccessful() {
	suite.rawExtension.Install = func(context *ExtensionContext, version string) {
	}
	err := suite.extension.Install(createMockContext(), "version")
	suite.NoError(err)
}

func (suite *ErrorHandlingExtensionSuite) TestInstallFailure() {
	suite.rawExtension.Install = func(context *ExtensionContext, version string) {
		panic("mock error")
	}
	err := suite.extension.Install(createMockContext(), "version")
	suite.EqualError(err, "failed to install extension \"id\": mock error")
}

func (suite *ErrorHandlingExtensionSuite) TestAddInstanceSuccessful() {
	suite.rawExtension.AddInstance = func(context *ExtensionContext, version string, params *ParameterValues) *JsExtInstance {
		return &JsExtInstance{Name: "newInstance"}
	}
	instance, err := suite.extension.AddInstance(createMockContext(), "version", &ParameterValues{})
	suite.NoError(err)
	suite.Equal(&JsExtInstance{Name: "newInstance"}, instance)
}

func (suite *ErrorHandlingExtensionSuite) TestAddInstanceFails() {
	suite.rawExtension.AddInstance = func(context *ExtensionContext, version string, params *ParameterValues) *JsExtInstance {
		panic("mock error")
	}
	instance, err := suite.extension.AddInstance(createMockContext(), "version", &ParameterValues{})
	suite.EqualError(err, "failed to add instance for extension \"id\": mock error")
	suite.Nil(instance)
}
