package backend

import (
	"database/sql"
	"fmt"
	"log"
)

type ExasolSqlClient struct {
	transaction *sql.Tx
}

func NewSqlClient(tx *sql.Tx) *ExasolSqlClient {
	return &ExasolSqlClient{transaction: tx}
}

func (c ExasolSqlClient) RunQuery(query string) {
	result, err := c.transaction.Exec(query)
	if err != nil {
		// Panic to signal a failed query to the JavaScript extension code.
		panic(fmt.Sprintf("error executing statement %q: %v", query, err))
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		panic(fmt.Sprintf("error getting rows affected for statement %q: %v", query, err))
	}
	log.Printf("Executed statement %q: rows affected: %d", query, rowsAffected)
}
