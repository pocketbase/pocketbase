package models

import "github.com/pocketbase/pocketbase/models/schema"

var _ Model = (*Collection)(nil)
var _ FilesManager = (*Collection)(nil)

type Collection struct {
	BaseModel

	Name       string        `db:"name" json:"name"`
	System     bool          `db:"system" json:"system"`
	Schema     schema.Schema `db:"schema" json:"schema"`
	ListRule   *string       `db:"listRule" json:"listRule"`
	ViewRule   *string       `db:"viewRule" json:"viewRule"`
	CreateRule *string       `db:"createRule" json:"createRule"`
	UpdateRule *string       `db:"updateRule" json:"updateRule"`
	DeleteRule *string       `db:"deleteRule" json:"deleteRule"`
}

func (m *Collection) TableName() string {
	return "_collections"
}

// BaseFilesPath returns the storage dir path used by the collection.
func (m *Collection) BaseFilesPath() string {
	return m.Id
}
