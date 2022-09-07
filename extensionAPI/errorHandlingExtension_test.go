package extensionAPI

import (
	"fmt"
	"testing"

	"github.com/dop251/goja"
	"github.com/exasol/extension-manager/apiErrors"
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
	suite.extension = wrapExtension(suite.rawExtension, "id", newJavaScriptVm())
}

func (suite *ErrorHandlingExtensionSuite) TestProperties() {
	suite.Equal(&JsExtension{
		Id:                  "id",
		Name:                "name",
		Description:         "desc",
		InstallableVersions: []string{"v1", "v2"},
		BucketFsUploads:     []BucketFsUpload{{Name: "uploadName"}},
		extension:           suite.rawExtension,
		vm:                  suite.extension.vm},
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

func (suite *ErrorHandlingExtensionSuite) TestConvertError_nonErrorObject() {
	err := suite.extension.convertError("msg", "dummyError")
	suite.EqualError(err, "msg: dummyError")
	suite.Equal("*errors.errorString", fmt.Sprintf("%T", err))
}

func (suite *ErrorHandlingExtensionSuite) TestConvertError_errorObject() {
	err := suite.extension.convertError("msg", fmt.Errorf("dummyError"))
	suite.EqualError(err, "msg: dummyError")
	suite.Equal("*errors.errorString", fmt.Sprintf("%T", err))
}

func (suite *ErrorHandlingExtensionSuite) TestConvertError_nilGojaException() {
	var exception goja.Exception
	err := suite.extension.convertError("msg", &exception)
	suite.Equal("*errors.errorString", fmt.Sprintf("%T", err))
	suite.EqualError(err, "msg: <nil>")
}

func (suite *ErrorHandlingExtensionSuite) TestConvertError_genericJavaScriptError() {
	exception := suite.getGojaException("throw Error('jsError')")
	err := suite.extension.convertError("msg", exception)
	suite.Equal("*errors.errorString", fmt.Sprintf("%T", err))
	suite.EqualError(err, "msg: Error: jsError at <eval>:1:1(3)")
}

func (suite *ErrorHandlingExtensionSuite) TestConvertError_genericNewJavaScriptError() {
	exception := suite.getGojaException("throw new Error('jsError')")
	err := suite.extension.convertError("msg", exception)
	suite.Equal("*errors.errorString", fmt.Sprintf("%T", err))
	suite.EqualError(err, "msg: Error: jsError at <eval>:1:7(2)")
}

func (suite *ErrorHandlingExtensionSuite) TestConvertError_JavaScriptString() {
	exception := suite.getGojaException("throw 'jsError'")
	err := suite.extension.convertError("msg", exception)
	suite.Equal("*errors.errorString", fmt.Sprintf("%T", err))
	suite.EqualError(err, "msg: jsError at <eval>:1:1(1)")
}

func (suite *ErrorHandlingExtensionSuite) TestConvertError_JavaScriptErrorWithStatus() {
	exception := suite.getGojaException("const err = new Error('jsError'); err.status = 400; throw err")
	err := suite.extension.convertError("msg", exception)
	suite.Equal("*apiErrors.APIError", fmt.Sprintf("%T", err))
	suite.EqualError(err, "jsError")
	apiErr := err.(*apiErrors.APIError)
	suite.Equal(apiErrors.NewAPIError(400, "jsError"), apiErr)
}

func (suite *ErrorHandlingExtensionSuite) getGojaException(javaScript string) *goja.Exception {
	_, err := suite.extension.vm.RunString(javaScript)
	suite.Error(err)
	suite.Equal("*goja.Exception", fmt.Sprintf("%T", err))
	exception := err.(*goja.Exception)
	suite.NotNil(exception)
	return exception
}
