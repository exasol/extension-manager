package extensionApi

import (
	"database/sql"
	testSetupAbstraction "github.com/exasol/exasol-test-setup-abstraction-server/go-client"
	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite
	Exasol     *testSetupAbstraction.TestSetupAbstraction
	Connection *sql.DB
}

func (suite *IntegrationTestSuite) SetupSuite() {
	exasol := testSetupAbstraction.Create("./exasol-test-setup-config.json") // file does not exist. --> we use the testcontainer test setup
	suite.Exasol = &exasol
	suite.Connection = exasol.CreateConnection()
}

func (suite *IntegrationTestSuite) ExecSQL(query string) {
	_, err := suite.Connection.Exec(query)
	suite.NoError(err)
}
