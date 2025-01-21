package cli_test

import (
	"bytes"
	"testing"

	"github.com/pocketbase/pocketbase/tools/cli"
)

func TestPromptMessageWithDefaultYes(t *testing.T) {
	var writer, reader bytes.Buffer
	cli := cli.NewWithIO(&writer, &reader)

	_ = cli.Confirm("Do you want to proceed with the update?", true)

	expected := "Do you want to proceed with the update? (Y/n) "
	result := writer.String()
	if result != expected {
		t.Fatalf("Expected prompt message \"%q\", got \"%q\"", expected, result)
	}
}

func TestPromptMessageWithDefaultNo(t *testing.T) {
	var writer, reader bytes.Buffer
	cli := cli.NewWithIO(&writer, &reader)

	_ = cli.Confirm("Do you want to proceed with the update?", false)

	expected := "Do you want to proceed with the update? (y/N) "
	result := writer.String()
	if result != expected {
		t.Fatalf("Expected prompt message \"%q\", got \"%q\"", expected, result)
	}
}

func TestSubmit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		def      bool
		expected bool
	}{
		{name: "TestSubmitBlankWithPromptDefaultYes", input: "", def: true, expected: true},
		{name: "TestSubmitBlankWithPromptDefaultNo", input: "", def: false, expected: false},
		{name: "TestSubmitY", input: "Y", def: false, expected: true},
		{name: "TestSubmitLowerY", input: "y", def: false, expected: true},
		{name: "TestSubmitYeS", input: "YeS", def: false, expected: true},
		{name: "TestSubmitN", input: "N", def: true, expected: false},
		{name: "TestSubmitLowerN", input: "n", def: true, expected: false},
		{name: "TestSubmitNO", input: "NO", def: true, expected: false},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			var writer, reader bytes.Buffer
			cli := cli.NewWithIO(&writer, &reader)

			reader.WriteString(tt.input)
			result := cli.Confirm("Do you want to proceed with the update?", tt.def)

			if result != tt.expected {
				t.Fatalf("Expected prompt result to be %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestSubmitSomethingElse(t *testing.T) {
	var writer, reader bytes.Buffer
	cli := cli.NewWithIO(&writer, &reader)

	reader.WriteString("YeSNoOk")
	_ = cli.Confirm("Do you want to proceed with the update?", false)
	prompt := writer.String()

	expected := "Do you want to proceed with the update? (y/N) Do you want to proceed with the update? (y/N) "
	if prompt != expected {
		t.Fatalf("Expected prompt result to be %q, got %q", expected, prompt)
	}

}
