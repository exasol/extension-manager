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
	log.Printf("Executing statement %q", query)
	_, err := client.Connection.Exec(query)
	if err != nil {
		// Panic to signal a failed query to the JavaScript extension code.
		panic(fmt.Sprintf("error executing statement %q: %v", query, err))
	}
}
