package extensionAPI

import (
	"fmt"

	"github.com/exasol/extension-manager/pkg/extensionAPI/context"
	"github.com/exasol/extension-manager/pkg/extensionAPI/exaMetadata"
	log "github.com/sirupsen/logrus"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
)

// LoadExtension loads an extension from the given file content.
/* [impl -> dsn~extension-definition~1]. */
func LoadExtension(id, content string) (*JsExtension, error) {
	logPrefix := fmt.Sprintf("JS:%s>", id)
	vm := newJavaScriptVm(logPrefix)
	extensionJs, err := loadExtension(vm, id, content)
	if err != nil {
		return nil, err
	}
	err = validateExtensionIsCompatibleWithApiVersion(id, extensionJs.APIVersion)
	if err != nil {
		return nil, err
	}
	wrappedExtension := wrapExtension(&extensionJs.Extension, id, vm)
	log.Debugf("Extension %q with id %q using API version %q loaded successfully", wrappedExtension.Name, wrappedExtension.Id, extensionJs.APIVersion)
	return wrappedExtension, nil
}

func newJavaScriptVm(logPrefix string) *goja.Runtime {
	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
	registry := new(require.Registry)
	registry.Enable(vm)
	configureLogging(registry, vm, logPrefix)
	return vm
}

func configureLogging(registry *require.Registry, vm *goja.Runtime, logPrefix string) {
	printer := createJavaScriptLogger(logPrefix)
	registry.RegisterNativeModule(console.ModuleName, console.RequireWithPrinter(printer))
	console.Enable(vm)
}

func loadExtension(vm *goja.Runtime, id, content string) (*installedExtension, error) {
	globalJsObj := vm.NewObject()
	err := vm.Set("global", globalJsObj)
	if err != nil {
		return nil, fmt.Errorf("failed to set global to a new object. Cause: %w", err)
	}
	_, err = vm.RunScript(id, content)
	if err != nil {
		return nil, fmt.Errorf("failed to run extension %q with content %q: %w", id, content, err)
	}

	const extensionVariableName = "installedExtension"
	extensionVariable := globalJsObj.Get(extensionVariableName)
	if extensionVariable == nil {
		return nil, fmt.Errorf("extension %q did not set global.%s", id, extensionVariableName)
	}
	var extension installedExtension
	err = vm.ExportTo(extensionVariable, &extension)
	if err != nil {
		return nil, fmt.Errorf("failed to read installedExtension variable for extension %q. Cause: %w", id, err)
	}
	return &extension, nil
}

// installedExtension allows deserializing extension definitions that implement the extension-manager-interface (https://github.com/exasol/extension-manager-interface/).
/* [impl -> dsn~extension-api~1]. */
type installedExtension struct {
	Extension  rawJsExtension `json:"extension"`
	APIVersion string         `json:"apiVersion"`
}

type rawJsExtension struct {
	Name                string                  `json:"name"`
	Category            string                  `json:"category"`
	Description         string                  `json:"description"`
	BucketFsUploads     []BucketFsUpload        `json:"bucketFsUploads"`
	InstallableVersions []rawJsExtensionVersion `json:"installableVersions"`
	// [impl -> dsn~parameter-versioning~1]
	// [impl -> dsn~configuration-parameters~1]
	GetParameterDefinitions func(context *context.ExtensionContext, version string) []interface{}                           `json:"getInstanceParameters"`
	Install                 func(context *context.ExtensionContext, version string)                                         `json:"install"`
	Uninstall               func(context *context.ExtensionContext, version string)                                         `json:"uninstall"`
	Upgrade                 func(context *context.ExtensionContext) *JsUpgradeResult                                        `json:"upgrade"`
	FindInstallations       func(context *context.ExtensionContext, metadata *exaMetadata.ExaMetadata) []*JsExtInstallation `json:"findInstallations"`
	AddInstance             func(context *context.ExtensionContext, version string, params *ParameterValues) *JsExtInstance `json:"addInstance"`
	FindInstances           func(context *context.ExtensionContext, version string) []*JsExtInstance                        `json:"findInstances"`
	DeleteInstance          func(context *context.ExtensionContext, version, instanceId string)                             `json:"deleteInstance"`
}

type rawJsExtensionVersion struct {
	Name       string `json:"name"`
	Latest     bool   `json:"latest"`
	Deprecated bool   `json:"deprecated"`
}

type BucketFsUpload struct {
	Name             string `json:"name"`             // Human-readable name or short description of the file
	DownloadURL      string `json:"downloadUrl"`      // Optional
	LicenseURL       string `json:"licenseUrl"`       // Optional
	FileSize         int    `json:"fileSize"`         // File size in bytes. Negative if EM should ignore the file size
	BucketFsFilename string `json:"bucketFsFilename"` // File name in BucketFS
}

type JsExtInstallation struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

type JsUpgradeResult struct {
	PreviousVersion string `json:"previousVersion"`
	NewVersion      string `json:"newVersion"`
}

type JsExtInstance struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ParameterValues struct {
	Values []ParameterValue `json:"values"`
}

// Find returns the parameter with the given ID and true if the parameter exists
// or an empty parameter and false in case none exists.
func (pv ParameterValues) Find(id string) (value ParameterValue, found bool) {
	for _, v := range pv.Values {
		if v.Name == id {
			return v, true
		}
	}
	return ParameterValue{Name: "", Value: ""}, false
}

type ParameterValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
