package apis

import (
	"context"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"
)

// bindSQLConsoleApi registers the SQL console api endpoint.
func bindSQLConsoleApi(app core.App, rg *router.RouterGroup[*core.RequestEvent]) {
	subGroup := rg.Group("/sql").Bind(RequireSuperuserAuth())
	subGroup.POST("", sqlExecute)
	subGroup.GET("/schema", sqlSchema)
}

type sqlExecuteRequest struct {
	SQL        string `json:"sql"`
	AllowWrite bool   `json:"allowWrite"`
}

type sqlExecuteResponse struct {
	Results      []map[string]any `json:"results,omitempty"`
	RowsAffected int64            `json:"rowsAffected,omitempty"`
	IsWrite      bool             `json:"isWrite"`
	ExecutionMs  int64            `json:"executionMs"`
}

var writeQueryPattern = regexp.MustCompile(`(?i)^\s*(INSERT|UPDATE|DELETE|CREATE|DROP|ALTER|TRUNCATE|REPLACE)\s`)
var readQueryPattern = regexp.MustCompile(`(?i)^\s*(SELECT|PRAGMA|EXPLAIN)\s`)

// isWriteQuery determines if a SQL query is a write operation
func isWriteQuery(sql string) bool {
	trimmed := strings.TrimSpace(sql)
	return writeQueryPattern.MatchString(trimmed) && !readQueryPattern.MatchString(trimmed)
}

func sqlExecute(e *core.RequestEvent) error {
	var req sqlExecuteRequest
	if err := e.BindBody(&req); err != nil {
		return e.BadRequestError("Invalid request body.", err)
	}

	if req.SQL == "" {
		return e.BadRequestError("SQL query is required.", nil)
	}

	// Check if it's a write query
	isWrite := isWriteQuery(req.SQL)

	// Block write queries if not explicitly allowed
	if isWrite && !req.AllowWrite {
		return e.ForbiddenError("Write queries are not allowed in read-only mode. Enable write mode to execute this query.", nil)
	}

	// Create a context with 30-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	startTime := time.Now()
	var resp sqlExecuteResponse
	resp.IsWrite = isWrite

	// Execute the query
	if isWrite {
		// Use nonconcurrent DB for write queries
		result, err := e.App.NonconcurrentDB().NewQuery(req.SQL).WithContext(ctx).Execute()
		if err != nil {
			return e.BadRequestError("SQL execution failed: "+err.Error(), err)
		}

		rowsAffected, _ := result.RowsAffected()
		resp.RowsAffected = rowsAffected
	} else {
		// Use concurrent DB for read queries
		rows, err := e.App.ConcurrentDB().NewQuery(req.SQL).WithContext(ctx).Rows()
		if err != nil {
			return e.BadRequestError("SQL execution failed: "+err.Error(), err)
		}
		defer rows.Close()

		// Get column names
		columns, err := rows.Columns()
		if err != nil {
			return e.BadRequestError("Failed to get column names: "+err.Error(), err)
		}

		// Parse results
		results := []map[string]any{}
		for rows.Next() {
			// Create a slice of interface{} to hold each column's value
			values := make([]any, len(columns))
			valuePtrs := make([]any, len(columns))
			for i := range columns {
				valuePtrs[i] = &values[i]
			}

			if err := rows.Scan(valuePtrs...); err != nil {
				return e.BadRequestError("Failed to scan row: "+err.Error(), err)
			}

			// Create a map for this row
			rowMap := make(map[string]any)
			for i, col := range columns {
				val := values[i]

				// Handle different types and convert byte arrays to strings
				switch v := val.(type) {
				case []byte:
					rowMap[col] = string(v)
				case nil:
					rowMap[col] = nil
				default:
					rowMap[col] = v
				}
			}
			results = append(results, rowMap)
		}

		if err := rows.Err(); err != nil {
			return e.BadRequestError("Error iterating rows: "+err.Error(), err)
		}

		resp.Results = results
	}

	resp.ExecutionMs = time.Since(startTime).Milliseconds()

	return e.JSON(http.StatusOK, resp)
}

type tableInfo struct {
	Name    string       `json:"name"`
	Type    string       `json:"type"` // "table" or "view"
	Columns []columnInfo `json:"columns"`
}

type columnInfo struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	NotNull      bool   `json:"notNull"`
	DefaultValue string `json:"defaultValue"`
	PrimaryKey   bool   `json:"primaryKey"`
}

func sqlSchema(e *core.RequestEvent) error {
	// Get all tables and views from sqlite_master
	type schemaRow struct {
		Type string `db:"type"`
		Name string `db:"name"`
	}

	var schemaRows []schemaRow
	err := e.App.DB().
		Select("type", "name").
		From("sqlite_master").
		AndWhere(dbx.In("type", "table", "view")).
		AndWhere(dbx.Not(dbx.Like("name", "sqlite_%"))).
		OrderBy("type ASC", "name ASC").
		All(&schemaRows)

	if err != nil {
		return e.BadRequestError("Failed to retrieve schema.", err)
	}

	tables := make([]tableInfo, 0, len(schemaRows))

	for _, row := range schemaRows {
		table := tableInfo{
			Name:    row.Name,
			Type:    row.Type,
			Columns: []columnInfo{},
		}

		// Get table info (columns)
		info, err := e.App.TableInfo(row.Name)
		if err != nil {
			// Skip tables we can't read
			continue
		}

		for _, col := range info {
			defaultVal := ""
			if col.DefaultValue.Valid {
				defaultVal = col.DefaultValue.String
			}

			table.Columns = append(table.Columns, columnInfo{
				Name:         col.Name,
				Type:         col.Type,
				NotNull:      col.NotNull,
				DefaultValue: defaultVal,
				PrimaryKey:   col.PK > 0,
			})
		}

		tables = append(tables, table)
	}

	return e.JSON(http.StatusOK, tables)
}
