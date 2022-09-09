package backend

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

type ExasolSqlClient struct {
	transaction *sql.Tx
	ctx         context.Context
}

func NewSqlClient(ctx context.Context, tx *sql.Tx) *ExasolSqlClient {
	return &ExasolSqlClient{ctx: ctx, transaction: tx}
}

func (c ExasolSqlClient) Execute(query string) {
	err := validateQuery(query)
	if err != nil {
		reportError(err)
	}
	result, err := c.transaction.ExecContext(c.ctx, query)
	if err != nil {
		reportError(fmt.Errorf("error executing statement %q: %v", query, err))
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		reportError(fmt.Errorf("error getting rows affected for statement %q: %v", query, err))
	}
	log.Printf("Executed statement %q: rows affected: %d", query, rowsAffected)
}

func (c ExasolSqlClient) Query(query string) Rows {
	// TODO
	return Rows{}
}

type Rows struct {
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
