package extensionAPI

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
)

const SupportedApiVersion = "0.1.8"

// GetExtensionFromFile loads an extension from a .js file.
func GetExtensionFromFile(extensionPath string) (*JsExtension, error) {
	vm := newJavaScriptVm()
	extensionJs, err := loadExtension(vm, extensionPath)
	if err != nil {
		return nil, err
	}
	if extensionJs.APIVersion != SupportedApiVersion {
		return nil, fmt.Errorf("incompatible extension API version %q. Please update the extension to use supported version %q", extensionJs.APIVersion, SupportedApiVersion)
	}
	wrappedExtension := wrapExtension(&extensionJs.Extension)
	_, fileName := path.Split(extensionPath)
	wrappedExtension.Id = fileName
	log.Printf("Extension %q with id %q loaded from file %q", wrappedExtension.Name, wrappedExtension.Id, extensionPath)
	return wrappedExtension, nil
}

func newJavaScriptVm() *goja.Runtime {
	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
	registry := new(require.Registry)
	registry.Enable(vm)
	console.Enable(vm)
	return vm
}

func loadExtension(vm *goja.Runtime, fileName string) (*installedExtension, error) {
	globalJsObj := vm.NewObject()
	err := vm.Set("global", globalJsObj)
	if err != nil {
		return nil, fmt.Errorf("failed to set global to a new object. Cause: %w", err)
	}
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open extension file %v. Cause: %w", fileName, err)
	}
	_, err = vm.RunScript(fileName, string(bytes))
	if err != nil {
		return nil, fmt.Errorf("failed to run extension file %v. Cause %w", fileName, err)
	}

	const extensionVariableName = "installedExtension"
	extensionVariable := globalJsObj.Get(extensionVariableName)
	if extensionVariable == nil {
		return nil, fmt.Errorf("extension did not set global.%s", extensionVariableName)
	}
	var extension installedExtension
	err = vm.ExportTo(extensionVariable, &extension)
	if err != nil {
		return nil, fmt.Errorf("failed to read installedExtension variable. Cause: %w", err)
	}
	return &extension, nil
}

type installedExtension struct {
	Extension  rawJsExtension `json:"extension"`
	APIVersion string         `json:"apiVersion"`
}

type rawJsExtension struct {
	Id                  string
	Name                string                                                                                  `json:"name"`
	Description         string                                                                                  `json:"description"`
	BucketFsUploads     []BucketFsUpload                                                                        `json:"bucketFsUploads"`
	InstallableVersions []string                                                                                `json:"installableVersions"`
	Install             func(context *ExtensionContext, version string)                                         `json:"install"`
	FindInstallations   func(context *ExtensionContext, metadata *ExaMetadata) []*JsExtInstallation             `json:"findInstallations"`
	AddInstance         func(context *ExtensionContext, version string, params *ParameterValues) *JsExtInstance `json:"addInstance"`
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
	// InstanceParameters is deserialized to a structure of []interface{} and maps.
	InstanceParameters []interface{} `json:"instanceParameters"`
}

type JsExtInstance struct {
	Name string `json:"name"`
}

type ParameterValues struct {
	Values []ParameterValue `json:"values"`
}

type ParameterValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type SimpleSQLClient interface {
	RunQuery(query string)
}

type LoggingSimpleSQLClient struct {
}

func (client LoggingSimpleSQLClient) RunQuery(query string) {
	fmt.Printf("sql: %v\n", query)
}
