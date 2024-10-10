package core_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/core"
)

func TestMigrationsList(t *testing.T) {
	l1 := core.MigrationsList{}
	l1.Add(&core.Migration{File: "5_test.go"})
	l1.Add(&core.Migration{ /* auto detect file name */ })
	l1.Register(nil, nil, "3_test.go")
	l1.Register(nil, nil, "1_test.go")
	l1.Register(nil, nil, "2_test.go")
	l1.Register(nil, nil /* auto detect file name */)

	l2 := core.MigrationsList{}
	l2.Register(nil, nil, "4_test.go")
	l2.Copy(l1)

	expected := []string{
		"1_test.go",
		"2_test.go",
		"3_test.go",
		"4_test.go",
		"5_test.go",
		// twice because there 2 test migrations with auto filename
		"migrations_list_test.go",
		"migrations_list_test.go",
	}

	items := l2.Items()
	if len(items) != len(expected) {
		names := make([]string, len(items))
		for i, item := range items {
			names[i] = item.File
		}
		t.Fatalf("Expected %d items, got %d:\n%v", len(expected), len(names), names)
	}

	for i, name := range expected {
		item := l2.Item(i)
		if item.File != name {
			t.Fatalf("Expected name %s for index %d, got %s", name, i, item.File)
		}
	}
}
