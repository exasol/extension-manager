package context

import "github.com/exasol/extension-manager/pkg/backend"

type ContextSqlClient interface {
	// Execute runs a query that does not return rows, e.g. INSERT or UPDATE.
	Execute(query string, args ...any)

	// Query runs a query that returns rows, typically a SELECT.
	Query(query string, args ...any) backend.QueryResult
}

type contextSqlClient struct {
	client backend.SimpleSQLClient
}

func (c *contextSqlClient) Execute(query string, args ...any) {
	_, err := c.client.Execute(query, args...)
	if err != nil {
		reportError(err)
	}
}

func (c *contextSqlClient) Query(query string, args ...any) backend.QueryResult {
	result, err := c.client.Query(query, args...)
	if err != nil {
		reportError(err)
	}
	return *result
}
