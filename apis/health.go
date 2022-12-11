package apis

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"
)

// bindHealthApi registers the health api endpoint.
func bindHealthApi(app core.App, rg *echo.Group) {
	api := healthApi{app: app}

	subGroup := rg.Group("/health")
	subGroup.GET("", api.healthCheck)
}

type healthApi struct {
	app core.App
}

// healthCheck returns a 200 OK response if the server is healthy.
func (api *healthApi) healthCheck(c echo.Context) error {
	payload := map[string]any{
		"code":    http.StatusOK,
		"message": "API is healthy.",
	}

	return c.JSON(http.StatusOK, payload)
}
