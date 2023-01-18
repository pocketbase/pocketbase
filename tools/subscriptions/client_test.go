package subscriptions_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/tools/subscriptions"
)

func TestNewDefaultClient(t *testing.T) {
	c := subscriptions.NewDefaultClient()

	if c.Channel() == nil {
		t.Errorf("Expected channel to be initialized")
	}

	if c.Subscriptions() == nil {
		t.Errorf("Expected subscriptions map to be initialized")
	}

	if c.Id() == "" {
		t.Errorf("Expected unique id to be set")
	}
}

func TestId(t *testing.T) {
	clients := []*subscriptions.DefaultClient{
		subscriptions.NewDefaultClient(),
		subscriptions.NewDefaultClient(),
		subscriptions.NewDefaultClient(),
		subscriptions.NewDefaultClient(),
	}

	ids := map[string]struct{}{}
	for i, c := range clients {
		// check uniqueness
		if _, ok := ids[c.Id()]; ok {
			t.Errorf("(%d) Expected unique id, got %v", i, c.Id())
		} else {
			ids[c.Id()] = struct{}{}
		}

		// check length
		if len(c.Id()) != 40 {
			t.Errorf("(%d) Expected unique id to have 40 chars length, got %v", i, c.Id())
		}
	}
}

func TestChannel(t *testing.T) {
	c := subscriptions.NewDefaultClient()

	if c.Channel() == nil {
		t.Errorf("Expected channel to be initialized, got")
	}
}

func TestSubscriptions(t *testing.T) {
	c := subscriptions.NewDefaultClient()

	if len(c.Subscriptions()) != 0 {
		t.Errorf("Expected subscriptions to be empty")
	}

	c.Subscribe("sub1", "sub2", "sub3")

	if len(c.Subscriptions()) != 3 {
		t.Errorf("Expected 3 subscriptions, got %v", c.Subscriptions())
	}
}

func TestSubscribe(t *testing.T) {
	c := subscriptions.NewDefaultClient()

	subs := []string{"", "sub1", "sub2", "sub3"}
	expected := []string{"sub1", "sub2", "sub3"}

	c.Subscribe(subs...) // empty string should be skipped

	if len(c.Subscriptions()) != 3 {
		t.Errorf("Expected 3 subscriptions, got %v", c.Subscriptions())
	}

	for i, s := range expected {
		if !c.HasSubscription(s) {
			t.Errorf("(%d) Expected sub %s", i, s)
		}
	}
}

func TestUnsubscribe(t *testing.T) {
	c := subscriptions.NewDefaultClient()

	c.Subscribe("sub1", "sub2", "sub3")

	c.Unsubscribe("sub1")

	if c.HasSubscription("sub1") {
		t.Error("Expected sub1 to be removed")
	}

	c.Unsubscribe( /* all */ )
	if len(c.Subscriptions()) != 0 {
		t.Errorf("Expected all subscriptions to be removed, got %v", c.Subscriptions())
	}
}

func TestHasSubscription(t *testing.T) {
	c := subscriptions.NewDefaultClient()

	if c.HasSubscription("missing") {
		t.Error("Expected false, got true")
	}

	c.Subscribe("sub")

	if !c.HasSubscription("sub") {
		t.Error("Expected true, got false")
	}
}

func TestSetAndGet(t *testing.T) {
	c := subscriptions.NewDefaultClient()

	c.Set("demo", 1)

	result, _ := c.Get("demo").(int)

	if result != 1 {
		t.Errorf("Expected 1, got %v", result)
	}
}

func TestDiscard(t *testing.T) {
	c := subscriptions.NewDefaultClient()

	if v := c.IsDiscarded(); v {
		t.Fatal("Expected false, got true")
	}

	c.Discard()

	if v := c.IsDiscarded(); !v {
		t.Fatal("Expected true, got false")
	}
}
