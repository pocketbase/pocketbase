package apis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/resolvers"
	"github.com/pocketbase/pocketbase/tools/search"
	"github.com/pocketbase/pocketbase/tools/subscriptions"
)

// bindRealtimeApi registers the realtime api endpoints.
func bindRealtimeApi(app core.App, rg *echo.Group) {
	api := realtimeApi{app: app}

	subGroup := rg.Group("/realtime", ActivityLogger(app))
	subGroup.GET("", api.connect)
	subGroup.POST("", api.setSubscriptions)

	api.bindEvents()
}

type realtimeApi struct {
	app core.App
}

func (api *realtimeApi) connect(c echo.Context) error {
	cancelCtx, cancelRequest := context.WithCancel(c.Request().Context())
	defer cancelRequest()
	c.SetRequest(c.Request().Clone(cancelCtx))

	// register new subscription client
	client := subscriptions.NewDefaultClient()
	api.app.SubscriptionsBroker().Register(client)
	defer api.app.SubscriptionsBroker().Unregister(client.Id())

	c.Response().Header().Set("Content-Type", "text/event-stream; charset=UTF-8")
	c.Response().Header().Set("Cache-Control", "no-store")
	c.Response().Header().Set("Connection", "keep-alive")
	// https://github.com/pocketbase/pocketbase/discussions/480#discussioncomment-3657640
	// https://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_buffering
	c.Response().Header().Set("X-Accel-Buffering", "no")

	event := &core.RealtimeConnectEvent{
		HttpContext: c,
		Client:      client,
	}

	if err := api.app.OnRealtimeConnectRequest().Trigger(event); err != nil {
		return err
	}

	if api.app.IsDebug() {
		log.Printf("Realtime connection established: %s\n", client.Id())
	}

	// signalize established connection (aka. fire "connect" message)
	fmt.Fprint(c.Response(), "id:"+client.Id()+"\n")
	fmt.Fprint(c.Response(), "event:PB_CONNECT\n")
	fmt.Fprint(c.Response(), "data:{\"clientId\":\""+client.Id()+"\"}\n\n")
	c.Response().Flush()

	// start an idle timer to keep track of inactive/forgotten connections
	idleDuration := 5 * time.Minute
	idleTimer := time.NewTimer(idleDuration)
	defer idleTimer.Stop()

	for {
		select {
		case <-idleTimer.C:
			cancelRequest()
		case msg, ok := <-client.Channel():
			if !ok {
				// channel is closed
				if api.app.IsDebug() {
					log.Println("Realtime connection closed (closed channel):", client.Id())
				}
				return nil
			}

			w := c.Response()
			fmt.Fprint(w, "id:"+client.Id()+"\n")
			fmt.Fprint(w, "event:"+msg.Name+"\n")
			fmt.Fprint(w, "data:"+msg.Data+"\n\n")
			w.Flush()

			idleTimer.Stop()
			idleTimer.Reset(idleDuration)
		case <-c.Request().Context().Done():
			// connection is closed
			if api.app.IsDebug() {
				log.Println("Realtime connection closed (cancelled request):", client.Id())
			}
			return nil
		}
	}
}

// note: in case of reconnect, clients will have to resubmit all subscriptions again
func (api *realtimeApi) setSubscriptions(c echo.Context) error {
	form := forms.NewRealtimeSubscribe()

	// read request data
	if err := c.Bind(form); err != nil {
		return NewBadRequestError("", err)
	}

	// validate request data
	if err := form.Validate(); err != nil {
		return NewBadRequestError("", err)
	}

	// find subscription client
	client, err := api.app.SubscriptionsBroker().ClientById(form.ClientId)
	if err != nil {
		return NewNotFoundError("Missing or invalid client id.", err)
	}

	// check if the previous request was authorized
	oldAuthId := extractAuthIdFromGetter(client)
	newAuthId := extractAuthIdFromGetter(c)
	if oldAuthId != "" && oldAuthId != newAuthId {
		return NewForbiddenError("The current and the previous request authorization don't match.", nil)
	}

	event := &core.RealtimeSubscribeEvent{
		HttpContext:   c,
		Client:        client,
		Subscriptions: form.Subscriptions,
	}

	handlerErr := api.app.OnRealtimeBeforeSubscribeRequest().Trigger(event, func(e *core.RealtimeSubscribeEvent) error {
		// update auth state
		e.Client.Set(ContextAdminKey, e.HttpContext.Get(ContextAdminKey))
		e.Client.Set(ContextAuthRecordKey, e.HttpContext.Get(ContextAuthRecordKey))

		// unsubscribe from any previous existing subscriptions
		e.Client.Unsubscribe()

		// subscribe to the new subscriptions
		e.Client.Subscribe(e.Subscriptions...)

		return e.HttpContext.NoContent(http.StatusNoContent)
	})

	if handlerErr == nil {
		api.app.OnRealtimeAfterSubscribeRequest().Trigger(event)
	}

	return handlerErr
}

// updateClientsAuthModel updates the existing clients auth model with the new one (matched by ID).
func (api *realtimeApi) updateClientsAuthModel(contextKey string, newModel models.Model) error {
	for _, client := range api.app.SubscriptionsBroker().Clients() {
		clientModel, _ := client.Get(contextKey).(models.Model)
		if clientModel != nil && clientModel.GetId() == newModel.GetId() {
			client.Set(contextKey, newModel)
		}
	}

	return nil
}

// unregisterClientsByAuthModel unregister all clients that has the provided auth model.
func (api *realtimeApi) unregisterClientsByAuthModel(contextKey string, model models.Model) error {
	for _, client := range api.app.SubscriptionsBroker().Clients() {
		clientModel, _ := client.Get(contextKey).(models.Model)
		if clientModel != nil && clientModel.GetId() == model.GetId() {
			api.app.SubscriptionsBroker().Unregister(client.Id())
		}
	}

	return nil
}

