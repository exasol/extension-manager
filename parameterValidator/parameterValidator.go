package parameterValidator

import (
	_ "embed"
	"fmt"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
)

// The dist.js file is built using go:generate in ../go-generate.go
//go:embed dist.js
var dependencyValidatorJs string

type ValidationResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ValidateParameters uses the given parameter definition
func ValidateParameters(definition interface{}, value string) (*ValidationResult, error) {
	vm := newJavaScriptVm()
	globalJsObj := vm.NewObject()
	err := vm.Set("global", globalJsObj)
	if err != nil {
		return nil, err
	}
	_, err = vm.RunString(dependencyValidatorJs)
	if err != nil {
		return nil, fmt.Errorf("failed to load validateParameter script. Cause: %w", err)
	}
	function := globalJsObj.Get("validateParameter")
	var validateParameterJs func(definition interface{}, value string) ValidationResult
	err = vm.ExportTo(function, &validateParameterJs)
	if err != nil {
		return nil, err
	}
	result := validateParameterJs(definition, value)
	return &result, nil
}

func newJavaScriptVm() *goja.Runtime {
	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
	registry := new(require.Registry)
	registry.Enable(vm)
	console.Enable(vm)
	return vm
}
