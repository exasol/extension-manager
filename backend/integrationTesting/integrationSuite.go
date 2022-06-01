package integrationTesting

import (
	"database/sql"
	"log"
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
	exasol, err := testSetupAbstraction.Create("./exasol-test-setup-config.json") // file does not exist. --> we use the testcontainer test setup
	if err != nil {
		log.Fatalf("failed to create test setup abstraction. Cause: %v", err)
	}
	suite.Exasol = exasol
	suite.Connection, err = exasol.CreateConnection()
	suite.NoError(err)
}

func (suite *IntegrationTestSuite) TearDownSuite() {
	suite.NoError(suite.Exasol.Stop())
}

func (suite *IntegrationTestSuite) ExecSQL(query string) {
	_, err := suite.Connection.Exec(query)
	suite.NoError(err)
}
