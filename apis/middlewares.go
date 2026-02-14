package apis

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/hook"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/router"
	"github.com/pocketbase/pocketbase/tools/routine"
	"github.com/spf13/cast"
)

// Common request event store keys used by the middlewares and api handlers.
const (
	RequestEventKeyLogMeta = "pbLogMeta" // extra data to store with the request activity log

	requestEventKeyExecStart              = "__execStart"                 // the value must be time.Time
	requestEventKeySkipSuccessActivityLog = "__skipSuccessActivityLogger" // the value must be bool
)

const (
	DefaultWWWRedirectMiddlewarePriority = -99999
	DefaultWWWRedirectMiddlewareId       = "pbWWWRedirect"

	DefaultActivityLoggerMiddlewarePriority   = DefaultRateLimitMiddlewarePriority - 40
	DefaultActivityLoggerMiddlewareId         = "pbActivityLogger"
	DefaultSkipSuccessActivityLogMiddlewareId = "pbSkipSuccessActivityLog"
	DefaultEnableAuthIdActivityLog            = "pbEnableAuthIdActivityLog"

	DefaultPanicRecoverMiddlewarePriority = DefaultRateLimitMiddlewarePriority - 30
	DefaultPanicRecoverMiddlewareId       = "pbPanicRecover"

	DefaultLoadAuthTokenMiddlewarePriority = DefaultRateLimitMiddlewarePriority - 20
	DefaultLoadAuthTokenMiddlewareId       = "pbLoadAuthToken"

	DefaultSecurityHeadersMiddlewarePriority = DefaultRateLimitMiddlewarePriority - 10
	DefaultSecurityHeadersMiddlewareId       = "pbSecurityHeaders"

	DefaultRequireGuestOnlyMiddlewareId                 = "pbRequireGuestOnly"
	DefaultRequireAuthMiddlewareId                      = "pbRequireAuth"
	DefaultRequireSuperuserAuthMiddlewareId             = "pbRequireSuperuserAuth"
	DefaultRequireSuperuserOrOwnerAuthMiddlewareId      = "pbRequireSuperuserOrOwnerAuth"
	DefaultRequireSameCollectionContextAuthMiddlewareId = "pbRequireSameCollectionContextAuth"
)

// RequireGuestOnly middleware requires a request to NOT have a valid
// Authorization header.
//
// This middleware is the opposite of [apis.RequireAuth()].
func RequireGuestOnly() *hook.Handler[*core.RequestEvent] {
	return &hook.Handler[*core.RequestEvent]{
		Id: DefaultRequireGuestOnlyMiddlewareId,
		Func: func(e *core.RequestEvent) error {
			if e.Auth != nil {
				return router.NewBadRequestError("The request can be accessed only by guests.", nil)
			}

			return e.Next()
		},
	}
}

// RequireAuth middleware requires a request to have a valid record Authorization header.
//
// The auth record could be from any collection.
// You can further filter the allowed record auth collections by specifying their names.
//
// Example:
//
//	apis.RequireAuth()                      // any auth collection
//	apis.RequireAuth("_superusers", "users") // only the listed auth collections
func RequireAuth(optCollectionNames ...string) *hook.Handler[*core.RequestEvent] {
	return &hook.Handler[*core.RequestEvent]{
		Id:   DefaultRequireAuthMiddlewareId,
		Func: requireAuth(optCollectionNames...),
	}
}

func requireAuth(optCollectionNames ...string) func(*core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		if e.Auth == nil {
			return e.UnauthorizedError("The request requires valid record authorization token.", nil)
		}

		// check record collection name
		if len(optCollectionNames) > 0 && !slices.Contains(optCollectionNames, e.Auth.Collection().Name) {
			return e.ForbiddenError("The authorized record is not allowed to perform this action.", nil)
		}

		return e.Next()
	}
}

// RequireSuperuserAuth middleware requires a request to have
// a valid superuser Authorization header.
func RequireSuperuserAuth() *hook.Handler[*core.RequestEvent] {
	return &hook.Handler[*core.RequestEvent]{
		Id:   DefaultRequireSuperuserAuthMiddlewareId,
		Func: requireAuth(core.CollectionNameSuperusers),
	}
}

// RequireSuperuserOrOwnerAuth middleware requires a request to have
// a valid superuser or regular record owner Authorization header set.
//
// This middleware is similar to [apis.RequireAuth()] but
// for the auth record token expects to have the same id as the path
// parameter ownerIdPathParam (default to "id" if empty).
func RequireSuperuserOrOwnerAuth(ownerIdPathParam string) *hook.Handler[*core.RequestEvent] {
	return &hook.Handler[*core.RequestEvent]{
		Id: DefaultRequireSuperuserOrOwnerAuthMiddlewareId,
		Func: func(e *core.RequestEvent) error {
			if e.Auth == nil {
				return e.UnauthorizedError("The request requires superuser or record authorization token.", nil)
			}

			if e.Auth.IsSuperuser() {
				return e.Next()
			}

			if ownerIdPathParam == "" {
				ownerIdPathParam = "id"
			}
			ownerId := e.Request.PathValue(ownerIdPathParam)

			// note: it is considered "safe" to compare only the record id
			// since the auth record ids are treated as unique across all auth collections
			if e.Auth.Id != ownerId {
				return e.ForbiddenError("You are not allowed to perform this request.", nil)
			}

			return e.Next()
		},
	}
}

