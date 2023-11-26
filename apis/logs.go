package apis

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/search"
)

// bindLogsApi registers the request logs api endpoints.
func bindLogsApi(app core.App, rg *echo.Group) {
	api := logsApi{app: app}

	subGroup := rg.Group("/logs", RequireAdminAuth())
	subGroup.GET("", api.list)
	subGroup.GET("/stats", api.stats)
	subGroup.GET("/:id", api.view)
}

type logsApi struct {
	app core.App
}

var logFilterFields = []string{
	"rowid", "id", "created", "updated",
	"level", "message", "data",
	`^data\.[\w\.\:]*\w+$`,
}

func (api *logsApi) list(c echo.Context) error {
	fieldResolver := search.NewSimpleFieldResolver(logFilterFields...)

	result, err := search.NewProvider(fieldResolver).
		Query(api.app.LogsDao().LogQuery()).
		ParseAndExec(c.QueryParams().Encode(), &[]*models.Log{})

	if err != nil {
		return NewBadRequestError("", err)
	}

	return c.JSON(http.StatusOK, result)
}

func (api *logsApi) stats(c echo.Context) error {
	fieldResolver := search.NewSimpleFieldResolver(logFilterFields...)

	filter := c.QueryParam(search.FilterQueryParam)

	var expr dbx.Expression
	if filter != "" {
		var err error
		expr, err = search.FilterData(filter).BuildExpr(fieldResolver)
		if err != nil {
			return NewBadRequestError("Invalid filter format.", err)
		}
	}

	stats, err := api.app.LogsDao().LogsStats(expr)
	if err != nil {
		return NewBadRequestError("Failed to generate logs stats.", err)
	}

	return c.JSON(http.StatusOK, stats)
}

func (api *logsApi) view(c echo.Context) error {
	id := c.PathParam("id")
	if id == "" {
		return NewNotFoundError("", nil)
	}

	log, err := api.app.LogsDao().FindLogById(id)
	if err != nil || log == nil {
		return NewNotFoundError("", err)
	}

	return c.JSON(http.StatusOK, log)
}
