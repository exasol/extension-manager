package extensionAPI

import (
	"database/sql"
	"fmt"
	"strings"
)

type ExaMetadataReader interface {
	ReadMetadataTables(tx *sql.Tx, schemaName string) (*ExaMetadata, error)
}

type ExaMetadata struct {
	AllScripts        ExaScriptTable         `json:"allScripts"`
	AllVirtualSchemas ExaVirtualSchemasTable `json:"allVirtualSchemas"`
}

func CreateExaMetaDataReader() ExaMetadataReader {
	return &metaDataReaderImpl{}
}

type metaDataReaderImpl struct {
}

// ReadMetadataTables reads the metadata tables of the given schema.
/* [impl -> dsn~extension-components~1] */
func (r *metaDataReaderImpl) ReadMetadataTables(tx *sql.Tx, schemaName string) (*ExaMetadata, error) {
	allScripts, err := readExaAllScriptTable(tx, schemaName)
	if err != nil {
		return nil, err
	}
	allVirtualSchemas, err := readExaAllVirtualSchemasTable(tx)
	if err != nil {
		return nil, err
	}
	return &ExaMetadata{AllScripts: *allScripts, AllVirtualSchemas: *allVirtualSchemas}, nil
}

func readExaAllScriptTable(tx *sql.Tx, schemaName string) (*ExaScriptTable, error) {
	result, err := tx.Query(`
SELECT SCRIPT_SCHEMA, SCRIPT_NAME, SCRIPT_TYPE, SCRIPT_INPUT_TYPE, SCRIPT_RESULT_TYPE, SCRIPT_TEXT, SCRIPT_COMMENT
FROM SYS.EXA_ALL_SCRIPTS
WHERE SCRIPT_SCHEMA=?`, schemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to read SYS.EXA_ALL_SCRIPTS: %w", err)
	}
	rows := make([]ExaScriptRow, 0)
	for result.Next() {
		var row ExaScriptRow
		var inputType sql.NullString
		var resultType sql.NullString
		var comment sql.NullString
		err := result.Scan(&row.Schema, &row.Name, &row.Type, &inputType, &resultType, &row.Text, &comment)
		if err != nil {
			return nil, fmt.Errorf("failed to read row of SYS.EXA_ALL_SCRIPTS: %w", err)
		}
		row.InputType = inputType.String
		row.ResultType = resultType.String
		row.Comment = comment.String
		rows = append(rows, row)
	}
	return &ExaScriptTable{Rows: rows}, nil
}

func readExaAllVirtualSchemasTable(tx *sql.Tx) (*ExaVirtualSchemasTable, error) {
	result, err := tx.Query(`
SELECT SCHEMA_NAME, SCHEMA_OWNER, ADAPTER_SCRIPT, ADAPTER_NOTES
FROM SYS.EXA_ALL_VIRTUAL_SCHEMAS`)
	if err != nil {
		return nil, fmt.Errorf("failed to read SYS.EXA_ALL_VIRTUAL_SCHEMAS: %w", err)
	}
	rows := make([]ExaVirtualSchemaRow, 0)
	for result.Next() {
		var row ExaVirtualSchemaRow
		var adapterScriptSchemaAndName string
		err := result.Scan(&row.Name, &row.Owner, &adapterScriptSchemaAndName, &row.AdapterNotes)
		if err != nil {
			return nil, fmt.Errorf("failed to read row of SYS.EXA_ALL_VIRTUAL_SCHEMAS: %w", err)
		}
		row.AdapterScriptSchema, row.AdapterScriptName, err = extractSchemaAndName(adapterScriptSchemaAndName)
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}
	return &ExaVirtualSchemasTable{Rows: rows}, nil
}

func extractSchemaAndName(schemaAndName string) (schema, name string, err error) {
	i := strings.Index(schemaAndName, ".")
	if i < 0 {
		return "", "", fmt.Errorf("invalid format for adapter script: %q", schemaAndName)
	}
	return schemaAndName[:i], schemaAndName[i+1:], nil
}

type ExaScriptTable struct {
	Rows []ExaScriptRow `json:"rows"`
}

type ExaScriptRow struct {
	Schema     string `json:"schema"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	InputType  string `json:"inputType"`
	ResultType string `json:"resultType"`
	Text       string `json:"text"`
	Comment    string `json:"comment"`
}

type ExaVirtualSchemasTable struct {
	Rows []ExaVirtualSchemaRow `json:"rows"`
}

type ExaVirtualSchemaRow struct {
	Name                string `json:"name"`
	Owner               string `json:"owner"`
	AdapterScriptSchema string `json:"adapterScriptSchema"`
	AdapterScriptName   string `json:"adapterScriptName"`
	AdapterNotes        string `json:"adapterNotes"`
}
