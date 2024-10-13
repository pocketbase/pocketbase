package hook

import (
	"github.com/pocketbase/pocketbase/tools/list"
)

// Tagger defines an interface for event data structs that support tags/groups/categories/etc.
// Usually used together with TaggedHook.
type Tagger interface {
	Resolver

	Tags() []string
}

// wrapped local Hook embedded struct to limit the public API surface.
type mainHook[T Tagger] struct {
	*Hook[T]
}

// NewTaggedHook creates a new TaggedHook with the provided main hook and optional tags.
func NewTaggedHook[T Tagger](hook *Hook[T], tags ...string) *TaggedHook[T] {
	return &TaggedHook[T]{
		mainHook[T]{hook},
		tags,
	}
}

// TaggedHook defines a proxy hook which register handlers that are triggered only
// if the TaggedHook.tags are empty or includes at least one of the event data tag(s).
type TaggedHook[T Tagger] struct {
	mainHook[T]

	tags []string
}

// CanTriggerOn checks if the current TaggedHook can be triggered with
// the provided event data tags.
//
// It returns always true if the hook doens't have any tags.
func (h *TaggedHook[T]) CanTriggerOn(tagsToCheck []string) bool {
	if len(h.tags) == 0 {
		return true // match all
	}

	for _, t := range tagsToCheck {
		if list.ExistInSlice(t, h.tags) {
			return true
		}
	}

	return false
}

// Bind registers the provided handler to the current hooks queue.
//
// It is similar to [Hook.Bind] with the difference that the handler
// function is invoked only if the event data tags satisfy h.CanTriggerOn.
func (h *TaggedHook[T]) Bind(handler *Handler[T]) string {
	fn := handler.Func

	handler.Func = func(e T) error {
		if h.CanTriggerOn(e.Tags()) {
			return fn(e)
		}

		return e.Next()
	}

	return h.mainHook.Bind(handler)
}

// BindFunc registers a new handler with the specified function.
//
// It is similar to [Hook.Bind] with the difference that the handler
// function is invoked only if the event data tags satisfy h.CanTriggerOn.
func (h *TaggedHook[T]) BindFunc(fn func(e T) error) string {
	return h.mainHook.BindFunc(func(e T) error {
		if h.CanTriggerOn(e.Tags()) {
			return fn(e)
		}

		return e.Next()
	})
}
