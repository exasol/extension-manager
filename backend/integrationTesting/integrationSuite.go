package integrationTesting

import (
	"database/sql"
	"testing"

	testSetupAbstraction "github.com/exasol/exasol-test-setup-abstraction-server/go-client"
	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite
	Exasol     *testSetupAbstraction.TestSetupAbstraction
	Connection *sql.DB
}

func (suite *IntegrationTestSuite) SetupSuite() {
	if testing.Short() {
		suite.T().Skip()
	}
	/** Since the  testSetupAbstraction reuses the container parallel use of this suite would cause conflicts. --> We make sure it's not used in parallel using a mutex */
	exasol := testSetupAbstraction.Create("./exasol-test-setup-config.json") // file does not exist. --> we use the testcontainer test setup
	suite.Exasol = &exasol
	suite.Connection = exasol.CreateConnection()
}

func (suite *IntegrationTestSuite) TearDownSuite() {
	suite.Exasol.Stop()
}

func (suite *IntegrationTestSuite) ExecSQL(query string) {
	_, err := suite.Connection.Exec(query)
	suite.NoError(err)
}
