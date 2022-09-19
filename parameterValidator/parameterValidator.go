package parameterValidator

import (
	_ "embed" // Required to embed the validator JS code
	"fmt"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
	"github.com/exasol/extension-manager/extensionAPI"
)

//go:generate npm ci
//go:generate npm run build
//go:embed parameterValidator.js
var dependencyValidatorJs string

type ParameterDefinition struct {
	Id            string
	Name          string
	RawDefinition map[string]interface{}
}

func ConvertDefinitions(rawDefinitions []interface{}) ([]ParameterDefinition, error) {
	definitions := make([]ParameterDefinition, 0, len(rawDefinitions))
	for _, d := range rawDefinitions {
		if def, ok := d.(map[string]interface{}); ok {
			convertedDef, err := convertDefinition(def)
			if err != nil {
				return nil, err
			}
			definitions = append(definitions, convertedDef)
		} else {
			return nil, fmt.Errorf("unexpected type %T of definition: %v, expected map[string]interface{}", d, d)
		}
	}
	return definitions, nil
}

func convertDefinition(rawDefinition map[string]interface{}) (ParameterDefinition, error) {
	id, name, err := extractValues(rawDefinition)
	if err != nil {
		return ParameterDefinition{}, err
	}
	return ParameterDefinition{Id: id, Name: name, RawDefinition: rawDefinition}, nil
}

func extractValues(def map[string]interface{}) (id, name string, err error) {
	id, err = extractStringValue(def, "id")
	if err != nil {
		return "", "", err
	}
	name, err = extractStringValue(def, "name")
	if err != nil {
		return "", "", err
	}
	return id, name, nil
}

func extractStringValue(def map[string]interface{}, key string) (string, error) {
	if _, ok := def[key]; !ok {
		return "", fmt.Errorf("entry %q missing in parameter definition %v", key, def)
	}
	 if value, ok := def[key].(string); !ok {
		return value, nil
	} 
        return "", fmt.Errorf("unexpected type of key %q in parameter definition: %t, expected string", key, def[key])
}

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
func (v *Validator) ValidateParameters(definitions []ParameterDefinition, params extensionAPI.ParameterValues) (failedValidations []ValidationResult, err error) {
	result := make([]ValidationResult, 0)
	for _, def := range definitions {
		name, r, err := v.validateParameter(def, params)
		if err != nil {
			return nil, err
		}
		if !r.Success {
			result = append(result, ValidationResult{Success: false, Message: fmt.Sprintf("Failed to validate parameter '%s': %s", name, r.Message)})
		}
	}
	return result, nil
}

func (v *Validator) validateParameter(def ParameterDefinition, params extensionAPI.ParameterValues) (string, *ValidationResult, error) {
	paramValue := findParamValue(params, def.Id)
	result, err := v.ValidateParameter(def, paramValue)
	if err != nil {
		return "", nil, fmt.Errorf("failed to validate parameter value %q with id %q using definition %v", paramValue, def.Id, def.RawDefinition)
	}
	return def.Name, result, nil
}

func findParamValue(params extensionAPI.ParameterValues, id string) string {
	if param, found := params.Find(id); found {
		return param.Value
	}
	return ""
}

// ValidateParameters uses the given parameter definition to validate a single value
func (v *Validator) ValidateParameter(def ParameterDefinition, value string) (validationResult *ValidationResult, errorResult error) {
	defer func() {
		if err := recover(); err != nil {
			errorResult = fmt.Errorf("failed to validate parameter value %q using definition %v: %v", value, def, err)
		}
	}()
	result := v.validate(def.RawDefinition, value)
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
