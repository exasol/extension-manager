package backend

import (
	"database/sql"
	"fmt"
	"log"
)

type ExasolSqlClient struct {
	Connection *sql.DB
}

func (client ExasolSqlClient) RunQuery(query string) {
	result, err := client.Connection.Exec(query)
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
