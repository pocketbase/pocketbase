package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type Payload struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
}

func validatePayload(c echo.Context) (code int, errorMsg string, payload Payload) {
	// This function checks if the payload in request body is valid and returns the payload
	// Password and PasswordConfirm are compulsory fields and should match
	var r = c.Request()
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		return 400, "Invalid request payload", payload
	}
	if payload.Password == "" {
		return 400, "Password is required", payload
	}
	if payload.PasswordConfirm != payload.Password {
		return 400, "PasswordConfirm and Password don't match", payload
	}
	return 200, "", payload
}

func addRecordInUsersCollection(c echo.Context, payload Payload) (code int, errorMsg string) {
	// This function adds the payload to the users collection
	// This function returns error code and error message for the operation performed
	// The payload is added to the users collection using the pocketbase API
	recordURL := url.URL{
		Scheme: "http",
		Host:   c.Request().Host,
		Path:   "/api/collections/users/records",
	}
	recordEndpoint := recordURL.String()
	requestBody, err := json.Marshal(payload)
	if err != nil {
		return 500, "Internal server error"
	}
	req, err := http.NewRequest("POST", recordEndpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return 500, "Internal server error"
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 500, "Internal server error"
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 500, "Internal server error"
	}
	if resp.StatusCode != 200 {
		return resp.StatusCode, string(responseBody)
	}
	return 200, fmt.Sprintf("Valid payload received and record added %s ", string(responseBody))
}

func main() {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.AddRoute(echo.Route{
			Method: http.MethodPost,
			Path:   "/validatePayloadAndAddRecord",
			Handler: func(c echo.Context) error {
				// Validate payload in request body
				code, errorMsg, payload := validatePayload(c)
				if code != 200 {
					return c.String(code, errorMsg)
				}
				// Add record in users collection
				code, errorMsg = addRecordInUsersCollection(c, payload)
				return c.String(code, errorMsg)
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
