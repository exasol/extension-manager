package parameterValidator

import (
	"github.com/stretchr/testify/suite"
	"testing"
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
	for _, testCase := range cases {
		result, err := validateParameter(testCase.definition, "test")
		suite.NoError(err)
		suite.Assert().Equal(testCase.expected, *result)
	}
}
