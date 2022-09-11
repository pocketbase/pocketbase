package models

var _ Model = (*ExternalAuth)(nil)

type ExternalAuth struct {
	BaseModel

	UserId     string `db:"userId" json:"userId"`
	Provider   string `db:"provider" json:"provider"`
	ProviderId string `db:"providerId" json:"providerId"`
}

func (m *ExternalAuth) TableName() string {
	return "_externalAuths"
}
