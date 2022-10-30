package models_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/models"
)

func TestBaseModelHasId(t *testing.T) {
	scenarios := []struct {
		model    models.BaseModel
		expected bool
	}{
		{
			models.BaseModel{},
			false,
		},
		{
			models.BaseModel{Id: ""},
			false,
		},
		{
			models.BaseModel{Id: "abc"},
			true,
		},
	}

	for i, s := range scenarios {
		result := s.model.HasId()
		if result != s.expected {
			t.Errorf("(%d) Expected %v, got %v", i, s.expected, result)
		}
	}
}

func TestBaseModelId(t *testing.T) {
	m := models.BaseModel{}

	if m.GetId() != "" {
		t.Fatalf("Expected empty id value, got %v", m.GetId())
	}

	m.SetId("test")

	if m.GetId() != "test" {
		t.Fatalf("Expected %q id, got %v", "test", m.GetId())
	}

	m.RefreshId()

	if len(m.GetId()) != 15 {
		t.Fatalf("Expected 15 chars id, got %v", m.GetId())
	}
}

func TestBaseModelIsNew(t *testing.T) {
	m0 := models.BaseModel{}
	m1 := models.BaseModel{Id: ""}
	m2 := models.BaseModel{Id: "test"}
	m3 := models.BaseModel{}
	m3.MarkAsNew()
	m4 := models.BaseModel{Id: "test"}
	m4.MarkAsNew()
	m5 := models.BaseModel{Id: "test"}
	m5.MarkAsNew()
	m5.UnmarkAsNew()
	// check if MarkAsNew will be called on initial RefreshId()
	m6 := models.BaseModel{}
	m6.RefreshId()

	scenarios := []struct {
		model    models.BaseModel
		expected bool
	}{
		{m0, true},
		{m1, true},
		{m2, false},
		{m3, true},
		{m4, true},
		{m5, false},
		{m6, true},
	}

	for i, s := range scenarios {
		result := s.model.IsNew()
		if result != s.expected {
			t.Errorf("(%d) Expected IsNew %v, got %v", i, s.expected, result)
		}
	}
}

func TestBaseModelCreated(t *testing.T) {
	m := models.BaseModel{}

	if !m.GetCreated().IsZero() {
		t.Fatalf("Expected zero datetime, got %v", m.GetCreated())
	}

	m.RefreshCreated()

	if m.GetCreated().IsZero() {
		t.Fatalf("Expected non-zero datetime, got %v", m.GetCreated())
	}
}

func TestBaseModelUpdated(t *testing.T) {
	m := models.BaseModel{}

	if !m.GetUpdated().IsZero() {
		t.Fatalf("Expected zero datetime, got %v", m.GetUpdated())
	}

	m.RefreshUpdated()

	if m.GetUpdated().IsZero() {
		t.Fatalf("Expected non-zero datetime, got %v", m.GetUpdated())
	}
}
