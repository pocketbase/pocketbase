package core

import (
	"context"
	"net/http"
	"time"

	"github.com/pocketbase/pocketbase/tools/auth"
	"github.com/pocketbase/pocketbase/tools/hook"
	"github.com/pocketbase/pocketbase/tools/mailer"
	"github.com/pocketbase/pocketbase/tools/router"
	"github.com/pocketbase/pocketbase/tools/search"
	"github.com/pocketbase/pocketbase/tools/subscriptions"
	"golang.org/x/crypto/acme/autocert"
)

type HookTagger interface {
	HookTags() []string
}

// -------------------------------------------------------------------

type baseModelEventData struct {
	Model Model
}

func (e *baseModelEventData) Tags() []string {
	if e.Model == nil {
		return nil
	}

	if ht, ok := e.Model.(HookTagger); ok {
		return ht.HookTags()
	}

	return []string{e.Model.TableName()}
}

// -------------------------------------------------------------------

type baseRecordEventData struct {
	Record *Record
}

func (e *baseRecordEventData) Tags() []string {
	if e.Record == nil {
		return nil
	}

	return e.Record.HookTags()
}

// -------------------------------------------------------------------

type baseCollectionEventData struct {
	Collection *Collection
}

func (e *baseCollectionEventData) Tags() []string {
	if e.Collection == nil {
		return nil
	}

	tags := make([]string, 0, 2)

	if e.Collection.Id != "" {
		tags = append(tags, e.Collection.Id)
	}

	if e.Collection.Name != "" {
		tags = append(tags, e.Collection.Name)
	}

	return tags
}

// -------------------------------------------------------------------
// App events data
// -------------------------------------------------------------------

type BootstrapEvent struct {
	hook.Event
	App App
}

type TerminateEvent struct {
	hook.Event
	App       App
	IsRestart bool
}

type BackupEvent struct {
	hook.Event
	App     App
	Context context.Context
	Name    string   // the name of the backup to create/restore.
	Exclude []string // list of dir entries to exclude from the backup create/restore.
}

type ServeEvent struct {
	hook.Event
	App         App
	Router      *router.Router[*RequestEvent]
	Server      *http.Server
	CertManager *autocert.Manager

	// InstallerFunc is the "installer" function that is called after
	// successful server tcp bind but only if there is no explicit
	// superuser record created yet.
	//
	// It runs in a separate goroutine and its default value is [apis.DefaultInstallerFunc].
	//
	// It receives a system superuser record as argument that you can use to generate
	// a short-lived auth token (e.g. systemSuperuser.NewStaticAuthToken(30 * time.Minute))
	// and concatenate it as query param for your installer page
	// (if you are using the client-side SDKs, you can then load the
	// token with pb.authStore.save(token) and perform any Web API request
	// e.g. creating a new superuser).
	//
	// Set it to nil if you want to skip the installer.
	InstallerFunc func(app App, systemSuperuser *Record, baseURL string) error
}

// -------------------------------------------------------------------
// Settings events data
// -------------------------------------------------------------------

type SettingsListRequestEvent struct {
	hook.Event
	*RequestEvent

	Settings *Settings
}

type SettingsUpdateRequestEvent struct {
	hook.Event
	*RequestEvent

	OldSettings *Settings
	NewSettings *Settings
}

type SettingsReloadEvent struct {
	hook.Event
	App App
}

// -------------------------------------------------------------------
// Mailer events data
// -------------------------------------------------------------------

type MailerEvent struct {
	hook.Event
	App App

	Mailer  mailer.Mailer
	Message *mailer.Message
}

type MailerRecordEvent struct {
	MailerEvent
	baseRecordEventData
	Meta map[string]any
}

// -------------------------------------------------------------------
// Model events data
// -------------------------------------------------------------------

const (
	ModelEventTypeCreate   = "create"
	ModelEventTypeUpdate   = "update"
	ModelEventTypeDelete   = "delete"
	ModelEventTypeValidate = "validate"
)

type ModelEvent struct {
	hook.Event
	App App
	baseModelEventData
	Context context.Context

	// Could be any of the ModelEventType* constants, like:
	// - create
	// - update
	// - delete
	// - validate
	Type string
}

type ModelErrorEvent struct {
	Error error
	ModelEvent
}

// -------------------------------------------------------------------
// Record events data
// -------------------------------------------------------------------

