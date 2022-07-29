package extensionAPI

import (
	"fmt"
)

type JsExtension struct {
	extension           *rawJsExtension
	Id                  string
	Name                string
	Description         string
	InstallableVersions []string
	BucketFsUploads     []BucketFsUpload
}

func wrapExtension(ext *rawJsExtension) *JsExtension {
	return &JsExtension{
		extension:           ext,
		Id:                  ext.Id,
		Name:                ext.Name,
		Description:         ext.Description,
		InstallableVersions: ext.InstallableVersions,
		BucketFsUploads:     ext.BucketFsUploads}
}

func (e *JsExtension) Install(context *ExtensionContext, version string) (errorResult error) {
	defer func() {
		if err := recover(); err != nil {
			errorResult = fmt.Errorf("failed to install extension %q: %v", e.Id, err)
		}
	}()
	e.extension.Install(context, version)
	return nil
}

func (e *JsExtension) FindInstallations(context *ExtensionContext, metadata *ExaMetadata) (installations []*JsExtInstallation, errorResult error) {
	defer func() {
		if err := recover(); err != nil {
			errorResult = fmt.Errorf("failed to find installations for extension %q: %v", e.Id, err)
		}
	}()
	return e.extension.FindInstallations(context, metadata), nil
}

func (e *JsExtension) AddInstance(context *ExtensionContext, version string, params *ParameterValues) (instance *JsExtInstance, errorResult error) {
	defer func() {
		if err := recover(); err != nil {
			errorResult = fmt.Errorf("failed to add instance for extension %q: %v", e.Id, err)
		}
	}()
	return e.extension.AddInstance(context, version, params), nil
}
