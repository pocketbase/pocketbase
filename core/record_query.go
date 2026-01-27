package core

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/tools/dbutils"
	"github.com/pocketbase/pocketbase/tools/inflector"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/search"
	"github.com/pocketbase/pocketbase/tools/security"
)

var recordProxyType = reflect.TypeOf((*RecordProxy)(nil)).Elem()

// RecordQuery returns a new Record select query from a collection model, id or name.
//
// In case a collection id or name is provided and that collection doesn't
// actually exists, the generated query will be created with a cancelled context
// and will fail once an executor (Row(), One(), All(), etc.) is called.
func (app *BaseApp) RecordQuery(collectionModelOrIdentifier any) *dbx.SelectQuery {
	var tableName string

	collection, collectionErr := getCollectionByModelOrIdentifier(app, collectionModelOrIdentifier)
	if collection != nil {
		tableName = collection.Name
	}
	if tableName == "" {
		// update with some fake table name for easier debugging
		tableName = "@@__invalidCollectionModelOrIdentifier"
	}

	query := app.ConcurrentDB().Select(app.ConcurrentDB().QuoteSimpleColumnName(tableName) + ".*").From(tableName)

	// in case of an error attach a new context and cancel it immediately with the error
	if collectionErr != nil {
		ctx, cancelFunc := context.WithCancelCause(context.Background())
		query.WithContext(ctx)
		cancelFunc(collectionErr)
	}

	return query.WithBuildHook(func(q *dbx.Query) {
		q.WithExecHook(execLockRetry(app.config.QueryTimeout, defaultMaxLockRetries)).
			WithOneHook(func(q *dbx.Query, a any, op func(b any) error) error {
				if a == nil {
					return op(a)
				}

				switch v := a.(type) {
				case *Record:
					record, err := resolveRecordOneHook(collection, op)
					if err != nil {
						return err
					}

					*v = *record

					return nil
				case RecordProxy:
					record, err := resolveRecordOneHook(collection, op)
					if err != nil {
						return err
					}

					v.SetProxyRecord(record)
					return nil
				default:
					return op(a)
				}
			}).
			WithAllHook(func(q *dbx.Query, sliceA any, op func(sliceB any) error) error {
				if sliceA == nil {
					return op(sliceA)
				}

				switch v := sliceA.(type) {
				case *[]*Record:
					records, err := resolveRecordAllHook(collection, op)
					if err != nil {
						return err
					}

					*v = records

					return nil
				case *[]Record:
					records, err := resolveRecordAllHook(collection, op)
					if err != nil {
						return err
					}

					nonPointers := make([]Record, len(records))
					for i, r := range records {
						nonPointers[i] = *r
					}

					*v = nonPointers

					return nil
				default: // expects []RecordProxy slice
					rv := reflect.ValueOf(v)
					if rv.Kind() != reflect.Ptr || rv.IsNil() {
						return errors.New("must be a pointer")
					}

					rv = dereference(rv)

					if rv.Kind() != reflect.Slice {
						return errors.New("must be a slice of RecordSetters")
					}

					et := rv.Type().Elem()

					var isSliceOfPointers bool
					if et.Kind() == reflect.Ptr {
						isSliceOfPointers = true
						et = et.Elem()
					}

					if !reflect.PointerTo(et).Implements(recordProxyType) {
						return op(sliceA)
					}

					records, err := resolveRecordAllHook(collection, op)
					if err != nil {
						return err
					}

					// create an empty slice
					if rv.IsNil() {
						rv.Set(reflect.MakeSlice(rv.Type(), 0, len(records)))
					}

					for _, record := range records {
						ev := reflect.New(et)

						if !ev.CanInterface() {
							continue
						}

						ps, ok := ev.Interface().(RecordProxy)
						if !ok {
							continue
						}

						ps.SetProxyRecord(record)

						ev = ev.Elem()
						if isSliceOfPointers {
							ev = ev.Addr()
						}

						rv.Set(reflect.Append(rv, ev))
					}

					return nil
				}
			})
	})
}

func resolveRecordOneHook(collection *Collection, op func(dst any) error) (*Record, error) {
	data := dbx.NullStringMap{}
	if err := op(&data); err != nil {
		return nil, err
	}
	return newRecordFromNullStringMap(collection, data)
}

