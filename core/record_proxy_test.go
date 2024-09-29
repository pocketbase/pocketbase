package core_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/core"
)

func TestBaseRecordProxy(t *testing.T) {
	p := core.BaseRecordProxy{}

	record := core.NewRecord(core.NewBaseCollection("test"))
	record.Id = "test"

	p.SetProxyRecord(record)

	if p.ProxyRecord() == nil || p.ProxyRecord().Id != p.Id || p.Id != "test" {
		t.Fatalf("Expected proxy record to be set")
	}
}
