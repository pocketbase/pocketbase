package models

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/tools/types"
)

var _ Model = (*User)(nil)

const (
	// ProfileCollectionName is the name of the system user profiles collection.
	ProfileCollectionName = "profiles"

	// ProfileCollectionUserFieldName is the name of the user field from the system user profiles collection.
	ProfileCollectionUserFieldName = "userId"
)

type User struct {
	BaseAccount

	Verified               bool           `db:"verified" json:"verified"`
	LastVerificationSentAt types.DateTime `db:"lastVerificationSentAt" json:"lastVerificationSentAt"`

	// profile rel
	Profile *Record `db:"-" json:"profile"`
}

func (m *User) TableName() string {
	return "_users"
}

// AsMap returns the current user data as a plain map
// (including the profile relation, if loaded).
func (m *User) AsMap() (map[string]any, error) {
	userBytes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	result := map[string]any{}
	if err := json.Unmarshal(userBytes, &result); err != nil {
		return nil, err
	}

	return result, nil
}
