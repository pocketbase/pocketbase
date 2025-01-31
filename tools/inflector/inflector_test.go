package inflector_test

import (
	"fmt"
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

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%#v", i, s.val), func(t *testing.T) {
			result := inflector.UcFirst(s.val)
			if result != s.expected {
				t.Fatalf("Expected %q, got %q", s.expected, result)
			}
		})
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

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%#v", i, s.val), func(t *testing.T) {
			result := inflector.Columnify(s.val)
			if result != s.expected {
				t.Fatalf("Expected %q, got %q", s.expected, result)
			}
		})
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

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%#v", i, s.val), func(t *testing.T) {
			result := inflector.Sentenize(s.val)
			if result != s.expected {
				t.Fatalf("Expected %q, got %q", s.expected, result)
			}
		})
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

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%#v", i, s.val), func(t *testing.T) {
			result, err := inflector.Sanitize(s.val, s.pattern)
			hasErr := err != nil

			if s.expectErr != hasErr {
				t.Fatalf("Expected hasErr %v, got %v (%v)", s.expectErr, hasErr, err)
			}

			if result != s.expected {
				t.Fatalf("Expected %q, got %q", s.expected, result)
			}
		})
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

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%#v", i, s.val), func(t *testing.T) {
			result := inflector.Snakecase(s.val)
			if result != s.expected {
				t.Fatalf("Expected %q, got %q", s.expected, result)
			}
		})
	}
}

func TestCamelize(t *testing.T) {
	scenarios := []struct {
		val      string
		expected string
	}{
		{"", ""},
		{" ", ""},
		{"Test", "Test"},
		{"test", "Test"},
		{"testTest2", "TestTest2"},
		{"TestTest2", "TestTest2"},
		{"test test2", "TestTest2"},
		{"test-test2", "TestTest2"},
		{"test'test2", "TestTest2"},
		{"test1test2", "Test1test2"},
		{"1test-test2", "1testTest2"},
		{"123", "123"},
		{"123a", "123a"},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%#v", i, s.val), func(t *testing.T) {
			result := inflector.Camelize(s.val)
			if result != s.expected {
				t.Fatalf("Expected %q, got %q", s.expected, result)
			}
		})
	}
}
