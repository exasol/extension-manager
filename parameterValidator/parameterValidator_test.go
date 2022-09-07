package parameterValidator

import (
	"testing"

	"github.com/exasol/extension-manager/extensionAPI"
	"github.com/stretchr/testify/suite"
)

type ParameterValidatorSuite struct {
	suite.Suite
	validator *Validator
}

func TestParameterValidatorSuite(t *testing.T) {
	suite.Run(t, new(ParameterValidatorSuite))
}

func (suite *ParameterValidatorSuite) SetupSuite() {
	v, err := New()
	if err != nil {
		suite.Fail(err.Error())
	}
	suite.validator = v
}

func (suite *ParameterValidatorSuite) TestValidateParameter() {
	var cases = []struct {
		definition map[string]interface{}
		expected   ValidationResult
	}{
		{definition: map[string]interface{}{"type": "string", "id": "my-value", "required": true, "regex": ".*"},
			expected: ValidationResult{Success: true, Message: ""}},
		{definition: map[string]interface{}{"type": "string", "id": "my-value", "required": true, "regex": "a+"},
			expected: ValidationResult{Success: false, Message: "The value has an invalid format."}},
	}

	for _, testCase := range cases {
		result, err := suite.validator.ValidateParameter(testCase.definition, "test")
		suite.NoError(err)
		suite.Equal(testCase.expected, *result)
	}
}

func (suite *ParameterValidatorSuite) TestValidateParameters() {
	var tests = []struct {
		name        string
		definitions []interface{}
		params      []extensionAPI.ParameterValue
		expected    []ValidationResult
	}{
		{name: "success",
			definitions: []interface{}{map[string]interface{}{"id": "param1", "name": "My param", "type": "string"}},
			params:      []extensionAPI.ParameterValue{{Name: "param1", Value: "value"}},
			expected:    []ValidationResult{}},
		{name: "missing optional parameter",
			definitions: []interface{}{map[string]interface{}{"id": "param1", "name": "My param", "type": "string"}},
			params:      []extensionAPI.ParameterValue{},
			expected:    []ValidationResult{}},
		{name: "empty required parameter",
			definitions: []interface{}{map[string]interface{}{"id": "param1", "name": "My param", "type": "string", "required": true}},
			params:      []extensionAPI.ParameterValue{{Name: "param1", Value: ""}},
			expected:    []ValidationResult{{Success: false, Message: `Failed to validate parameter 'My param': This is a required parameter.`}}},
		{name: "empty non-required parameter",
			definitions: []interface{}{map[string]interface{}{"id": "param1", "name": "My param", "type": "string", "required": false}},
			params:      []extensionAPI.ParameterValue{{Name: "param1", Value: ""}},
			expected:    []ValidationResult{}},
		{name: "missing non-required parameter",
			definitions: []interface{}{map[string]interface{}{"id": "param1", "name": "My param", "type": "string", "required": false}},
			params:      []extensionAPI.ParameterValue{},
			expected:    []ValidationResult{}},
		{name: "valid regex parameter",
			definitions: []interface{}{map[string]interface{}{"id": "param1", "name": "My param", "type": "string", "regex": "^a+$"}},
			params:      []extensionAPI.ParameterValue{{Name: "param1", Value: "aaa"}},
			expected:    []ValidationResult{}},
		{name: "invalid regex parameter",
			definitions: []interface{}{map[string]interface{}{"id": "param1", "name": "My param", "type": "string", "regex": "^a+$"}},
			params:      []extensionAPI.ParameterValue{{Name: "param1", Value: "ab"}},
			expected:    []ValidationResult{{Success: false, Message: `Failed to validate parameter 'My param': The value has an invalid format.`}}},
		{name: "invalid boolean parameter",
			definitions: []interface{}{map[string]interface{}{"id": "param1", "name": "My param", "type": "boolean"}},
			params:      []extensionAPI.ParameterValue{{Name: "param1", Value: "invalid"}},
			expected:    []ValidationResult{{Success: false, Message: `Failed to validate parameter 'My param': Boolean value must be 'true' or 'false'.`}}},
		{name: "valid boolean parameter",
			definitions: []interface{}{map[string]interface{}{"id": "param1", "name": "My param", "type": "boolean"}},
			params:      []extensionAPI.ParameterValue{{Name: "param1", Value: "true"}},
			expected:    []ValidationResult{}},
	}

	for _, t := range tests {
		suite.Run(t.name, func() {
			result, err := suite.validator.ValidateParameters(t.definitions, extensionAPI.ParameterValues{Values: t.params})
			suite.NoError(err)
			suite.Len(result, len(t.expected))
			suite.Equal(t.expected, result)
		})
	}
}

func (suite *ParameterValidatorSuite) TestInvalidDefinitionIgnored() {
	result, err := suite.validator.ValidateParameters([]interface{}{map[string]interface{}{"id": "param1", "name": "My param", "type": "invalidType"}}, extensionAPI.ParameterValues{})
	suite.NoError(err)
	suite.Empty(result)
}
