package exaMetadata

import (
	"database/sql"
	"fmt"
)

// ExaMetadataReader allows accessing the Exasol metadata tables.
type ExaMetadataReader interface {
	// ReadMetadataTables reads all metadata tables.
	ReadMetadataTables(tx *sql.Tx, schemaName string) (*ExaMetadata, error)

	// GetScriptByName gets a row from the SYS.EXA_ALL_SCRIPTS table for the given schema and script name.
	//
	// Returns `(nil, nil)` when no script exists with the given name.
	GetScriptByName(tx *sql.Tx, schemaName, scriptName string) (*ExaScriptRow, error)
}

type ExaMetadata struct {
	AllScripts        ExaScriptTable         `json:"allScripts"`
	AllVirtualSchemas ExaVirtualSchemasTable `json:"allVirtualSchemas"`
}

// CreateExaMetaDataReader creates a new ExaMetadataReader for the Exasol meta data schema SYS.
func CreateExaMetaDataReader() ExaMetadataReader {
	return CreateExaMetaDataReaderForCustomMetadataSchema("SYS")
}

// CreateExaMetaDataReaderForCustomMetadataSchema creates a new ExaMetadataReader for the given Exasol meta data schema.
// This is only used for integration tests.
func CreateExaMetaDataReaderForCustomMetadataSchema(metaDataSchema string) ExaMetadataReader {
	return &metaDataReaderImpl{metaDataSchema: metaDataSchema}
}

type metaDataReaderImpl struct {
	metaDataSchema string
}

// ReadMetadataTables reads the metadata tables of the given schema.
/* [impl -> dsn~extension-components~1]. */
func (r *metaDataReaderImpl) ReadMetadataTables(tx *sql.Tx, schemaName string) (*ExaMetadata, error) {
	allScripts, err := r.readExaAllScriptTable(tx, schemaName)
	if err != nil {
		return nil, err
	}
	allVirtualSchemas, err := r.readExaAllVirtualSchemasTable(tx)
	if err != nil {
		return nil, err
	}
	return &ExaMetadata{AllScripts: *allScripts, AllVirtualSchemas: *allVirtualSchemas}, nil
}

/* [impl -> dsn~extension-context-metadata~1]. */
func (r *metaDataReaderImpl) GetScriptByName(tx *sql.Tx, schemaName, scriptName string) (*ExaScriptRow, error) {
	query := fmt.Sprintf(`
SELECT SCRIPT_SCHEMA, SCRIPT_NAME, SCRIPT_TYPE, SCRIPT_INPUT_TYPE, SCRIPT_RESULT_TYPE, SCRIPT_TEXT, SCRIPT_COMMENT
FROM %s.EXA_ALL_SCRIPTS
WHERE SCRIPT_SCHEMA=? AND SCRIPT_NAME=?`, r.metaDataSchema)
	result, err := tx.Query(query, schemaName, scriptName)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s.EXA_ALL_SCRIPTS: %w", r.metaDataSchema, err)
	}
	defer result.Close()
	if !result.Next() {
		return nil, nil
	}
	row, err := readScriptRow(result)
	if err != nil {
		return nil, err
	}
	return row, nil
}

func (r *metaDataReaderImpl) readExaAllScriptTable(tx *sql.Tx, schemaName string) (*ExaScriptTable, error) {
	query := fmt.Sprintf(`
SELECT SCRIPT_SCHEMA, SCRIPT_NAME, SCRIPT_TYPE, SCRIPT_INPUT_TYPE, SCRIPT_RESULT_TYPE, SCRIPT_TEXT, SCRIPT_COMMENT
FROM %s.EXA_ALL_SCRIPTS
WHERE SCRIPT_SCHEMA=?`, r.metaDataSchema)
	result, err := tx.Query(query, schemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s.EXA_ALL_SCRIPTS: %w", r.metaDataSchema, err)
	}
	defer result.Close()
	rows := make([]ExaScriptRow, 0)
	for result.Next() {
		if result.Err() != nil {
			return nil, fmt.Errorf("failed to iterate %s.EXA_ALL_SCRIPTS: %w", r.metaDataSchema, result.Err())
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
	var schema sql.NullString
	var name sql.NullString
	var scriptType sql.NullString
	var inputType sql.NullString
	var resultType sql.NullString
	var text sql.NullString
	var comment sql.NullString
	err := result.Scan(&schema, &name, &scriptType, &inputType, &resultType, &text, &comment)
	if err != nil {
		return nil, fmt.Errorf("failed to read row of EXA_ALL_SCRIPTS: %w", err)
	}
	row := ExaScriptRow{
		Schema:     schema.String,
		Name:       name.String,
		Type:       scriptType.String,
		InputType:  inputType.String,
		ResultType: resultType.String,
		Text:       text.String,
		Comment:    comment.String,
	}
	return &row, nil
}

func (r *metaDataReaderImpl) readExaAllVirtualSchemasTable(tx *sql.Tx) (*ExaVirtualSchemasTable, error) {
	query := fmt.Sprintf(`
SELECT SCHEMA_NAME, SCHEMA_OWNER, ADAPTER_SCRIPT_SCHEMA, ADAPTER_SCRIPT_NAME, ADAPTER_NOTES
FROM %s.EXA_ALL_VIRTUAL_SCHEMAS`, r.metaDataSchema)
	result, err := tx.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s.EXA_ALL_VIRTUAL_SCHEMAS: %w", r.metaDataSchema, err)
	}
	defer result.Close()
	rows := make([]ExaVirtualSchemaRow, 0)
	for result.Next() {
		if result.Err() != nil {
			return nil, fmt.Errorf("failed to iterate %s.EXA_ALL_VIRTUAL_SCHEMAS: %w", r.metaDataSchema, result.Err())
		}
		row, err := readVirtualSchemaRow(result)
		if err != nil {
			return nil, err
		}
		rows = append(rows, *row)
	}
	return &ExaVirtualSchemasTable{Rows: rows}, nil
}

func readVirtualSchemaRow(result *sql.Rows) (*ExaVirtualSchemaRow, error) {
	var name sql.NullString
	var owner sql.NullString
	var adapterScriptSchema sql.NullString
	var adapterScriptName sql.NullString
	var adapterNotes sql.NullString
	err := result.Scan(&name, &owner, &adapterScriptSchema, &adapterScriptName, &adapterNotes)
	if err != nil {
		return nil, fmt.Errorf("failed to read row of EXA_ALL_VIRTUAL_SCHEMAS: %w", err)
	}
	row := ExaVirtualSchemaRow{
		Name:                name.String,
		Owner:               owner.String,
		AdapterScriptSchema: adapterScriptSchema.String,
		AdapterScriptName:   adapterScriptName.String,
		AdapterNotes:        adapterNotes.String,
	}
	return &row, nil
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
