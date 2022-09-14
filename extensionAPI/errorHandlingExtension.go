package extensionAPI

import (
	"fmt"

	"github.com/dop251/goja"
	"github.com/exasol/extension-manager/apiErrors"
)

type JsExtension struct {
	extension           *rawJsExtension
	vm                  *goja.Runtime
	Id                  string
	Name                string
	Description         string
	InstallableVersions []string
	BucketFsUploads     []BucketFsUpload
}

func wrapExtension(ext *rawJsExtension, id string, vm *goja.Runtime) *JsExtension {
	return &JsExtension{
		extension:           ext,
		Id:                  id,
		vm:                  vm,
		Name:                ext.Name,
		Description:         ext.Description,
		InstallableVersions: ext.InstallableVersions,
		BucketFsUploads:     ext.BucketFsUploads,
	}
}

func (e *JsExtension) Install(context *ExtensionContext, version string) (errorResult error) {
	defer func() {
		if err := recover(); err != nil {
			errorResult = e.convertError(fmt.Sprintf("failed to install extension %q", e.Id), err)
		}
	}()
	e.extension.Install(context, version)
	return nil
}

func (e *JsExtension) Uninstall(context *ExtensionContext, version string) (errorResult error) {
	defer func() {
		if err := recover(); err != nil {
			errorResult = e.convertError(fmt.Sprintf("failed to uninstall extension %q", e.Id), err)
		}
	}()
	e.extension.Uninstall(context, version)
	return nil
}

func (e *JsExtension) FindInstallations(context *ExtensionContext, metadata *ExaMetadata) (installations []*JsExtInstallation, errorResult error) {
	defer func() {
		if err := recover(); err != nil {
			errorResult = e.convertError(fmt.Sprintf("failed to find installations for extension %q", e.Id), err)
		}
	}()
	return e.extension.FindInstallations(context, metadata), nil
}

func (e *JsExtension) AddInstance(context *ExtensionContext, version string, params *ParameterValues) (instance *JsExtInstance, errorResult error) {
	defer func() {
		if err := recover(); err != nil {
			errorResult = e.convertError(fmt.Sprintf("failed to add instance for extension %q", e.Id), err)
		}
	}()
	return e.extension.AddInstance(context, version, params), nil
}

func (e *JsExtension) ListInstances(context *ExtensionContext, version string) (instances []*JsExtInstance, errorResult error) {
	defer func() {
		if err := recover(); err != nil {
			errorResult = e.convertError(fmt.Sprintf("failed to list instances for extension %q in version %q", e.Id, version), err)
		}
	}()
	return e.extension.FindInstances(context, version), nil
}

func (e *JsExtension) DeleteInstance(context *ExtensionContext, extensionVersion, instanceId string) (errorResult error) {
	defer func() {
		if err := recover(); err != nil {
			errorResult = e.convertError(fmt.Sprintf("failed to delete instance %q for extension %q", instanceId, e.Id), err)
		}
	}()
	e.extension.DeleteInstance(context, instanceId)
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

type jsApiError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
