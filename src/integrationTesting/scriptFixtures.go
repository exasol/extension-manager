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

func execSQL(db *sql.DB, sql string) {
	_, err := db.Exec(sql)
	if err != nil {
		panic(fmt.Sprintf("error executing SQL %q: %v", sql, err))
	}
}

func (f ScriptFixture) GetSchemaName() string {
	return "TEST"
}

func (f ScriptFixture) Cleanup(t *testing.T) {
	t.Cleanup(func() {
		execSQL(f.db, "DROP SCHEMA IF EXISTS TEST CASCADE")
	})
}
