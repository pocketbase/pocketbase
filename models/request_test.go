package models_test

import (
	"testing"

	"github.com/lilysnc/pocketbasepg/models"
)

func TestRequestTableName(t *testing.T) {
	m := models.Request{}
	if m.TableName() != "_requests" {
		t.Fatalf("Unexpected table name, got %q", m.TableName())
	}
}
