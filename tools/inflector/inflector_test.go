package inflector_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/tools/inflector"
)

func TestUcFirst(t *testing.T) {
	scenarios := []struct {
		val      string
		expected string
	}{
		{"", ""},
		{" ", " "},
		{"Test", "Test"},
		{"test", "Test"},
		{"test test2", "Test test2"},
	}

	for i, scenario := range scenarios {
		if result := inflector.UcFirst(scenario.val); result != scenario.expected {
			t.Errorf("(%d) Expected %q, got %q", i, scenario.expected, result)
		}
	}
}

func TestColumnify(t *testing.T) {
	scenarios := []struct {
		val      string
		expected string
	}{
		{"", ""},
		{"   ", ""},
		{"123", "123"},
		{"Test.", "Test."},
		{" test ", "test"},
		{"test1.test2", "test1.test2"},
		{"@test!abc", "@testabc"},
		{"#test?abc", "#testabc"},
		{"123test(123)#", "123test123#"},
		{"test1--test2", "test1--test2"},
	}

	for i, scenario := range scenarios {
		if result := inflector.Columnify(scenario.val); result != scenario.expected {
			t.Errorf("(%d) Expected %q, got %q", i, scenario.expected, result)
		}
	}
}

func TestSentenize(t *testing.T) {
	scenarios := []struct {
		val      string
		expected string
	}{
		{"", ""},
		{"   ", ""},
		{".", "."},
		{"?", "?"},
		{"!", "!"},
		{"Test", "Test."},
		{" test ", "Test."},
		{"hello world", "Hello world."},
		{"hello world.", "Hello world."},
		{"hello world!", "Hello world!"},
		{"hello world?", "Hello world?"},
	}

	for i, scenario := range scenarios {
		if result := inflector.Sentenize(scenario.val); result != scenario.expected {
			t.Errorf("(%d) Expected %q, got %q", i, scenario.expected, result)
		}
	}
}

func TestSanitize(t *testing.T) {
	scenarios := []struct {
		val       string
		pattern   string
		expected  string
		expectErr bool
	}{
		{"", ``, "", false},
		{" ", ``, " ", false},
		{" ", ` `, "", false},
		{"", `[A-Z]`, "", false},
		{"abcABC", `[A-Z]`, "abc", false},
		{"abcABC", `[A-Z`, "", true}, // invalid pattern
	}

	for i, scenario := range scenarios {
		result, err := inflector.Sanitize(scenario.val, scenario.pattern)
		hasErr := err != nil

		if scenario.expectErr != hasErr {
			if scenario.expectErr {
				t.Errorf("(%d) Expected error, got nil", i)
			} else {
				t.Errorf("(%d) Didn't expect error, got", err)
			}
		}

		if result != scenario.expected {
			t.Errorf("(%d) Expected %q, got %q", i, scenario.expected, result)
		}
	}
}

func TestSnakecase(t *testing.T) {
	scenarios := []struct {
		val      string
		expected string
	}{
		{"", ""},
		{"  ", ""},
		{"!@#$%^", ""},
		{"...", ""},
		{"_", ""},
		{"John Doe", "john_doe"},
		{"John_Doe", "john_doe"},
		{".a!b@c#d$e%123. ", "a_b_c_d_e_123"},
		{"HelloWorld", "hello_world"},
		{"HelloWorld1HelloWorld2", "hello_world1_hello_world2"},
		{"TEST", "test"},
		{"testABR", "test_abr"},
	}

	for i, scenario := range scenarios {
		if result := inflector.Snakecase(scenario.val); result != scenario.expected {
			t.Errorf("(%d) Expected %q, got %q", i, scenario.expected, result)
		}
	}
}
