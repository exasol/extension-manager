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
	/** TestSetupAbstraction reuses the container, so parallel use of this suite would cause conflicts. We make sure it's not used in parallel using a mutex. */
	exasol, err := testSetupAbstraction.Create("./exasol-test-setup-config.json") // file does not exist --> we use the testcontainer test setup
	if err != nil {
		suite.FailNowf("failed to create test setup abstraction: %v", err.Error())
	}
	suite.Exasol = exasol
}

func (suite *IntegrationTestSuite) TearDownSuite() {
	suite.NoError(suite.Exasol.Stop())
}

func (suite *IntegrationTestSuite) ExecSQL(query string) {
	_, err := suite.Connection.Exec(query)
	suite.NoError(err)
}

func (suite *IntegrationTestSuite) BeforeTest(suiteName, testName string) {
	if suite.Connection != nil {
		suite.FailNow("previous connection was not closed")
	}
	db, err := suite.Exasol.CreateConnectionWithConfig(false)
	if err != nil {
		suite.FailNowf("failed to connect to db: %v", err.Error())
	}
	suite.Connection = db
}

func (suite *IntegrationTestSuite) AfterTest(suiteName, testName string) {
	err := suite.Connection.Close()
	suite.NoError(err)
	suite.Connection = nil
}
