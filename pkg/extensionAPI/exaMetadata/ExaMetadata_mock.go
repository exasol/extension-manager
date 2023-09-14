package exaMetadata

import (
	"database/sql"

	"github.com/stretchr/testify/mock"
)

type ExaMetaDataReaderMock struct {
	mock.Mock
	extensionSchema string
}

func CreateExaMetaDataReaderMock(extensionSchema string) *ExaMetaDataReaderMock {
	//nolint:exhaustruct // Default value for Mock is fine
	return &ExaMetaDataReaderMock{extensionSchema: extensionSchema}
}

func (m *ExaMetaDataReaderMock) SimulateExaAllScripts(scripts []ExaScriptRow) {
	m.SimulateExaMetaData(ExaMetadata{
		AllScripts:        ExaScriptTable{Rows: scripts},
		AllVirtualSchemas: ExaVirtualSchemasTable{Rows: []ExaVirtualSchemaRow{}}})
}

func (m *ExaMetaDataReaderMock) SimulateExaMetaData(metaData ExaMetadata) {
	m.On("ReadMetadataTables", mock.Anything, m.extensionSchema).Return(&metaData, nil)
}

func (mock *ExaMetaDataReaderMock) ReadMetadataTables(tx *sql.Tx, schemaName string) (*ExaMetadata, error) {
	args := mock.Called(tx, schemaName)
	if buckets, ok := args.Get(0).(*ExaMetadata); ok {
		return buckets, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ExaMetaDataReaderMock) SimulateGetScriptByNameScriptText(scriptName string, scriptText string) {
	script := &ExaScriptRow{
		Schema:     "?",
		Name:       scriptName,
		Text:       scriptText,
		Type:       "",
		InputType:  "",
		ResultType: "",
		Comment:    "",
	}
	m.SimulateGetScriptByName(scriptName, script)
}

func (m *ExaMetaDataReaderMock) SimulateGetScriptByName(scriptName string, script *ExaScriptRow) {
	m.On("GetScriptByName", mock.Anything, m.extensionSchema, scriptName).Return(script, nil)
}

func (m *ExaMetaDataReaderMock) SimulateGetScriptByNameFails(scriptName string, err error) {
	m.On("GetScriptByName", mock.Anything, m.extensionSchema, scriptName).Return(nil, err)
}

func (mock *ExaMetaDataReaderMock) GetScriptByName(tx *sql.Tx, schemaName, scriptName string) (*ExaScriptRow, error) {
	args := mock.Called(tx, schemaName, scriptName)
	if script, ok := args.Get(0).(*ExaScriptRow); ok {
		return script, args.Error(1)
	}
	return nil, args.Error(1)
}
