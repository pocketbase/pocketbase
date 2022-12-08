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
	subGroup.GET("/requests", api.requestsList)
	subGroup.GET("/requests/stats", api.requestsStats)
	subGroup.GET("/requests/:id", api.requestView)
}

type logsApi struct {
	app core.App
}

var requestFilterFields = []string{
	"rowid", "id", "created", "updated",
	"url", "method", "status", "auth",
	"remoteIp", "userIp", "referer", "userAgent",
}

func (api *logsApi) requestsList(c echo.Context) error {
	fieldResolver := search.NewSimpleFieldResolver(requestFilterFields...)

	result, err := search.NewProvider(fieldResolver).
		Query(api.app.LogsDao().RequestQuery()).
		ParseAndExec(c.QueryParams().Encode(), &[]*models.Request{})

	if err != nil {
		return NewBadRequestError("", err)
	}

	return c.JSON(http.StatusOK, result)
}

func (api *logsApi) requestsStats(c echo.Context) error {
	fieldResolver := search.NewSimpleFieldResolver(requestFilterFields...)

	filter := c.QueryParam(search.FilterQueryParam)

	var expr dbx.Expression
	if filter != "" {
		var err error
		expr, err = search.FilterData(filter).BuildExpr(fieldResolver)
		if err != nil {
			return NewBadRequestError("Invalid filter format.", err)
		}
	}

	stats, err := api.app.LogsDao().RequestsStats(expr)
	if err != nil {
		return NewBadRequestError("Failed to generate requests stats.", err)
	}

	return c.JSON(http.StatusOK, stats)
}

func (api *logsApi) requestView(c echo.Context) error {
	id := c.PathParam("id")
	if id == "" {
		return NewNotFoundError("", nil)
	}

	request, err := api.app.LogsDao().FindRequestById(id)
	if err != nil || request == nil {
		return NewNotFoundError("", err)
	}

	return c.JSON(http.StatusOK, request)
}
