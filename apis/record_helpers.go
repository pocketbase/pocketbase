package apis

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/mails"
	"github.com/pocketbase/pocketbase/tools/router"
	"github.com/pocketbase/pocketbase/tools/routine"
	"github.com/pocketbase/pocketbase/tools/search"
	"github.com/pocketbase/pocketbase/tools/security"
)

const (
	expandQueryParam = "expand"
	fieldsQueryParam = "fields"
)

var ErrMFA = errors.New("mfa required")

// RecordAuthResponse writes standardized json record auth response
// into the specified request context.
//
// The authMethod argument specify the name of the current authentication method (eg. password, oauth2, etc.)
// that it is used primarily as an auth identifier during MFA and for login alerts.
//
// Set authMethod to empty string if you want to ignore the MFA checks and the login alerts
// (can be also adjusted additionally via the OnRecordAuthRequest hook).
func RecordAuthResponse(e *core.RequestEvent, authRecord *core.Record, authMethod string, meta any) error {
	token, tokenErr := authRecord.NewAuthToken()
	if tokenErr != nil {
		return e.InternalServerError("Failed to create auth token.", tokenErr)
	}

	return recordAuthResponse(e, authRecord, token, authMethod, meta)
}

func recordAuthResponse(e *core.RequestEvent, authRecord *core.Record, token string, authMethod string, meta any) error {
	originalRequestInfo, err := e.RequestInfo()
	if err != nil {
		return err
	}

	ok, err := e.App.CanAccessRecord(authRecord, originalRequestInfo, authRecord.Collection().AuthRule)
	if !ok {
		return firstApiError(err, e.ForbiddenError("The request doesn't satisfy the collection requirements to authenticate.", err))
	}

	event := new(core.RecordAuthRequestEvent)
	event.RequestEvent = e
	event.Collection = authRecord.Collection()
	event.Record = authRecord
	event.Token = token
	event.Meta = meta
	event.AuthMethod = authMethod

	return e.App.OnRecordAuthRequest().Trigger(event, func(e *core.RecordAuthRequestEvent) error {
		if e.Written() {
			return nil
		}

		// MFA
		// ---
		mfaId, err := checkMFA(e.RequestEvent, e.Record, e.AuthMethod)
		if err != nil {
			return err
		}

		// require additional authentication
		if mfaId != "" {
			// eagerly write the mfa response and return an err so that
			// external middlewars are aware that the auth response requires an extra step
			e.JSON(http.StatusUnauthorized, map[string]string{
				"mfaId": mfaId,
			})
			return ErrMFA
		}
		// ---

		// create a shallow copy of the cached request data and adjust it to the current auth record
		requestInfo := *originalRequestInfo
		requestInfo.Auth = e.Record

		err = triggerRecordEnrichHooks(e.App, &requestInfo, []*core.Record{e.Record}, func() error {
			if e.Record.IsSuperuser() {
				e.Record.Unhide(e.Record.Collection().Fields.FieldNames()...)
			}

			// allow always returning the email address of the authenticated model
			e.Record.IgnoreEmailVisibility(true)

			// expand record relations
			expands := strings.Split(e.Request.URL.Query().Get(expandQueryParam), ",")
			if len(expands) > 0 {
				failed := e.App.ExpandRecord(e.Record, expands, expandFetch(e.App, &requestInfo))
				if len(failed) > 0 {
					e.App.Logger().Warn("[recordAuthResponse] Failed to expand relations", "error", failed)
				}
			}

			return nil
		})
		if err != nil {
			return err
		}

		if e.AuthMethod != "" && authRecord.Collection().AuthAlert.Enabled {
			if err = authAlert(e.RequestEvent, e.Record); err != nil {
				e.App.Logger().Warn("[recordAuthResponse] Failed to send login alert", "error", err)
			}
		}

		result := struct {
			Meta   any          `json:"meta,omitempty"`
			Record *core.Record `json:"record"`
			Token  string       `json:"token"`
		}{
			Token:  e.Token,
			Record: e.Record,
		}

		if e.Meta != nil {
			result.Meta = e.Meta
		}

		return e.JSON(http.StatusOK, result)
	})
}

// wantsMFA checks whether to enable MFA for the specified auth record based on its MFA rule
// (note: returns true even in case of an error as a safer default).
func wantsMFA(e *core.RequestEvent, record *core.Record) (bool, error) {
	rule := record.Collection().MFA.Rule
	if rule == "" {
		return true, nil
	}

	requestInfo, err := e.RequestInfo()
	if err != nil {
		return true, err
	}

	var exists int

	query := e.App.RecordQuery(record.Collection()).
		Select("(1)").
		AndWhere(dbx.HashExp{record.Collection().Name + ".id": record.Id})

	// parse and apply the access rule filter
	resolver := core.NewRecordFieldResolver(e.App, record.Collection(), requestInfo, true)
	expr, err := search.FilterData(rule).BuildExpr(resolver)
	if err != nil {
		return true, err
	}
	resolver.UpdateQuery(query)

	err = query.AndWhere(expr).Limit(1).Row(&exists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return true, err
	}

	return exists > 0, nil
}

