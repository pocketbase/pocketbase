package apis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/hook"
	"github.com/pocketbase/pocketbase/tools/picker"
	"github.com/pocketbase/pocketbase/tools/router"
	"github.com/pocketbase/pocketbase/tools/routine"
	"github.com/pocketbase/pocketbase/tools/search"
	"github.com/pocketbase/pocketbase/tools/subscriptions"
	"golang.org/x/sync/errgroup"
)

// note: the chunk size is arbitrary chosen and may change in the future
const clientsChunkSize = 150

// RealtimeClientAuthKey is the name of the realtime client store key that holds its auth state.
const RealtimeClientAuthKey = "auth"

// bindRealtimeApi registers the realtime api endpoints.
func bindRealtimeApi(app core.App, rg *router.RouterGroup[*core.RequestEvent]) {
	sub := rg.Group("/realtime")
	sub.GET("", realtimeConnect).Bind(SkipSuccessActivityLog())
	sub.POST("", realtimeSetSubscriptions)

	bindRealtimeEvents(app)
}

func realtimeConnect(e *core.RequestEvent) error {
	// disable global write deadline for the SSE connection
	rc := http.NewResponseController(e.Response)
	writeDeadlineErr := rc.SetWriteDeadline(time.Time{})
	if writeDeadlineErr != nil {
		if !errors.Is(writeDeadlineErr, http.ErrNotSupported) {
			return e.InternalServerError("Failed to initialize SSE connection.", writeDeadlineErr)
		}

		// only log since there are valid cases where it may not be implement (e.g. httptest.ResponseRecorder)
		e.App.Logger().Warn("SetWriteDeadline is not supported, fallback to the default server WriteTimeout")
	}

	// create cancellable request
	cancelCtx, cancelRequest := context.WithCancel(e.Request.Context())
	defer cancelRequest()
	e.Request = e.Request.Clone(cancelCtx)

	e.Response.Header().Set("Content-Type", "text/event-stream")
	e.Response.Header().Set("Cache-Control", "no-store")
	// https://github.com/pocketbase/pocketbase/discussions/480#discussioncomment-3657640
	// https://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_buffering
	e.Response.Header().Set("X-Accel-Buffering", "no")

	connectEvent := new(core.RealtimeConnectRequestEvent)
	connectEvent.RequestEvent = e
	connectEvent.Client = subscriptions.NewDefaultClient()
	connectEvent.IdleTimeout = 5 * time.Minute

	return e.App.OnRealtimeConnectRequest().Trigger(connectEvent, func(ce *core.RealtimeConnectRequestEvent) error {
		// register new subscription client
		ce.App.SubscriptionsBroker().Register(ce.Client)
		defer func() {
			e.App.SubscriptionsBroker().Unregister(ce.Client.Id())
		}()

		ce.App.Logger().Debug("Realtime connection established.", slog.String("clientId", ce.Client.Id()))

		// signalize established connection (aka. fire "connect" message)
		connectMsgEvent := new(core.RealtimeMessageEvent)
		connectMsgEvent.RequestEvent = ce.RequestEvent
		connectMsgEvent.Client = ce.Client
		connectMsgEvent.Message = &subscriptions.Message{
			Name: "PB_CONNECT",
			Data: []byte(`{"clientId":"` + ce.Client.Id() + `"}`),
		}
		connectMsgErr := ce.App.OnRealtimeMessageSend().Trigger(connectMsgEvent, func(me *core.RealtimeMessageEvent) error {
			err := me.Message.WriteSSE(me.Response, me.Client.Id())
			if err != nil {
				return err
			}
			return me.Flush()
		})
		if connectMsgErr != nil {
			ce.App.Logger().Debug(
				"Realtime connection closed (failed to deliver PB_CONNECT)",
				slog.String("clientId", ce.Client.Id()),
				slog.String("error", connectMsgErr.Error()),
			)
			return nil
		}

		// start an idle timer to keep track of inactive/forgotten connections
		idleTimer := time.NewTimer(ce.IdleTimeout)
		defer idleTimer.Stop()

		for {
			select {
			case <-idleTimer.C:
				cancelRequest()
			case msg, ok := <-ce.Client.Channel():
				if !ok {
					// channel is closed
					ce.App.Logger().Debug(
						"Realtime connection closed (closed channel)",
						slog.String("clientId", ce.Client.Id()),
					)
					return nil
				}

				msgEvent := new(core.RealtimeMessageEvent)
				msgEvent.RequestEvent = ce.RequestEvent
				msgEvent.Client = ce.Client
				msgEvent.Message = &msg
				msgErr := ce.App.OnRealtimeMessageSend().Trigger(msgEvent, func(me *core.RealtimeMessageEvent) error {
					err := me.Message.WriteSSE(me.Response, me.Client.Id())
					if err != nil {
						return err
					}
					return me.Flush()
				})
				if msgErr != nil {
					ce.App.Logger().Debug(
						"Realtime connection closed (failed to deliver message)",
						slog.String("clientId", ce.Client.Id()),
						slog.String("error", msgErr.Error()),
					)
					return nil
				}

				idleTimer.Stop()
				idleTimer.Reset(ce.IdleTimeout)
			case <-ce.Request.Context().Done():
				// connection is closed
				ce.App.Logger().Debug(
					"Realtime connection closed (cancelled request)",
					slog.String("clientId", ce.Client.Id()),
				)
				return nil
			}
		}
	})
}

