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
	"github.com/pocketbase/pocketbase/tokens"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/routine"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/spf13/cast"
)

// Common request context keys used by the middlewares and api handlers.
const (
	ContextAdminKey      string = "admin"
	ContextAuthRecordKey string = "authRecord"
	ContextCollectionKey string = "collection"
)

// RequireGuestOnly middleware requires a request to NOT have a valid
// Authorization header.
//
// This middleware is the opposite of [apis.RequireAdminOrRecordAuth()].
func RequireGuestOnly() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := NewBadRequestError("The request can be accessed only by guests.", nil)

			record, _ := c.Get(ContextAuthRecordKey).(*models.Record)
			if record != nil {
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

// RequireRecordAuth middleware requires a request to have
// a valid record auth Authorization header.
//
// The auth record could be from any collection.
//
// You can further filter the allowed record auth collections by
// specifying their names.
//
// Example:
//
//	apis.RequireRecordAuth()
//
// Or:
//
//	apis.RequireRecordAuth("users", "supervisors")
//
// To restrict the auth record only to the loaded context collection,
// use [apis.RequireSameContextRecordAuth()] instead.
func RequireRecordAuth(optCollectionNames ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			record, _ := c.Get(ContextAuthRecordKey).(*models.Record)
			if record == nil {
				return NewUnauthorizedError("The request requires valid record authorization token to be set.", nil)
			}

			// check record collection name
			if len(optCollectionNames) > 0 && !list.ExistInSlice(record.Collection().Name, optCollectionNames) {
				return NewForbiddenError("The authorized record model is not allowed to perform this action.", nil)
			}

			return next(c)
		}
	}
}

// RequireSameContextRecordAuth middleware requires a request to have
// a valid record Authorization header.
//
// The auth record must be from the same collection already loaded in the context.
func RequireSameContextRecordAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			record, _ := c.Get(ContextAuthRecordKey).(*models.Record)
			if record == nil {
				return NewUnauthorizedError("The request requires valid record authorization token to be set.", nil)
			}

			collection, _ := c.Get(ContextCollectionKey).(*models.Collection)
			if collection == nil || record.Collection().Id != collection.Id {
				return NewForbiddenError(fmt.Sprintf("The request requires auth record from %s collection.", record.Collection().Name), nil)
			}

			return next(c)
		}
	}
}

// RequireAdminAuth middleware requires a request to have
// a valid admin Authorization header.
func RequireAdminAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			admin, _ := c.Get(ContextAdminKey).(*models.Admin)
			if admin == nil {
				return NewUnauthorizedError("The request requires valid admin authorization token to be set.", nil)
			}

			return next(c)
		}
	}
}

// RequireAdminAuthOnlyIfAny middleware requires a request to have
// a valid admin Authorization header ONLY if the application has
// at least 1 existing Admin model.
func RequireAdminAuthOnlyIfAny(app core.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			admin, _ := c.Get(ContextAdminKey).(*models.Admin)
			if admin != nil {
				return next(c)
			}

			totalAdmins, err := app.Dao().TotalAdmins()
			if err != nil {
				return NewBadRequestError("Failed to fetch admins info.", err)
			}

			if totalAdmins == 0 {
				return next(c)
			}

			return NewUnauthorizedError("The request requires valid admin authorization token to be set.", nil)
		}
	}
}

// RequireAdminOrRecordAuth middleware requires a request to have
// a valid admin or record Authorization header set.
//
// You can further filter the allowed auth record collections by providing their names.
//
// This middleware is the opposite of [apis.RequireGuestOnly()].
func RequireAdminOrRecordAuth(optCollectionNames ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			admin, _ := c.Get(ContextAdminKey).(*models.Admin)
			record, _ := c.Get(ContextAuthRecordKey).(*models.Record)

			if admin == nil && record == nil {
				return NewUnauthorizedError("The request requires admin or record authorization token to be set.", nil)
			}

			if record != nil && len(optCollectionNames) > 0 && !list.ExistInSlice(record.Collection().Name, optCollectionNames) {
				return NewForbiddenError("The authorized record model is not allowed to perform this action.", nil)
			}

			return next(c)
		}
	}
}