// checkMFA handles any MFA auth checks that needs to be performed for the specified request event.
// Returns the mfaId that needs to be written as response to the user.
//
// (note: all auth methods are treated as equal and there is no requirement for "pairing").
func checkMFA(e *core.RequestEvent, authRecord *core.Record, currentAuthMethod string) (string, error) {
	if !authRecord.Collection().MFA.Enabled || currentAuthMethod == "" {
		return "", nil
	}

	ok, err := wantsMFA(e, authRecord)
	if err != nil {
		return "", e.BadRequestError("Failed to authenticate.", fmt.Errorf("MFA rule failure: %w", err))
	}
	if !ok {
		return "", nil // no mfa needed for this auth record
	}

	// read the mfaId either from the qyery params or request body
	mfaId := e.Request.URL.Query().Get("mfaId")
	if mfaId == "" {
		// check the body
		data := struct {
			MfaId string `form:"mfaId" json:"mfaId" xml:"mfaId"`
		}{}
		if err := e.BindBody(&data); err != nil {
			return "", firstApiError(err, e.BadRequestError("Failed to read MFA Id", err))
		}
		mfaId = data.MfaId
	}

	// first-time auth
	// ---
	if mfaId == "" {
		mfa := core.NewMFA(e.App)
		mfa.SetCollectionRef(authRecord.Collection().Id)
		mfa.SetRecordRef(authRecord.Id)
		mfa.SetMethod(currentAuthMethod)
		if err := e.App.Save(mfa); err != nil {
			return "", firstApiError(err, e.InternalServerError("Failed to create MFA record", err))
		}

		return mfa.Id, nil
	}

	// second-time auth
	// ---
	mfa, err := e.App.FindMFAById(mfaId)
	deleteMFA := func() {
		// try to delete the expired mfa
		if mfa != nil {
			if deleteErr := e.App.Delete(mfa); deleteErr != nil {
				e.App.Logger().Warn("Failed to delete expired MFA record", "error", deleteErr, "mfaId", mfa.Id)
			}
		}
	}
	if err != nil || mfa.HasExpired(authRecord.Collection().MFA.DurationTime()) {
		deleteMFA()
		return "", e.BadRequestError("Invalid or expired MFA session.", err)
	}

	if mfa.RecordRef() != authRecord.Id || mfa.CollectionRef() != authRecord.Collection().Id {
		return "", e.BadRequestError("Invalid MFA session.", nil)
	}

	if mfa.Method() == currentAuthMethod {
		return "", e.BadRequestError("A different authentication method is required.", nil)
	}

	deleteMFA()

	return "", nil
}

// EnrichRecord parses the request context and enrich the provided record:
//   - expands relations (if defaultExpands and/or ?expand query param is set)
//   - ensures that the emails of the auth record and its expanded auth relations
//     are visible only for the current logged superuser, record owner or record with manage access
func EnrichRecord(e *core.RequestEvent, record *core.Record, defaultExpands ...string) error {
	return EnrichRecords(e, []*core.Record{record}, defaultExpands...)
}

// EnrichRecords parses the request context and enriches the provided records:
//   - expands relations (if defaultExpands and/or ?expand query param is set)
//   - ensures that the emails of the auth records and their expanded auth relations
//     are visible only for the current logged superuser, record owner or record with manage access
//
// Note: Expects all records to be from the same collection!
func EnrichRecords(e *core.RequestEvent, records []*core.Record, defaultExpands ...string) error {
	if len(records) == 0 {
		return nil
	}

	info, err := e.RequestInfo()
	if err != nil {
		return err
	}

	return triggerRecordEnrichHooks(e.App, info, records, func() error {
		expands := defaultExpands
		if param := info.Query[expandQueryParam]; param != "" {
			expands = append(expands, strings.Split(param, ",")...)
		}

		err := defaultEnrichRecords(e.App, info, records, expands...)
		if err != nil {
			// only log because it is not critical
			e.App.Logger().Warn("failed to apply default enriching", "error", err)
		}

		return nil
	})
}

type iterator[T any] struct {
	items []T
	index int
}

func (ri *iterator[T]) next() T {
	var item T

	if ri.index < len(ri.items) {
		item = ri.items[ri.index]
		ri.index++
	}

	return item
}

