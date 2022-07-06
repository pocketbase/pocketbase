package models_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/models"
)

func TestCollectionTableName(t *testing.T) {
	m := models.Collection{}
	if m.TableName() != "_collections" {
		t.Fatalf("Unexpected table name, got %q", m.TableName())
	}
}

func TestCollectionBaseFilesPath(t *testing.T) {
	m := models.Collection{}

	m.RefreshId()

	expected := m.Id
	if m.BaseFilesPath() != expected {
		t.Fatalf("Expected path %s, got %s", expected, m.BaseFilesPath())
	}
}
