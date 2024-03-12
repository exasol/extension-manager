package exaMetadata

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
)

type ExaMetadataUTestSuite struct {
	suite.Suite
	db     *sql.DB
	dbMock sqlmock.Sqlmock
}

func TestExaMetadataUTestSuite(t *testing.T) {
	suite.Run(t, new(ExaMetadataUTestSuite))
}

func (suite *ExaMetadataUTestSuite) SetupTest() {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	suite.Require().NoError(err)
	suite.db = db
	suite.dbMock = mock
	suite.dbMock.MatchExpectationsInOrder(true)
}

func (suite *ExaMetadataUTestSuite) AfterTest(suiteName, testName string) {
	suite.NoError(suite.dbMock.ExpectationsWereMet())
}

func (suite *ExaMetadataUTestSuite) TestCreateExaMetaDataReader() {
	suite.NotNil(CreateExaMetaDataReader())
}

const SCHEMA_NAME = "EXA_SCHEMA_NAME"

func (suite *ExaMetadataUTestSuite) TestReadMetadataTables() {
	tx := suite.beginTransaction()
	suite.dbMock.ExpectQuery("(?m)SELECT .*FROM SYS.EXA_ALL_SCRIPTS .*").WithArgs(SCHEMA_NAME).
		WillReturnRows(sqlmock.
			NewRows([]string{"SCRIPT_SCHEMA", "SCRIPT_NAME", "SCRIPT_TYPE", "SCRIPT_INPUT_TYPE", "SCRIPT_RESULT_TYPE", "SCRIPT_TEXT", "SCRIPT_COMMENT"}).
			AddRow("schema1", "script1", "type1", "input_type1", "result_type1", "text1", "comment1").
			AddRow("schema2", "script2", "type2", "input_type2", "result_type2", "text2", "comment2")).
		RowsWillBeClosed()
	suite.dbMock.ExpectQuery("SELECT SCHEMA_NAME, SCHEMA_OWNER, ADAPTER_SCRIPT_SCHEMA, ADAPTER_SCRIPT_NAME, ADAPTER_NOTES\\s+FROM SYS.EXA_ALL_VIRTUAL_SCHEMAS").WillReturnRows(sqlmock.
		NewRows([]string{"SCHEMA_NAME", "SCHEMA_OWNER", "ADAPTER_SCRIPT_SCHEMA", "ADAPTER_SCRIPT_NAME", "ADAPTER_NOTES"}).
		AddRow("schema1", "owner1", "scriptSchema1", "script1", "notes1").
		AddRow("schema2", "owner2", "scriptSchema2", "script2", "notes2")).
		RowsWillBeClosed()

	metadata, err := CreateExaMetaDataReader().ReadMetadataTables(tx, SCHEMA_NAME)
	suite.Require().NoError(err)
	suite.Equal(&ExaMetadata{AllScripts: ExaScriptTable{Rows: []ExaScriptRow{
		{Schema: "schema1", Name: "script1", Type: "type1", InputType: "input_type1", ResultType: "result_type1", Text: "text1", Comment: "comment1"},
		{Schema: "schema2", Name: "script2", Type: "type2", InputType: "input_type2", ResultType: "result_type2", Text: "text2", Comment: "comment2"},
	}}, AllVirtualSchemas: ExaVirtualSchemasTable{Rows: []ExaVirtualSchemaRow{
		{Name: "schema1", Owner: "owner1", AdapterScriptSchema: "scriptSchema1", AdapterScriptName: "script1", AdapterNotes: "notes1"},
		{Name: "schema2", Owner: "owner2", AdapterScriptSchema: "scriptSchema2", AdapterScriptName: "script2", AdapterNotes: "notes2"},
	}}}, metadata)
}

func (suite *ExaMetadataUTestSuite) TestReadMetadataTablesAllScriptsFails() {
	tx := suite.beginTransaction()
	suite.dbMock.ExpectQuery("(?m)SELECT .*FROM SYS.EXA_ALL_SCRIPTS .*").WithArgs(SCHEMA_NAME).WillReturnError(errors.New("mock error"))

	metadata, err := CreateExaMetaDataReader().ReadMetadataTables(tx, SCHEMA_NAME)
	suite.Require().EqualError(err, "failed to read SYS.EXA_ALL_SCRIPTS: mock error")
	suite.Nil(metadata)
}

