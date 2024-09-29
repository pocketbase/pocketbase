package hook

// Resolver defines a common interface for a Hook event (see [Event]).
type Resolver interface {
	// Next triggers the next handler in the hook's chain (if any).
	Next() error

	// note: kept only for the generic interface; may get removed in the future
	nextFunc() func() error
	setNextFunc(f func() error)
}

var _ Resolver = (*Event)(nil)

// Event implements [Resolver] and it is intended to be used as a base
// Hook event that you can embed in your custom typed event structs.
//
// Example:
//
//	type CustomEvent struct {
//		hook.Event
//
//		SomeField int
//	}
type Event struct {
	next func() error
}

// Next calls the next hook handler.
func (e *Event) Next() error {
	if e.next != nil {
		return e.next()
	}
	return nil
}

// nextFunc returns the function that Next calls.
func (e *Event) nextFunc() func() error {
	return e.next
}

// setNextFunc sets the function that Next calls.
func (e *Event) setNextFunc(f func() error) {
	e.next = f
}
