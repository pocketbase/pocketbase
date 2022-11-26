package migrate

import (
	"testing"
)

func TestMigrationsList(t *testing.T) {
	l := MigrationsList{}

	l.Register(nil, nil, "3_test.go")
	l.Register(nil, nil, "1_test.go")
	l.Register(nil, nil, "2_test.go")
	l.Register(nil, nil /* auto detect file name */)

	expected := []string{
		"1_test.go",
		"2_test.go",
		"3_test.go",
		"list_test.go",
	}

	items := l.Items()
	if len(items) != len(expected) {
		t.Fatalf("Expected %d items, got %d: \n%#v", len(expected), len(items), items)
	}

	for i, name := range expected {
		item := l.Item(i)
		if item.File != name {
			t.Fatalf("Expected name %s for index %d, got %s", name, i, item.File)
		}
	}
}
