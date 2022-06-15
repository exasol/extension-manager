package parameterValidator

import (
	_ "embed"
	"fmt"
	"github.com/dop251/goja"
)

// The dist.js file is built using go:generate in ../go-generate.go
//go:embed dist.js
var dependencyValidatorJs string

var validateParameterJs func(definition interface{}, value string) ValidationResult

type ValidationResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func validateParameter(definition interface{}, value string) (*ValidationResult, error) {
	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
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
	err = vm.ExportTo(function, &validateParameterJs)
	if err != nil {
		return nil, err
	}
	result := validateParameterJs(definition, value)
	return &result, nil
}
