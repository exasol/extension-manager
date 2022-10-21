package extensionAPI

import (
	"fmt"

	"github.com/exasol/extension-manager/pkg/backend"
	log "github.com/sirupsen/logrus"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
)

const SupportedApiVersion = "0.1.15"

// LoadExtension loads an extension from the given file content.
func LoadExtension(id, content string) (*JsExtension, error) {
	logPrefix := fmt.Sprintf("JS:%s>", id)
	vm := newJavaScriptVm(logPrefix)
	extensionJs, err := loadExtension(vm, id, content)
	if err != nil {
		return nil, err
	}
	if extensionJs.APIVersion != SupportedApiVersion {
		return nil, fmt.Errorf("incompatible extension API version %q. Please update the extension to use supported version %q", extensionJs.APIVersion, SupportedApiVersion)
	}
	wrappedExtension := wrapExtension(&extensionJs.Extension, id, vm)
	log.Debugf("Extension %q with id %q loaded", wrappedExtension.Name, wrappedExtension.Id)
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
	var printer console.Printer = console.PrinterFunc(func(s string) { log.Print(logPrefix + s) })
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

type installedExtension struct {
	Extension  rawJsExtension `json:"extension"`
	APIVersion string         `json:"apiVersion"`
}

type rawJsExtension struct {
	Name                    string                                                                                  `json:"name"`
	Description             string                                                                                  `json:"description"`
	BucketFsUploads         []BucketFsUpload                                                                        `json:"bucketFsUploads"`
	InstallableVersions     []rawJsExtensionVersion                                                                 `json:"installableVersions"`
	GetParameterDefinitions func(context *ExtensionContext, version string) []interface{}                           `json:"getInstanceParameters"`
	Install                 func(context *ExtensionContext, version string)                                         `json:"install"`
	Uninstall               func(context *ExtensionContext, version string)                                         `json:"uninstall"`
	FindInstallations       func(context *ExtensionContext, metadata *ExaMetadata) []*JsExtInstallation             `json:"findInstallations"`
	AddInstance             func(context *ExtensionContext, version string, params *ParameterValues) *JsExtInstance `json:"addInstance"`
	FindInstances           func(context *ExtensionContext, version string) []*JsExtInstance                        `json:"findInstances"`
	DeleteInstance          func(context *ExtensionContext, version, instanceId string)                             `json:"deleteInstance"`
}

type rawJsExtensionVersion struct {
	Name       string `json:"name"`
	Latest     bool   `json:"latest"`
	Deprecated bool   `json:"deprecated"`
}

type BucketFsUpload struct {
	Name             string `json:"name"`
	DownloadURL      string `json:"downloadUrl"`
	LicenseURL       string `json:"licenseUrl"`
	FileSize         int    `json:"fileSize"`
	BucketFsFilename string `json:"bucketFsFilename"`
}

type JsExtInstallation struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type JsExtInstance struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ParameterValues struct {
	Values []ParameterValue `json:"values"`
}

// Find returns the parameter with the given ID or nil in case none exists.
func (pv ParameterValues) Find(id string) (value ParameterValue, found bool) {
	for _, v := range pv.Values {
		if v.Name == id {
			return v, true
		}
	}
	return ParameterValue{}, false
}

type ParameterValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Extensions use this SQL client to execute queries.
type SimpleSQLClient interface {
	// Execute runs a query that does not return rows, e.g. INSERT or UPDATE.
	Execute(query string, args ...any)

	// Query runs a query that returns rows, typically a SELECT.
	Query(query string, args ...any) backend.QueryResult
}
