package exaMetadata

import (
	"database/sql"

	"github.com/stretchr/testify/mock"
)

type ExaMetaDataReaderMock struct {
	mock.Mock
	extensionSchema string
}

func CreateExaMetaDataReaderMock(extensionSchema string) ExaMetaDataReaderMock {
	var _ ExaMetadataReader = &ExaMetaDataReaderMock{}
	return ExaMetaDataReaderMock{extensionSchema: extensionSchema}
}

func (m *ExaMetaDataReaderMock) SimulateExaAllScripts(scripts []ExaScriptRow) {
	m.SimulateExaMetaData(ExaMetadata{AllScripts: ExaScriptTable{Rows: scripts}})
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
