package backend

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type ExasolSqlClient struct {
	transaction *sql.Tx
	ctx         context.Context
}

func NewSqlClient(ctx context.Context, tx *sql.Tx) *ExasolSqlClient {
	return &ExasolSqlClient{ctx: ctx, transaction: tx}
}

func (c *ExasolSqlClient) Execute(query string, args ...any) {
	err := validateQuery(query)
	if err != nil {
		reportError(err)
	}

	_, err = c.transaction.ExecContext(c.ctx, query, args...)
	if err != nil {
		reportError(fmt.Errorf("error executing statement %q: %v", query, err))
	}
}

func (c *ExasolSqlClient) Query(query string, args ...any) QueryResult {
	err := validateQuery(query)
	if err != nil {
		reportError(err)
	}
	rows, err := c.transaction.QueryContext(c.ctx, query, args...)
	if err != nil {
		reportError(fmt.Errorf("error executing statement %q: %v", query, err))
	}
	defer closeRows(rows)
	result, err := c.extractResult(rows)
	if err != nil || result == nil {
		reportError(fmt.Errorf("error reading result from statement %q: %v", query, err))
	}
	return *result
}

func closeRows(rows *sql.Rows) {
	err := rows.Close()
	if err != nil {
		reportError(fmt.Errorf("error closing result: %v", err))
	}
	err = rows.Err()
	if err != nil {
		reportError(fmt.Errorf("error while iterating result: %v", err))
	}
}

func (c ExasolSqlClient) extractResult(rows *sql.Rows) (*QueryResult, error) {
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

type QueryResult struct {
	Columns []Column `json:"columns"`
	Rows    []Row    `json:"rows"`
}

type Column struct {
	Name     string `json:"name"`
	TypeName string `json:"typeName"`
}

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

func reportError(err error) {
	// Panic to signal a failure to the JavaScript extension code.
	panic(err)
}
