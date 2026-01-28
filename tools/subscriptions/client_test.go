package subscriptions_test

import (
	"encoding/json"
	"strings"
	"sync"
	"testing"
	"time"

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
		t.Fatalf("Expected channel to be initialized, got")
	}
}

func TestSubscriptions(t *testing.T) {
	c := subscriptions.NewDefaultClient()

	if len(c.Subscriptions()) != 0 {
		t.Fatalf("Expected subscriptions to be empty")
	}

	c.Subscribe("sub1", "sub11", "sub2")

	scenarios := []struct {
		prefixes []string
		expected []string
	}{
		{nil, []string{"sub1", "sub11", "sub2"}},
		{[]string{"missing"}, nil},
		{[]string{"sub1"}, []string{"sub1", "sub11"}},
		{[]string{"sub2"}, []string{"sub2"}}, // with extra query start char
	}

	for _, s := range scenarios {
		t.Run(strings.Join(s.prefixes, ","), func(t *testing.T) {
			subs := c.Subscriptions(s.prefixes...)

			if len(subs) != len(s.expected) {
				t.Fatalf("Expected %d subscriptions, got %d", len(s.expected), len(subs))
			}

			for _, s := range s.expected {
				if _, ok := subs[s]; !ok {
					t.Fatalf("Missing subscription %q in \n%v", s, subs)
				}
			}
		})
	}
}

func TestSubscribe(t *testing.T) {
	c := subscriptions.NewDefaultClient()

	subs := []string{"", "sub1", "sub2", "sub3"}
	expected := []string{"sub1", "sub2", "sub3"}

	c.Subscribe(subs...) // empty string should be skipped

	if len(c.Subscriptions()) != 3 {
		t.Fatalf("Expected 3 subscriptions, got %v", c.Subscriptions())
	}

	for i, s := range expected {
		if !c.HasSubscription(s) {
			t.Errorf("(%d) Expected sub %s", i, s)
		}
	}
}

func TestSubscribeOptions(t *testing.T) {
	c := subscriptions.NewDefaultClient()

	sub1 := "test1"
	sub2 := `test2?options={"query":{"name":123},"headers":{"X-Token":456}}`

	c.Subscribe(sub1, sub2)

	subs := c.Subscriptions()

	scenarios := []struct {
		name            string
		expectedOptions string
	}{
		{sub1, `{"query":{},"headers":{}}`},
		{sub2, `{"query":{"name":"123"},"headers":{"x_token":"456"}}`},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			options, ok := subs[s.name]
			if !ok {
				t.Fatalf("Missing subscription \n%q \nin \n%v", s.name, subs)
			}

			rawBytes, err := json.Marshal(options)
			if err != nil {
				t.Fatal(err)
			}
			rawStr := string(rawBytes)

			if rawStr != s.expectedOptions {
				t.Fatalf("Expected options\n%v\ngot\n%v", s.expectedOptions, rawStr)
			}
		})
	}
}

func TestUnsubscribe(t *testing.T) {
	c := subscriptions.NewDefaultClient()

	c.Subscribe("sub1", "sub2", "sub3")

	c.Unsubscribe("sub1")

	if c.HasSubscription("sub1") {
		t.Fatalf("Expected sub1 to be removed")
	}

	c.Unsubscribe( /* all */ )
	if len(c.Subscriptions()) != 0 {
		t.Fatalf("Expected all subscriptions to be removed, got %v", c.Subscriptions())
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

func TestSend(t *testing.T) {
	var mu sync.Mutex

	c := subscriptions.NewDefaultClient()

	received := []string{}
	go func() {
		for m := range c.Channel() {
			mu.Lock()
			received = append(received, m.Name)
			mu.Unlock()
		}
	}()

	c.Send(subscriptions.Message{Name: "m1"})
	c.Send(subscriptions.Message{Name: "m2"})
	c.Discard()
	c.Send(subscriptions.Message{Name: "m3"})
	c.Send(subscriptions.Message{Name: "m4"})
	time.Sleep(5 * time.Millisecond)

	expected := []string{"m1", "m2"}

	mu.Lock()
	defer mu.Unlock()
	if len(received) != len(expected) {
		t.Fatalf("Expected %d messages, got %d", len(expected), len(received))
	}
	for _, name := range expected {
		var exists bool
		for _, n := range received {
			if n == name {
				exists = true
				break
			}
		}
		if !exists {
			t.Fatalf("Missing expected %q message, got %v", name, received)
		}
	}
}
