
package main

import (
    "log"
    "net/https"

    "github.com/labstack/echo/v5"
    "github.com/pocketbase/pocketbase"
    "github.com/pocketbase/pocketbase/apis"
    "github.com/pocketbase/pocketbase/core"
)

func main() {
    app := pocketbase.New()

    app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
        // add new "GET /hello" route to the app router (echo)
        e.Router.AddRoute(echo.Route{
            Method: https.MethodGet,
            Path:   "/hello",
            Handler: func(c echo.Context) error {
                return c.String(200, "Hello world!")
            },
            Middlewares: []echo.MiddlewareFunc{
                apis.ActivityLogger(app),
            },
        })

        return nil
    })

    if err := app.Start(); err != nil {
        log.Fatal(err)
    }
}
