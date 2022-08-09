package parameterValidator

import (
	_ "embed" // Required to embed the validator JS code
	"fmt"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
	"github.com/exasol/extension-manager/extensionAPI"
)

// The dist.js file is built using go:generate in ../go-generate.go
//
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

// ValidateParameters validates parameter values against the parameter definition and returns a list of failed validations.
// If all parameters are valid, this returns an empty slice.
func (v *Validator) ValidateParameters(definitions []interface{}, params extensionAPI.ParameterValues) (failedValidations []ValidationResult, err error) {
	result := make([]ValidationResult, 0)
	for _, def := range definitions {
		name, r, err := v.validateParameter(def, params)
		if err != nil {
			return nil, err
		}
		if !r.Success {
			result = append(result, ValidationResult{Success: false, Message: fmt.Sprintf("Failed to validate parameter %q: %s", name, r.Message)})
		}
	}
	return result, nil
}

func (v *Validator) validateParameter(def interface{}, params extensionAPI.ParameterValues) (string, *ValidationResult, error) {
	id, name, err := extractFromDefinition(def)
	if err != nil {
		return "", nil, err
	}
	paramValue := findParamValue(params, id)
	result, err := v.ValidateParameter(def, paramValue)
	if err != nil {
		return "", nil, fmt.Errorf("failed to validate parameter value %q with id %q using definition %v", paramValue, id, def)
	}
	return name, result, nil
}

func findParamValue(params extensionAPI.ParameterValues, id string) string {
	if param, found := params.Find(id); found {
		return param.Value
	}
	return ""
}

func extractFromDefinition(d interface{}) (id string, name string, err error) {
	if def, ok := d.(map[string]interface{}); ok {
		return extractIdAndName(def)
	} else {
		return "", "", fmt.Errorf("unexpected type of definition: %t", d)
	}
}

func extractIdAndName(def map[string]interface{}) (id string, name string, err error) {
	if id, ok := def["id"].(string); ok {
		name, err := extractName(def, id)
		return id, name, err
	} else {
		return "", "", fmt.Errorf("unexpected type of id in parameter definition: %t", def["id"])
	}
}

func extractName(def map[string]interface{}, id string) (name string, err error) {
	if name, ok := def["name"].(string); ok {
		return name, nil
	} else {
		return "", fmt.Errorf("unexpected type of name in parameter definition: %t", def["name"])
	}
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