func triggerRecordEnrichHooks(app core.App, requestInfo *core.RequestInfo, records []*core.Record, finalizer func() error) error {
	it := iterator[*core.Record]{items: records}

	enrichHook := app.OnRecordEnrich()

	event := new(core.RecordEnrichEvent)
	event.App = app
	event.RequestInfo = requestInfo

	var iterate func(record *core.Record) error
	iterate = func(record *core.Record) error {
		if record == nil {
			return nil
		}

		event.Record = record

		return enrichHook.Trigger(event, func(ee *core.RecordEnrichEvent) error {
			next := it.next()
			if next == nil {
				if finalizer != nil {
					return finalizer()
				}
				return nil
			}

			event.App = ee.App // in case it was replaced with a transaction
			event.Record = next

			err := iterate(next)

			event.App = app
			event.Record = record

			return err
		})
	}

	return iterate(it.next())
}

func defaultEnrichRecords(app core.App, requestInfo *core.RequestInfo, records []*core.Record, expands ...string) error {
	err := autoResolveRecordsFlags(app, records, requestInfo)
	if err != nil {
		return fmt.Errorf("failed to resolve records flags: %w", err)
	}

	if len(expands) > 0 {
		expandErrs := app.ExpandRecords(records, expands, expandFetch(app, requestInfo))
		if len(expandErrs) > 0 {
			errsSlice := make([]error, 0, len(expandErrs))
			for key, err := range expandErrs {
				errsSlice = append(errsSlice, fmt.Errorf("failed to expand %q: %w", key, err))
			}
			return fmt.Errorf("failed to expand records: %w", errors.Join(errsSlice...))
		}
	}

	return nil
}

// expandFetch is the records fetch function that is used to expand related records.
func expandFetch(app core.App, originalRequestInfo *core.RequestInfo) core.ExpandFetchFunc {
	// shallow clone the provided request info to set an "expand" context
	requestInfoClone := *originalRequestInfo
	requestInfoPtr := &requestInfoClone
	requestInfoPtr.Context = core.RequestInfoContextExpand

	return func(relCollection *core.Collection, relIds []string) ([]*core.Record, error) {
		records, findErr := app.FindRecordsByIds(relCollection.Id, relIds, func(q *dbx.SelectQuery) error {
			if requestInfoPtr.Auth != nil && requestInfoPtr.Auth.IsSuperuser() {
				return nil // superusers can access everything
			}

			if relCollection.ViewRule == nil {
				return fmt.Errorf("only superusers can view collection %q records", relCollection.Name)
			}

			if *relCollection.ViewRule != "" {
				resolver := core.NewRecordFieldResolver(app, relCollection, requestInfoPtr, true)
				expr, err := search.FilterData(*(relCollection.ViewRule)).BuildExpr(resolver)
				if err != nil {
					return err
				}
				resolver.UpdateQuery(q)
				q.AndWhere(expr)
			}

			return nil
		})
		if findErr != nil {
			return nil, findErr
		}

		enrichErr := triggerRecordEnrichHooks(app, requestInfoPtr, records, func() error {
			if err := autoResolveRecordsFlags(app, records, requestInfoPtr); err != nil {
				// non-critical error
				app.Logger().Warn("Failed to apply autoResolveRecordsFlags for the expanded records", "error", err)
			}

			return nil
		})
		if enrichErr != nil {
			return nil, enrichErr
		}

		return records, nil
	}
}

// autoResolveRecordsFlags resolves various visibility flags of the provided records.
//
// Currently it enables:
// - export of hidden fields if the current auth model is a superuser
// - email export ignoring the emailVisibity checks if the current auth model is superuser, owner or a "manager".
//
// Note: Expects all records to be from the same collection!
func autoResolveRecordsFlags(app core.App, records []*core.Record, requestInfo *core.RequestInfo) error {
	if len(records) == 0 {
		return nil // nothing to resolve
	}

	if requestInfo.HasSuperuserAuth() {
		hiddenFields := records[0].Collection().Fields.FieldNames()
		for _, rec := range records {
			rec.Unhide(hiddenFields...)
			rec.IgnoreEmailVisibility(true)
		}
	}

	// additional emailVisibility checks
	// ---------------------------------------------------------------
	if !records[0].Collection().IsAuth() {
		return nil // not auth collection records
	}

	collection := records[0].Collection()

	mappedRecords := make(map[string]*core.Record, len(records))
	recordIds := make([]any, len(records))
	for i, rec := range records {
		mappedRecords[rec.Id] = rec
		recordIds[i] = rec.Id
	}

	if requestInfo.Auth != nil && mappedRecords[requestInfo.Auth.Id] != nil {
		mappedRecords[requestInfo.Auth.Id].IgnoreEmailVisibility(true)
	}

	if collection.ManageRule == nil || *collection.ManageRule == "" {
		return nil // no manage rule to check
	}

	// fetch the ids of the managed records
	// ---
	managedIds := []string{}

	query := app.RecordQuery(collection).
		Select(app.DB().QuoteSimpleColumnName(collection.Name) + ".id").
		AndWhere(dbx.In(app.DB().QuoteSimpleColumnName(collection.Name)+".id", recordIds...))

	resolver := core.NewRecordFieldResolver(app, collection, requestInfo, true)
	expr, err := search.FilterData(*collection.ManageRule).BuildExpr(resolver)
	if err != nil {
		return err
	}
	resolver.UpdateQuery(query)
	query.AndWhere(expr)

	if err := query.Column(&managedIds); err != nil {
		return err
	}
	// ---

	// ignore the email visibility check for the managed records
	for _, id := range managedIds {
		if rec, ok := mappedRecords[id]; ok {
			rec.IgnoreEmailVisibility(true)
		}
	}

	return nil
}