func (suite *ExaMetadataUTestSuite) TestReadMetadataTablesAllVirtualSchemasFails() {
	tx := suite.beginTransaction()
	suite.dbMock.ExpectQuery("(?m)SELECT .*FROM SYS.EXA_ALL_SCRIPTS .*").WithArgs(SCHEMA_NAME).
		WillReturnRows(sqlmock.
			NewRows([]string{"SCRIPT_SCHEMA", "SCRIPT_NAME", "SCRIPT_TYPE", "SCRIPT_INPUT_TYPE", "SCRIPT_RESULT_TYPE", "SCRIPT_TEXT", "SCRIPT_COMMENT"}).
			AddRow("schema1", "script1", "type1", "input_type1", "result_type1", "text1", "comment1").
			AddRow("schema2", "script2", "type2", "input_type2", "result_type2", "text2", "comment2")).
		RowsWillBeClosed()
	suite.dbMock.ExpectQuery("SELECT SCHEMA_NAME, SCHEMA_OWNER, ADAPTER_SCRIPT_SCHEMA, ADAPTER_SCRIPT_NAME, ADAPTER_NOTES\\s+FROM SYS.EXA_ALL_VIRTUAL_SCHEMAS").WillReturnError(errors.New("mock error"))

	metadata, err := CreateExaMetaDataReader().ReadMetadataTables(tx, SCHEMA_NAME)
	suite.Require().EqualError(err, "failed to read SYS.EXA_ALL_VIRTUAL_SCHEMAS: mock error")
	suite.Nil(metadata)
}

func (suite *ExaMetadataUTestSuite) TestReadExaAllScriptTableQueryFails() {
	tx := suite.beginTransaction()
	suite.dbMock.ExpectQuery("(?m)SELECT .*FROM SYS.EXA_ALL_SCRIPTS .*").WithArgs(SCHEMA_NAME).WillReturnError(errors.New("mock error"))
	result, err := readExaAllScriptTable(tx, SCHEMA_NAME)
	suite.Require().EqualError(err, "failed to read SYS.EXA_ALL_SCRIPTS: mock error")
	suite.Nil(result)
}

func (suite *ExaMetadataUTestSuite) TestReadExaAllScriptTableScanFails() {
	tx := suite.beginTransaction()
	suite.dbMock.ExpectQuery("(?m)SELECT .*FROM SYS.EXA_ALL_SCRIPTS .*").WithArgs(SCHEMA_NAME).
		WillReturnRows(sqlmock.NewRows([]string{"WRONG_COL"}).AddRow("Wrong")).
		RowsWillBeClosed()
	result, err := readExaAllScriptTable(tx, SCHEMA_NAME)
	suite.Require().EqualError(err, "failed to read row of SYS.EXA_ALL_SCRIPTS: sql: expected 1 destination arguments in Scan, not 7")
	suite.Nil(result)
}

func (suite *ExaMetadataUTestSuite) TestReadExaAllVirtualSchemasTableFails() {
	tx := suite.beginTransaction()
	suite.dbMock.ExpectQuery("(?m)SELECT .* FROM SYS.EXA_ALL_VIRTUAL_SCHEMAS").WillReturnError(errors.New("mock error"))
	result, err := readExaAllVirtualSchemasTable(tx)
	suite.Require().EqualError(err, "failed to read SYS.EXA_ALL_VIRTUAL_SCHEMAS: mock error")
	suite.Nil(result)
}

func (suite *ExaMetadataUTestSuite) TestReadExaAllVirtualSchemasTableScanFails() {
	tx := suite.beginTransaction()
	suite.dbMock.ExpectQuery("(?m)SELECT .* FROM SYS.EXA_ALL_VIRTUAL_SCHEMAS").WillReturnRows(sqlmock.NewRows([]string{"wrong"}).AddRow("wrong"))
	result, err := readExaAllVirtualSchemasTable(tx)
	suite.Require().EqualError(err, "failed to read row of SYS.EXA_ALL_VIRTUAL_SCHEMAS: sql: expected 1 destination arguments in Scan, not 5")
	suite.Nil(result)
}

// GetScriptByName

/* [utest -> dsn~extension-context-metadata~1]. */
func (suite *ExaMetadataUTestSuite) TestGetScriptByName() {
	tx := suite.beginTransaction()
	suite.dbMock.ExpectQuery("(?m)SELECT .*FROM SYS.EXA_ALL_SCRIPTS .*WHERE SCRIPT_SCHEMA=\\? AND SCRIPT_NAME=\\?").WithArgs(SCHEMA_NAME, "script").
		WillReturnRows(sqlmock.
			NewRows([]string{"SCRIPT_SCHEMA", "SCRIPT_NAME", "SCRIPT_TYPE", "SCRIPT_INPUT_TYPE", "SCRIPT_RESULT_TYPE", "SCRIPT_TEXT", "SCRIPT_COMMENT"}).
			AddRow("schema1", "script1", "type1", "input_type1", "result_type1", "text1", "comment1")).
		RowsWillBeClosed()
	result, err := CreateExaMetaDataReader().GetScriptByName(tx, SCHEMA_NAME, "script")
	suite.Require().NoError(err)
	suite.Equal(&ExaScriptRow{Schema: "schema1", Name: "script1", Type: "type1", InputType: "input_type1", ResultType: "result_type1", Text: "text1", Comment: "comment1"}, result)
}

