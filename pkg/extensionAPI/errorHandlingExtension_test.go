package extensionAPI

import (
	"fmt"
	"testing"

	"github.com/dop251/goja"
	"github.com/exasol/extension-manager/pkg/apiErrors"
	"github.com/exasol/extension-manager/pkg/backend"
	"github.com/exasol/extension-manager/pkg/extensionAPI/context"
	"github.com/exasol/extension-manager/pkg/extensionAPI/exaMetadata"
	"github.com/exasol/extension-manager/pkg/extensionController/transaction"
	"github.com/stretchr/testify/suite"
)

const mockErrorMessage = "mock error"

type ErrorHandlingExtensionSuite struct {
	suite.Suite
	rawExtension *rawJsExtension
	extension    *JsExtension
}

func TestJsExtensionSuite(t *testing.T) {
	suite.Run(t, new(ErrorHandlingExtensionSuite))
}

func (suite *ErrorHandlingExtensionSuite) SetupSuite() {
	suite.rawExtension = &rawJsExtension{
		Name:                    "name",
		Category:                "category",
		Description:             "desc",
		InstallableVersions:     []rawJsExtensionVersion{{Name: "v1", Deprecated: true, Latest: false}, {Name: "v2", Deprecated: false, Latest: true}},
		BucketFsUploads:         []BucketFsUpload{{Name: "uploadName", DownloadURL: "url", LicenseURL: "license", FileSize: 123, BucketFsFilename: "filename"}},
		GetParameterDefinitions: nil,
		Install:                 nil,
		Uninstall:               nil,
		Upgrade:                 nil,
		FindInstallations:       nil,
		AddInstance:             nil,
		FindInstances:           nil,
		DeleteInstance:          nil,
	}
	suite.extension = wrapExtension(suite.rawExtension, "id", newJavaScriptVm("logPrefix>"))
}

func (suite *ErrorHandlingExtensionSuite) TestProperties() {
	suite.Equal(&JsExtension{
		Id:                  "id",
		Category:            "category",
		Name:                "name",
		Description:         "desc",
		InstallableVersions: []JsExtensionVersion{{Name: "v1", Deprecated: true, Latest: false}, {Name: "v2", Deprecated: false, Latest: true}},
		BucketFsUploads:     []BucketFsUpload{{Name: "uploadName", DownloadURL: "url", LicenseURL: "license", FileSize: 123, BucketFsFilename: "filename"}},
		extension:           suite.rawExtension,
		vm:                  suite.extension.vm},
		suite.extension)
}

const EXTENSION_SCHEMA = "extension_schema"

func createMockContextWithClients(
	sqlClient backend.SimpleSQLClient,
	bucketFsContext context.BucketFsContext,
	metadataReader exaMetadata.ExaMetadataReader,
) *context.ExtensionContext {
	txCtx := &transaction.TransactionContext{}
	return context.CreateContextWithClient(EXTENSION_SCHEMA, txCtx, sqlClient, bucketFsContext, metadataReader)
}

func createMockContext() *context.ExtensionContext {
	var sqlClientMock backend.SimpleSQLClient = backend.CreateSimpleSqlClientMock()
	var bucketFsClientMock context.BucketFsContext = context.CreateBucketFsContextMock()
	var metadataReader exaMetadata.ExaMetadataReader = exaMetadata.CreateExaMetaDataReaderMock(EXTENSION_SCHEMA)
	return createMockContextWithClients(sqlClientMock, bucketFsClientMock, metadataReader)
}

// FindInstallations

func (suite *ErrorHandlingExtensionSuite) TestFindInstallationsSuccessful() {
	expectedInstallations := []*JsExtInstallation{{Name: "instName"}}
	suite.rawExtension.FindInstallations = func(context *context.ExtensionContext, metadata *exaMetadata.ExaMetadata) []*JsExtInstallation {
		return expectedInstallations
	}
	installations, err := suite.extension.FindInstallations(createMockContext(), createMetaData())
	suite.NoError(err)
	suite.Equal(expectedInstallations, installations)
}

func (suite *ErrorHandlingExtensionSuite) TestFindInstallationsFailure() {
	suite.rawExtension.FindInstallations = func(context *context.ExtensionContext, metadata *exaMetadata.ExaMetadata) []*JsExtInstallation {
		panic(mockErrorMessage)
	}
	installations, err := suite.extension.FindInstallations(createMockContext(), createMetaData())
	suite.EqualError(err, `failed to find installations for extension "id": `+mockErrorMessage)
	suite.Nil(installations)
}

func (suite *ErrorHandlingExtensionSuite) TestFindInstallationsUnsupported() {
	suite.rawExtension.FindInstallations = nil
	installations, err := suite.extension.FindInstallations(createMockContext(), createMetaData())
	suite.EqualError(err, `extension "id" does not support operation "findInstallations"`)
	suite.Nil(installations)
}

// GetParameterDefinitions