var ruleQueryParams = []string{search.FilterQueryParam, search.SortQueryParam}
var superuserOnlyRuleFields = []string{"@collection.", "@request."}

// checkForSuperuserOnlyRuleFields loosely checks and returns an error if
// the provided RequestInfo contains rule fields that only the superuser can use.
func checkForSuperuserOnlyRuleFields(requestInfo *core.RequestInfo) error {
	if len(requestInfo.Query) == 0 || requestInfo.HasSuperuserAuth() {
		return nil // superuser or nothing to check
	}

	for _, param := range ruleQueryParams {
		v := requestInfo.Query[param]
		if v == "" {
			continue
		}

		for _, field := range superuserOnlyRuleFields {
			if strings.Contains(v, field) {
				return router.NewForbiddenError("Only superusers can filter by "+field, nil)
			}
		}
	}

	return nil
}

// firstApiError returns the first ApiError from the errors list
// (this is used usually to prevent unnecessary wraping and to allow bubling ApiError from nested hooks)
//
// If no ApiError is found, returns a default "Internal server" error.
func firstApiError(errs ...error) *router.ApiError {
	var apiErr *router.ApiError
	var ok bool

	for _, err := range errs {
		if err == nil {
			continue
		}

		// quick assert to avoid the reflection checks
		apiErr, ok = err.(*router.ApiError)
		if ok {
			return apiErr
		}

		// nested/wrapped errors
		if errors.As(err, &apiErr) {
			return apiErr
		}
	}

	return router.NewInternalServerError("", errors.Join(errs...))
}

// -------------------------------------------------------------------

const maxAuthOrigins = 5

func authAlert(e *core.RequestEvent, authRecord *core.Record) error {
	// generating fingerprint
	// ---
	userAgent := e.Request.UserAgent()
	if len(userAgent) > 300 {
		userAgent = userAgent[:300]
	}
	fingerprint := security.MD5(e.RealIP() + userAgent)
	// ---

	origins, err := e.App.FindAllAuthOriginsByRecord(authRecord)
	if err != nil {
		return err
	}

	isFirstLogin := len(origins) == 0

	var currentOrigin *core.AuthOrigin
	for _, origin := range origins {
		if origin.Fingerprint() == fingerprint {
			currentOrigin = origin
			break
		}
	}
	if currentOrigin == nil {
		currentOrigin = core.NewAuthOrigin(e.App)
		currentOrigin.SetCollectionRef(authRecord.Collection().Id)
		currentOrigin.SetRecordRef(authRecord.Id)
		currentOrigin.SetFingerprint(fingerprint)
	}

	// send email alert for the new origin auth (skip first login)
	//
	// Note: The "fake" timeout is a temp solution to avoid blocking
	//       for too long when the SMTP server is not accessible, due
	//       to the lack of context cancellation support in the underlying
	//       mailer and net/smtp package.
	//       The goroutine technically "leaks" but we assume that the OS will
	//       terminate the connection after some time (usually after 3-4 mins).
	if !isFirstLogin && currentOrigin.IsNew() && authRecord.Email() != "" {
		mailSent := make(chan error, 1)

		timer := time.AfterFunc(15*time.Second, func() {
			mailSent <- errors.New("auth alert mail send wait timeout reached")
		})

		routine.FireAndForget(func() {
			err := mails.SendRecordAuthAlert(e.App, authRecord)
			timer.Stop()
			mailSent <- err
		})

		err = <-mailSent
		if err != nil {
			return err
		}
	}

	// try to keep only up to maxAuthOrigins
	// (pop the last used ones; it is not executed in a transaction to avoid unnecessary locks)
	if currentOrigin.IsNew() && len(origins) >= maxAuthOrigins {
		for i := len(origins) - 1; i >= maxAuthOrigins-1; i-- {
			if err := e.App.Delete(origins[i]); err != nil {
				// treat as non-critical error, just log for now
				e.App.Logger().Warn("Failed to delete old AuthOrigin record", "error", err, "authOriginId", origins[i].Id)
			}
		}
	}

	// create/update the origin fingerprint
	return e.App.Save(currentOrigin)
}
