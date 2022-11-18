package types_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/tools/types"
)

func TestPointer(t *testing.T) {
	s1 := types.Pointer("")
	if s1 == nil || *s1 != "" {
		t.Fatalf("Expected empty string pointer, got %#v", s1)
	}

	s2 := types.Pointer("test")
	if s2 == nil || *s2 != "test" {
		t.Fatalf("Expected 'test' string pointer, got %#v", s2)
	}

	s3 := types.Pointer(123)
	if s3 == nil || *s3 != 123 {
		t.Fatalf("Expected 123 string pointer, got %#v", s3)
	}
}
