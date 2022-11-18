package core

import (
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/mailer"
	"github.com/pocketbase/pocketbase/tools/search"
	"github.com/pocketbase/pocketbase/tools/subscriptions"

	"github.com/labstack/echo/v5"
)

// -------------------------------------------------------------------
// Serve events data
// -------------------------------------------------------------------

type ServeEvent struct {
	App    App
	Router *echo.Echo
}

// -------------------------------------------------------------------
// Model DAO events data
// -------------------------------------------------------------------

type ModelEvent struct {
	Dao   *daos.Dao
	Model models.Model
}

// -------------------------------------------------------------------
// Mailer events data
// -------------------------------------------------------------------

type MailerRecordEvent struct {
	MailClient mailer.Mailer
	Record     *models.Record
	Meta       map[string]any
}

type MailerAdminEvent struct {
	MailClient mailer.Mailer
	Admin      *models.Admin
	Meta       map[string]any
}

// -------------------------------------------------------------------
// Realtime API events data
// -------------------------------------------------------------------

type RealtimeConnectEvent struct {
	HttpContext echo.Context
	Client      subscriptions.Client
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
	RedactedSettings *Settings
}

type SettingsUpdateEvent struct {
	HttpContext echo.Context
	OldSettings *Settings
	NewSettings *Settings
}

// -------------------------------------------------------------------
// Record API events data
// -------------------------------------------------------------------

type RecordsListEvent struct {
	HttpContext echo.Context
	Collection  *models.Collection
	Records     []*models.Record
	Result      *search.Result
}

type RecordViewEvent struct {
	HttpContext echo.Context
	Record      *models.Record
}

type RecordCreateEvent struct {
	HttpContext echo.Context
	Record      *models.Record
}

type RecordUpdateEvent struct {
	HttpContext echo.Context
	Record      *models.Record
}

type RecordDeleteEvent struct {
	HttpContext echo.Context
	Record      *models.Record
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

// -------------------------------------------------------------------
// Auth Record API events data
// -------------------------------------------------------------------

type RecordAuthEvent struct {
	HttpContext echo.Context
	Record      *models.Record
	Token       string
	Meta        any
}

type RecordListExternalAuthsEvent struct {
	HttpContext   echo.Context
	Record        *models.Record
	ExternalAuths []*models.ExternalAuth
}

type RecordUnlinkExternalAuthEvent struct {
	HttpContext  echo.Context
	Record       *models.Record
	ExternalAuth *models.ExternalAuth
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
	HttpContext echo.Context
	Collection  *models.Collection
}

type CollectionCreateEvent struct {
	HttpContext echo.Context
	Collection  *models.Collection
}

type CollectionUpdateEvent struct {
	HttpContext echo.Context
	Collection  *models.Collection
}

type CollectionDeleteEvent struct {
	HttpContext echo.Context
	Collection  *models.Collection
}

type CollectionsImportEvent struct {
	HttpContext echo.Context
	Collections []*models.Collection
}

// -------------------------------------------------------------------
// View API events data
// -------------------------------------------------------------------

type ViewListEvent struct {
	HttpContext echo.Context
	Result      *search.Result
}

type ViewViewEvent struct {
	HttpContext echo.Context
	View        *models.View
}
type ViewCreateEvent struct {
	HttpContext echo.Context
	View        *models.View
}
type ViewUpdateEvent struct {
	HttpContext echo.Context
	View        *models.View
}

type ViewDeleteEvent struct {
	HttpContext echo.Context
	View        *models.View
}
type RecordsFromViewListEvent struct {
	HttpContext echo.Context
	View        *models.View
	Result      *search.Result
}

// -------------------------------------------------------------------
// File API events data
// -------------------------------------------------------------------

type FileDownloadEvent struct {
	HttpContext echo.Context
	Collection  *models.Collection
	Record      *models.Record
	FileField   *schema.SchemaField
	ServedPath  string
	ServedName  string
}
