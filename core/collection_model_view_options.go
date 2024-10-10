package core

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var _ optionsValidator = (*collectionViewOptions)(nil)

// collectionViewOptions defines the options for the "view" type collection.
type collectionViewOptions struct {
	ViewQuery string `form:"viewQuery" json:"viewQuery"`
}

func (o *collectionViewOptions) validate(cv *collectionValidator) error {
	return validation.ValidateStruct(o,
		validation.Field(&o.ViewQuery, validation.Required, validation.By(cv.checkViewQuery)),
	)
}
