// Package models implements all PocketBase DB models and DTOs.
package models

import (
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/pocketbase/pocketbase/tools/types"
)

const (
	// DefaultIdLength is the default length of the generated model id.
	DefaultIdLength = 15

	// DefaultIdAlphabet is the default characters set used for generating the model id.
	DefaultIdAlphabet = "abcdefghijklmnopqrstuvwxyz0123456789"
)

// ColumnValueMapper defines an interface for custom db model data serialization.
type ColumnValueMapper interface {
	// ColumnValueMap returns the data to be used when persisting the model.
	ColumnValueMap() map[string]any
}

// FilesManager defines an interface with common methods that files manager models should implement.
type FilesManager interface {
	// BaseFilesPath returns the storage dir path used by the interface instance.
	BaseFilesPath() string
}

// Model defines an interface with common methods that all db models should have.
type Model interface {
	TableName() string
	IsNew() bool
	MarkAsNew()
	MarkAsNotNew()
	HasId() bool
	GetId() string
	SetId(id string)
	GetCreated() types.DateTime
	GetUpdated() types.DateTime
	RefreshId()
	RefreshCreated()
	RefreshUpdated()
}

// -------------------------------------------------------------------
// BaseModel
// -------------------------------------------------------------------

// BaseModel defines common fields and methods used by all other models.
type BaseModel struct {
	isNotNew bool

	Id      string         `db:"id" json:"id"`
	Created types.DateTime `db:"created" json:"created"`
	Updated types.DateTime `db:"updated" json:"updated"`
}

// HasId returns whether the model has a nonzero id.
func (m *BaseModel) HasId() bool {
	return m.GetId() != ""
}

// GetId returns the model id.
func (m *BaseModel) GetId() string {
	return m.Id
}

// SetId sets the model id to the provided string value.
func (m *BaseModel) SetId(id string) {
	m.Id = id
}

// MarkAsNew marks the model as "new" (aka. enforces m.IsNew() to be true).
func (m *BaseModel) MarkAsNew() {
	m.isNotNew = false
}

// MarkAsNotNew marks the model as "not new" (aka. enforces m.IsNew() to be false)
func (m *BaseModel) MarkAsNotNew() {
	m.isNotNew = true
}

// IsNew indicates what type of db query (insert or update)
// should be used with the model instance.
func (m *BaseModel) IsNew() bool {
	return !m.isNotNew
}

// GetCreated returns the model Created datetime.
func (m *BaseModel) GetCreated() types.DateTime {
	return m.Created
}

// GetUpdated returns the model Updated datetime.
func (m *BaseModel) GetUpdated() types.DateTime {
	return m.Updated
}

// RefreshId generates and sets a new model id.
//
// The generated id is a cryptographically random 15 characters length string.
func (m *BaseModel) RefreshId() {
	m.Id = security.RandomStringWithAlphabet(DefaultIdLength, DefaultIdAlphabet)
}

// RefreshCreated updates the model Created field with the current datetime.
func (m *BaseModel) RefreshCreated() {
	m.Created = types.NowDateTime()
}

// RefreshUpdated updates the model Updated field with the current datetime.
func (m *BaseModel) RefreshUpdated() {
	m.Updated = types.NowDateTime()
}

// PostScan implements the [dbx.PostScanner] interface.
//
// It is executed right after the model was populated with the db row values.
func (m *BaseModel) PostScan() error {
	m.MarkAsNotNew()
	return nil
}