type RecordEvent struct {
	hook.Event
	App App
	baseRecordEventData
	Context context.Context

	// Could be any of the ModelEventType* constants, like:
	// - create
	// - update
	// - delete
	// - validate
	Type string
}

type RecordErrorEvent struct {
	Error error
	RecordEvent
}

func syncModelEventWithRecordEvent(me *ModelEvent, re *RecordEvent) {
	me.App = re.App
	me.Context = re.Context
	me.Type = re.Type

	// @todo enable if after profiling doesn't have significant impact
	// 		 skip for now to avoid excessive checks and assume that the
	// 		 Model and the Record fields still points to the same instance
	//
	// if _, ok := me.Model.(*Record); ok {
	// 	me.Model = re.Record
	// } else if proxy, ok := me.Model.(RecordProxy); ok {
	// 	proxy.SetProxyRecord(re.Record)
	// }
}

func syncRecordEventWithModelEvent(re *RecordEvent, me *ModelEvent) {
	re.App = me.App
	re.Context = me.Context
	re.Type = me.Type
}

func newRecordEventFromModelEvent(me *ModelEvent) (*RecordEvent, bool) {
	record, ok := me.Model.(*Record)
	if !ok {
		proxy, ok := me.Model.(RecordProxy)
		if !ok {
			return nil, false
		}
		record = proxy.ProxyRecord()
	}

	re := new(RecordEvent)
	re.App = me.App
	re.Context = me.Context
	re.Type = me.Type
	re.Record = record

	return re, true
}

func newRecordErrorEventFromModelErrorEvent(me *ModelErrorEvent) (*RecordErrorEvent, bool) {
	recordEvent, ok := newRecordEventFromModelEvent(&me.ModelEvent)
	if !ok {
		return nil, false
	}

	re := new(RecordErrorEvent)
	re.RecordEvent = *recordEvent
	re.Error = me.Error

	return re, true
}

func syncModelErrorEventWithRecordErrorEvent(me *ModelErrorEvent, re *RecordErrorEvent) {
	syncModelEventWithRecordEvent(&me.ModelEvent, &re.RecordEvent)
	me.Error = re.Error
}

func syncRecordErrorEventWithModelErrorEvent(re *RecordErrorEvent, me *ModelErrorEvent) {
	syncRecordEventWithModelEvent(&re.RecordEvent, &me.ModelEvent)
	re.Error = me.Error
}

// -------------------------------------------------------------------
// Collection events data
// -------------------------------------------------------------------

type CollectionEvent struct {
	hook.Event
	App App
	baseCollectionEventData
	Context context.Context

	// Could be any of the ModelEventType* constants, like:
	// - create
	// - update
	// - delete
	// - validate
	Type string
}

type CollectionErrorEvent struct {
	Error error
	CollectionEvent
}

func syncModelEventWithCollectionEvent(me *ModelEvent, ce *CollectionEvent) {
	me.App = ce.App
	me.Context = ce.Context
	me.Type = ce.Type
	me.Model = ce.Collection
}

func syncCollectionEventWithModelEvent(ce *CollectionEvent, me *ModelEvent) {
	ce.App = me.App
	ce.Context = me.Context
	ce.Type = me.Type
	if c, ok := me.Model.(*Collection); ok {
		ce.Collection = c
	}
}

func newCollectionEventFromModelEvent(me *ModelEvent) (*CollectionEvent, bool) {
	record, ok := me.Model.(*Collection)
	if !ok {
		return nil, false
	}

	ce := new(CollectionEvent)
	ce.App = me.App
	ce.Context = me.Context
	ce.Type = me.Type
	ce.Collection = record

	return ce, true
}

func newCollectionErrorEventFromModelErrorEvent(me *ModelErrorEvent) (*CollectionErrorEvent, bool) {
	collectionevent, ok := newCollectionEventFromModelEvent(&me.ModelEvent)
	if !ok {
		return nil, false
	}

	ce := new(CollectionErrorEvent)
	ce.CollectionEvent = *collectionevent
	ce.Error = me.Error

	return ce, true
}

func syncModelErrorEventWithCollectionErrorEvent(me *ModelErrorEvent, ce *CollectionErrorEvent) {
	syncModelEventWithCollectionEvent(&me.ModelEvent, &ce.CollectionEvent)
	me.Error = ce.Error
}

func syncCollectionErrorEventWithModelErrorEvent(ce *CollectionErrorEvent, me *ModelErrorEvent) {
	syncCollectionEventWithModelEvent(&ce.CollectionEvent, &me.ModelEvent)
	ce.Error = me.Error
}

