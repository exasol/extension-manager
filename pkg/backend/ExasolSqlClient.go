package backend

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// SimpleSQLClient allows extensions to execute statements and queries against the database.
type SimpleSQLClient interface {
	// Execute runs a query that does not return rows, e.g. INSERT or UPDATE.
	Execute(query string, args ...any) (sql.Result, error)

	// Query runs a query that returns rows, typically a SELECT.
	Query(query string, args ...any) (*QueryResult, error)
}

type exasolSqlClient struct {
	transaction *sql.Tx
	ctx         context.Context
}

// NewSqlClient creates a new [SimpleSQLClient].
func NewSqlClient(ctx context.Context, tx *sql.Tx) SimpleSQLClient {
	return &exasolSqlClient{ctx: ctx, transaction: tx}
}

// Execute executes a statement like `CREATE VIRTUAL SCHEMA`.
/* [impl -> dsn~extension-context-sql-client~1]. */
func (c *exasolSqlClient) Execute(query string, args ...any) (sql.Result, error) {
	err := validateQuery(query)
	if err != nil {
		return nil, err
	}

	result, err := c.transaction.ExecContext(c.ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing statement %q: %w", query, err)
	}
	return result, nil
}

// Query runs a query like `SELECT` and returns the result.
/* [impl -> dsn~extension-context-sql-client~1]. */
func (c *exasolSqlClient) Query(query string, args ...any) (result *QueryResult, errResult error) {
	err := validateQuery(query)
	if err != nil {
		return nil, err
	}
	rows, err := c.transaction.QueryContext(c.ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing statement %q: %w", query, err)
	}
	defer func() {
		if err := closeRows(rows); err != nil {
			errResult = err
			result = nil
		}
	}()
	result, err = c.extractResult(rows)
	if err != nil {
		return nil, fmt.Errorf("error reading result from statement %q: %w", query, err)
	}
	return result, nil
}

func closeRows(rows *sql.Rows) error {
	err := rows.Close()
	if err != nil {
		return fmt.Errorf("error closing result: %w", err)
	}
	if err = rows.Err(); err != nil {
		return fmt.Errorf("error while iterating result: %w", err)
	}
	return nil
}

func (c exasolSqlClient) extractResult(rows *sql.Rows) (*QueryResult, error) {
	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	cols := make([]Column, 0, len(colTypes))
	for _, c := range colTypes {
		cols = append(cols, Column{Name: c.Name(), TypeName: c.DatabaseTypeName()})
	}

	resultRows, err := extractRows(rows, len(cols))
	if err != nil {
		return nil, err
	}
	return &QueryResult{Columns: cols, Rows: resultRows}, nil
}

func extractRows(rows *sql.Rows, columnCount int) ([]Row, error) {
	resultRows := make([]Row, 0)
	values := make([]interface{}, columnCount)
	for rows.Next() {
		for i := range values {
			values[i] = new(interface{})
		}
		if err := rows.Scan(values...); err != nil {
			return nil, err
		}
		row := make([]any, 0, len(values))
		for _, v := range values {
			row = append(row, *v.(*interface{}))
		}
		resultRows = append(resultRows, row)
	}
	return resultRows, nil
}

// Result of a database query.
type QueryResult struct {
	Columns []Column `json:"columns"` // Column definitions of the query result
	Rows    []Row    `json:"rows"`    // The result rows
}

// Column definition of a query result.
type Column struct {
	Name     string `json:"name"`     // Column name
	TypeName string `json:"typeName"` // Column type
}

// Row of a database query result.
type Row []any

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
