package hook

import (
	"github.com/pocketbase/pocketbase/tools/list"
)

// Tagger defines an interface for event data structs that support tags/groups/categories/etc.
// Usually used together with TaggedHook.
type Tagger interface {
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
func (h *TaggedHook[T]) CanTriggerOn(tags []string) bool {
	if len(h.tags) == 0 {
		return true // match all
	}

	for _, t := range tags {
		if list.ExistInSlice(t, h.tags) {
			return true
		}
	}

	return false
}

// PreAdd registers a new handler to the hook by prepending it to the existing queue.
//
// The fn handler will be called only if the event data tags satisfy h.CanTriggerOn.
func (h *TaggedHook[T]) PreAdd(fn Handler[T]) string {
	return h.mainHook.PreAdd(func(e T) error {
		if h.CanTriggerOn(e.Tags()) {
			return fn(e)
		}

		return nil
	})
}

// Add registers a new handler to the hook by appending it to the existing queue.
//
// The fn handler will be called only if the event data tags satisfy h.CanTriggerOn.
func (h *TaggedHook[T]) Add(fn Handler[T]) string {
	return h.mainHook.Add(func(e T) error {
		if h.CanTriggerOn(e.Tags()) {
			return fn(e)
		}

		return nil
	})
}
