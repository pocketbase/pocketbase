package models

var _ Model = (*ExternalAuth)(nil)

type ExternalAuth struct {
	BaseModel

	CollectionId string `db:"collectionId" json:"collectionId"`
	RecordId     string `db:"recordId" json:"recordId"`
	Provider     string `db:"provider" json:"provider"`
	ProviderId   string `db:"providerId" json:"providerId"`
}

func (m *ExternalAuth) TableName() string {
	return "_externalAuths"
}
