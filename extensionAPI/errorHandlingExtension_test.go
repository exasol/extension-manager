package extensionAPI

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type JsExtensionSuite struct {
	suite.Suite
	rawExtension *rawJsExtension
	extension    *JsExtension
}

func TestJsExtensionSuite(t *testing.T) {
	suite.Run(t, new(JsExtensionSuite))
}

func (suite *JsExtensionSuite) SetupSuite() {
	suite.rawExtension = &rawJsExtension{Id: "id", Name: "name", Description: "desc", InstallableVersions: []string{"v1", "v2"}, BucketFsUploads: []BucketFsUpload{{Name: "uploadName"}}}
	suite.extension = wrapExtension(suite.rawExtension)
}

func (suite *JsExtensionSuite) TestProperties() {
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
	return &ExtensionContext{ExtensionSchemaName: "ext_schema", SqlClient: sqlClient}
}

func createMockContext() *ExtensionContext {
	var client SimpleSQLClient = &MockSimpleSQLClient{}
	return createMockContextWithSqlClient(client)
}

func (suite *JsExtensionSuite) TestFindInstallationsSuccessful() {
	expectedInstallations := []*JsExtInstallation{{Name: "instName"}}
	suite.rawExtension.FindInstallations = func(context *ExtensionContext, metadata *ExaMetadata) []*JsExtInstallation {
		return expectedInstallations
	}
	installations, err := suite.extension.FindInstallations(createMockContext(), &ExaMetadata{})
	suite.NoError(err)
	suite.Equal(expectedInstallations, installations)
}

func (suite *JsExtensionSuite) TestFindInstallationsFailure() {
	suite.rawExtension.FindInstallations = func(context *ExtensionContext, metadata *ExaMetadata) []*JsExtInstallation {
		panic("mock error")
	}
	installations, err := suite.extension.FindInstallations(createMockContext(), &ExaMetadata{})
	suite.EqualError(err, "failed to find installations for extension \"id\": mock error")
	suite.Nil(installations)
}

func (suite *JsExtensionSuite) TestInstallSuccessful() {
	suite.rawExtension.Install = func(context *ExtensionContext, version string) {
	}
	err := suite.extension.Install(createMockContext(), "version")
	suite.NoError(err)
}

func (suite *JsExtensionSuite) TestInstallFailure() {
	suite.rawExtension.Install = func(context *ExtensionContext, version string) {
		panic("mock error")
	}
	err := suite.extension.Install(createMockContext(), "version")
	suite.EqualError(err, "failed to install extension \"id\": mock error")
}
