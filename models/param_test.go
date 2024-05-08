package models_test

import (
	"testing"

	"github.com/lilysnc/pocketbasepg/models"
)

func TestParamTableName(t *testing.T) {
	t.Parallel()

	m := models.Param{}
	if m.TableName() != "_params" {
		t.Fatalf("Unexpected table name, got %q", m.TableName())
	}
}
