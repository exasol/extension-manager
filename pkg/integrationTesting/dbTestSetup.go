package integrationTesting

import (
	"database/sql"
	"os"
	"testing"

	testSetupAbstraction "github.com/exasol/exasol-test-setup-abstraction-server/go-client"
	"github.com/stretchr/testify/suite"
)

const defaultExasolDbVersion = "7.1.22"

type DbTestSetup struct {
	suite          *suite.Suite
	Exasol         *testSetupAbstraction.TestSetupAbstraction
	connection     *sql.DB
	ConnectionInfo *testSetupAbstraction.ConnectionInfo
}

func StartDbSetup(suite *suite.Suite) *DbTestSetup {
	if testing.Short() {
		suite.T().Skip()
	}
	exasolDbVersion := getDbVersion()
	suite.T().Logf("Starting Exasol %s...", exasolDbVersion)
	exasol, err := testSetupAbstraction.New().DockerDbVersion(exasolDbVersion).Start()
	if err != nil {
		suite.FailNowf("failed to create test setup abstraction: %v", err.Error())
	}
	connectionInfo, err := exasol.GetConnectionInfo()
	if err != nil {
		suite.FailNowf("error getting connection info: %v", err.Error())
	}
	setup := DbTestSetup{suite: suite, Exasol: exasol, ConnectionInfo: connectionInfo}
	return &setup
}

func getDbVersion() string {
	dbVersion := os.Getenv("EXASOL_VERSION")
	if dbVersion != "" {
		return dbVersion
	}
	return defaultExasolDbVersion
}

func (setup *DbTestSetup) StopDb() {
	setup.suite.NoError(setup.Exasol.Stop())
}

func (setup *DbTestSetup) ExecSQL(query string) {
	_, err := setup.connection.Exec(query)
	setup.suite.NoError(err)
}

func (setup *DbTestSetup) GetConnection() *sql.DB {
	if setup.connection == nil {
		setup.suite.FailNow("no db connection. CreateConnection() in BeforeTest(suiteName, testName string).")
	}
	return setup.connection
}

func (setup *DbTestSetup) CreateConnection() {
	if setup.connection != nil {
		setup.suite.FailNow("previous connection was not closed. Run CloseConnection() in AfterTest(suiteName, testName string).")
	}
	db, err := setup.Exasol.CreateConnectionWithConfig(false)
	if err != nil {
		setup.suite.FailNowf("failed to connect to db: %v", err.Error())
	}
	setup.connection = db
}

func (setup *DbTestSetup) CloseConnection() {
	if setup.connection == nil {
		setup.suite.FailNow("no connection to close after test. Run CreateConnection() in BeforeTest(suiteName, testName string).")
	}
	err := setup.connection.Close()
	setup.suite.NoError(err)
	setup.connection = nil
}
