package core

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/models/settings"
	"github.com/pocketbase/pocketbase/tools/auth"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/pocketbase/pocketbase/tools/hook"
	"github.com/pocketbase/pocketbase/tools/mailer"
	"github.com/pocketbase/pocketbase/tools/search"
	"github.com/pocketbase/pocketbase/tools/subscriptions"
	"golang.org/x/crypto/acme/autocert"
)

var (
	_ hook.Tagger = (*BaseModelEvent)(nil)
	_ hook.Tagger = (*BaseCollectionEvent)(nil)
)

type BaseModelEvent struct {
	Model models.Model
}

func (e *BaseModelEvent) Tags() []string {
	if e.Model == nil {
		return nil
	}

	if r, ok := e.Model.(*models.Record); ok && r.Collection() != nil {
		return []string{r.Collection().Id, r.Collection().Name}
	}

	return []string{e.Model.TableName()}
}

type BaseCollectionEvent struct {
	Collection *models.Collection
}

func (e *BaseCollectionEvent) Tags() []string {
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
// Serve events data
// -------------------------------------------------------------------

type BootstrapEvent struct {
	App App
}

type TerminateEvent struct {
	App       App
	IsRestart bool
}

type ServeEvent struct {
	App         App
	Router      *echo.Echo
	Server      *http.Server
	CertManager *autocert.Manager
}

type ApiErrorEvent struct {
	HttpContext echo.Context
	Error       error
}

// -------------------------------------------------------------------
// Model DAO events data
// -------------------------------------------------------------------

type ModelEvent struct {
	BaseModelEvent

	Dao *daos.Dao
}

// -------------------------------------------------------------------
// Mailer events data
// -------------------------------------------------------------------

type MailerRecordEvent struct {
	BaseCollectionEvent

	MailClient mailer.Mailer
	Message    *mailer.Message
	Record     *models.Record
	Meta       map[string]any
}

type MailerAdminEvent struct {
	MailClient mailer.Mailer
	Message    *mailer.Message
	Admin      *models.Admin
	Meta       map[string]any
}

// -------------------------------------------------------------------
// Realtime API events data
// -------------------------------------------------------------------

type RealtimeConnectEvent struct {
	HttpContext echo.Context
	Client      subscriptions.Client
	IdleTimeout time.Duration
}

type RealtimeDisconnectEvent struct {
	HttpContext echo.Context
	Client      subscriptions.Client
}

type RealtimeMessageEvent struct {
	HttpContext echo.Context
	Client      subscriptions.Client
	Message     *subscriptions.Message
}

type RealtimeSubscribeEvent struct {
	HttpContext   echo.Context
	Client        subscriptions.Client
	Subscriptions []string
}

// -------------------------------------------------------------------
// Settings API events data
// -------------------------------------------------------------------

type SettingsListEvent struct {
	HttpContext      echo.Context
	RedactedSettings *settings.Settings
}

type SettingsUpdateEvent struct {
	HttpContext echo.Context
	OldSettings *settings.Settings
	NewSettings *settings.Settings
}

// -------------------------------------------------------------------
// Record CRUD API events data
// -------------------------------------------------------------------

type RecordsListEvent struct {
	BaseCollectionEvent

	HttpContext echo.Context
	Records     []*models.Record
	Result      *search.Result
}

type RecordViewEvent struct {
	BaseCollectionEvent

	HttpContext echo.Context
	Record      *models.Record
}

type RecordCreateEvent struct {
	BaseCollectionEvent

	HttpContext   echo.Context
	Record        *models.Record
	UploadedFiles map[string][]*filesystem.File
}

type RecordUpdateEvent struct {
	BaseCollectionEvent

	HttpContext   echo.Context
	Record        *models.Record
	UploadedFiles map[string][]*filesystem.File
}

type RecordDeleteEvent struct {
	BaseCollectionEvent

	HttpContext echo.Context
	Record      *models.Record
}

// -------------------------------------------------------------------
// Auth Record API events data
// -------------------------------------------------------------------

type RecordAuthEvent struct {
	BaseCollectionEvent

	HttpContext echo.Context
	Record      *models.Record
	Token       string
	Meta        any
}

type RecordAuthWithPasswordEvent struct {
	BaseCollectionEvent

	HttpContext echo.Context
	Record      *models.Record
	Identity    string
	Password    string
}

type RecordAuthWithOAuth2Event struct {
	BaseCollectionEvent

	HttpContext    echo.Context
	ProviderName   string
	ProviderClient auth.Provider
	Record         *models.Record
	OAuth2User     *auth.AuthUser
	IsNewRecord    bool
}

type RecordAuthRefreshEvent struct {
	BaseCollectionEvent

	HttpContext echo.Context
	Record      *models.Record
}

type RecordRequestPasswordResetEvent struct {
	BaseCollectionEvent

	HttpContext echo.Context
	Record      *models.Record
}

type RecordConfirmPasswordResetEvent struct {
	BaseCollectionEvent

	HttpContext echo.Context
	Record      *models.Record
}

type RecordRequestVerificationEvent struct {
	BaseCollectionEvent

	HttpContext echo.Context
	Record      *models.Record
}

type RecordConfirmVerificationEvent struct {
	BaseCollectionEvent

	HttpContext echo.Context
	Record      *models.Record
}

type RecordRequestEmailChangeEvent struct {
	BaseCollectionEvent

	HttpContext echo.Context
	Record      *models.Record
}

type RecordConfirmEmailChangeEvent struct {
	BaseCollectionEvent

	HttpContext echo.Context
	Record      *models.Record
}

type RecordListExternalAuthsEvent struct {
	BaseCollectionEvent

	HttpContext   echo.Context
	Record        *models.Record
	ExternalAuths []*models.ExternalAuth
}

type RecordUnlinkExternalAuthEvent struct {
	BaseCollectionEvent

	HttpContext  echo.Context
	Record       *models.Record
	ExternalAuth *models.ExternalAuth
}

// -------------------------------------------------------------------
// Admin API events data
// -------------------------------------------------------------------

type AdminsListEvent struct {
	HttpContext echo.Context
	Admins      []*models.Admin
	Result      *search.Result
}

type AdminViewEvent struct {
	HttpContext echo.Context
	Admin       *models.Admin
}

type AdminCreateEvent struct {
	HttpContext echo.Context
	Admin       *models.Admin
}

type AdminUpdateEvent struct {
	HttpContext echo.Context
	Admin       *models.Admin
}

type AdminDeleteEvent struct {
	HttpContext echo.Context
	Admin       *models.Admin
}

type AdminAuthEvent struct {
	HttpContext echo.Context
	Admin       *models.Admin
	Token       string
}

type AdminAuthWithPasswordEvent struct {
	HttpContext echo.Context
	Admin       *models.Admin
	Identity    string
	Password    string
}

type AdminAuthRefreshEvent struct {
	HttpContext echo.Context
	Admin       *models.Admin
}

type AdminRequestPasswordResetEvent struct {
	HttpContext echo.Context
	Admin       *models.Admin
}

type AdminConfirmPasswordResetEvent struct {
	HttpContext echo.Context
	Admin       *models.Admin
}

// -------------------------------------------------------------------
// Collection API events data
// -------------------------------------------------------------------

type CollectionsListEvent struct {
	HttpContext echo.Context
	Collections []*models.Collection
	Result      *search.Result
}

type CollectionViewEvent struct {
	BaseCollectionEvent

	HttpContext echo.Context
}

type CollectionCreateEvent struct {
	BaseCollectionEvent

	HttpContext echo.Context
}

type CollectionUpdateEvent struct {
	BaseCollectionEvent

	HttpContext echo.Context
}

type CollectionDeleteEvent struct {
	BaseCollectionEvent

	HttpContext echo.Context
}

type CollectionsImportEvent struct {
	HttpContext echo.Context
	Collections []*models.Collection
}

// -------------------------------------------------------------------
// File API events data
// -------------------------------------------------------------------

type FileTokenEvent struct {
	BaseModelEvent

	HttpContext echo.Context
	Token       string
}

type FileDownloadEvent struct {
	BaseCollectionEvent

	HttpContext echo.Context
	Record      *models.Record
	FileField   *schema.SchemaField
	ServedPath  string
	ServedName  string
}
