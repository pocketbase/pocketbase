package core

var _ optionsValidator = (*collectionBaseOptions)(nil)

// collectionBaseOptions defines the options for the "base" type collection.
type collectionBaseOptions struct {
}

func (o *collectionBaseOptions) validate(cv *collectionValidator) error {
	return nil
}
