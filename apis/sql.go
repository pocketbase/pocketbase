package apis

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"
)

const (
	runSQLMaxRows    = 1000
	runSQLMaxTimeout = 3 * time.Minute
)

// bindSQLApi registers the SQL api endpoints.
func bindSQLApi(app core.App, rg *router.RouterGroup[*core.RequestEvent]) {
	subGroup := rg.Group("/sql").Bind(RequireSuperuserAuth())
	subGroup.POST("", runSQL)
}

func runSQL(e *core.RequestEvent) error {
	// extra precaution in case manually invoked from somewhere else
	if !e.HasSuperuserAuth() {
		return e.ForbiddenError("", nil)
	}

	form := runSQLForm{}

	err := e.BindBody(&form)
	if err != nil {
		return firstApiError(err, e.BadRequestError("An error occurred while loading the submitted data.", err))
	}

	err = form.validate()
	if err != nil {
		return firstApiError(err, e.BadRequestError("An error occurred while validating the submitted data.", err))
	}

	result, err := executeQuery(e.App, form.Query, runSQLMaxRows)
	if err != nil {
		return firstApiError(err, e.BadRequestError("Failed to execute query. Raw error:\n"+err.Error(), nil))
	}

	return e.JSON(http.StatusOK, result)
}

type runSQLForm struct {
	Query string `form:"query" json:"query"`
}

func (form *runSQLForm) validate() error {
	return validation.ValidateStruct(form,
		validation.Field(&form.Query, validation.Required, validation.Length(0, 3000)),
	)
}

type runSQLResultColumn struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Nullable bool   `json:"nullable"`
}

type runSQLResult struct {
	ExecTime     int64                `json:"execTime"`
	AffectedRows int64                `json:"affectedRows"`
	Columns      []runSQLResultColumn `json:"columns"`
	Rows         [][]any              `json:"rows"`
}

var knownWriteQueryPrefixes = []string{"INSERT", "CREATE", "UPDATE", "DELETE", "DROP", "DETACH"}

func executeQuery(app core.App, query string, maxRows int) (*runSQLResult, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		// see https://github.com/mattn/go-sqlite3/issues/950
		return nil, errors.New("empty query")
	}

	var isPossibleWriteQuery bool

	// loosely check the query type
	ucQuery := strings.ToUpper(query)
	if !strings.HasPrefix(ucQuery, "SELECT") {
		for _, prefix := range knownWriteQueryPrefixes {
			if strings.HasPrefix(ucQuery, prefix) {
				isPossibleWriteQuery = true
				break
			}
		}
	}

	// note: don't extend the request context to minimize the risk of
	// causing integrity issues with custom non-transaction mutations
	ctx, cancelFunc := context.WithTimeout(context.Background(), runSQLMaxTimeout)
	defer cancelFunc()

	result := &runSQLResult{
		// init empty slices to ensure "[]" serialization
		Columns: []runSQLResultColumn{},
		Rows:    [][]any{},
	}

	now := time.Now()
	defer func() {
		result.ExecTime = time.Since(now).Milliseconds()
	}()

	// assume write/mutation query
	// ---------------------------------------------------------------
	if isPossibleWriteQuery {
		// auto wrap in transaction in case there are multiple inline queries
		txErr := app.RunInTransaction(func(txApp core.App) error {
			execResult, err := txApp.NonconcurrentDB().NewQuery(query).WithContext(ctx).Execute()
			if err != nil {
				return err
			}

			result.AffectedRows, err = execResult.RowsAffected()
			if err != nil {
				// non-critical error (e.g. not supported by the driver)
				txApp.Logger().Debug("Unable to fetch affected rows", slog.String("error", err.Error()))
			}

			return nil
		})
		if txErr != nil {
			return nil, txErr
		}

		return result, nil
	}

	// assume query returning rows
	// ---------------------------------------------------------------
	rows, err := app.ConcurrentDB().NewQuery(query).WithContext(ctx).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// populate columns info
	// ---
	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	for _, colType := range colTypes {
		col := runSQLResultColumn{
			Name: colType.Name(),
			Type: colType.DatabaseTypeName(),
		}
		col.Nullable, _ = colType.Nullable()

		result.Columns = append(result.Columns, col)
	}

	// populate rows
	// ---
	for rows.Next() {
		if len(result.Rows) >= maxRows {
			break
		}

		rowData := make([]any, len(colTypes))
		for i := 0; i < len(colTypes); i++ {
			var v *string
			rowData[i] = &v
		}

		err := rows.Scan(rowData...)
		if err != nil {
			return nil, err
		}

		result.Rows = append(result.Rows, rowData)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, nil
}
