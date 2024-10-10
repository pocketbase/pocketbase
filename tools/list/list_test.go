package list_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/types"
)

func TestSubtractSliceString(t *testing.T) {
	scenarios := []struct {
		base     []string
		subtract []string
		expected string
	}{
		{
			[]string{},
			[]string{},
			`[]`,
		},
		{
			[]string{},
			[]string{"1", "2", "3", "4"},
			`[]`,
		},
		{
			[]string{"1", "2", "3", "4"},
			[]string{},
			`["1","2","3","4"]`,
		},
		{
			[]string{"1", "2", "3", "4"},
			[]string{"1", "2", "3", "4"},
			`[]`,
		},
		{
			[]string{"1", "2", "3", "4", "7"},
			[]string{"2", "4", "5", "6"},
			`["1","3","7"]`,
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%s", i, s.expected), func(t *testing.T) {
			result := list.SubtractSlice(s.base, s.subtract)

			raw, err := json.Marshal(result)
			if err != nil {
				t.Fatalf("Failed to serialize: %v", err)
			}

			strResult := string(raw)

			if strResult != s.expected {
				t.Fatalf("Expected %v, got %v", s.expected, strResult)
			}
		})
	}
}

func TestSubtractSliceInt(t *testing.T) {
	scenarios := []struct {
		base     []int
		subtract []int
		expected string
	}{
		{
			[]int{},
			[]int{},
			`[]`,
		},
		{
			[]int{},
			[]int{1, 2, 3, 4},
			`[]`,
		},
		{
			[]int{1, 2, 3, 4},
			[]int{},
			`[1,2,3,4]`,
		},
		{
			[]int{1, 2, 3, 4},
			[]int{1, 2, 3, 4},
			`[]`,
		},
		{
			[]int{1, 2, 3, 4, 7},
			[]int{2, 4, 5, 6},
			`[1,3,7]`,
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%s", i, s.expected), func(t *testing.T) {
			result := list.SubtractSlice(s.base, s.subtract)

			raw, err := json.Marshal(result)
			if err != nil {
				t.Fatalf("Failed to serialize: %v", err)
			}

			strResult := string(raw)

			if strResult != s.expected {
				t.Fatalf("Expected %v, got %v", s.expected, strResult)
			}
		})
	}
}

func TestExistInSliceString(t *testing.T) {
	scenarios := []struct {
		item     string
		list     []string
		expected bool
	}{
		{"", []string{""}, true},
		{"", []string{"1", "2", "test 123"}, false},
		{"test", []string{}, false},
		{"test", []string{"TEST"}, false},
		{"test", []string{"1", "2", "test 123"}, false},
		{"test", []string{"1", "2", "test"}, true},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%s", i, s.item), func(t *testing.T) {
			result := list.ExistInSlice(s.item, s.list)
			if result != s.expected {
				t.Fatalf("Expected %v, got %v", s.expected, result)
			}
		})
	}
}

func TestExistInSliceInt(t *testing.T) {
	scenarios := []struct {
		item     int
		list     []int
		expected bool
	}{
		{0, []int{}, false},
		{0, []int{0}, true},
		{4, []int{1, 2, 3}, false},
		{1, []int{1, 2, 3}, true},
		{-1, []int{0, 1, 2, 3}, false},
		{-1, []int{0, -1, -2, -3, -4}, true},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%d", i, s.item), func(t *testing.T) {
			result := list.ExistInSlice(s.item, s.list)
			if result != s.expected {
				t.Fatalf("Expected %v, got %v", s.expected, result)
			}
		})
	}
}

func TestExistInSliceWithRegex(t *testing.T) {
	scenarios := []struct {
		item     string
		list     []string
		expected bool
	}{
		{"", []string{``}, true},
		{"", []string{`^\W+$`}, false},
		{" ", []string{`^\W+$`}, true},
		{"test", []string{`^\invalid[+$`}, false}, // invalid regex
		{"test", []string{`^\W+$`, "test"}, true},
		{`^\W+$`, []string{`^\W+$`, "test"}, false}, // direct match shouldn't work for this case
		{`\W+$`, []string{`\W+$`, "test"}, true},    // direct match should work for this case because it is not an actual supported pattern format
		{"!?@", []string{`\W+$`, "test"}, false},    // the method requires the pattern elems to start with '^'
		{"!?@", []string{`^\W+`, "test"}, false},    // the method requires the pattern elems to end with '$'
		{"!?@", []string{`^\W+$`, "test"}, true},
		{"!?@test", []string{`^\W+$`, "test"}, false},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%s", i, s.item), func(t *testing.T) {
			result := list.ExistInSliceWithRegex(s.item, s.list)
			if result != s.expected {
				t.Fatalf("Expected %v, got %v", s.expected, result)
			}
		})
	}
}