// RequireSameCollectionContextAuth middleware requires a request to have
// a valid record Authorization header and the auth record's collection to
// match the one from the route path parameter (default to "collection" if collectionParam is empty).
func RequireSameCollectionContextAuth(collectionPathParam string) *hook.Handler[*core.RequestEvent] {
	return &hook.Handler[*core.RequestEvent]{
		Id: DefaultRequireSameCollectionContextAuthMiddlewareId,
		Func: func(e *core.RequestEvent) error {
			if e.Auth == nil {
				return e.UnauthorizedError("The request requires valid record authorization token.", nil)
			}

			if collectionPathParam == "" {
				collectionPathParam = "collection"
			}

			collection, _ := e.App.FindCachedCollectionByNameOrId(e.Request.PathValue(collectionPathParam))
			if collection == nil || e.Auth.Collection().Id != collection.Id {
				return e.ForbiddenError(fmt.Sprintf("The request requires auth record from %s collection.", e.Auth.Collection().Name), nil)
			}

			return e.Next()
		},
	}
}

// loadAuthToken attempts to load the auth context based on the "Authorization: TOKEN" header value.
//
// This middleware does nothing in case of:
//   - missing, invalid or expired token
//   - e.Auth is already loaded by another middleware
//
// This middleware is registered by default for all routes.
//
// Note: We don't throw an error on invalid or expired token to allow
// users to extend with their own custom handling in external middleware(s).
func loadAuthToken() *hook.Handler[*core.RequestEvent] {
	return &hook.Handler[*core.RequestEvent]{
		Id:       DefaultLoadAuthTokenMiddlewareId,
		Priority: DefaultLoadAuthTokenMiddlewarePriority,
		Func: func(e *core.RequestEvent) error {
			// already loaded by another middleware
			if e.Auth != nil {
				return e.Next()
			}

			token := getAuthTokenFromRequest(e)
			if token == "" {
				return e.Next()
			}

			record, err := e.App.FindAuthRecordByToken(token, core.TokenTypeAuth)
			if err != nil {
				e.App.Logger().Debug("loadAuthToken failure", "error", err)
			} else if record != nil {
				e.Auth = record
			}

			return e.Next()
		},
	}
}

func getAuthTokenFromRequest(e *core.RequestEvent) string {
	token := e.Request.Header.Get("Authorization")

	// the "Bearer" schema prefix is not required by PocketBase and it is
	// supported only for compatibility with the defaults of some HTTP clients
	if len(token) > 7 && strings.EqualFold(token[:7], "Bearer ") {
		return token[7:]
	}

	return token
}

// wwwRedirect performs www->non-www redirect(s) if the request host
// matches with one of the values in redirectHosts.
//
// This middleware is registered by default on Serve for all routes.
func wwwRedirect(redirectHosts []string) *hook.Handler[*core.RequestEvent] {
	return &hook.Handler[*core.RequestEvent]{
		Id:       DefaultWWWRedirectMiddlewareId,
		Priority: DefaultWWWRedirectMiddlewarePriority,
		Func: func(e *core.RequestEvent) error {
			host := e.Request.Host

			if strings.HasPrefix(host, "www.") && list.ExistInSlice(host, redirectHosts) {
				// note: e.Request.URL.Scheme would be empty
				schema := "http://"
				if e.IsTLS() {
					schema = "https://"
				}

				return e.Redirect(
					http.StatusTemporaryRedirect,
					(schema + host[4:] + e.Request.RequestURI),
				)
			}

			return e.Next()
		},
	}
}

// panicRecover returns a default panic-recover handler.
func panicRecover() *hook.Handler[*core.RequestEvent] {
	return &hook.Handler[*core.RequestEvent]{
		Id:       DefaultPanicRecoverMiddlewareId,
		Priority: DefaultPanicRecoverMiddlewarePriority,
		Func: func(e *core.RequestEvent) (err error) {
			// panic-recover
			defer func() {
				recoverResult := recover()
				if recoverResult == nil {
					return
				}

				recoverErr, ok := recoverResult.(error)
				if !ok {
					recoverErr = fmt.Errorf("%v", recoverResult)
				} else if errors.Is(recoverErr, http.ErrAbortHandler) {
					// don't recover ErrAbortHandler so the response to the client can be aborted
					panic(recoverResult)
				}

				stack := make([]byte, 2<<10) // 2 KB
				length := runtime.Stack(stack, true)
				err = e.InternalServerError("", fmt.Errorf("[PANIC RECOVER] %w %s", recoverErr, stack[:length]))
			}()

			err = e.Next()

			return err
		},
	}
}

