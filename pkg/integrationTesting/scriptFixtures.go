package integrationTesting

import (
	"database/sql"
	"fmt"
	"testing"
)

type ScriptFixture struct {
	db *sql.DB
}

func CreateLuaScriptFixture(db *sql.DB) ScriptFixture {
	execSQL(db, "CREATE SCHEMA TEST")
	execSQL(db, `
CREATE LUA SET SCRIPT test.my_script (a DOUBLE)
RETURNS DOUBLE AS
function run(ctx) return 1 end`)
	execSQL(db, "COMMENT ON SCRIPT test.my_script IS 'my comment'")
	return ScriptFixture{db: db}
}

func CreateJavaAdapterScriptFixture(db *sql.DB) ScriptFixture {
	execSQL(db, "CREATE SCHEMA TEST")
	execSQL(db, `
CREATE OR REPLACE JAVA ADAPTER SCRIPT TEST.VS_ADAPTER AS
%scriptclass com.exasol.adapter.RequestDispatcher;
%jar /buckets/bfsdefault/default/vs.jar;`)
	return ScriptFixture{db: db}
}

func CreateJavaSetScriptFixture(db *sql.DB) ScriptFixture {
	execSQL(db, "CREATE SCHEMA TEST")
	execSQL(db, `
CREATE OR REPLACE JAVA SET SCRIPT TEST.IMPORT_FROM_S3_DOCUMENT_FILES(
DATA_LOADER VARCHAR(2000000),
SCHEMA_MAPPING_REQUEST VARCHAR(2000000),
CONNECTION_NAME VARCHAR(500))
EMITS(...) AS
%scriptclass com.exasol.adapter.document.UdfEntryPoint;
%jar /buckets/bfsdefault/default/vs.jar;`)
	return ScriptFixture{db: db}
}

func CreateVirtualSchemaFixture(db *sql.DB) ScriptFixture {
	execSQL(db, "CREATE SCHEMA TEST_META_DATA")
	createMetaDataTable(db, "EXA_ALL_VIRTUAL_SCHEMAS")
	createMetaDataTable(db, "EXA_ALL_SCRIPTS")
	execSQL(db, `INSERT INTO TEST_META_DATA.EXA_ALL_VIRTUAL_SCHEMAS VALUES ('schema1', 'owner1', 1, 'TEST', 'script1', '2024-04-25 09:31:08', 'user1', 'notes1')`)
	return ScriptFixture{db: db}
}

func CreateVirtualSchemaFixtureNullValues(db *sql.DB) ScriptFixture {
	execSQL(db, "CREATE SCHEMA TEST_META_DATA")
	createMetaDataTable(db, "EXA_ALL_VIRTUAL_SCHEMAS")
	createMetaDataTable(db, "EXA_ALL_SCRIPTS")
	// nolint:dupword // Duplicating the word "NULL" is required here
	execSQL(db, `INSERT INTO TEST_META_DATA.EXA_ALL_VIRTUAL_SCHEMAS VALUES (NULL, NULL, NULL, 'TEST', NULL, NULL, NULL, NULL)`)
	return ScriptFixture{db: db}
}

func CreateScriptFixtureNullValues(db *sql.DB) ScriptFixture {
	execSQL(db, "CREATE SCHEMA TEST_META_DATA")
	createMetaDataTable(db, "EXA_ALL_VIRTUAL_SCHEMAS")
	createMetaDataTable(db, "EXA_ALL_SCRIPTS")
	// nolint:dupword // Duplicating the word "NULL" is required here
	execSQL(db, `INSERT INTO TEST_META_DATA.EXA_ALL_SCRIPTS VALUES ('TEST', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL)`)
	return ScriptFixture{db: db}
}

func createMetaDataTable(db *sql.DB, tableName string) {
	execSQL(db, fmt.Sprintf(`CREATE TABLE TEST_META_DATA.%s AS SELECT * FROM SYS.%s WHERE 1=2`, tableName, tableName))
}

func execSQL(db *sql.DB, sql string) {
	_, err := db.Exec(sql)
	if err != nil {
		panic(fmt.Sprintf("error executing SQL %q: %v", sql, err))
	}
}

func (f ScriptFixture) GetSchemaName() string {
	return "TEST"
}

func (f ScriptFixture) GetMetaDataSchemaName() string {
	return "TEST_META_DATA"
}

func (f ScriptFixture) Cleanup(t *testing.T) {
	t.Cleanup(func() {
		execSQL(f.db, "DROP SCHEMA IF EXISTS TEST CASCADE")
		execSQL(f.db, "DROP SCHEMA IF EXISTS TEST_META_DATA CASCADE")
	})
}