type realtimeSubscribeForm struct {
	ClientId      string   `form:"clientId" json:"clientId"`
	Subscriptions []string `form:"subscriptions" json:"subscriptions"`
}

func (form *realtimeSubscribeForm) validate() error {
	return validation.ValidateStruct(form,
		validation.Field(&form.ClientId, validation.Required, validation.Length(1, 255)),
		validation.Field(&form.Subscriptions,
			validation.Length(0, 1000),
			validation.Each(validation.Length(0, 2500)),
		),
	)
}

// note: in case of reconnect, clients will have to resubmit all subscriptions again
func realtimeSetSubscriptions(e *core.RequestEvent) error {
	form := new(realtimeSubscribeForm)

	err := e.BindBody(form)
	if err != nil {
		return e.BadRequestError("", err)
	}

	err = form.validate()
	if err != nil {
		return e.BadRequestError("", err)
	}

	// find subscription client
	client, err := e.App.SubscriptionsBroker().ClientById(form.ClientId)
	if err != nil {
		return e.NotFoundError("Missing or invalid client id.", err)
	}

	// for now allow only guest->auth upgrades and any other auth change is forbidden
	clientAuth, _ := client.Get(RealtimeClientAuthKey).(*core.Record)
	if clientAuth != nil && !isSameAuth(clientAuth, e.Auth) {
		return e.ForbiddenError("The current and the previous request authorization don't match.", nil)
	}

	event := new(core.RealtimeSubscribeRequestEvent)
	event.RequestEvent = e
	event.Client = client
	event.Subscriptions = form.Subscriptions

	return e.App.OnRealtimeSubscribeRequest().Trigger(event, func(e *core.RealtimeSubscribeRequestEvent) error {
		// update auth state
		e.Client.Set(RealtimeClientAuthKey, e.Auth)

		// unsubscribe from any previous existing subscriptions
		e.Client.Unsubscribe()

		// subscribe to the new subscriptions
		e.Client.Subscribe(e.Subscriptions...)

		e.App.Logger().Debug(
			"Realtime subscriptions updated.",
			slog.String("clientId", e.Client.Id()),
			slog.Any("subscriptions", e.Subscriptions),
		)

		return execAfterSuccessTx(true, e.App, func() error {
			return e.NoContent(http.StatusNoContent)
		})
	})
}

// updateClientsAuth updates the existing clients auth record with the new one (matched by ID).
func realtimeUpdateClientsAuth(app core.App, newAuthRecord *core.Record) error {
	chunks := app.SubscriptionsBroker().ChunkedClients(clientsChunkSize)

	group := new(errgroup.Group)

	for _, chunk := range chunks {
		group.Go(func() error {
			for _, client := range chunk {
				clientAuth, _ := client.Get(RealtimeClientAuthKey).(*core.Record)
				if clientAuth != nil &&
					clientAuth.Id == newAuthRecord.Id &&
					clientAuth.Collection().Name == newAuthRecord.Collection().Name {
					client.Set(RealtimeClientAuthKey, newAuthRecord)
				}
			}

			return nil
		})
	}

	return group.Wait()
}

