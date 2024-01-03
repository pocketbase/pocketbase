package models

import "github.com/pocketbase/pocketbase/tools/types"

var _ Model = (*Request)(nil)

// list with the supported values for `Request.Auth`
const (
	RequestAuthGuest  = "guest"
	RequestAuthAdmin  = "admin"
	RequestAuthRecord = "authRecord"
)

// Deprecated: Replaced by the Log model and will be removed in a future version.
type Request struct {
	BaseModel

	Url       string        `db:"url" json:"url"`
	Method    string        `db:"method" json:"method"`
	Status    int           `db:"status" json:"status"`
	Auth      string        `db:"auth" json:"auth"`
	UserIp    string        `db:"userIp" json:"userIp"`
	RemoteIp  string        `db:"remoteIp" json:"remoteIp"`
	Referer   string        `db:"referer" json:"referer"`
	UserAgent string        `db:"userAgent" json:"userAgent"`
	Meta      types.JsonMap `db:"meta" json:"meta"`
}

func (m *Request) TableName() string {
	return "_requests"
}
