package exaMetadata

import (
	"database/sql"
	"fmt"
	"strings"
)

type ExaMetadataReader interface {
	ReadMetadataTables(tx *sql.Tx, schemaName string) (*ExaMetadata, error)
	GetScriptByName(tx *sql.Tx, schemaName, scriptName string) (*ExaScriptRow, error)
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
/* [impl -> dsn~extension-components~1]. */
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

/* [impl -> dsn~extension-context-metadata~1] */
func (r *metaDataReaderImpl) GetScriptByName(tx *sql.Tx, schemaName, scriptName string) (*ExaScriptRow, error) {
	result, err := tx.Query(`
SELECT SCRIPT_SCHEMA, SCRIPT_NAME, SCRIPT_TYPE, SCRIPT_INPUT_TYPE, SCRIPT_RESULT_TYPE, SCRIPT_TEXT, SCRIPT_COMMENT
FROM SYS.EXA_ALL_SCRIPTS
WHERE SCRIPT_SCHEMA=? AND SCRIPT_NAME=?`, schemaName, scriptName)
	if err != nil {
		return nil, fmt.Errorf("failed to read SYS.EXA_ALL_SCRIPTS: %w", err)
	}
	defer result.Close()
	if !result.Next() {
		return nil, fmt.Errorf("no script found in schema %q for name %q", schemaName, scriptName)
	}
	row, err := readScriptRow(result)
	if err != nil {
		return nil, err
	}
	return row, nil
}

func readExaAllScriptTable(tx *sql.Tx, schemaName string) (*ExaScriptTable, error) {
	result, err := tx.Query(`
SELECT SCRIPT_SCHEMA, SCRIPT_NAME, SCRIPT_TYPE, SCRIPT_INPUT_TYPE, SCRIPT_RESULT_TYPE, SCRIPT_TEXT, SCRIPT_COMMENT
FROM SYS.EXA_ALL_SCRIPTS
WHERE SCRIPT_SCHEMA=?`, schemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to read SYS.EXA_ALL_SCRIPTS: %w", err)
	}
	defer result.Close()
	rows := make([]ExaScriptRow, 0)
	for result.Next() {
		if result.Err() != nil {
			return nil, fmt.Errorf("failed to iterate SYS.EXA_ALL_SCRIPTS: %w", result.Err())
		}
		row, err := readScriptRow(result)
		if err != nil {
			return nil, err
		}
		rows = append(rows, *row)
	}
	return &ExaScriptTable{Rows: rows}, nil
}

func readScriptRow(result *sql.Rows) (*ExaScriptRow, error) {
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
	return &row, nil
}

func readExaAllVirtualSchemasTable(tx *sql.Tx) (*ExaVirtualSchemasTable, error) {
	majorVersion, err := getExasolMajorVersion(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to find db version: %w", err)
	}
	if majorVersion == "8" {
		return readExaAllVirtualSchemasTableV8(tx)
	} else {
		return readExaAllVirtualSchemasTableV7(tx)
	}
}

func getExasolMajorVersion(tx *sql.Tx) (string, error) {
	result, err := tx.Query("SELECT PARAM_VALUE FROM SYS.EXA_METADATA WHERE PARAM_NAME='databaseMajorVersion'")
	if err != nil {
		return "", fmt.Errorf("querying exasol version failed: %w", err)
	}
	defer result.Close()
	if !result.Next() {
		if result.Err() != nil {
			return "", fmt.Errorf("failed to iterate exasol version: %w", result.Err())
		}
		return "", fmt.Errorf("no result found for exasol version query")
	}
	var majorVersion string
	err = result.Scan(&majorVersion)
	if err != nil {
		return "", fmt.Errorf("failed to read exasol version result: %w", err)
	}
	return majorVersion, nil
}

// This reads virtual schemas from the metadata tables of Exasol version 8.
func readExaAllVirtualSchemasTableV8(tx *sql.Tx) (*ExaVirtualSchemasTable, error) {
	result, err := tx.Query(`
SELECT SCHEMA_NAME, SCHEMA_OWNER, ADAPTER_SCRIPT_SCHEMA, ADAPTER_SCRIPT_NAME, ADAPTER_NOTES
FROM SYS.EXA_ALL_VIRTUAL_SCHEMAS`)
	if err != nil {
		return nil, fmt.Errorf("failed to read SYS.EXA_ALL_VIRTUAL_SCHEMAS: %w", err)
	}
	defer result.Close()
	rows := make([]ExaVirtualSchemaRow, 0)
	for result.Next() {
		if result.Err() != nil {
			return nil, fmt.Errorf("failed to iterate SYS.EXA_ALL_VIRTUAL_SCHEMAS: %w", result.Err())
		}
		var row ExaVirtualSchemaRow
		err := result.Scan(&row.Name, &row.Owner, &row.AdapterScriptSchema, &row.AdapterScriptName, &row.AdapterNotes)
		if err != nil {
			return nil, fmt.Errorf("failed to read row of SYS.EXA_ALL_VIRTUAL_SCHEMAS: %w", err)
		}
		rows = append(rows, row)
	}
	return &ExaVirtualSchemasTable{Rows: rows}, nil
}

// This reads virtual schemas from the metadata tables of Exasol version 7.1.
func readExaAllVirtualSchemasTableV7(tx *sql.Tx) (*ExaVirtualSchemasTable, error) {
	result, err := tx.Query(`
SELECT SCHEMA_NAME, SCHEMA_OWNER, ADAPTER_SCRIPT, ADAPTER_NOTES
FROM SYS.EXA_ALL_VIRTUAL_SCHEMAS`)
	if err != nil {
		return nil, fmt.Errorf("failed to read SYS.EXA_ALL_VIRTUAL_SCHEMAS: %w", err)
	}
	defer result.Close()
	rows := make([]ExaVirtualSchemaRow, 0)
	for result.Next() {
		if result.Err() != nil {
			return nil, fmt.Errorf("failed to iterate SYS.EXA_ALL_VIRTUAL_SCHEMAS: %w", result.Err())
		}
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