func resolveRecordAllHook(collection *Collection, op func(dst any) error) ([]*Record, error) {
	data := []dbx.NullStringMap{}
	if err := op(&data); err != nil {
		return nil, err
	}
	return newRecordsFromNullStringMaps(collection, data)
}

// dereference returns the underlying value v points to.
func dereference(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			// initialize with a new value and continue searching
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}
	return v
}

func getCollectionByModelOrIdentifier(app App, collectionModelOrIdentifier any) (*Collection, error) {
	switch c := collectionModelOrIdentifier.(type) {
	case *Collection:
		return c, nil
	case Collection:
		return &c, nil
	case string:
		return app.FindCachedCollectionByNameOrId(c)
	default:
		return nil, errors.New("unknown collection identifier - must be collection model, id or name")
	}
}

// -------------------------------------------------------------------

// FindRecordById finds the Record model by its id.
func (app *BaseApp) FindRecordById(
	collectionModelOrIdentifier any,
	recordId string,
	optFilters ...func(q *dbx.SelectQuery) error,
) (*Record, error) {
	collection, err := getCollectionByModelOrIdentifier(app, collectionModelOrIdentifier)
	if err != nil {
		return nil, err
	}

	record := &Record{}

	query := app.RecordQuery(collection).
		AndWhere(dbx.HashExp{collection.Name + ".id": recordId})

	// apply filter funcs (if any)
	for _, filter := range optFilters {
		if filter == nil {
			continue
		}
		if err = filter(query); err != nil {
			return nil, err
		}
	}

	err = query.Limit(1).One(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

// FindRecordsByIds finds all records by the specified ids.
// If no records are found, returns an empty slice.
func (app *BaseApp) FindRecordsByIds(
	collectionModelOrIdentifier any,
	recordIds []string,
	optFilters ...func(q *dbx.SelectQuery) error,
) ([]*Record, error) {
	collection, err := getCollectionByModelOrIdentifier(app, collectionModelOrIdentifier)
	if err != nil {
		return nil, err
	}

	query := app.RecordQuery(collection).
		AndWhere(dbx.In(
			collection.Name+".id",
			list.ToInterfaceSlice(recordIds)...,
		))

	for _, filter := range optFilters {
		if filter == nil {
			continue
		}
		if err = filter(query); err != nil {
			return nil, err
		}
	}

	records := make([]*Record, 0, len(recordIds))

	err = query.All(&records)
	if err != nil {
		return nil, err
	}

	return records, nil
}

// FindAllRecords finds all records matching specified db expressions.
//
// Returns all collection records if no expression is provided.
//
// Returns an empty slice if no records are found.
//
// Example:
//
//	// no extra expressions
//	app.FindAllRecords("example")
//
//	// with extra expressions
//	expr1 := dbx.HashExp{"email": "test@example.com"}
//	expr2 := dbx.NewExp("LOWER(username) = {:username}", dbx.Params{"username": "test"})
//	app.FindAllRecords("example", expr1, expr2)
func (app *BaseApp) FindAllRecords(collectionModelOrIdentifier any, exprs ...dbx.Expression) ([]*Record, error) {
	query := app.RecordQuery(collectionModelOrIdentifier)

	for _, expr := range exprs {
		if expr != nil { // add only the non-nil expressions
			query.AndWhere(expr)
		}
	}

	var records []*Record

	if err := query.All(&records); err != nil {
		return nil, err
	}

	return records, nil
}

// FindFirstRecordByData returns the first found record matching
// the provided key-value pair.
func (app *BaseApp) FindFirstRecordByData(collectionModelOrIdentifier any, key string, value any) (*Record, error) {
	collection, err := getCollectionByModelOrIdentifier(app, collectionModelOrIdentifier)
	if err != nil {
		return nil, err
	}

	field := collection.Fields.GetByName(key)
	if field == nil {
		return nil, errors.New("invalid or missing field " + key)
	}

	record := &Record{}

	err = app.RecordQuery(collection).
		AndWhere(dbx.HashExp{inflector.Columnify(key): value}).
		Limit(1).
		One(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

// FindRecordsByFilter returns limit number of records matching the
// provided string filter.
//
// NB! Use the last "params" argument to bind untrusted user variables!
//
// The filter argument is optional and can be empty string to target
// all available records.
//
// The sort argument is optional and can be empty string OR the same format
// used in the web APIs, ex. "-created,title".
//
// If the limit argument is <= 0, no limit is applied to the query and
// all matching records are returned.
//
// Returns an empty slice if no records are found.
//
// Example:
//
//	app.FindRecordsByFilter(
//		"posts",
//		"title ~ {:title} && visible = {:visible}",
//		"-created",
//		10,
//		0,
//		dbx.Params{"title": "lorem ipsum", "visible": true}
//	)
func (app *BaseApp) FindRecordsByFilter(
	collectionModelOrIdentifier any,
	filter string,
	sort string,
	limit int,
	offset int,
	params ...dbx.Params,
) ([]*Record, error) {
	collection, err := getCollectionByModelOrIdentifier(app, collectionModelOrIdentifier)
	if err != nil {
		return nil, err
	}

	q := app.RecordQuery(collection)

	// build a fields resolver and attach the generated conditions to the query
	// ---
	resolver := NewRecordFieldResolver(
		app,
		collection, // the base collection
		nil,        // no request data
		true,       // allow searching hidden/protected fields like "email"
	)

	if filter != "" {
		expr, err := search.FilterData(filter).BuildExpr(resolver, params...)
		if err != nil {
			return nil, fmt.Errorf("invalid filter expression: %w", err)
		}
		q.AndWhere(expr)
	}

	if sort != "" {
		for _, sortField := range search.ParseSortFromString(sort) {
			expr, err := sortField.BuildExpr(resolver)
			if err != nil {
				return nil, err
			}
			if expr != "" {
				q.AndOrderBy(expr)
			}
		}
	}

	err = resolver.UpdateQuery(q) // attaches any adhoc joins and aliases
	if err != nil {
		return nil, err
	}
	// ---

	if offset > 0 {
		q.Offset(int64(offset))
	}

	if limit > 0 {
		q.Limit(int64(limit))
	}

	records := []*Record{}

	if err := q.All(&records); err != nil {
		return nil, err
	}

	return records, nil
}

// FindFirstRecordByFilter returns the first available record matching the provided filter (if any).
//
// NB! Use the last params argument to bind untrusted user variables!
//
// Returns sql.ErrNoRows if no record is found.
//
// Example:
//
//	app.FindFirstRecordByFilter("posts", "")
//	app.FindFirstRecordByFilter("posts", "slug={:slug} && status='public'", dbx.Params{"slug": "test"})
func (app *BaseApp) FindFirstRecordByFilter(
	collectionModelOrIdentifier any,
	filter string,
	params ...dbx.Params,
) (*Record, error) {
	result, err := app.FindRecordsByFilter(collectionModelOrIdentifier, filter, "", 1, 0, params...)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, sql.ErrNoRows
	}

	return result[0], nil
}

// CountRecords returns the total number of records in a collection.
func (app *BaseApp) CountRecords(collectionModelOrIdentifier any, exprs ...dbx.Expression) (int64, error) {
	var total int64

	q := app.RecordQuery(collectionModelOrIdentifier).Select("count(*)")

	for _, expr := range exprs {
		if expr != nil { // add only the non-nil expressions
			q.AndWhere(expr)
		}
	}

	err := q.Row(&total)

	return total, err
}

// FindAuthRecordByToken finds the auth record associated with the provided JWT
// (auth, file, verifyEmail, changeEmail, passwordReset types).
//
// Optionally specify a list of validTypes to check tokens only from those types.
//
// Returns an error if the JWT is invalid, expired or not associated to an auth collection record.
func (app *BaseApp) FindAuthRecordByToken(token string, validTypes ...string) (*Record, error) {
	if token == "" {
		return nil, errors.New("missing token")
	}

	unverifiedClaims, err := security.ParseUnverifiedJWT(token)
	if err != nil {
		return nil, err
	}

	// check required claims
	id, _ := unverifiedClaims[TokenClaimId].(string)
	collectionId, _ := unverifiedClaims[TokenClaimCollectionId].(string)
	tokenType, _ := unverifiedClaims[TokenClaimType].(string)
	if id == "" || collectionId == "" || tokenType == "" {
		return nil, errors.New("missing or invalid token claims")
	}

	// check types (if explicitly set)
	if len(validTypes) > 0 && !list.ExistInSlice(tokenType, validTypes) {
		return nil, fmt.Errorf("invalid token type %q, expects %q", tokenType, strings.Join(validTypes, ","))
	}

	record, err := app.FindRecordById(collectionId, id)
	if err != nil {
		return nil, err
	}

	if !record.Collection().IsAuth() {
		return nil, errors.New("the token is not associated to an auth collection record")
	}

	var baseTokenKey string
	switch tokenType {
	case TokenTypeAuth:
		baseTokenKey = record.Collection().AuthToken.Secret
	case TokenTypeFile:
		baseTokenKey = record.Collection().FileToken.Secret
	case TokenTypeVerification:
		baseTokenKey = record.Collection().VerificationToken.Secret
	case TokenTypePasswordReset:
		baseTokenKey = record.Collection().PasswordResetToken.Secret
	case TokenTypeEmailChange:
		baseTokenKey = record.Collection().EmailChangeToken.Secret
	default:
		return nil, errors.New("unknown token type " + tokenType)
	}

	secret := record.TokenKey() + baseTokenKey

	// verify token signature
	_, err = security.ParseJWT(token, secret)
	if err != nil {
		return nil, err
	}

	return record, nil
}

// FindAuthRecordByEmail finds the auth record associated with the provided email.
//
// The email check would be case-insensitive if the related collection
// email unique index has COLLATE NOCASE specified for the email column.
//
// Returns an error if it is not an auth collection or the record is not found.
func (app *BaseApp) FindAuthRecordByEmail(collectionModelOrIdentifier any, email string) (*Record, error) {
	collection, err := getCollectionByModelOrIdentifier(app, collectionModelOrIdentifier)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch auth collection: %w", err)
	}

	if !collection.IsAuth() {
		return nil, fmt.Errorf("%q is not an auth collection", collection.Name)
	}

	record := &Record{}

	var expr dbx.Expression

	index, ok := dbutils.FindSingleColumnUniqueIndex(collection.Indexes, FieldNameEmail)
	if ok && strings.EqualFold(index.Columns[0].Collate, "nocase") {
		// case-insensitive search
		expr = dbx.NewExp("[["+FieldNameEmail+"]] = {:email} COLLATE NOCASE", dbx.Params{"email": email})
	} else {
		expr = dbx.HashExp{FieldNameEmail: email}
	}

	err = app.RecordQuery(collection).
		AndWhere(expr).
		Limit(1).
		One(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

// CanAccessRecord checks if a record is allowed to be accessed by the
// specified requestInfo and accessRule.
//
// Rule and db checks are ignored in case requestInfo.Auth is a superuser.
//
// The returned error indicate that something unexpected happened during
// the check (eg. invalid rule or db query error).
//
// The method always return false on invalid rule or db query error.
//
// Example:
//
//	requestInfo, _ := e.RequestInfo()
//	record, _ := app.FindRecordById("example", "RECORD_ID")
//	rule := types.Pointer("@request.auth.id != '' || status = 'public'")
//	// ... or use one of the record collection's rule, eg. record.Collection().ViewRule
//
//	if ok, _ := app.CanAccessRecord(record, requestInfo, rule); ok { ... }
func (app *BaseApp) CanAccessRecord(record *Record, requestInfo *RequestInfo, accessRule *string) (bool, error) {
	// superusers can access everything
	if requestInfo.HasSuperuserAuth() {
		return true, nil
	}

	// only superusers can access this record
	if accessRule == nil {
		return false, nil
	}

	// empty public rule, aka. everyone can access
	if *accessRule == "" {
		return true, nil
	}

	var exists int

	query := app.RecordQuery(record.Collection()).
		Select("(1)").
		AndWhere(dbx.HashExp{record.Collection().Name + ".id": record.Id})

	// parse and apply the access rule filter
	resolver := NewRecordFieldResolver(app, record.Collection(), requestInfo, true)
	expr, err := search.FilterData(*accessRule).BuildExpr(resolver)
	if err != nil {
		return false, err
	}

	err = resolver.UpdateQuery(query)
	if err != nil {
		return false, err
	}

	err = query.AndWhere(expr).Limit(1).Row(&exists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}

	return exists > 0, nil
}