// RequireAdminOrOwnerAuth middleware requires a request to have
// a valid admin or auth record owner Authorization header set.
//
// This middleware is similar to [apis.RequireAdminOrRecordAuth()] but
// for the auth record token expects to have the same id as the path
// parameter ownerIdParam (default to "id" if empty).
func RequireAdminOrOwnerAuth(ownerIdParam string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			admin, _ := c.Get(ContextAdminKey).(*models.Admin)
			if admin != nil {
				return next(c)
			}

			record, _ := c.Get(ContextAuthRecordKey).(*models.Record)
			if record == nil {
				return NewUnauthorizedError("The request requires admin or record authorization token to be set.", nil)
			}

			if ownerIdParam == "" {
				ownerIdParam = "id"
			}
			ownerId := c.PathParam(ownerIdParam)

			// note: it is "safe" to compare only the record id since the auth
			// record ids are treated as unique across all auth collections
			if record.Id != ownerId {
				return NewForbiddenError("You are not allowed to perform this request.", nil)
			}

			return next(c)
		}
	}
}

// LoadAuthContext middleware reads the Authorization request header
// and loads the token related record or admin instance into the
// request's context.
//
// This middleware is expected to be already registered by default for all routes.
func LoadAuthContext(app core.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")
			if token == "" {
				return next(c)
			}

			// the schema is not required and it is only for
			// compatibility with the defaults of some HTTP clients
			token = strings.TrimPrefix(token, "Bearer ")

			claims, _ := security.ParseUnverifiedJWT(token)
			tokenType := cast.ToString(claims["type"])

			switch tokenType {
			case tokens.TypeAdmin:
				admin, err := app.Dao().FindAdminByToken(
					token,
					app.Settings().AdminAuthToken.Secret,
				)
				if err == nil && admin != nil {
					c.Set(ContextAdminKey, admin)
				}
			case tokens.TypeAuthRecord:
				record, err := app.Dao().FindAuthRecordByToken(
					token,
					app.Settings().RecordAuthToken.Secret,
				)
				if err == nil && record != nil {
					c.Set(ContextAuthRecordKey, record)
				}
			}

			return next(c)
		}
	}
}

// LoadCollectionContext middleware finds the collection with related
// path identifier and loads it into the request context.
//
// Set optCollectionTypes to further filter the found collection by its type.
func LoadCollectionContext(app core.App, optCollectionTypes ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if param := c.PathParam("collection"); param != "" {
				collection, err := app.Dao().FindCollectionByNameOrId(param)
				if err != nil || collection == nil {
					return NewNotFoundError("", err)
				}

				if len(optCollectionTypes) > 0 && !list.ExistInSlice(collection.Type, optCollectionTypes) {
					return NewBadRequestError("Unsupported collection type.", nil)
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
				case *ApiError:
					status = v.Code
					meta["errorMessage"] = v.Message
					meta["errorDetails"] = fmt.Sprint(v.RawData())
				default:
					status = http.StatusBadRequest
					meta["errorMessage"] = v.Error()
				}
			}

			requestAuth := models.RequestAuthGuest
			if c.Get(ContextAuthRecordKey) != nil {
				requestAuth = models.RequestAuthRecord
			} else if c.Get(ContextAdminKey) != nil {
				requestAuth = models.RequestAuthAdmin
			}

			ip, _, _ := net.SplitHostPort(httpRequest.RemoteAddr)

			model := &models.Request{
				Url:       httpRequest.URL.RequestURI(),
				Method:    strings.ToUpper(httpRequest.Method),
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
				if err := app.LogsDao().SaveRequest(model); err != nil && app.IsDebug() {
					log.Println("Log save failed:", err)
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

	if ip := r.Header.Get("Fly-Client-IP"); ip != "" {
		return ip
	}

	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}

	if ipsList := r.Header.Get("X-Forwarded-For"); ipsList != "" {
		// extract the first non-empty leftmost-ish ip
		ips := strings.Split(ipsList, ",")
		for _, ip := range ips {
			ip = strings.TrimSpace(ip)
			if ip != "" {
				return ip
			}
		}
	}

	return fallbackIp
}

// eagerRequestDataCache ensures that the request data is cached in the request
// context to allow reading for example the json request body data more than once.
func eagerRequestDataCache(app core.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			switch c.Request().Method {
			// currently we are eagerly caching only the requests with body
			case "POST", "PUT", "PATCH", "DELETE":
				RequestData(c)
			}

			return next(c)
		}
	}
}
