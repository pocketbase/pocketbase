package subscriptions

import (
	"fmt"
	"sync"
)

// Broker defines a struct for managing subscriptions clients.
type Broker struct {
	clients map[string]Client
	mux     sync.RWMutex
}

// NewBroker initializes and returns a new Broker instance.
func NewBroker() *Broker {
	return &Broker{
		clients: make(map[string]Client),
	}
}

// Clients returns a shallow copy of all registered clients indexed
// with their connection id.
func (b *Broker) Clients() map[string]Client {
	b.mux.RLock()
	defer b.mux.RUnlock()

	copy := make(map[string]Client, len(b.clients))

	for id, c := range b.clients {
		copy[id] = c
	}

	return copy
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

	if client, ok := b.clients[clientId]; ok {
		client.Discard()
		delete(b.clients, clientId)
	}
}
