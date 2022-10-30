package subscriptions

import (
	"sync"

	"github.com/pocketbase/pocketbase/tools/security"
)

// Message defines a client's channel data.
type Message struct {
	Name string
	Data string
}

// Client is an interface for a generic subscription client.
type Client interface {
	// Id Returns the unique id of the client.
	Id() string

	// Channel returns the client's communication channel.
	Channel() chan Message

	// Subscriptions returns all subscriptions to which the client has subscribed to.
	Subscriptions() map[string]struct{}

	// Subscribe subscribes the client to the provided subscriptions list.
	Subscribe(subs ...string)

	// Unsubscribe unsubscribes the client from the provided subscriptions list.
	Unsubscribe(subs ...string)

	// HasSubscription checks if the client is subscribed to `sub`.
	HasSubscription(sub string) bool

	// Set stores any value to the client's context.
	Set(key string, value any)

	// Get retrieves the key value from the client's context.
	Get(key string) any
}

// ensures that DefaultClient satisfies the Client interface
var _ Client = (*DefaultClient)(nil)

// DefaultClient defines a generic subscription client.
type DefaultClient struct {
	mux           sync.RWMutex
	id            string
	store         map[string]any
	channel       chan Message
	subscriptions map[string]struct{}
}

// NewDefaultClient creates and returns a new DefaultClient instance.
func NewDefaultClient() *DefaultClient {
	return &DefaultClient{
		id:            security.RandomString(40),
		store:         map[string]any{},
		channel:       make(chan Message),
		subscriptions: make(map[string]struct{}),
	}
}

// Id implements the [Client.Id] interface method.
func (c *DefaultClient) Id() string {
	return c.id
}

// Channel implements the [Client.Channel] interface method.
func (c *DefaultClient) Channel() chan Message {
	return c.channel
}

// Subscriptions implements the [Client.Subscriptions] interface method.
func (c *DefaultClient) Subscriptions() map[string]struct{} {
	c.mux.RLock()
	defer c.mux.RUnlock()

	return c.subscriptions
}

// Subscribe implements the [Client.Subscribe] interface method.
//
// Empty subscriptions (aka. "") are ignored.
func (c *DefaultClient) Subscribe(subs ...string) {
	c.mux.Lock()
	defer c.mux.Unlock()

	for _, s := range subs {
		if s == "" {
			continue // skip empty
		}

		c.subscriptions[s] = struct{}{}
	}
}

// Unsubscribe implements the [Client.Unsubscribe] interface method.
//
// If subs is not set, this method removes all registered client's subscriptions.
func (c *DefaultClient) Unsubscribe(subs ...string) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if len(subs) > 0 {
		for _, s := range subs {
			delete(c.subscriptions, s)
		}
	} else {
		// unsubscribe all
		for s := range c.subscriptions {
			delete(c.subscriptions, s)
		}
	}
}

// HasSubscription implements the [Client.HasSubscription] interface method.
func (c *DefaultClient) HasSubscription(sub string) bool {
	c.mux.RLock()
	defer c.mux.RUnlock()

	_, ok := c.subscriptions[sub]

	return ok
}

// Get implements the [Client.Get] interface method.
func (c *DefaultClient) Get(key string) any {
	c.mux.RLock()
	defer c.mux.RUnlock()

	return c.store[key]
}

// Set implements the [Client.Set] interface method.
func (c *DefaultClient) Set(key string, value any) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.store[key] = value
}