func (suite *ExaMetadataUTestSuite) TestGetScriptByNameIgnoresSecondLine() {
	tx := suite.beginTransaction()
	suite.dbMock.ExpectQuery("(?m)SELECT .*FROM SYS.EXA_ALL_SCRIPTS .*WHERE SCRIPT_SCHEMA=\\? AND SCRIPT_NAME=\\?").WithArgs(SCHEMA_NAME, "script").
		WillReturnRows(sqlmock.
			NewRows([]string{"SCRIPT_SCHEMA", "SCRIPT_NAME", "SCRIPT_TYPE", "SCRIPT_INPUT_TYPE", "SCRIPT_RESULT_TYPE", "SCRIPT_TEXT", "SCRIPT_COMMENT"}).
			AddRow("schema1", "script1", "type1", "input_type1", "result_type1", "text1", "comment1").
			AddRow("ignored-schema1", "ignored-script1", "ignored-type1", "ignored-input_type1", "ignored-result_type1", "ignored-text1", "ignored-comment1")).
		RowsWillBeClosed()
	result, err := CreateExaMetaDataReader().GetScriptByName(tx, SCHEMA_NAME, "script")
	suite.Require().NoError(err)
	suite.Equal(&ExaScriptRow{Schema: "schema1", Name: "script1", Type: "type1", InputType: "input_type1", ResultType: "result_type1", Text: "text1", Comment: "comment1"}, result)
}

func (suite *ExaMetadataUTestSuite) TestGetScriptByNameQueryFails() {
	tx := suite.beginTransaction()
	suite.dbMock.ExpectQuery("(?m)SELECT .*FROM SYS.EXA_ALL_SCRIPTS .*WHERE SCRIPT_SCHEMA=\\? AND SCRIPT_NAME=\\?").WithArgs(SCHEMA_NAME, "script").
		WillReturnError(errors.New("mock error"))
	result, err := CreateExaMetaDataReader().GetScriptByName(tx, SCHEMA_NAME, "script")
	suite.Require().EqualError(err, "failed to read SYS.EXA_ALL_SCRIPTS: mock error")
	suite.Nil(result)
}

func (suite *ExaMetadataUTestSuite) TestGetScriptByNameNoResult() {
	tx := suite.beginTransaction()
	suite.dbMock.ExpectQuery("(?m)SELECT .*FROM SYS.EXA_ALL_SCRIPTS .*WHERE SCRIPT_SCHEMA=\\? AND SCRIPT_NAME=\\?").WithArgs(SCHEMA_NAME, "script").
		WillReturnRows(sqlmock.
			NewRows([]string{"SCRIPT_SCHEMA", "SCRIPT_NAME", "SCRIPT_TYPE", "SCRIPT_INPUT_TYPE", "SCRIPT_RESULT_TYPE", "SCRIPT_TEXT", "SCRIPT_COMMENT"})).
		RowsWillBeClosed()
	result, err := CreateExaMetaDataReader().GetScriptByName(tx, SCHEMA_NAME, "script")
	suite.Require().NoError(err)
	suite.Nil(result)
}

func (suite *ExaMetadataUTestSuite) TestGetScriptByNameReadingFails() {
	tx := suite.beginTransaction()
	suite.dbMock.ExpectQuery("(?m)SELECT .*FROM SYS.EXA_ALL_SCRIPTS .*WHERE SCRIPT_SCHEMA=\\? AND SCRIPT_NAME=\\?").WithArgs(SCHEMA_NAME, "script").
		WillReturnRows(sqlmock.NewRows([]string{"invalid"}).AddRow("invalid")).
		RowsWillBeClosed()
	result, err := CreateExaMetaDataReader().GetScriptByName(tx, SCHEMA_NAME, "script")
	suite.Require().EqualError(err, `failed to read row of SYS.EXA_ALL_SCRIPTS: sql: expected 1 destination arguments in Scan, not 7`)
	suite.Nil(result)
}

func (suite *ExaMetadataUTestSuite) beginTransaction() *sql.Tx {
	suite.dbMock.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.Require().NoError(err)
	return tx
}