func (suite *ErrorHandlingExtensionSuite) GetParameterDefinitionsSuccessful() {
	expectedDefinitions := []interface{}{map[string]interface{}{"id": "param1", "name": "My param", "type": "string"}}
	suite.rawExtension.GetParameterDefinitions = func(context *context.ExtensionContext, version string) []interface{} {
		return expectedDefinitions
	}
	definitions, err := suite.extension.GetParameterDefinitions(createMockContext(), "ext-version")
	suite.NoError(err)
	suite.Equal(expectedDefinitions, definitions)
}

func (suite *ErrorHandlingExtensionSuite) GetParameterDefinitionsFailure() {
	suite.rawExtension.GetParameterDefinitions = func(context *context.ExtensionContext, version string) []interface{} {
		panic(mockErrorMessage)
	}
	installations, err := suite.extension.GetParameterDefinitions(createMockContext(), "ext-version")
	suite.EqualError(err, "failed to get parameter definitions for extension \"id\": "+mockErrorMessage)
	suite.Nil(installations)
}

func (suite *ErrorHandlingExtensionSuite) GetParameterDefinitionsUnsupported() {
	suite.rawExtension.GetParameterDefinitions = nil
	installations, err := suite.extension.GetParameterDefinitions(createMockContext(), "ext-version")
	suite.EqualError(err, `extension "id" does not support operation "getInstanceParameters"`)
	suite.Nil(installations)
}

// Install

func (suite *ErrorHandlingExtensionSuite) TestInstallSuccessful() {
	suite.rawExtension.Install = func(context *context.ExtensionContext, version string) {
		// empty mocked function
	}
	err := suite.extension.Install(createMockContext(), "version")
	suite.NoError(err)
}

func (suite *ErrorHandlingExtensionSuite) TestInstallFailure() {
	suite.rawExtension.Install = func(context *context.ExtensionContext, version string) {
		panic(mockErrorMessage)
	}
	err := suite.extension.Install(createMockContext(), "version")
	suite.EqualError(err, `failed to install extension "id": `+mockErrorMessage)
}

func (suite *ErrorHandlingExtensionSuite) TestInstallUnsupported() {
	suite.rawExtension.Install = nil
	err := suite.extension.Install(createMockContext(), "version")
	suite.EqualError(err, `extension "id" does not support operation "install"`)
}

// Uninstall

func (suite *ErrorHandlingExtensionSuite) TestUninstallSuccessful() {
	suite.rawExtension.Uninstall = func(context *context.ExtensionContext, version string) {
		// empty mocked function
	}
	err := suite.extension.Uninstall(createMockContext(), "version")
	suite.NoError(err)
}

func (suite *ErrorHandlingExtensionSuite) TestUninstallFailure() {
	suite.rawExtension.Uninstall = func(context *context.ExtensionContext, version string) {
		panic(mockErrorMessage)
	}
	err := suite.extension.Uninstall(createMockContext(), "version")
	suite.EqualError(err, `failed to uninstall extension "id": `+mockErrorMessage)
}

func (suite *ErrorHandlingExtensionSuite) TestUninstallUnsupported() {
	suite.rawExtension.Uninstall = nil
	err := suite.extension.Uninstall(createMockContext(), "version")
	suite.EqualError(err, `extension "id" does not support operation "uninstall"`)
}

// Upgrade

func (suite *ErrorHandlingExtensionSuite) TestUpgradeSuccessful() {
	suite.rawExtension.Upgrade = func(context *context.ExtensionContext) *JsUpgradeResult {
		return &JsUpgradeResult{PreviousVersion: "old", NewVersion: "new"}
	}
	result, err := suite.extension.Upgrade(createMockContext())
	suite.NoError(err)
	suite.Equal(&JsUpgradeResult{PreviousVersion: "old", NewVersion: "new"}, result)
}

func (suite *ErrorHandlingExtensionSuite) TestUpgradeFails() {
	suite.rawExtension.Upgrade = func(context *context.ExtensionContext) *JsUpgradeResult {
		panic(mockErrorMessage)
	}
	result, err := suite.extension.Upgrade(createMockContext())
	suite.EqualError(err, `failed to upgrade extension "id": `+mockErrorMessage)
	suite.Nil(result)
}

func (suite *ErrorHandlingExtensionSuite) TestUpgradeUnsupported() {
	suite.rawExtension.Upgrade = nil
	instance, err := suite.extension.Upgrade(createMockContext())
	suite.EqualError(err, `extension "id" does not support operation "upgrade"`)
	suite.Nil(instance)
}

// AddInstance

func (suite *ErrorHandlingExtensionSuite) TestAddInstanceSuccessful() {
	suite.rawExtension.AddInstance = func(context *context.ExtensionContext, version string, params *ParameterValues) *JsExtInstance {
		return &JsExtInstance{Id: "inst", Name: "newInstance"}
	}
	instance, err := suite.extension.AddInstance(createMockContext(), "version", &ParameterValues{Values: []ParameterValue{}})
	suite.NoError(err)
	suite.Equal(&JsExtInstance{Id: "inst", Name: "newInstance"}, instance)
}

