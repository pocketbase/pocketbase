package models_test

import (
	"encoding/json"
	"testing"

	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

func TestUserTableName(t *testing.T) {
	m := models.User{}
	if m.TableName() != "_users" {
		t.Fatalf("Unexpected table name, got %q", m.TableName())
	}
}

func TestUserAsMap(t *testing.T) {
	date, _ := types.ParseDateTime("2022-01-01 01:12:23.456")

	m := models.User{}
	m.Id = "210a896c-1e32-4c94-ae06-90c25fcf6791"
	m.Email = "test@example.com"
	m.PasswordHash = "test"
	m.LastResetSentAt = date
	m.Updated = date
	m.RefreshTokenKey()

	result, err := m.AsMap()
	if err != nil {
		t.Fatal(err)
	}

	encoded, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}

	expected := `{"created":"","email":"test@example.com","id":"210a896c-1e32-4c94-ae06-90c25fcf6791","lastResetSentAt":"2022-01-01 01:12:23.456","lastVerificationSentAt":"","profile":null,"updated":"2022-01-01 01:12:23.456","verified":false}`
	if string(encoded) != expected {
		t.Errorf("Expected %s, got %s", expected, string(encoded))
	}
}
