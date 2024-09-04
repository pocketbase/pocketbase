package hook

import (
	"errors"
	"fmt"
	"sync"

	"github.com/thinkonmay/pocketbase/tools/security"
)

var StopPropagation = errors.New("Event hook propagation stopped")

// Handler defines a hook handler function.
type Handler[T any] func(e T) error

// handlerPair defines a pair of string id and Handler.
type handlerPair[T any] struct {
	id      string
	handler Handler[T]
}

// Hook defines a concurrent safe structure for handling event hooks
// (aka. callbacks propagation).
type Hook[T any] struct {
	mux      sync.RWMutex
	handlers []*handlerPair[T]
}

// PreAdd registers a new handler to the hook by prepending it to the existing queue.
//
// Returns an autogenerated hook id that could be used later to remove the hook with Hook.Remove(id).
func (h *Hook[T]) PreAdd(fn Handler[T]) string {
	h.mux.Lock()
	defer h.mux.Unlock()

	id := generateHookId()

	// minimize allocations by shifting the slice
	h.handlers = append(h.handlers, nil)
	copy(h.handlers[1:], h.handlers)
	h.handlers[0] = &handlerPair[T]{id, fn}

	return id
}

// Add registers a new handler to the hook by appending it to the existing queue.
//
// Returns an autogenerated hook id that could be used later to remove the hook with Hook.Remove(id).
func (h *Hook[T]) Add(fn Handler[T]) string {
	h.mux.Lock()
	defer h.mux.Unlock()

	id := generateHookId()

	h.handlers = append(h.handlers, &handlerPair[T]{id, fn})

	return id
}

// Remove removes a single hook handler by its id.
func (h *Hook[T]) Remove(id string) {
	h.mux.Lock()
	defer h.mux.Unlock()

	for i := len(h.handlers) - 1; i >= 0; i-- {
		if h.handlers[i].id == id {
			h.handlers = append(h.handlers[:i], h.handlers[i+1:]...)
			return
		}
	}
}

// RemoveAll removes all registered handlers.
func (h *Hook[T]) RemoveAll() {
	h.mux.Lock()
	defer h.mux.Unlock()

	h.handlers = nil
}

// Trigger executes all registered hook handlers one by one
// with the specified `data` as an argument.
//
// Optionally, this method allows also to register additional one off
// handlers that will be temporary appended to the handlers queue.
//
// The execution stops when:
// - hook.StopPropagation is returned in one of the handlers
// - any non-nil error is returned in one of the handlers
func (h *Hook[T]) Trigger(data T, oneOffHandlers ...Handler[T]) error {
	h.mux.RLock()

	handlers := make([]*handlerPair[T], 0, len(h.handlers)+len(oneOffHandlers))
	handlers = append(handlers, h.handlers...)

	// append the one off handlers
	for i, oneOff := range oneOffHandlers {
		handlers = append(handlers, &handlerPair[T]{
			id:      fmt.Sprintf("@%d", i),
			handler: oneOff,
		})
	}

	// unlock is not deferred to avoid deadlocks in case Trigger
	// is called recursively by the handlers
	h.mux.RUnlock()

	for _, item := range handlers {
		err := item.handler(data)
		if err == nil {
			continue
		}

		if errors.Is(err, StopPropagation) {
			return nil
		}

		return err
	}

	return nil
}

func generateHookId() string {
	return security.PseudorandomString(8)
}
