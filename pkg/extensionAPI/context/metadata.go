package context

import (
	"database/sql"

	"github.com/exasol/extension-manager/pkg/extensionAPI/exaMetadata"
)

type MetadataContext interface {
	GetScriptByName(name string) exaMetadata.ExaScriptRow
}

type metadataContextImpl struct {
	metadataReader exaMetadata.ExaMetadataReader
	schemaName     string
	transaction    *sql.Tx
}

/* [impl -> dsn~extension-context-metadata~1] */
func (m *metadataContextImpl) GetScriptByName(name string) exaMetadata.ExaScriptRow {
	script, err := m.metadataReader.GetScriptByName(m.transaction, m.schemaName, name)
	if err != nil {
		reportError(err)
	}
	return *script
}
