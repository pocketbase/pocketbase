package subscriptions

import (
	"encoding/json"
	"net/url"
	"strings"
	"sync"

	"github.com/pocketbase/pocketbase/tools/inflector"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/spf13/cast"
)

const optionsParam = "options"

// Message defines a client's channel data.
type Message struct {
	Name string `json:"name"`
	Data []byte `json:"data"`
}

// SubscriptionOptions defines the request options (query params, headers, etc.)
// for a single subscription topic.
type SubscriptionOptions struct {
	// @todo after the requests handling refactoring consider
	// changing to map[string]string or map[string][]string

	Query   map[string]any `json:"query"`
	Headers map[string]any `json:"headers"`
}

// Client is an interface for a generic subscription client.
type Client interface {
	// Id Returns the unique id of the client.
	Id() string

	// Channel returns the client's communication channel.
	Channel() chan Message

	// Subscriptions returns a shallow copy of the client subscriptions matching the prefixes.
	// If no prefix is specified, returns all subscriptions.
	Subscriptions(prefixes ...string) map[string]SubscriptionOptions

	// Subscribe subscribes the client to the provided subscriptions list.
	//
	// Each subscription can also have "options" (json serialized SubscriptionOptions) as query parameter.
	//
	// Example:
	//
	// 	Subscribe(
	// 	    "subscriptionA",
	// 	    `subscriptionB?options={"query":{"a":1},"headers":{"x_token":"abc"}}`,
	// 	)
	Subscribe(subs ...string)

	// Unsubscribe unsubscribes the client from the provided subscriptions list.
	Unsubscribe(subs ...string)

	// HasSubscription checks if the client is subscribed to `sub`.
	HasSubscription(sub string) bool

	// Set stores any value to the client's context.
	Set(key string, value any)

	// Unset removes a single value from the client's context.
	Unset(key string)

	// Get retrieves the key value from the client's context.
	Get(key string) any

	// Discard marks the client as "discarded", meaning that it
	// shouldn't be used anymore for sending new messages.
	//
	// It is safe to call Discard() multiple times.
	Discard()

	// IsDiscarded indicates whether the client has been "discarded"
	// and should no longer be used.
	IsDiscarded() bool

	// Send sends the specified message to the client's channel (if not discarded).
	Send(m Message)
}

// ensures that DefaultClient satisfies the Client interface
var _ Client = (*DefaultClient)(nil)

// DefaultClient defines a generic subscription client.
type DefaultClient struct {
	store         map[string]any
	subscriptions map[string]SubscriptionOptions
	channel       chan Message
	id            string
	mux           sync.RWMutex
	isDiscarded   bool
}

// NewDefaultClient creates and returns a new DefaultClient instance.
func NewDefaultClient() *DefaultClient {
	return &DefaultClient{
		id:            security.RandomString(40),
		store:         map[string]any{},
		channel:       make(chan Message),
		subscriptions: map[string]SubscriptionOptions{},
	}
}

// Id implements the [Client.Id] interface method.
func (c *DefaultClient) Id() string {
	c.mux.RLock()
	defer c.mux.RUnlock()

	return c.id
}

// Channel implements the [Client.Channel] interface method.
func (c *DefaultClient) Channel() chan Message {
	c.mux.RLock()
	defer c.mux.RUnlock()

	return c.channel
}

// Subscriptions implements the [Client.Subscriptions] interface method.
//
// It returns a shallow copy of the client subscriptions matching the prefixes.
// If no prefix is specified, returns all subscriptions.
func (c *DefaultClient) Subscriptions(prefixes ...string) map[string]SubscriptionOptions {
	c.mux.RLock()
	defer c.mux.RUnlock()

	// no prefix -> return copy of all subscriptions
	if len(prefixes) == 0 {
		result := make(map[string]SubscriptionOptions, len(c.subscriptions))

		for s, options := range c.subscriptions {
			result[s] = options
		}

		return result
	}

	result := make(map[string]SubscriptionOptions)

	for _, prefix := range prefixes {
		for s, options := range c.subscriptions {
			// "?" ensures that the options query start character is always there
			// so that it can be used as an end separator when looking only for the main subscription topic
			if strings.HasPrefix(s+"?", prefix) {
				result[s] = options
			}
		}
	}

	return result
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

		// extract subscription options (if any)
		options := SubscriptionOptions{}
		u, err := url.Parse(s)
		if err == nil {
			rawOptions := u.Query().Get(optionsParam)
			if rawOptions != "" {
				json.Unmarshal([]byte(rawOptions), &options)
			}
		}

		// normalize query
		// (currently only single string values are supported for consistency with the default routes handling)
		for k, v := range options.Query {
			options.Query[k] = cast.ToString(v)
		}

		// normalize headers name and values, eg. "X-Token" is converted to "x_token"
		// (currently only single string values are supported for consistency with the default routes handling)
		for k, v := range options.Headers {
			delete(options.Headers, k)
			options.Headers[inflector.Snakecase(k)] = cast.ToString(v)
		}

		c.subscriptions[s] = options
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

// Unset implements the [Client.Unset] interface method.
func (c *DefaultClient) Unset(key string) {
	c.mux.Lock()
	defer c.mux.Unlock()

	delete(c.store, key)
}

// Discard implements the [Client.Discard] interface method.
func (c *DefaultClient) Discard() {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.isDiscarded = true
}

// IsDiscarded implements the [Client.IsDiscarded] interface method.
func (c *DefaultClient) IsDiscarded() bool {
	c.mux.RLock()
	defer c.mux.RUnlock()

	return c.isDiscarded
}

// Send sends the specified message to the client's channel (if not discarded).
func (c *DefaultClient) Send(m Message) {
	if c.IsDiscarded() {
		return
	}

	c.Channel() <- m
}
