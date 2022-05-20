package backend

import (
	"database/sql"
	"fmt"
)

type ExasolSqlClient struct {
	Connection *sql.DB
}

func (client ExasolSqlClient) RunQuery(query string) {
	_, err := client.Connection.Exec(query)
	if err != nil {
		panic(fmt.Sprintf("sql error in extension: %v", err.Error()))
	}
}