// realtimeUnsetClientsAuthState unsets the auth state of all clients that have the provided auth model.
func realtimeUnsetClientsAuthState(app core.App, authModel core.Model) error {
	chunks := app.SubscriptionsBroker().ChunkedClients(clientsChunkSize)

	group := new(errgroup.Group)

	for _, chunk := range chunks {
		group.Go(func() error {
			for _, client := range chunk {
				clientAuth, _ := client.Get(RealtimeClientAuthKey).(*core.Record)
				if clientAuth != nil &&
					clientAuth.Id == authModel.PK() &&
					clientAuth.Collection().Name == authModel.TableName() {
					client.Unset(RealtimeClientAuthKey)
				}
			}

			return nil
		})
	}

	return group.Wait()
}

func bindRealtimeEvents(app core.App) {
	// update the clients that has auth record association
	app.OnModelAfterUpdateSuccess().Bind(&hook.Handler[*core.ModelEvent]{
		Func: func(e *core.ModelEvent) error {
			authRecord := realtimeResolveRecord(e.App, e.Model, core.CollectionTypeAuth)
			if authRecord != nil {
				if err := realtimeUpdateClientsAuth(e.App, authRecord); err != nil {
					app.Logger().Warn(
						"Failed to update client(s) associated to the updated auth record",
						slog.Any("id", authRecord.Id),
						slog.String("collectionName", authRecord.Collection().Name),
						slog.String("error", err.Error()),
					)
				}
			}

			return e.Next()
		},
		Priority: -99,
	})

	// remove the client(s) associated to the deleted auth model
	// (note: works also with custom model for backward compatibility)
	app.OnModelAfterDeleteSuccess().Bind(&hook.Handler[*core.ModelEvent]{
		Func: func(e *core.ModelEvent) error {
			collection := realtimeResolveRecordCollection(e.App, e.Model)
			if collection != nil && collection.IsAuth() {
				if err := realtimeUnsetClientsAuthState(e.App, e.Model); err != nil {
					app.Logger().Warn(
						"Failed to remove client(s) associated to the deleted auth model",
						slog.Any("id", e.Model.PK()),
						slog.String("collectionName", e.Model.TableName()),
						slog.String("error", err.Error()),
					)
				}
			}

			return e.Next()
		},
		Priority: -99,
	})

	app.OnModelAfterCreateSuccess().Bind(&hook.Handler[*core.ModelEvent]{
		Func: func(e *core.ModelEvent) error {
			record := realtimeResolveRecord(e.App, e.Model, "")
			if record != nil {
				err := realtimeBroadcastRecord(e.App, "create", record, false)
				if err != nil {
					app.Logger().Debug(
						"Failed to broadcast record create",
						slog.String("id", record.Id),
						slog.String("collectionName", record.Collection().Name),
						slog.String("error", err.Error()),
					)
				}
			}

			return e.Next()
		},
		Priority: -99,
	})

	app.OnModelAfterUpdateSuccess().Bind(&hook.Handler[*core.ModelEvent]{
		Func: func(e *core.ModelEvent) error {
			record := realtimeResolveRecord(e.App, e.Model, "")
			if record != nil {
				err := realtimeBroadcastRecord(e.App, "update", record, false)
				if err != nil {
					app.Logger().Debug(
						"Failed to broadcast record update",
						slog.String("id", record.Id),
						slog.String("collectionName", record.Collection().Name),
						slog.String("error", err.Error()),
					)
				}
			}

			return e.Next()
		},
		Priority: -99,
	})

	// delete: dry cache
	app.OnModelDelete().Bind(&hook.Handler[*core.ModelEvent]{
		Func: func(e *core.ModelEvent) error {
			record := realtimeResolveRecord(e.App, e.Model, "")
			if record != nil {
				// note: use the outside scoped app instance for the access checks so that the API rules
				// are performed out of the delete transaction ensuring that they would still work even if
				// a cascade-deleted record's API rule relies on an already deleted parent record
				err := realtimeBroadcastRecord(e.App, "delete", record, true, app)
				if err != nil {
					app.Logger().Debug(
						"Failed to dry cache record delete",
						slog.String("id", record.Id),
						slog.String("collectionName", record.Collection().Name),
						slog.String("error", err.Error()),
					)
				}
			}

			return e.Next()
		},
		Priority: 99, // execute as later as possible
	})

	// delete: broadcast
	app.OnModelAfterDeleteSuccess().Bind(&hook.Handler[*core.ModelEvent]{
		Func: func(e *core.ModelEvent) error {
			// note: only ensure that it is a collection record
			// and don't use realtimeResolveRecord because in case of a
			// custom model it'll fail to resolve since the record is already deleted
			collection := realtimeResolveRecordCollection(e.App, e.Model)
			if collection != nil {
				err := realtimeBroadcastDryCacheKey(e.App, getDryCacheKey("delete", e.Model))
				if err != nil {
					app.Logger().Debug(
						"Failed to broadcast record delete",
						slog.Any("id", e.Model.PK()),
						slog.String("collectionName", collection.Name),
						slog.String("error", err.Error()),
					)
				}
			}

			return e.Next()
		},
		Priority: -99,
	})

	// delete: failure
	app.OnModelAfterDeleteError().Bind(&hook.Handler[*core.ModelErrorEvent]{
		Func: func(e *core.ModelErrorEvent) error {
			record := realtimeResolveRecord(e.App, e.Model, "")
			if record != nil {
				err := realtimeUnsetDryCacheKey(e.App, getDryCacheKey("delete", record))
				if err != nil {
					app.Logger().Debug(
						"Failed to cleanup after broadcast record delete failure",
						slog.String("id", record.Id),
						slog.String("collectionName", record.Collection().Name),
						slog.String("error", err.Error()),
					)
				}
			}

			return e.Next()
		},
		Priority: -99,
	})
}

