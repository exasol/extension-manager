package extensionApi

import (
	"database/sql"
	"fmt"
)

func ReadExaAllScriptTable(connection *sql.DB) (*ExaAllScriptTable, error) {
	result, err := connection.Query("SELECT SCRIPT_SCHEMA, SCRIPT_NAME, SCRIPT_TEXT FROM SYS.EXA_ALL_SCRIPTS")
	if err != nil {
		return nil, fmt.Errorf("failed to read SYS.EXA_ALL_SCRIPTS. Cause: %v", err.Error())
	}
	var rows []ExaAllScriptRow
	for result.Next() {
		var row ExaAllScriptRow
		var schema string
		var name string
		err := result.Scan(&schema, &name, &row.Text)
		if err != nil {
			return nil, fmt.Errorf("failed to read row of SYS.EXA_ALL_SCRIPTS. Cause: %w", err)
		}
		row.Name = schema + "." + name
		rows = append(rows, row)
	}
	return &ExaAllScriptTable{Rows: rows}, nil
}

type ExaAllScriptTable struct {
	Rows []ExaAllScriptRow `json:"rows"`
}

type ExaAllScriptRow struct {
	Name string `json:"name"`
	Text string `json:"text"`
}
