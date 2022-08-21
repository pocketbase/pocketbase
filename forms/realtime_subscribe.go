package forms

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// RealtimeSubscribe is a realtime subscriptions request form.
type RealtimeSubscribe struct {
	ClientId      string   `form:"clientId" json:"clientId"`
	Subscriptions []string `form:"subscriptions" json:"subscriptions"`
}

// NewRealtimeSubscribe creates new RealtimeSubscribe request form.
func NewRealtimeSubscribe() *RealtimeSubscribe {
	return &RealtimeSubscribe{}
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *RealtimeSubscribe) Validate() error {
	return validation.ValidateStruct(form,
		validation.Field(&form.ClientId, validation.Required, validation.Length(1, 255)),
	)
}