// resolveRecord converts *if possible* the provided model interface to a Record.
// This is usually helpful if the provided model is a custom Record model struct.
func realtimeResolveRecord(app core.App, model core.Model, optCollectionType string) *core.Record {
	var record *core.Record
	switch m := model.(type) {
	case *core.Record:
		record = m
	case core.RecordProxy:
		record = m.ProxyRecord()
	}

	if record != nil {
		if optCollectionType == "" || record.Collection().Type == optCollectionType {
			return record
		}
		return nil
	}

	tblName := model.TableName()

	// skip Log model checks
	if tblName == core.LogsTableName {
		return nil
	}

	// check if it is custom Record model struct
	collection, _ := app.FindCachedCollectionByNameOrId(tblName)
	if collection != nil && (optCollectionType == "" || collection.Type == optCollectionType) {
		if id, ok := model.PK().(string); ok {
			record, _ = app.FindRecordById(collection, id)
		}
	}

	return record
}

// realtimeResolveRecordCollection extracts *if possible* the Collection model from the provided model interface.
// This is usually helpful if the provided model is a custom Record model struct.
func realtimeResolveRecordCollection(app core.App, model core.Model) (collection *core.Collection) {
	switch m := model.(type) {
	case *core.Record:
		return m.Collection()
	case core.RecordProxy:
		return m.ProxyRecord().Collection()
	default:
		// check if it is custom Record model struct
		collection, err := app.FindCachedCollectionByNameOrId(model.TableName())
		if err == nil {
			return collection
		}
	}

	return nil
}

// recordData represents the broadcasted record subscrition message data.
type recordData struct {
	Record any    `json:"record"` /* map or core.Record */
	Action string `json:"action"`
}