func (api *realtimeApi) bindEvents() {
	// update the clients that has admin or auth record association
	api.app.OnModelAfterUpdate().PreAdd(func(e *core.ModelEvent) error {
		if record, ok := e.Model.(*models.Record); ok && record != nil && record.Collection().IsAuth() {
			return api.updateClientsAuthModel(ContextAuthRecordKey, record)
		}

		if admin, ok := e.Model.(*models.Admin); ok && admin != nil {
			return api.updateClientsAuthModel(ContextAdminKey, admin)
		}

		return nil
	})

	// remove the client(s) associated to the deleted admin or auth record
	api.app.OnModelAfterDelete().PreAdd(func(e *core.ModelEvent) error {
		if record, ok := e.Model.(*models.Record); ok && record != nil && record.Collection().IsAuth() {
			return api.unregisterClientsByAuthModel(ContextAuthRecordKey, record)
		}

		if admin, ok := e.Model.(*models.Admin); ok && admin != nil {
			return api.unregisterClientsByAuthModel(ContextAdminKey, admin)
		}

		return nil
	})

	api.app.OnModelAfterCreate().PreAdd(func(e *core.ModelEvent) error {
		if record, ok := e.Model.(*models.Record); ok {
			api.broadcastRecord("create", record)
		}
		return nil
	})

	api.app.OnModelAfterUpdate().PreAdd(func(e *core.ModelEvent) error {
		if record, ok := e.Model.(*models.Record); ok {
			api.broadcastRecord("update", record)
		}
		return nil
	})

	api.app.OnModelBeforeDelete().Add(func(e *core.ModelEvent) error {
		if record, ok := e.Model.(*models.Record); ok {
			api.broadcastRecord("delete", record)
		}
		return nil
	})
}

func (api *realtimeApi) canAccessRecord(client subscriptions.Client, record *models.Record, accessRule *string) bool {
	admin, _ := client.Get(ContextAdminKey).(*models.Admin)
	if admin != nil {
		// admins can access everything
		return true
	}

	if accessRule == nil {
		// only admins can access this record
		return false
	}

	ruleFunc := func(q *dbx.SelectQuery) error {
		if *accessRule == "" {
			return nil // empty public rule
		}

		// emulate request data
		requestData := &models.RequestData{
			Method: "GET",
		}
		requestData.AuthRecord, _ = client.Get(ContextAuthRecordKey).(*models.Record)

		resolver := resolvers.NewRecordFieldResolver(api.app.Dao(), record.Collection(), requestData, true)
		expr, err := search.FilterData(*accessRule).BuildExpr(resolver)
		if err != nil {
			return err
		}
		resolver.UpdateQuery(q)
		q.AndWhere(expr)

		return nil
	}

	foundRecord, err := api.app.Dao().FindRecordById(record.Collection().Id, record.Id, ruleFunc)
	if err == nil && foundRecord != nil {
		return true
	}

	return false
}

type recordData struct {
	Action string         `json:"action"`
	Record *models.Record `json:"record"`
}

func (api *realtimeApi) broadcastRecord(action string, record *models.Record) error {
	collection := record.Collection()
	if collection == nil {
		return errors.New("Record collection not set.")
	}

	clients := api.app.SubscriptionsBroker().Clients()
	if len(clients) == 0 {
		return nil // no subscribers
	}

	// remove the expand from the broadcasted record because we don't
	// know if the clients have access to view the expanded records
	cleanRecord := *record
	cleanRecord.SetExpand(nil)
	cleanRecord.WithUnkownData(false)
	cleanRecord.IgnoreEmailVisibility(false)

	subscriptionRuleMap := map[string]*string{
		(collection.Name + "/" + cleanRecord.Id): collection.ViewRule,
		(collection.Id + "/" + cleanRecord.Id):   collection.ViewRule,
		(collection.Name + "/*"):                 collection.ListRule,
		(collection.Id + "/*"):                   collection.ListRule,
		// @deprecated: the same as the wildcard topic but kept for backward compatibility
		collection.Name: collection.ListRule,
		collection.Id:   collection.ListRule,
	}

	data := &recordData{
		Action: action,
		Record: &cleanRecord,
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		if api.app.IsDebug() {
			log.Println(err)
		}
		return err
	}

	encodedData := string(dataBytes)

	for _, client := range clients {
		for subscription, rule := range subscriptionRuleMap {
			if !client.HasSubscription(subscription) {
				continue
			}

			if !api.canAccessRecord(client, data.Record, rule) {
				continue
			}

			msg := subscriptions.Message{
				Name: subscription,
				Data: encodedData,
			}

			// ignore the auth record email visibility checks for
			// auth owner, admin or manager
			if collection.IsAuth() {
				authId := extractAuthIdFromGetter(client)
				if authId == data.Record.Id ||
					api.canAccessRecord(client, data.Record, collection.AuthOptions().ManageRule) {
					data.Record.IgnoreEmailVisibility(true) // ignore
					if newData, err := json.Marshal(data); err == nil {
						msg.Data = string(newData)
					}
					data.Record.IgnoreEmailVisibility(false) // restore
				}
			}

			client.Channel() <- msg
		}
	}

	return nil
}

type getter interface {
	Get(string) any
}

func extractAuthIdFromGetter(val getter) string {
	record, _ := val.Get(ContextAuthRecordKey).(*models.Record)
	if record != nil {
		return record.Id
	}

	admin, _ := val.Get(ContextAdminKey).(*models.Admin)
	if admin != nil {
		return admin.Id
	}

	return ""
}
