package extensionAPI

import (
	"database/sql"
	"fmt"
)

type ExaMetadata struct {
	AllScripts ExaAllScriptTable `json:"allScripts"`
}

func ReadMetadataTables(connection *sql.DB, schemaName string) (*ExaMetadata, error) {
	allScripts, err := readExaAllScriptTable(connection, schemaName)
	if err != nil {
		return nil, err
	}
	return &ExaMetadata{AllScripts: *allScripts}, nil
}

func readExaAllScriptTable(connection *sql.DB, schemaName string) (*ExaAllScriptTable, error) {
	result, err := connection.Query(`
SELECT SCRIPT_SCHEMA, SCRIPT_NAME, SCRIPT_TYPE, SCRIPT_INPUT_TYPE, SCRIPT_RESULT_TYPE, SCRIPT_TEXT, SCRIPT_COMMENT
FROM SYS.EXA_ALL_SCRIPTS
WHERE SCRIPT_SCHEMA=?`, schemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to read SYS.EXA_ALL_SCRIPTS. Cause: %w", err)
	}
	var rows []ExaAllScriptRow
	for result.Next() {
		var row ExaAllScriptRow
		var inputType sql.NullString
		var resultType sql.NullString
		var comment sql.NullString
		err := result.Scan(&row.Schema, &row.Name, &row.Type, &inputType, &resultType, &row.Text, &comment)
		if err != nil {
			return nil, fmt.Errorf("failed to read row of SYS.EXA_ALL_SCRIPTS. Cause: %w", err)
		}
		row.InputType = inputType.String
		row.ResultType = resultType.String
		row.Comment = comment.String
		rows = append(rows, row)
	}
	return &ExaAllScriptTable{Rows: rows}, nil
}

type ExaAllScriptTable struct {
	Rows []ExaAllScriptRow `json:"rows"`
}

type ExaAllScriptRow struct {
	Schema     string `json:"schema"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	InputType  string `json:"inputType"`
	ResultType string `json:"resultType"`
	Text       string `json:"text"`
	Comment    string `json:"comment"`
}
