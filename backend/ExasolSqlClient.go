package backend

import (
	"database/sql"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

type ExasolSqlClient struct {
	transaction *sql.Tx
}

func NewSqlClient(tx *sql.Tx) *ExasolSqlClient {
	return &ExasolSqlClient{transaction: tx}
}

func (c ExasolSqlClient) RunQuery(query string) {
	err := validateQuery(query)
	if err != nil {
		reportError(err)
	}
	result, err := c.transaction.Exec(query)
	if err != nil {
		reportError(fmt.Errorf("error executing statement %q: %v", query, err))
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		reportError(fmt.Errorf("error getting rows affected for statement %q: %v", query, err))
	}
	log.Printf("Executed statement %q: rows affected: %d", query, rowsAffected)
}

var transactionStatements = []string{"commit", "rollback"}

func validateQuery(originalQuery string) error {
	query := strings.ToLower(originalQuery)
	query = strings.Trim(query, "\t\r\n ;")
	for _, forbiddenStatement := range transactionStatements {
		if strings.ToLower(forbiddenStatement) == query {
			return fmt.Errorf("statement %q contains forbidden command %q. Transaction handling is done by extension manager", originalQuery, forbiddenStatement)
		}
	}
	return nil
}

func reportError(err error) {
	// Panic to signal a failure to the JavaScript extension code.
	panic(err)
}