// securityHeaders middleware adds common security headers to the response.
//
// This middleware is registered by default for all routes.
func securityHeaders() *hook.Handler[*core.RequestEvent] {
	return &hook.Handler[*core.RequestEvent]{
		Id:       DefaultSecurityHeadersMiddlewareId,
		Priority: DefaultSecurityHeadersMiddlewarePriority,
		Func: func(e *core.RequestEvent) error {
			e.Response.Header().Set("X-XSS-Protection", "1; mode=block")
			e.Response.Header().Set("X-Content-Type-Options", "nosniff")
			e.Response.Header().Set("X-Frame-Options", "SAMEORIGIN")

			// @todo consider a default HSTS?
			// (see also https://webkit.org/blog/8146/protecting-against-hsts-abuse/)

			return e.Next()
		},
	}
}

// SkipSuccessActivityLog is a helper middleware that instructs the global
// activity logger to log only requests that have failed/returned an error.
func SkipSuccessActivityLog() *hook.Handler[*core.RequestEvent] {
	return &hook.Handler[*core.RequestEvent]{
		Id: DefaultSkipSuccessActivityLogMiddlewareId,
		Func: func(e *core.RequestEvent) error {
			e.Set(requestEventKeySkipSuccessActivityLog, true)
			return e.Next()
		},
	}
}

// activityLogger middleware takes care to save the request information
// into the logs database.
//
// This middleware is registered by default for all routes.
//
// The middleware does nothing if the app logs retention period is zero
// (aka. app.Settings().Logs.MaxDays = 0).
//
// Users can attach the [apis.SkipSuccessActivityLog()] middleware if
// you want to log only the failed requests.
func activityLogger() *hook.Handler[*core.RequestEvent] {
	return &hook.Handler[*core.RequestEvent]{
		Id:       DefaultActivityLoggerMiddlewareId,
		Priority: DefaultActivityLoggerMiddlewarePriority,
		Func: func(e *core.RequestEvent) error {
			e.Set(requestEventKeyExecStart, time.Now())

			err := e.Next()

			logRequest(e, err)

			return err
		},
	}
}

func logRequest(event *core.RequestEvent, err error) {
	// no logs retention
	if event.App.Settings().Logs.MaxDays == 0 {
		return
	}

	// the non-error route has explicitly disabled the activity logger
	if err == nil && event.Get(requestEventKeySkipSuccessActivityLog) != nil {
		return
	}

	attrs := make([]any, 0, 15)

	attrs = append(attrs, slog.String("type", "request"))

	started := cast.ToTime(event.Get(requestEventKeyExecStart))
	if !started.IsZero() {
		attrs = append(attrs, slog.Float64("execTime", float64(time.Since(started))/float64(time.Millisecond)))
	}

	if meta := event.Get(RequestEventKeyLogMeta); meta != nil {
		attrs = append(attrs, slog.Any("meta", meta))
	}

	status := event.Status()
	method := cutStr(strings.ToUpper(event.Request.Method), 50)
	requestUri := cutStr(event.Request.URL.RequestURI(), 3000)

	// parse the request error
	if err != nil {
		apiErr, isPlainApiError := err.(*router.ApiError)
		if isPlainApiError || errors.As(err, &apiErr) {
			// the status header wasn't written yet
			if status == 0 {
				status = apiErr.Status
			}

			var errMsg string
			if isPlainApiError {
				errMsg = apiErr.Message
			} else {
				// wrapped ApiError -> add the full serialized version
				// of the original error since it could contain more information
				errMsg = err.Error()
			}

			attrs = append(
				attrs,
				slog.String("error", errMsg),
				slog.Any("details", apiErr.RawData()),
			)
		} else {
			attrs = append(attrs, slog.String("error", err.Error()))
		}
	}

	attrs = append(
		attrs,
		slog.String("url", requestUri),
		slog.String("method", method),
		slog.Int("status", status),
		slog.String("referer", cutStr(event.Request.Referer(), 2000)),
		slog.String("userAgent", cutStr(event.Request.UserAgent(), 2000)),
	)

	if event.Auth != nil {
		attrs = append(attrs, slog.String("auth", event.Auth.Collection().Name))

		if event.App.Settings().Logs.LogAuthId {
			attrs = append(attrs, slog.String("authId", event.Auth.Id))
		}
	} else {
		attrs = append(attrs, slog.String("auth", ""))
	}

	if event.App.Settings().Logs.LogIP {
		attrs = append(
			attrs,
			slog.String("userIP", event.RealIP()),
			slog.String("remoteIP", event.RemoteIP()),
		)
	}

	// don't block on logs write
	routine.FireAndForget(func() {
		message := method + " "

		if escaped, unescapeErr := url.PathUnescape(requestUri); unescapeErr == nil {
			message += escaped
		} else {
			message += requestUri
		}

		if err != nil {
			event.App.Logger().Error(message, attrs...)
		} else {
			event.App.Logger().Info(message, attrs...)
		}
	})
}

func cutStr(str string, max int) string {
	if len(str) > max {
		return str[:max] + "..."
	}
	return str
}
