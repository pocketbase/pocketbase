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

type MailerUserEvent struct {
	MailClient mailer.Mailer
	User       *models.User
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
// User API events data
// -------------------------------------------------------------------

type UsersListEvent struct {
	HttpContext echo.Context
	Users       []*models.User
	Result      *search.Result
}

type UserViewEvent struct {
	HttpContext echo.Context
	User        *models.User
}

type UserCreateEvent struct {
	HttpContext echo.Context
	User        *models.User
}

type UserUpdateEvent struct {
	HttpContext echo.Context
	User        *models.User
}

type UserDeleteEvent struct {
	HttpContext echo.Context
	User        *models.User
}

type UserAuthEvent struct {
	HttpContext echo.Context
	User        *models.User
	Token       string
	Meta        any
}

type UserListExternalAuthsEvent struct {
	HttpContext   echo.Context
	User          *models.User
	ExternalAuths []*models.ExternalAuth
}

type UserUnlinkExternalAuthEvent struct {
	HttpContext  echo.Context
	User         *models.User
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
