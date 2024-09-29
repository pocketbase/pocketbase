package router_test

import (
	"io"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/tools/router"
)

func TestRereadableReadCloser(t *testing.T) {
	content := "test"

	rereadable := &router.RereadableReadCloser{
		ReadCloser: io.NopCloser(strings.NewReader(content)),
	}

	// read multiple times
	for i := 0; i < 3; i++ {
		result, err := io.ReadAll(rereadable)
		if err != nil {
			t.Fatalf("[read:%d] %v", i, err)
		}
		if str := string(result); str != content {
			t.Fatalf("[read:%d] Expected %q, got %q", i, content, result)
		}
	}
}