func TestToInterfaceSlice(t *testing.T) {
	scenarios := []struct {
		items []string
	}{
		{[]string{}},
		{[]string{""}},
		{[]string{"1", "test"}},
		{[]string{"test1", "test1", "test2", "test3"}},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%#v", i, s.items), func(t *testing.T) {
			result := list.ToInterfaceSlice(s.items)

			if len(result) != len(s.items) {
				t.Fatalf("Expected length %d, got %d", len(s.items), len(result))
			}

			for j, v := range result {
				if v != s.items[j] {
					t.Fatalf("Result list item doesn't match with the original list item, got %v VS %v", v, s.items[j])
				}
			}
		})
	}
}

func TestNonzeroUniquesString(t *testing.T) {
	scenarios := []struct {
		items    []string
		expected []string
	}{
		{[]string{}, []string{}},
		{[]string{""}, []string{}},
		{[]string{"1", "test"}, []string{"1", "test"}},
		{[]string{"test1", "", "test2", "Test2", "test1", "test3"}, []string{"test1", "test2", "Test2", "test3"}},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%#v", i, s.items), func(t *testing.T) {
			result := list.NonzeroUniques(s.items)

			if len(result) != len(s.expected) {
				t.Fatalf("Expected length %d, got %d", len(s.expected), len(result))
			}

			for j, v := range result {
				if v != s.expected[j] {
					t.Fatalf("Result list item doesn't match with the expected list item, got %v VS %v", v, s.expected[j])
				}
			}
		})
	}
}

func TestToUniqueStringSlice(t *testing.T) {
	scenarios := []struct {
		value    any
		expected []string
	}{
		{nil, []string{}},
		{"", []string{}},
		{[]any{}, []string{}},
		{[]int{}, []string{}},
		{"test", []string{"test"}},
		{[]int{1, 2, 3}, []string{"1", "2", "3"}},
		{[]any{0, 1, "test", ""}, []string{"0", "1", "test"}},
		{[]string{"test1", "test2", "test1"}, []string{"test1", "test2"}},
		{`["test1", "test2", "test2"]`, []string{"test1", "test2"}},
		{types.JSONArray[string]{"test1", "test2", "test1"}, []string{"test1", "test2"}},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%#v", i, s.value), func(t *testing.T) {
			result := list.ToUniqueStringSlice(s.value)

			if len(result) != len(s.expected) {
				t.Fatalf("Expected length %d, got %d", len(s.expected), len(result))
			}

			for j, v := range result {
				if v != s.expected[j] {
					t.Fatalf("Result list item doesn't match with the expected list item, got %v vs %v", v, s.expected[j])
				}
			}
		})
	}
}

func TestToChunks(t *testing.T) {
	scenarios := []struct {
		items     []any
		chunkSize int
		expected  string
	}{
		{nil, 2, "[]"},
		{[]any{}, 2, "[]"},
		{[]any{1, 2, 3, 4}, -1, "[[1],[2],[3],[4]]"},
		{[]any{1, 2, 3, 4}, 0, "[[1],[2],[3],[4]]"},
		{[]any{1, 2, 3, 4}, 2, "[[1,2],[3,4]]"},
		{[]any{1, 2, 3, 4, 5}, 2, "[[1,2],[3,4],[5]]"},
		{[]any{1, 2, 3, 4, 5}, 10, "[[1,2,3,4,5]]"},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%#v", i, s.items), func(t *testing.T) {
			result := list.ToChunks(s.items, s.chunkSize)

			raw, err := json.Marshal(result)
			if err != nil {
				t.Fatal(err)
			}
			rawStr := string(raw)

			if rawStr != s.expected {
				t.Fatalf("Expected %v, got %v", s.expected, rawStr)
			}
		})
	}
}
