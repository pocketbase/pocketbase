package core_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/core"
)

func TestBaseModel(t *testing.T) {
	id := "test_id"

	m := core.BaseModel{Id: id}

	if m.PK() != id {
		t.Fatalf("[before PostScan] Expected PK %q, got %q", "", m.PK())
	}

	if m.LastSavedPK() != "" {
		t.Fatalf("[before PostScan] Expected LastSavedPK %q, got %q", "", m.LastSavedPK())
	}

	if !m.IsNew() {
		t.Fatalf("[before PostScan] Expected IsNew %v, got %v", true, m.IsNew())
	}

	if err := m.PostScan(); err != nil {
		t.Fatal(err)
	}

	if m.PK() != id {
		t.Fatalf("[after PostScan] Expected PK %q, got %q", "", m.PK())
	}

	if m.LastSavedPK() != id {
		t.Fatalf("[after PostScan] Expected LastSavedPK %q, got %q", id, m.LastSavedPK())
	}

	if m.IsNew() {
		t.Fatalf("[after PostScan] Expected IsNew %v, got %v", false, m.IsNew())
	}

	m.MarkAsNew()

	if m.PK() != id {
		t.Fatalf("[after MarkAsNew] Expected PK %q, got %q", id, m.PK())
	}

	if m.LastSavedPK() != "" {
		t.Fatalf("[after MarkAsNew] Expected LastSavedPK %q, got %q", "", m.LastSavedPK())
	}

	if !m.IsNew() {
		t.Fatalf("[after MarkAsNew] Expected IsNew %v, got %v", true, m.IsNew())
	}

	// mark as not new without id
	m.MarkAsNotNew()

	if m.PK() != id {
		t.Fatalf("[after MarkAsNotNew] Expected PK %q, got %q", id, m.PK())
	}

	if m.LastSavedPK() != id {
		t.Fatalf("[after MarkAsNotNew] Expected LastSavedPK %q, got %q", id, m.LastSavedPK())
	}

	if m.IsNew() {
		t.Fatalf("[after MarkAsNotNew] Expected IsNew %v, got %v", false, m.IsNew())
	}
}
