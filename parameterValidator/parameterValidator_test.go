package parameterValidator

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ParameterValidatorSuite struct {
	suite.Suite
}

func TestParameterValidatorSuite(t *testing.T) {
	suite.Run(t, new(ParameterValidatorSuite))
}

func (suite *ParameterValidatorSuite) TestValidate() {
	var cases = []struct {
		definition map[string]interface{}
		expected   ValidationResult
	}{
		{definition: map[string]interface{}{"type": "string", "id": "my-value", "required": true, "regex": ".*"},
			expected: ValidationResult{Success: true, Message: ""}},
		{definition: map[string]interface{}{"type": "string", "id": "my-value", "required": true, "regex": "a+"},
			expected: ValidationResult{Success: false, Message: "The value has an invalid format."}},
	}

	validator, err := New()
	if err != nil {
		suite.Fail(err.Error())
	}
	for _, testCase := range cases {
		result, err := validator.ValidateParameter(testCase.definition, "test")
		suite.NoError(err)
		suite.Assert().Equal(testCase.expected, *result)
	}
}
