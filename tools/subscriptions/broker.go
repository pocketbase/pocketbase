package subscriptions

import (
	"fmt"
	"sync"
)

// Broker defines a struct for managing subscriptions clients.
type Broker struct {
	mux     sync.RWMutex
	clients map[string]Client
}

// NewBroker initializes and returns a new Broker instance.
func NewBroker() *Broker {
	return &Broker{
		clients: make(map[string]Client),
	}
}

// Clients returns all registered clients.
func (b *Broker) Clients() map[string]Client {
	b.mux.RLock()
	defer b.mux.RUnlock()

	return b.clients
}

// ClientById finds a registered client by its id.
//
// Returns non-nil error when client with clientId is not registered.
func (b *Broker) ClientById(clientId string) (Client, error) {
	b.mux.RLock()
	defer b.mux.RUnlock()

	client, ok := b.clients[clientId]
	if !ok {
		return nil, fmt.Errorf("No client associated with connection ID %q", clientId)
	}

	return client, nil
}

// Register adds a new client to the broker instance.
func (b *Broker) Register(client Client) {
	b.mux.Lock()
	defer b.mux.Unlock()

	b.clients[client.Id()] = client
}

// Unregister removes a single client by its id.
//
// If client with clientId doesn't exist, this method does nothing.
func (b *Broker) Unregister(clientId string) {
	b.mux.Lock()
	defer b.mux.Unlock()

	// Note:
	// There is no need to explicitly close the client's channel since it will be GC-ed anyway.
	// Addinitionally, closing the channel explicitly could panic when there are several
	// subscriptions attached to the client that needs to receive the same event.
	delete(b.clients, clientId)
}
