package models_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/models"
)

func TestAdminTableName(t *testing.T) {
	m := models.Admin{}
	if m.TableName() != "_admins" {
		t.Fatalf("Unexpected table name, got %q", m.TableName())
	}
}
