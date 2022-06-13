package extensionAPI

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
	"io/ioutil"
)

// GetExtensionFromFile loads an extension from a .js file.
func GetExtensionFromFile(fileName string) (*JsExtension, error) {
	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
	registry := new(require.Registry)
	registry.Enable(vm)
	console.Enable(vm)
	extensionJs, err := loadExtension(vm, fileName)
	if err != nil {
		return nil, err
	}
	if extensionJs.APIVersion != "0.1.0" {
		return nil, fmt.Errorf("incompatible extension API version %v. Please update the extension to use a supported version of the extension API", extensionJs.APIVersion)
	}
	return &extensionJs.Extension, nil
}

func loadExtension(vm *goja.Runtime, fileName string) (*installedExtension, error) {
	const extensionVariable = "installedExtension"
	err := vm.Set(extensionVariable, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to set installedExtension = null. Cause: %v", err.Error())
	}
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open extension file %v. Cause: %v", fileName, err.Error())
	}
	_, err = vm.RunScript(fileName, string(bytes))
	if err != nil {
		return nil, fmt.Errorf("failed to run extension file %v. Cause %v", fileName, err.Error())
	}
	var extension installedExtension
	err = vm.ExportTo(vm.Get(extensionVariable), &extension)
	if err != nil {
		return nil, fmt.Errorf("failed to read installedExtension variable. Cause: %v", err.Error())
	}
	return &extension, nil
}

type installedExtension struct {
	Extension  JsExtension `json:"extension"`
	APIVersion string      `json:"apiVersion"`
}

type JsExtension struct {
	Name                string                                                                                 `json:"name"`
	Description         string                                                                                 `json:"description"`
	BucketFsUploads     []BucketFsUpload                                                                       `json:"bucketFsUploads"`
	InstallableVersions []string                                                                               `json:"installableVersions"`
	Install             func(client SimpleSQLClient)                                                           `json:"install"`
	FindInstallations   func(sqlClient SimpleSQLClient, exaAllScripts *ExaAllScriptTable) []*JsExtInstallation `json:"findInstallations"`
}

type BucketFsUpload struct {
	Name             string `json:"name"`
	DownloadURL      string `json:"downloadUrl"`
	LicenseURL       string `json:"licenseUrl"`
	FileSize         int    `json:"fileSize"`
	BucketFsFilename string `json:"bucketFsFilename"`
}

type JsExtInstallation struct {
	Name string `json:"name"`
}

type SimpleSQLClient interface {
	RunQuery(query string)
}

type LoggingSimpleSQLClient struct {
}

func (client LoggingSimpleSQLClient) RunQuery(query string) {
	fmt.Printf("sql: %v\n", query)
}
