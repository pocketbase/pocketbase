package models_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
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

// -------------------------------------------------------------------
// BaseAccount tests
// -------------------------------------------------------------------

func TestBaseAccountValidatePassword(t *testing.T) {
	scenarios := []struct {
		account  models.BaseAccount
		password string
		expected bool
	}{
		{
			// empty passwordHash + empty pass
			models.BaseAccount{},
			"",
			false,
		},
		{
			// empty passwordHash + nonempty pass
			models.BaseAccount{},
			"123456",
			false,
		},
		{
			// nonempty passwordHash + empty pass
			models.BaseAccount{PasswordHash: "$2a$10$SKk/Y/Yc925PBtsSYBvq3Ous9Jy18m4KTn6b/PQQ.Y9QVjy3o/Fv."},
			"",
			false,
		},
		{
			// nonempty passwordHash + wrong pass
			models.BaseAccount{PasswordHash: "$2a$10$SKk/Y/Yc925PBtsSYBvq3Ous9Jy18m4KTn6b/PQQ.Y9QVjy3o/Fv."},
			"654321",
			false,
		},
		{
			// nonempty passwordHash + correct pass
			models.BaseAccount{PasswordHash: "$2a$10$SKk/Y/Yc925PBtsSYBvq3Ous9Jy18m4KTn6b/PQQ.Y9QVjy3o/Fv."},
			"123456",
			true,
		},
	}

	for i, s := range scenarios {
		result := s.account.ValidatePassword(s.password)
		if result != s.expected {
			t.Errorf("(%d) Expected %v, got %v", i, s.expected, result)
		}
	}
}

func TestBaseAccountSetPassword(t *testing.T) {
	m := models.BaseAccount{
		// 123456
		PasswordHash:    "$2a$10$SKk/Y/Yc925PBtsSYBvq3Ous9Jy18m4KTn6b/PQQ.Y9QVjy3o/Fv.",
		LastResetSentAt: types.NowDateTime(),
		TokenKey:        "test",
	}

	// empty pass
	err1 := m.SetPassword("")
	if err1 == nil {
		t.Fatal("Expected empty password error")
	}

	err2 := m.SetPassword("654321")
	if err2 != nil {
		t.Fatalf("Expected nil, got error %v", err2)
	}

	if !m.ValidatePassword("654321") {
		t.Fatalf("Password is invalid")
	}

	if m.TokenKey == "test" {
		t.Fatalf("Expected TokenKey to change, got %v", m.TokenKey)
	}

	if !m.LastResetSentAt.IsZero() {
		t.Fatalf("Expected LastResetSentAt to be zero datetime, got %v", m.LastResetSentAt)
	}
}

func TestBaseAccountRefreshTokenKey(t *testing.T) {
	m := models.BaseAccount{TokenKey: "test"}

	m.RefreshTokenKey()

	// empty pass
	if m.TokenKey == "" || m.TokenKey == "test" {
		t.Fatalf("Expected TokenKey to change, got %q", m.TokenKey)
	}
}