// Note: the optAccessCheckApp is there in case you want the access check
// to be performed against different db app context (e.g. out of a transaction).
// If set, it is expected that optAccessCheckApp instance is used for read-only operations to avoid deadlocks.
// If not set, it fallbacks to app.
func realtimeBroadcastRecord(app core.App, action string, record *core.Record, dryCache bool, optAccessCheckApp ...core.App) error {
	collection := record.Collection()
	if collection == nil {
		return errors.New("[broadcastRecord] Record collection not set")
	}

	chunks := app.SubscriptionsBroker().ChunkedClients(clientsChunkSize)
	if len(chunks) == 0 {
		return nil // no subscribers
	}

	subscriptionRuleMap := map[string]*string{
		(collection.Name + "/" + record.Id + "?"): collection.ViewRule,
		(collection.Id + "/" + record.Id + "?"):   collection.ViewRule,
		(collection.Name + "/*?"):                 collection.ListRule,
		(collection.Id + "/*?"):                   collection.ListRule,

		// @deprecated: the same as the wildcard topic but kept for backward compatibility
		(collection.Name + "?"): collection.ListRule,
		(collection.Id + "?"):   collection.ListRule,
	}

	dryCacheKey := getDryCacheKey(action, record)

	group := new(errgroup.Group)

	accessCheckApp := app
	if len(optAccessCheckApp) > 0 {
		accessCheckApp = optAccessCheckApp[0]
	}

	for _, chunk := range chunks {
		group.Go(func() error {
			var clientAuth *core.Record

			for _, client := range chunk {
				// note: not executed concurrently to avoid races and to ensure
				// that the access checks are applied for the current record db state
				for prefix, rule := range subscriptionRuleMap {
					subs := client.Subscriptions(prefix)
					if len(subs) == 0 {
						continue
					}

					clientAuth, _ = client.Get(RealtimeClientAuthKey).(*core.Record)

					for sub, options := range subs {
						// mock request data
						requestInfo := &core.RequestInfo{
							Context: core.RequestInfoContextRealtime,
							Method:  "GET",
							Query:   options.Query,
							Headers: options.Headers,
							Auth:    clientAuth,
						}

						if !realtimeCanAccessRecord(accessCheckApp, record, requestInfo, rule) {
							continue
						}

						// create a clean record copy without expand and unknown fields because we don't know yet
						// which exact fields the client subscription requested or has permissions to access
						cleanRecord := record.Fresh()

						// trigger the enrich hooks
						enrichErr := triggerRecordEnrichHooks(app, requestInfo, []*core.Record{cleanRecord}, func() error {
							// apply expand
							rawExpand := options.Query[expandQueryParam]
							if rawExpand != "" {
								expandErrs := app.ExpandRecord(cleanRecord, strings.Split(rawExpand, ","), expandFetch(app, requestInfo))
								if len(expandErrs) > 0 {
									app.Logger().Debug(
										"[broadcastRecord] expand errors",
										slog.String("id", cleanRecord.Id),
										slog.String("collectionName", cleanRecord.Collection().Name),
										slog.String("sub", sub),
										slog.String("expand", rawExpand),
										slog.Any("errors", expandErrs),
									)
								}
							}

							// ignore the auth record email visibility checks
							// for auth owner, superuser or manager
							if collection.IsAuth() {
								if isSameAuth(clientAuth, cleanRecord) ||
									realtimeCanAccessRecord(accessCheckApp, cleanRecord, requestInfo, collection.ManageRule) {
									cleanRecord.IgnoreEmailVisibility(true)
								}
							}

							return nil
						})
						if enrichErr != nil {
							app.Logger().Debug(
								"[broadcastRecord] record enrich error",
								slog.String("id", cleanRecord.Id),
								slog.String("collectionName", cleanRecord.Collection().Name),
								slog.String("sub", sub),
								slog.Any("error", enrichErr),
							)
							continue
						}

						data := &recordData{
							Action: action,
							Record: cleanRecord,
						}

						// check fields
						rawFields := options.Query[fieldsQueryParam]
						if rawFields != "" {
							decoded, err := picker.Pick(cleanRecord, rawFields)
							if err == nil {
								data.Record = decoded
							} else {
								app.Logger().Debug(
									"[broadcastRecord] pick fields error",
									slog.String("id", cleanRecord.Id),
									slog.String("collectionName", cleanRecord.Collection().Name),
									slog.String("sub", sub),
									slog.String("fields", rawFields),
									slog.String("error", err.Error()),
								)
							}
						}

						dataBytes, err := json.Marshal(data)
						if err != nil {
							app.Logger().Debug(
								"[broadcastRecord] data marshal error",
								slog.String("id", cleanRecord.Id),
								slog.String("collectionName", cleanRecord.Collection().Name),
								slog.String("error", err.Error()),
							)
							continue
						}

						msg := subscriptions.Message{
							Name: sub,
							Data: dataBytes,
						}

						if dryCache {
							messages, ok := client.Get(dryCacheKey).([]subscriptions.Message)
							if !ok {
								messages = []subscriptions.Message{msg}
							} else {
								messages = append(messages, msg)
							}
							client.Set(dryCacheKey, messages)
						} else {
							routine.FireAndForget(func() {
								client.Send(msg)
							})
						}
					}
				}
			}

			return nil
		})
	}

	return group.Wait()
}

