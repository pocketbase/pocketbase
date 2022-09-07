package apis

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/rest"
	"github.com/pocketbase/pocketbase/tools/routine"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/spf13/cast"
)

// Common request context keys used by the middlewares and api handlers.
const (
	ContextUserKey       string = "user"
	ContextAdminKey      string = "admin"
	ContextCollectionKey string = "collection"
)

// RequireGuestOnly middleware requires a request to NOT have a valid
// Authorization header set.
//
// This middleware is the opposite of [apis.RequireAdminOrUserAuth()].
func RequireGuestOnly() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := rest.NewBadRequestError("The request can be accessed only by guests.", nil)

			user, _ := c.Get(ContextUserKey).(*models.User)
			if user != nil {
				return err
			}

			admin, _ := c.Get(ContextAdminKey).(*models.Admin)
			if admin != nil {
				return err
			}

			return next(c)
		}
	}
}

// RequireUserAuth middleware requires a request to have
// a valid user Authorization header set (aka. `Authorization: User ...`).
func RequireUserAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, _ := c.Get(ContextUserKey).(*models.User)
			if user == nil {
				return rest.NewUnauthorizedError("The request requires valid user authorization token to be set.", nil)
			}

			return next(c)
		}
	}
}

// RequireAdminAuth middleware requires a request to have
// a valid admin Authorization header set (aka. `Authorization: Admin ...`).
func RequireAdminAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			admin, _ := c.Get(ContextAdminKey).(*models.Admin)
			if admin == nil {
				return rest.NewUnauthorizedError("The request requires admin authorization token to be set.", nil)
			}

			return next(c)
		}
	}
}

// RequireAdminAuthOnlyIfAny middleware requires a request to have
// a valid admin Authorization header set (aka. `Authorization: Admin ...`)
// ONLY if the application has at least 1 existing Admin model.
func RequireAdminAuthOnlyIfAny(app core.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			totalAdmins, err := app.Dao().TotalAdmins()
			if err != nil {
				return rest.NewBadRequestError("Failed to fetch admins info.", err)
			}

			admin, _ := c.Get(ContextAdminKey).(*models.Admin)

			if admin != nil || totalAdmins == 0 {
				return next(c)
			}

			return rest.NewUnauthorizedError("The request requires admin authorization token to be set.", nil)
		}
	}
}

// RequireAdminOrUserAuth middleware requires a request to have
// a valid admin or user Authorization header set
// (aka. `Authorization: Admin ...` or `Authorization: User ...`).
//
// This middleware is the opposite of [apis.RequireGuestOnly()].
func RequireAdminOrUserAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			admin, _ := c.Get(ContextAdminKey).(*models.Admin)
			user, _ := c.Get(ContextUserKey).(*models.User)

			if admin == nil && user == nil {
				return rest.NewUnauthorizedError("The request requires admin or user authorization token to be set.", nil)
			}

			return next(c)
		}
	}
}

// RequireAdminOrOwnerAuth middleware requires a request to have
// a valid admin or user owner Authorization header set
// (aka. `Authorization: Admin ...` or `Authorization: User ...`).
//
// This middleware is similar to [apis.RequireAdminOrUserAuth()] but
// for the user token expects to have the same id as the path parameter
// `ownerIdParam` (default to "id").
func RequireAdminOrOwnerAuth(ownerIdParam string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if ownerIdParam == "" {
				ownerIdParam = "id"
			}

			ownerId := c.PathParam(ownerIdParam)
			admin, _ := c.Get(ContextAdminKey).(*models.Admin)
			loggedUser, _ := c.Get(ContextUserKey).(*models.User)

			if admin == nil && loggedUser == nil {
				return rest.NewUnauthorizedError("The request requires admin or user authorization token to be set.", nil)
			}

			if admin == nil && loggedUser.Id != ownerId {
				return rest.NewForbiddenError("You are not allowed to perform this request.", nil)
			}

			return next(c)
		}
	}
}

// LoadAuthContext middleware reads the Authorization request header
// and loads the token related user or admin instance into the
// request's context.
//
// This middleware is expected to be registered by default for all routes.
func LoadAuthContext(app core.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")

			if token != "" {
				if strings.HasPrefix(token, "User ") {
					user, err := app.Dao().FindUserByToken(
						token[5:],
						app.Settings().UserAuthToken.Secret,
					)
					if err == nil && user != nil {
						c.Set(ContextUserKey, user)
					}
				} else if strings.HasPrefix(token, "Admin ") {
					admin, err := app.Dao().FindAdminByToken(
						token[6:],
						app.Settings().AdminAuthToken.Secret,
					)
					if err == nil && admin != nil {
						c.Set(ContextAdminKey, admin)
					}
				}
			}

			return next(c)
		}
	}
}

