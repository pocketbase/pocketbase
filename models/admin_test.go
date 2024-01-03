package models_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

func TestAdminTableName(t *testing.T) {
	t.Parallel()

	m := models.Admin{}
	if m.TableName() != "_admins" {
		t.Fatalf("Unexpected table name, got %q", m.TableName())
	}
}

func TestAdminValidatePassword(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		admin    models.Admin
		password string
		expected bool
	}{
		{
			// empty passwordHash + empty pass
			models.Admin{},
			"",
			false,
		},
		{
			// empty passwordHash + nonempty pass
			models.Admin{},
			"123456",
			false,
		},
		{
			// nonempty passwordHash + empty pass
			models.Admin{PasswordHash: "$2a$10$SKk/Y/Yc925PBtsSYBvq3Ous9Jy18m4KTn6b/PQQ.Y9QVjy3o/Fv."},
			"",
			false,
		},
		{
			// nonempty passwordHash + wrong pass
			models.Admin{PasswordHash: "$2a$10$SKk/Y/Yc925PBtsSYBvq3Ous9Jy18m4KTn6b/PQQ.Y9QVjy3o/Fv."},
			"654321",
			false,
		},
		{
			// nonempty passwordHash + correct pass
			models.Admin{PasswordHash: "$2a$10$SKk/Y/Yc925PBtsSYBvq3Ous9Jy18m4KTn6b/PQQ.Y9QVjy3o/Fv."},
			"123456",
			true,
		},
	}

	for i, s := range scenarios {
		result := s.admin.ValidatePassword(s.password)
		if result != s.expected {
			t.Errorf("(%d) Expected %v, got %v", i, s.expected, result)
		}
	}
}

func TestAdminSetPassword(t *testing.T) {
	t.Parallel()

	m := models.Admin{
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

func TestAdminRefreshTokenKey(t *testing.T) {
	t.Parallel()

	m := models.Admin{TokenKey: "test"}

	m.RefreshTokenKey()

	// empty pass
	if m.TokenKey == "" || m.TokenKey == "test" {
		t.Fatalf("Expected TokenKey to change, got %q", m.TokenKey)
	}
}
