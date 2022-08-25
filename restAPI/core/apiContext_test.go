package core

import (
	"encoding/base64"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ApiContextSuite struct {
	suite.Suite
}

func TestApiContextSuite(t *testing.T) {
	suite.Run(t, new(ApiContextSuite))
}

func (suite *ApiContextSuite) TestExtractUserPasswordExample() {
	// Example from https://datatracker.ietf.org/doc/html/rfc7617#page-5
	user, password, err := extractUserPassword("QWxhZGRpbjpvcGVuIHNlc2FtZQ==")
	suite.NoError(err)
	suite.Equal(user, "Aladdin")
	suite.Equal(password, "open sesame")
}

func (suite *ApiContextSuite) TestExtractUserPasswordInvalidBase64() {
	user, password, err := extractUserPassword("invalid base64")
	suite.EqualError(err, "invalid basic auth header \"invalid base64\": illegal base64 data at input byte 7")
	suite.Equal(user, "")
	suite.Equal(password, "")
}

func (suite *ApiContextSuite) TestExtractUserPassword() {

	tests := []struct {
		input            string
		expectedUser     string
		expectedPassword string
		expectedError    string
	}{
		{input: "user:password", expectedUser: "user", expectedPassword: "password"},
		{input: "user:pass:word", expectedUser: "user", expectedPassword: "pass:word"},
		{input: "öäü!µ:`«@≠", expectedUser: "öäü!µ", expectedPassword: "`«@≠"},
		{input: "nocolon", expectedError: "colon missing in basic auth header"},
	}
	for _, test := range tests {
		suite.Run(test.input, func() {
			encoded := base64.StdEncoding.EncodeToString([]byte(test.input))
			user, password, err := extractUserPassword(encoded)
			if test.expectedError == "" {
				suite.NoError(err)
				suite.Equal(test.expectedUser, user)
				suite.Equal(test.expectedPassword, password)
			} else {
				suite.EqualError(err, test.expectedError)
				suite.Equal("", user)
				suite.Equal("", password)
			}
		})
	}
}