// LoadCollectionContext middleware finds the collection with related
// path identifier and loads it into the request context.
func LoadCollectionContext(app core.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if param := c.PathParam("collection"); param != "" {
				collection, err := app.Dao().FindCollectionByNameOrId(param)
				if err != nil || collection == nil {
					return rest.NewNotFoundError("", err)
				}

				c.Set(ContextCollectionKey, collection)
			}

			return next(c)
		}
	}
}

// ActivityLogger middleware takes care to save the request information
// into the logs database.
//
// The middleware does nothing if the app logs retention period is zero
// (aka. app.Settings().Logs.MaxDays = 0).
func ActivityLogger(app core.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)

			// no logs retention
			if app.Settings().Logs.MaxDays == 0 {
				return err
			}

			httpRequest := c.Request()
			httpResponse := c.Response()
			status := httpResponse.Status
			meta := types.JsonMap{}

			if err != nil {
				switch v := err.(type) {
				case *echo.HTTPError:
					status = v.Code
					meta["errorMessage"] = v.Message
					meta["errorDetails"] = fmt.Sprint(v.Internal)
				case *rest.ApiError:
					status = v.Code
					meta["errorMessage"] = v.Message
					meta["errorDetails"] = fmt.Sprint(v.RawData())
				default:
					status = http.StatusBadRequest
					meta["errorMessage"] = v.Error()
				}
			}

			requestAuth := models.RequestAuthGuest
			if c.Get(ContextUserKey) != nil {
				requestAuth = models.RequestAuthUser
			} else if c.Get(ContextAdminKey) != nil {
				requestAuth = models.RequestAuthAdmin
			}

			ip, _, _ := net.SplitHostPort(httpRequest.RemoteAddr)

			model := &models.Request{
				Url:       httpRequest.URL.RequestURI(),
				Method:    strings.ToLower(httpRequest.Method),
				Status:    status,
				Auth:      requestAuth,
				UserIp:    realUserIp(httpRequest, ip),
				RemoteIp:  ip,
				Referer:   httpRequest.Referer(),
				UserAgent: httpRequest.UserAgent(),
				Meta:      meta,
			}
			// set timestamp fields before firing a new go routine
			model.RefreshCreated()
			model.RefreshUpdated()

			routine.FireAndForget(func() {
				attempts := 1

			BeginSave:
				logErr := app.LogsDao().SaveRequest(model)
				if logErr != nil {
					// try one more time after 10s in case of SQLITE_BUSY or "database is locked" error
					if attempts <= 2 {
						attempts++
						time.Sleep(10 * time.Second)
						goto BeginSave
					} else if app.IsDebug() {
						log.Println("Log save failed:", logErr)
					}
				}

				// Delete old request logs
				// ---
				now := time.Now()
				lastLogsDeletedAt := cast.ToTime(app.Cache().Get("lastLogsDeletedAt"))
				daysDiff := now.Sub(lastLogsDeletedAt).Hours() * 24

				if daysDiff > float64(app.Settings().Logs.MaxDays) {
					deleteErr := app.LogsDao().DeleteOldRequests(now.AddDate(0, 0, -1*app.Settings().Logs.MaxDays))
					if deleteErr == nil {
						app.Cache().Set("lastLogsDeletedAt", now)
					} else if app.IsDebug() {
						log.Println("Logs delete failed:", deleteErr)
					}
				}
			})

			return err
		}
	}
}

// Returns the "real" user IP from common proxy headers (or fallbackIp if none is found).
//
// The returned IP value shouldn't be trusted if not behind a trusted reverse proxy!
func realUserIp(r *http.Request, fallbackIp string) string {
	if ip := r.Header.Get("CF-Connecting-IP"); ip != "" {
		return ip
	}

	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}

	if ipsList := r.Header.Get("X-Forwarded-For"); ipsList != "" {
		ips := strings.Split(ipsList, ",")
		// extract the rightmost ip
		for i := len(ips) - 1; i >= 0; i-- {
			ip := strings.TrimSpace(ips[i])
			if ip != "" {
				return ip
			}
		}
	}

	return fallbackIp
}
