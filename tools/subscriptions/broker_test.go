package subscriptions_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/tools/subscriptions"
)

func TestNewBroker(t *testing.T) {
	b := subscriptions.NewBroker()

	if b.Clients() == nil {
		t.Fatal("Expected clients map to be initialized")
	}
}

func TestClients(t *testing.T) {
	b := subscriptions.NewBroker()

	if total := len(b.Clients()); total != 0 {
		t.Fatalf("Expected no clients, got %v", total)
	}

	b.Register(subscriptions.NewDefaultClient())
	b.Register(subscriptions.NewDefaultClient())

	// check if it is a shallow copy
	clients := b.Clients()
	for k := range clients {
		delete(clients, k)
	}

	// should return a new copy
	if total := len(b.Clients()); total != 2 {
		t.Fatalf("Expected 2 clients, got %v", total)
	}
}

func TestChunkedClients(t *testing.T) {
	b := subscriptions.NewBroker()

	chunks := b.ChunkedClients(2)
	if total := len(chunks); total != 0 {
		t.Fatalf("Expected %d chunks, got %d", 0, total)
	}

	b.Register(subscriptions.NewDefaultClient())
	b.Register(subscriptions.NewDefaultClient())
	b.Register(subscriptions.NewDefaultClient())

	chunks = b.ChunkedClients(2)
	if total := len(chunks); total != 2 {
		t.Fatalf("Expected %d chunks, got %d", 2, total)
	}

	if total := len(chunks[0]); total != 2 {
		t.Fatalf("Expected the first chunk to have 2 clients, got %d", total)
	}

	if total := len(chunks[1]); total != 1 {
		t.Fatalf("Expected the second chunk to have 1 client, got %d", total)
	}
}

func TestTotalClients(t *testing.T) {
	b := subscriptions.NewBroker()

	if total := b.TotalClients(); total != 0 {
		t.Fatalf("Expected no clients, got %d", total)
	}

	b.Register(subscriptions.NewDefaultClient())
	b.Register(subscriptions.NewDefaultClient())

	if total := b.TotalClients(); total != 2 {
		t.Fatalf("Expected %d clients, got %d", 2, total)
	}
}

func TestClientById(t *testing.T) {
	b := subscriptions.NewBroker()

	clientA := subscriptions.NewDefaultClient()
	clientB := subscriptions.NewDefaultClient()
	b.Register(clientA)
	b.Register(clientB)

	resultClient, err := b.ClientById(clientA.Id())
	if err != nil {
		t.Fatalf("Expected client with id %s, got error %v", clientA.Id(), err)
	}
	if resultClient.Id() != clientA.Id() {
		t.Fatalf("Expected client %s, got %s", clientA.Id(), resultClient.Id())
	}

	if c, err := b.ClientById("missing"); err == nil {
		t.Fatalf("Expected error, found client %v", c)
	}
}

func TestRegister(t *testing.T) {
	b := subscriptions.NewBroker()

	client := subscriptions.NewDefaultClient()
	b.Register(client)

	if _, err := b.ClientById(client.Id()); err != nil {
		t.Fatalf("Expected client with id %s, got error %v", client.Id(), err)
	}
}

func TestUnregister(t *testing.T) {
	b := subscriptions.NewBroker()

	clientA := subscriptions.NewDefaultClient()
	clientB := subscriptions.NewDefaultClient()
	b.Register(clientA)
	b.Register(clientB)

	if _, err := b.ClientById(clientA.Id()); err != nil {
		t.Fatalf("Expected client with id %s, got error %v", clientA.Id(), err)
	}

	b.Unregister(clientA.Id())

	if c, err := b.ClientById(clientA.Id()); err == nil {
		t.Fatalf("Expected error, found client %v", c)
	}

	// clientB shouldn't have been removed
	if _, err := b.ClientById(clientB.Id()); err != nil {
		t.Fatalf("Expected client with id %s, got error %v", clientB.Id(), err)
	}
}
