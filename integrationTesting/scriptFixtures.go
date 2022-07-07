package integrationTesting

import (
	"database/sql"
	"fmt"
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

func execSQL(db *sql.DB, sql string) {
	_, err := db.Exec(sql)
	if err != nil {
		panic(fmt.Sprintf("error executing SQL %q: %v", sql, err.Error()))
	}
}

func (fixture ScriptFixture) Close() {
	execSQL(fixture.db, "DROP SCHEMA IF EXISTS TEST CASCADE")
}
