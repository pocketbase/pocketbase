package subscriptions_test

import (
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/tools/subscriptions"
)

func TestMessageWrite(t *testing.T) {
	m := subscriptions.Message{
		Name: "test_name",
		Data: []byte("test_data"),
	}

	var sb strings.Builder

	m.WriteSSE(&sb, "test_id")

	expected := "id:test_id\nevent:test_name\ndata:test_data\n\n"

	if v := sb.String(); v != expected {
		t.Fatalf("Expected writer content\n%q\ngot\n%q", expected, v)
	}
}
