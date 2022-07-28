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

type Validator struct {
	validate func(definition interface{}, value string) ValidationResult
}

// New creates a new reusable validator
func New() (*Validator, error) {
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

	validator := Validator{}
	err = vm.ExportTo(function, &validator.validate)
	if err != nil {
		return nil, err
	}
	return &validator, nil
}

// ValidateParameters uses the given parameter definition to validate a single value
func (v *Validator) ValidateParameter(definition interface{}, value string) (validationResult *ValidationResult, errorResult error) {
	defer func() {
		if err := recover(); err != nil {
			errorResult = fmt.Errorf("failed to validate parameter value %q using definition %v: %v", value, definition, err)
		}
	}()
	result := v.validate(definition, value)
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