func (suite *ErrorHandlingExtensionSuite) TestAddInstanceFails() {
	suite.rawExtension.AddInstance = func(context *context.ExtensionContext, version string, params *ParameterValues) *JsExtInstance {
		panic(mockErrorMessage)
	}
	instance, err := suite.extension.AddInstance(createMockContext(), "version", &ParameterValues{Values: []ParameterValue{}})
	suite.EqualError(err, `failed to add instance for extension "id": `+mockErrorMessage)
	suite.Nil(instance)
}

func (suite *ErrorHandlingExtensionSuite) TestAddInstanceUnsupported() {
	suite.rawExtension.AddInstance = nil
	instance, err := suite.extension.AddInstance(createMockContext(), "version", &ParameterValues{Values: []ParameterValue{}})
	suite.EqualError(err, `extension "id" does not support operation "addInstance"`)
	suite.Nil(instance)
}

// DeleteInstance

func (suite *ErrorHandlingExtensionSuite) TestDeleteInstanceSuccessful() {
	suite.rawExtension.DeleteInstance = func(context *context.ExtensionContext, version, instanceId string) {
		// empty mocked function
	}
	err := suite.extension.DeleteInstance(createMockContext(), "version", "instance-id")
	suite.NoError(err)
}

func (suite *ErrorHandlingExtensionSuite) TestDeleteInstanceFails() {
	suite.rawExtension.DeleteInstance = func(context *context.ExtensionContext, version, instanceId string) {
		panic(mockErrorMessage)
	}
	err := suite.extension.DeleteInstance(createMockContext(), "version", "instance-id")
	suite.EqualError(err, `failed to delete instance "instance-id" for extension "id": `+mockErrorMessage)
}

func (suite *ErrorHandlingExtensionSuite) TestDeleteInstanceUnsupported() {
	suite.rawExtension.DeleteInstance = nil
	err := suite.extension.DeleteInstance(createMockContext(), "version", "instance-id")
	suite.EqualError(err, `extension "id" does not support operation "deleteInstance"`)
}

// convertError

func (suite *ErrorHandlingExtensionSuite) TestConvertErrorNonErrorObject() {
	err := suite.extension.convertError("msg", mockErrorMessage)
	suite.assertErrorStringError(err, "msg: "+mockErrorMessage)
}

func (suite *ErrorHandlingExtensionSuite) TestConvertErrorErrorObject() {
	err := suite.extension.convertError("msg", fmt.Errorf(mockErrorMessage))
	suite.assertErrorStringError(err, "msg: "+mockErrorMessage)
}

func (suite *ErrorHandlingExtensionSuite) TestConvertErrorNilGojaException() {
	var exception goja.Exception
	err := suite.extension.convertError("msg", &exception)
	suite.assertErrorStringError(err, "msg: <nil>")
}

func (suite *ErrorHandlingExtensionSuite) TestConvertErrorGenericJavaScriptError() {
	exception := suite.getGojaException("throw Error('jsError')")
	err := suite.extension.convertError("msg", exception)
	suite.assertErrorStringError(err, "msg: Error: jsError at Error (native)")
}

func (suite *ErrorHandlingExtensionSuite) TestConvertErrorGenericNewJavaScriptError() {
	exception := suite.getGojaException("throw new Error('jsError')")
	err := suite.extension.convertError("msg", exception)
	suite.assertErrorStringError(err, "msg: Error: jsError at <eval>:1:7(2)")
}

func (suite *ErrorHandlingExtensionSuite) TestConvertErrorJavaScriptString() {
	exception := suite.getGojaException("throw 'jsError'")
	err := suite.extension.convertError("msg", exception)
	suite.assertErrorStringError(err, "msg: jsError at <eval>:1:1(1)")
}

func (suite *ErrorHandlingExtensionSuite) TestConvertErrorJavaScriptErrorWithStatus() {
	exception := suite.getGojaException("const err = new Error('jsError'); err.status = 400; throw err")
	err := suite.extension.convertError("msg", exception)
	suite.Equal("*apiErrors.APIError", fmt.Sprintf("%T", err))
	suite.EqualError(err, "jsError")
	apiErr, ok := apiErrors.AsAPIError(err)
	suite.True(ok)
	suite.Equal(apiErrors.NewAPIError(400, "jsError"), apiErr)
}

func (suite *ErrorHandlingExtensionSuite) getGojaException(javaScript string) *goja.Exception {
	_, err := suite.extension.vm.RunString(javaScript)
	suite.Error(err)
	suite.Equal("*goja.Exception", fmt.Sprintf("%T", err))
	//nolint:errorlint // Type assertion is OK here because the error is not wrapped
	exception := err.(*goja.Exception)
	suite.NotNil(exception)
	return exception
}

func createMetaData() *exaMetadata.ExaMetadata {
	//nolint:exhaustruct // Not necessary for test data
	return &exaMetadata.ExaMetadata{}
}

func (suite *ErrorHandlingExtensionSuite) assertErrorStringError(err error, expectedMessage string) {
	suite.Equal("*errors.errorString", fmt.Sprintf("%T", err))
	suite.EqualError(err, expectedMessage)
}
