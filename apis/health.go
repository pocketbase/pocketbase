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
	subGroup.HEAD("", api.healthCheck)
	subGroup.GET("", api.healthCheck)
}

type healthApi struct {
	app core.App
}

type healthCheckResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Data    struct {
		CanBackup bool `json:"canBackup"`
	} `json:"data"`
}

// healthCheck returns a 200 OK response if the server is healthy.
func (api *healthApi) healthCheck(c echo.Context) error {
	if c.Request().Method == http.MethodHead {
		return c.NoContent(http.StatusOK)
	}

	resp := new(healthCheckResponse)
	resp.Code = http.StatusOK
	resp.Message = "API is healthy."
	resp.Data.CanBackup = !api.app.Store().Has(core.StoreKeyActiveBackup)

	return c.JSON(http.StatusOK, resp)
}
