package apis

import (
	"net/http"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"
	"github.com/pocketbase/pocketbase/tools/search"
)

// bindLogsApi registers the request logs api endpoints.
func bindLogsApi(app core.App, rg *router.RouterGroup[*core.RequestEvent]) {
	sub := rg.Group("/logs").Bind(
		RequireSuperuserAuth(),
		SkipSuccessActivityLog(),
	)

	sub.GET("", logsList)
	sub.DELETE("", logsTruncate)
	sub.GET("/stats", logsStats)
	sub.GET("/{id}", logsView)
}

var logFilterFields = []string{
	"id", "created", "level", "message", "data",
	`^data\.[\w\.\:]*\w+$`,
}

func logsList(e *core.RequestEvent) error {
	fieldResolver := search.NewSimpleFieldResolver(logFilterFields...)

	result, err := search.NewProvider(fieldResolver).
		Query(e.App.AuxModelQuery(&core.Log{})).
		ParseAndExec(e.Request.URL.Query().Encode(), &[]*core.Log{})

	if err != nil {
		return e.BadRequestError("", err)
	}

	return e.JSON(http.StatusOK, result)
}

func logsStats(e *core.RequestEvent) error {
	fieldResolver := search.NewSimpleFieldResolver(logFilterFields...)

	filter := e.Request.URL.Query().Get(search.FilterQueryParam)

	var expr dbx.Expression
	if filter != "" {
		var err error
		expr, err = search.FilterData(filter).BuildExpr(fieldResolver)
		if err != nil {
			return e.BadRequestError("Invalid filter format.", err)
		}
	}

	stats, err := e.App.LogsStats(expr)
	if err != nil {
		return e.BadRequestError("Failed to generate logs stats.", err)
	}

	return e.JSON(http.StatusOK, stats)
}

func logsView(e *core.RequestEvent) error {
	id := e.Request.PathValue("id")
	if id == "" {
		return e.NotFoundError("", nil)
	}

	log, err := e.App.FindLogById(id)
	if err != nil || log == nil {
		return e.NotFoundError("", err)
	}

	return e.JSON(http.StatusOK, log)
}

func logsTruncate(e *core.RequestEvent) error {
	// delete all rows directly (aka. no model hooks will be fired)
	_, err := e.App.AuxNonconcurrentDB().Delete((&core.Log{}).TableName(), nil).Execute()
	if err != nil {
		return e.InternalServerError("Failed to truncate all logs.", err)
	}

	// try to free the unused disk space
	err = e.App.AuxVacuum()
	if err != nil {
		e.App.Logger().Warn("Failed to VACUUM aux database", "error", err)
	}

	return e.NoContent(http.StatusNoContent)
}
