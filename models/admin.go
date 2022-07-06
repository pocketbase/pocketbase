package models

var _ Model = (*Admin)(nil)

type Admin struct {
	BaseAccount

	Avatar int `db:"avatar" json:"avatar"`
}

func (m *Admin) TableName() string {
	return "_admins"
}
