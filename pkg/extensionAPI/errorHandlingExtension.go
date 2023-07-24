package extensionAPI

import (
	"fmt"

	"github.com/dop251/goja"
	"github.com/exasol/extension-manager/pkg/apiErrors"
	"github.com/exasol/extension-manager/pkg/extensionAPI/context"
	"github.com/exasol/extension-manager/pkg/extensionAPI/exaMetadata"
)

type JsExtension struct {
	extension           *rawJsExtension
	vm                  *goja.Runtime
	Id                  string
	Name                string
	Category            string
	Description         string
	InstallableVersions []JsExtensionVersion
	BucketFsUploads     []BucketFsUpload
}

type JsExtensionVersion struct {
	Name       string
	Latest     bool
	Deprecated bool
}

func wrapExtension(ext *rawJsExtension, id string, vm *goja.Runtime) *JsExtension {
	return &JsExtension{
		extension:           ext,
		Id:                  id,
		vm:                  vm,
		Name:                ext.Name,
		Category:            ext.Category,
		Description:         ext.Description,
		InstallableVersions: convertVersions(ext.InstallableVersions),
		BucketFsUploads:     ext.BucketFsUploads,
	}
}

func convertVersions(versions []rawJsExtensionVersion) []JsExtensionVersion {
	result := make([]JsExtensionVersion, 0, len(versions))
	for _, v := range versions {
		result = append(result, JsExtensionVersion(v))

	}
	return result
}

func (e *JsExtension) GetParameterDefinitions(context *context.ExtensionContext, version string) (definitions []interface{}, errorResult error) {
	defer func() {
		if err := recover(); err != nil {
			errorResult = e.convertError(fmt.Sprintf("failed to get parameter definitions for extension %q", e.Id), err)
		}
	}()
	return e.extension.GetParameterDefinitions(context, version), nil
}

func (e *JsExtension) Install(context *context.ExtensionContext, version string) (errorResult error) {
	if e.extension.Install == nil {
		return e.unsupportedFunction("install")
	}
	defer func() {
		if err := recover(); err != nil {
			errorResult = e.convertError(fmt.Sprintf("failed to install extension %q", e.Id), err)
		}
	}()
	e.extension.Install(context, version)
	return nil
}

func (e *JsExtension) Uninstall(context *context.ExtensionContext, version string) (errorResult error) {
	if e.extension.Uninstall == nil {
		return e.unsupportedFunction("uninstall")
	}
	defer func() {
		if err := recover(); err != nil {
			errorResult = e.convertError(fmt.Sprintf("failed to uninstall extension %q", e.Id), err)
		}
	}()
	e.extension.Uninstall(context, version)
	return nil
}

func (e *JsExtension) FindInstallations(context *context.ExtensionContext, metadata *exaMetadata.ExaMetadata) (installations []*JsExtInstallation, errorResult error) {
	if e.extension.FindInstallations == nil {
		return nil, e.unsupportedFunction("findInstallations")
	}
	defer func() {
		if err := recover(); err != nil {
			errorResult = e.convertError(fmt.Sprintf("failed to find installations for extension %q", e.Id), err)
		}
	}()
	return e.extension.FindInstallations(context, metadata), nil
}

func (e *JsExtension) AddInstance(context *context.ExtensionContext, version string, params *ParameterValues) (instance *JsExtInstance, errorResult error) {
	if e.extension.AddInstance == nil {
		return nil, e.unsupportedFunction("addInstance")
	}
	defer func() {
		if err := recover(); err != nil {
			errorResult = e.convertError(fmt.Sprintf("failed to add instance for extension %q", e.Id), err)
		}
	}()
	return e.extension.AddInstance(context, version, params), nil
}

func (e *JsExtension) ListInstances(context *context.ExtensionContext, version string) (instances []*JsExtInstance, errorResult error) {
	if e.extension.FindInstances == nil {
		return nil, e.unsupportedFunction("findInstances")
	}
	defer func() {
		if err := recover(); err != nil {
			errorResult = e.convertError(fmt.Sprintf("failed to list instances for extension %q in version %q", e.Id, version), err)
		}
	}()
	return e.extension.FindInstances(context, version), nil
}

func (e *JsExtension) DeleteInstance(context *context.ExtensionContext, extensionVersion, instanceId string) (errorResult error) {
	if e.extension.DeleteInstance == nil {
		return e.unsupportedFunction("deleteInstance")
	}
	defer func() {
		if err := recover(); err != nil {
			errorResult = e.convertError(fmt.Sprintf("failed to delete instance %q for extension %q", instanceId, e.Id), err)
		}
	}()
	e.extension.DeleteInstance(context, extensionVersion, instanceId)
	return
}

func (e *JsExtension) convertError(message string, err any) error {
	if exception, ok := err.(*goja.Exception); ok {
		if exception.Value() == nil {
			return basicError(message, err)
		}
		statusField := exception.Value().ToObject(e.vm).Get("status")
		if statusField == nil {
			return basicError(message, err)
		}
		var apiError jsApiError
		exportErr := e.vm.ExportTo(exception.Value(), &apiError)
		if exportErr != nil {
			return fmt.Errorf("failed to convert error %v of type %T (message: %q) to ApiError: %w", err, err, message, exportErr)
		}
		return apiErrors.NewAPIError(apiError.Status, apiError.Message)
	}
	return basicError(message, err)
}

func basicError(message string, err any) error {
	return fmt.Errorf("%s: %v", message, err)
}

func (e *JsExtension) unsupportedFunction(functionName string) error {
	return fmt.Errorf("extension %q does not support operation %q", e.Id, functionName)
}

type jsApiError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
