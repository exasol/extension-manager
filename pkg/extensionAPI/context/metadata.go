package context

import (
	"database/sql"
	"fmt"

	"github.com/exasol/extension-manager/pkg/extensionAPI/exaMetadata"
)

// MetadataContext contains methods for reading Exasol metadata tables.
type MetadataContext interface {
	// Get a row from the SYS.EXA_ALL_SCRIPTS table for the given schema and script name.
	//
	// Returns `nil` when no script exists with the given name.
	// The JS runtime will convert `nil` to `null` in JavaScript code, so extensions can
	// check if a script was found by testing the result with `=== null`.
	GetScriptByName(name string) *exaMetadata.ExaScriptRow
}

type metadataContextImpl struct {
	metadataReader exaMetadata.ExaMetadataReader
	schemaName     string
	transaction    *sql.Tx
}

/* [impl -> dsn~extension-context-metadata~1]. */
func (m *metadataContextImpl) GetScriptByName(scriptName string) *exaMetadata.ExaScriptRow {
	script, err := m.metadataReader.GetScriptByName(m.transaction, m.schemaName, scriptName)
	if err != nil {
		reportError(fmt.Errorf("failed to find script %q.%q. Caused by: %w", m.schemaName, scriptName, err))
	}
	return script
}