// realtimeBroadcastDryCacheKey broadcasts the dry cached key related messages.
func realtimeBroadcastDryCacheKey(app core.App, key string) error {
	chunks := app.SubscriptionsBroker().ChunkedClients(clientsChunkSize)
	if len(chunks) == 0 {
		return nil // no subscribers
	}

	group := new(errgroup.Group)

	for _, chunk := range chunks {
		group.Go(func() error {
			for _, client := range chunk {
				messages, ok := client.Get(key).([]subscriptions.Message)
				if !ok {
					continue
				}

				client.Unset(key)

				client := client

				routine.FireAndForget(func() {
					for _, msg := range messages {
						client.Send(msg)
					}
				})
			}

			return nil
		})
	}

	return group.Wait()
}

// realtimeUnsetDryCacheKey removes the dry cached key related messages.
func realtimeUnsetDryCacheKey(app core.App, key string) error {
	chunks := app.SubscriptionsBroker().ChunkedClients(clientsChunkSize)
	if len(chunks) == 0 {
		return nil // no subscribers
	}

	group := new(errgroup.Group)

	for _, chunk := range chunks {
		group.Go(func() error {
			for _, client := range chunk {
				if client.Get(key) != nil {
					client.Unset(key)
				}
			}

			return nil
		})
	}

	return group.Wait()
}

func getDryCacheKey(action string, model core.Model) string {
	pkStr, ok := model.PK().(string)
	if !ok {
		pkStr = fmt.Sprintf("%v", model.PK())
	}

	return action + "/" + model.TableName() + "/" + pkStr
}

func isSameAuth(authA, authB *core.Record) bool {
	if authA == nil {
		return authB == nil
	}

	if authB == nil {
		return false
	}

	return authA.Id == authB.Id && authA.Collection().Id == authB.Collection().Id
}

// realtimeCanAccessRecord checks if the subscription client has access to the specified record model.
func realtimeCanAccessRecord(
	app core.App,
	record *core.Record,
	requestInfo *core.RequestInfo,
	accessRule *string,
) bool {
	// check the access rule
	// ---
	if ok, _ := app.CanAccessRecord(record, requestInfo, accessRule); !ok {
		return false
	}

	// check the subscription client-side filter (if any)
	// ---
	filter := requestInfo.Query[search.FilterQueryParam]
	if filter == "" {
		return true // no further checks needed
	}

	err := checkForSuperuserOnlyRuleFields(requestInfo)
	if err != nil {
		return false
	}

	var exists int

	q := app.ConcurrentDB().Select("(1)").
		From(record.Collection().Name).
		AndWhere(dbx.HashExp{record.Collection().Name + ".id": record.Id})

	resolver := core.NewRecordFieldResolver(app, record.Collection(), requestInfo, false)
	expr, err := search.FilterData(filter).BuildExpr(resolver)
	if err != nil {
		return false
	}

	q.AndWhere(expr)

	err = resolver.UpdateQuery(q)
	if err != nil {
		return false
	}

	err = q.Limit(1).Row(&exists)

	return err == nil && exists > 0
}
