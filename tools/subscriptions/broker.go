package subscriptions

import (
	"fmt"

	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/store"
)

// Broker defines a struct for managing subscriptions clients.
type Broker struct {
	store *store.Store[string, Client]
}

// NewBroker initializes and returns a new Broker instance.
func NewBroker() *Broker {
	return &Broker{
		store: store.New[string, Client](nil),
	}
}

// Clients returns a shallow copy of all registered clients indexed
// with their connection id.
func (b *Broker) Clients() map[string]Client {
	return b.store.GetAll()
}

// ChunkedClients splits the current clients into a chunked slice.
func (b *Broker) ChunkedClients(chunkSize int) [][]Client {
	return list.ToChunks(b.store.Values(), chunkSize)
}

// TotalClients returns the total number of registered clients.
func (b *Broker) TotalClients() int {
	return b.store.Length()
}

// ClientById finds a registered client by its id.
//
// Returns non-nil error when client with clientId is not registered.
func (b *Broker) ClientById(clientId string) (Client, error) {
	client, ok := b.store.GetOk(clientId)
	if !ok {
		return nil, fmt.Errorf("no client associated with connection ID %q", clientId)
	}

	return client, nil
}

// Register adds a new client to the broker instance.
func (b *Broker) Register(client Client) {
	b.store.Set(client.Id(), client)
}

// Unregister removes a single client by its id and marks it as discarded.
//
// If client with clientId doesn't exist, this method does nothing.
func (b *Broker) Unregister(clientId string) {
	client := b.store.Get(clientId)
	if client == nil {
		return
	}
	client.Discard()
	b.store.Remove(clientId)
}
