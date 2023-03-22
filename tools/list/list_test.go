package list_test

import (
	"encoding/json"
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
		result := list.SubtractSlice(s.base, s.subtract)

		raw, err := json.Marshal(result)
		if err != nil {
			t.Fatalf("(%d) Failed to serialize: %v", i, err)
		}

		strResult := string(raw)

		if strResult != s.expected {
			t.Fatalf("(%d) Expected %v, got %v", i, s.expected, strResult)
		}
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
		result := list.SubtractSlice(s.base, s.subtract)

		raw, err := json.Marshal(result)
		if err != nil {
			t.Fatalf("(%d) Failed to serialize: %v", i, err)
		}

		strResult := string(raw)

		if strResult != s.expected {
			t.Fatalf("(%d) Expected %v, got %v", i, s.expected, strResult)
		}
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

	for i, scenario := range scenarios {
		result := list.ExistInSlice(scenario.item, scenario.list)
		if result != scenario.expected {
			if scenario.expected {
				t.Errorf("(%d) Expected to exist in the list", i)
			} else {
				t.Errorf("(%d) Expected NOT to exist in the list", i)
			}
		}
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

	for i, scenario := range scenarios {
		result := list.ExistInSlice(scenario.item, scenario.list)
		if result != scenario.expected {
			if scenario.expected {
				t.Errorf("(%d) Expected to exist in the list", i)
			} else {
				t.Errorf("(%d) Expected NOT to exist in the list", i)
			}
		}
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

	for i, scenario := range scenarios {
		result := list.ExistInSliceWithRegex(scenario.item, scenario.list)
		if result != scenario.expected {
			if scenario.expected {
				t.Errorf("(%d) Expected the string to exist in the list", i)
			} else {
				t.Errorf("(%d) Expected the string NOT to exist in the list", i)
			}
		}
	}
}

func TestToInterfaceSlice(t *testing.T) {
	scenarios := []struct {
		items []string
	}{
		{[]string{}},
		{[]string{""}},
		{[]string{"1", "test"}},
		{[]string{"test1", "test2", "test3"}},
	}

	for i, scenario := range scenarios {
		result := list.ToInterfaceSlice(scenario.items)

		if len(result) != len(scenario.items) {
			t.Errorf("(%d) Result list length doesn't match with the original list", i)
		}

		for j, v := range result {
			if v != scenario.items[j] {
				t.Errorf("(%d:%d) Result list item should match with the original list item", i, j)
			}
		}
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

	for i, scenario := range scenarios {
		result := list.NonzeroUniques(scenario.items)

		if len(result) != len(scenario.expected) {
			t.Errorf("(%d) Result list length doesn't match with the expected list", i)
		}

		for j, v := range result {
			if v != scenario.expected[j] {
				t.Errorf("(%d:%d) Result list item should match with the expected list item", i, j)
			}
		}
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
		{types.JsonArray[string]{"test1", "test2", "test1"}, []string{"test1", "test2"}},
	}

	for i, scenario := range scenarios {
		result := list.ToUniqueStringSlice(scenario.value)

		if len(result) != len(scenario.expected) {
			t.Errorf("(%d) Result list length doesn't match with the expected list", i)
		}

		for j, v := range result {
			if v != scenario.expected[j] {
				t.Errorf("(%d:%d) Result list item should match with the expected list item", i, j)
			}
		}
	}
}