// -------------------------------------------------------------------
// File API events data
// -------------------------------------------------------------------

type FileTokenRequestEvent struct {
	hook.Event
	*RequestEvent
	baseRecordEventData

	Token string
}

type FileDownloadRequestEvent struct {
	hook.Event
	*RequestEvent
	baseCollectionEventData

	Record     *Record
	FileField  *FileField
	ServedPath string
	ServedName string
}

// -------------------------------------------------------------------
// Collection API events data
// -------------------------------------------------------------------

type CollectionsListRequestEvent struct {
	hook.Event
	*RequestEvent

	Collections []*Collection
	Result      *search.Result
}

type CollectionsImportRequestEvent struct {
	hook.Event
	*RequestEvent

	CollectionsData []map[string]any
	DeleteMissing   bool
}

type CollectionRequestEvent struct {
	hook.Event
	*RequestEvent
	baseCollectionEventData
}

// -------------------------------------------------------------------
// Realtime API events data
// -------------------------------------------------------------------

type RealtimeConnectRequestEvent struct {
	hook.Event
	*RequestEvent

	Client subscriptions.Client

	// note: modifying it after the connect has no effect
	IdleTimeout time.Duration
}

type RealtimeMessageEvent struct {
	hook.Event
	*RequestEvent

	Client  subscriptions.Client
	Message *subscriptions.Message
}

type RealtimeSubscribeRequestEvent struct {
	hook.Event
	*RequestEvent

	Client        subscriptions.Client
	Subscriptions []string
}

// -------------------------------------------------------------------
// Record CRUD API events data
// -------------------------------------------------------------------

type RecordsListRequestEvent struct {
	hook.Event
	*RequestEvent
	baseCollectionEventData

	// @todo consider removing and maybe add as generic to the search.Result?
	Records []*Record
	Result  *search.Result
}

type RecordRequestEvent struct {
	hook.Event
	*RequestEvent
	baseCollectionEventData

	Record *Record
}

type RecordEnrichEvent struct {
	hook.Event
	App App
	baseRecordEventData

	RequestInfo *RequestInfo
}

// -------------------------------------------------------------------
// Auth Record API events data
// -------------------------------------------------------------------

type RecordCreateOTPRequestEvent struct {
	hook.Event
	*RequestEvent
	baseCollectionEventData

	Record   *Record
	Password string
}

type RecordAuthWithOTPRequestEvent struct {
	hook.Event
	*RequestEvent
	baseCollectionEventData

	Record *Record
	OTP    *OTP
}

type RecordAuthRequestEvent struct {
	hook.Event
	*RequestEvent
	baseCollectionEventData

	Record     *Record
	Token      string
	Meta       any
	AuthMethod string
}

type RecordAuthWithPasswordRequestEvent struct {
	hook.Event
	*RequestEvent
	baseCollectionEventData

	Record        *Record
	Identity      string
	IdentityField string
	Password      string
}

type RecordAuthWithOAuth2RequestEvent struct {
	hook.Event
	*RequestEvent
	baseCollectionEventData

	ProviderName   string
	ProviderClient auth.Provider
	Record         *Record
	OAuth2User     *auth.AuthUser
	CreateData     map[string]any
	IsNewRecord    bool
}

type RecordAuthRefreshRequestEvent struct {
	hook.Event
	*RequestEvent
	baseCollectionEventData

	Record *Record
}

type RecordRequestPasswordResetRequestEvent struct {
	hook.Event
	*RequestEvent
	baseCollectionEventData

	Record *Record
}

type RecordConfirmPasswordResetRequestEvent struct {
	hook.Event
	*RequestEvent
	baseCollectionEventData

	Record *Record
}

type RecordRequestVerificationRequestEvent struct {
	hook.Event
	*RequestEvent
	baseCollectionEventData

	Record *Record
}

type RecordConfirmVerificationRequestEvent struct {
	hook.Event
	*RequestEvent
	baseCollectionEventData

	Record *Record
}

type RecordRequestEmailChangeRequestEvent struct {
	hook.Event
	*RequestEvent
	baseCollectionEventData

	Record   *Record
	NewEmail string
}

type RecordConfirmEmailChangeRequestEvent struct {
	hook.Event
	*RequestEvent
	baseCollectionEventData

	Record   *Record
	NewEmail string
}
