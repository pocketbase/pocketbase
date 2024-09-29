package router_test

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"strconv"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/tools/router"
)

func TestNewApiErrorWithRawData(t *testing.T) {
	t.Parallel()

	e := router.NewApiError(
		300,
		"message_test",
		"rawData_test",
	)

	result, _ := json.Marshal(e)
	expected := `{"data":{},"message":"Message_test.","status":300}`

	if string(result) != expected {
		t.Errorf("Expected\n%v\ngot\n%v", expected, string(result))
	}

	if e.Error() != "Message_test." {
		t.Errorf("Expected %q, got %q", "Message_test.", e.Error())
	}

	if e.RawData() != "rawData_test" {
		t.Errorf("Expected rawData\n%v\ngot\n%v", "rawData_test", e.RawData())
	}
}

func TestNewApiErrorWithValidationData(t *testing.T) {
	t.Parallel()

	e := router.NewApiError(
		300,
		"message_test",
		map[string]any{
			"err1": errors.New("test error"), // should be normalized
			"err2": validation.ErrRequired,
			"err3": validation.Errors{
				"err3.1": errors.New("test error"), // should be normalized
				"err3.2": validation.ErrRequired,
				"err3.3": validation.Errors{
					"err3.3.1": validation.ErrRequired,
				},
			},
			"err4": &mockSafeErrorItem{},
			"err5": map[string]error{
				"err5.1": validation.ErrRequired,
			},
		},
	)

	result, _ := json.Marshal(e)
	expected := `{"data":{"err1":{"code":"validation_invalid_value","message":"Invalid value."},"err2":{"code":"validation_required","message":"Cannot be blank."},"err3":{"err3.1":{"code":"validation_invalid_value","message":"Invalid value."},"err3.2":{"code":"validation_required","message":"Cannot be blank."},"err3.3":{"err3.3.1":{"code":"validation_required","message":"Cannot be blank."}}},"err4":{"code":"mock_code","message":"Mock_error.","mock_resolve":123},"err5":{"err5.1":{"code":"validation_required","message":"Cannot be blank."}}},"message":"Message_test.","status":300}`

	if string(result) != expected {
		t.Errorf("Expected \n%v, \ngot \n%v", expected, string(result))
	}

	if e.Error() != "Message_test." {
		t.Errorf("Expected %q, got %q", "Message_test.", e.Error())
	}

	if e.RawData() == nil {
		t.Error("Expected non-nil rawData")
	}
}

func TestNewNotFoundError(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		message  string
		data     any
		expected string
	}{
		{"", nil, `{"data":{},"message":"The requested resource wasn't found.","status":404}`},
		{"demo", "rawData_test", `{"data":{},"message":"Demo.","status":404}`},
		{"demo", validation.Errors{"err1": validation.NewError("test_code", "test_message")}, `{"data":{"err1":{"code":"test_code","message":"Test_message."}},"message":"Demo.","status":404}`},
	}

	for i, s := range scenarios {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			e := router.NewNotFoundError(s.message, s.data)
			result, _ := json.Marshal(e)

			if str := string(result); str != s.expected {
				t.Fatalf("Expected\n%v\ngot\n%v", s.expected, str)
			}
		})
	}
}

func TestNewBadRequestError(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		message  string
		data     any
		expected string
	}{
		{"", nil, `{"data":{},"message":"Something went wrong while processing your request.","status":400}`},
		{"demo", "rawData_test", `{"data":{},"message":"Demo.","status":400}`},
		{"demo", validation.Errors{"err1": validation.NewError("test_code", "test_message")}, `{"data":{"err1":{"code":"test_code","message":"Test_message."}},"message":"Demo.","status":400}`},
	}

	for i, s := range scenarios {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			e := router.NewBadRequestError(s.message, s.data)
			result, _ := json.Marshal(e)

			if str := string(result); str != s.expected {
				t.Fatalf("Expected\n%v\ngot\n%v", s.expected, str)
			}
		})
	}
}

func TestNewForbiddenError(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		message  string
		data     any
		expected string
	}{
		{"", nil, `{"data":{},"message":"You are not allowed to perform this request.","status":403}`},
		{"demo", "rawData_test", `{"data":{},"message":"Demo.","status":403}`},
		{"demo", validation.Errors{"err1": validation.NewError("test_code", "test_message")}, `{"data":{"err1":{"code":"test_code","message":"Test_message."}},"message":"Demo.","status":403}`},
	}

	for i, s := range scenarios {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			e := router.NewForbiddenError(s.message, s.data)
			result, _ := json.Marshal(e)

			if str := string(result); str != s.expected {
				t.Fatalf("Expected\n%v\ngot\n%v", s.expected, str)
			}
		})
	}
}

func TestNewUnauthorizedError(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		message  string
		data     any
		expected string
	}{
		{"", nil, `{"data":{},"message":"Missing or invalid authentication.","status":401}`},
		{"demo", "rawData_test", `{"data":{},"message":"Demo.","status":401}`},
		{"demo", validation.Errors{"err1": validation.NewError("test_code", "test_message")}, `{"data":{"err1":{"code":"test_code","message":"Test_message."}},"message":"Demo.","status":401}`},
	}

	for i, s := range scenarios {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			e := router.NewUnauthorizedError(s.message, s.data)
			result, _ := json.Marshal(e)

			if str := string(result); str != s.expected {
				t.Fatalf("Expected\n%v\ngot\n%v", s.expected, str)
			}
		})
	}
}

func TestNewInternalServerError(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		message  string
		data     any
		expected string
	}{
		{"", nil, `{"data":{},"message":"Something went wrong while processing your request.","status":500}`},
		{"demo", "rawData_test", `{"data":{},"message":"Demo.","status":500}`},
		{"demo", validation.Errors{"err1": validation.NewError("test_code", "test_message")}, `{"data":{"err1":{"code":"test_code","message":"Test_message."}},"message":"Demo.","status":500}`},
	}

	for i, s := range scenarios {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			e := router.NewInternalServerError(s.message, s.data)
			result, _ := json.Marshal(e)

			if str := string(result); str != s.expected {
				t.Fatalf("Expected\n%v\ngot\n%v", s.expected, str)
			}
		})
	}
}

func TestNewTooManyRequestsError(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		message  string
		data     any
		expected string
	}{
		{"", nil, `{"data":{},"message":"Too Many Requests.","status":429}`},
		{"demo", "rawData_test", `{"data":{},"message":"Demo.","status":429}`},
		{"demo", validation.Errors{"err1": validation.NewError("test_code", "test_message").SetParams(map[string]any{"test": 123})}, `{"data":{"err1":{"code":"test_code","message":"Test_message.","params":{"test":123}}},"message":"Demo.","status":429}`},
	}

	for i, s := range scenarios {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			e := router.NewTooManyRequestsError(s.message, s.data)
			result, _ := json.Marshal(e)

			if str := string(result); str != s.expected {
				t.Fatalf("Expected\n%v\ngot\n%v", s.expected, str)
			}
		})
	}
}

func TestApiErrorIs(t *testing.T) {
	t.Parallel()

	err0 := router.NewInternalServerError("", nil)
	err1 := router.NewInternalServerError("", nil)
	err2 := errors.New("test")
	err3 := fmt.Errorf("wrapped: %w", err0)

	scenarios := []struct {
		name     string
		err      error
		target   error
		expected bool
	}{
		{
			"nil error",
			err0,
			nil,
			false,
		},
		{
			"non ApiError",
			err0,
			err1,
			false,
		},
		{
			"different ApiError",
			err0,
			err2,
			false,
		},
		{
			"same ApiError",
			err0,
			err0,
			true,
		},
		{
			"wrapped ApiError",
			err3,
			err0,
			true,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			is := errors.Is(s.err, s.target)

			if is != s.expected {
				t.Fatalf("Expected %v, got %v", s.expected, is)
			}
		})
	}
}

func TestToApiError(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		name     string
		err      error
		expected string
	}{
		{
			"regular error",
			errors.New("test"),
			`{"data":{},"message":"Something went wrong while processing your request.","status":400}`,
		},
		{
			"fs.ErrNotExist",
			fs.ErrNotExist,
			`{"data":{},"message":"The requested resource wasn't found.","status":404}`,
		},
		{
			"sql.ErrNoRows",
			sql.ErrNoRows,
			`{"data":{},"message":"The requested resource wasn't found.","status":404}`,
		},
		{
			"ApiError",
			router.NewForbiddenError("test", nil),
			`{"data":{},"message":"Test.","status":403}`,
		},
		{
			"wrapped ApiError",
			fmt.Errorf("wrapped: %w", router.NewForbiddenError("test", nil)),
			`{"data":{},"message":"Test.","status":403}`,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			raw, err := json.Marshal(router.ToApiError(s.err))
			if err != nil {
				t.Fatal(err)
			}
			rawStr := string(raw)

			if rawStr != s.expected {
				t.Fatalf("Expected error\n%vgot\n%v", s.expected, rawStr)
			}
		})
	}
}

// -------------------------------------------------------------------

var (
	_ router.SafeErrorItem     = (*mockSafeErrorItem)(nil)
	_ router.SafeErrorResolver = (*mockSafeErrorItem)(nil)
)

type mockSafeErrorItem struct {
}

func (m *mockSafeErrorItem) Code() string {
	return "mock_code"
}

func (m *mockSafeErrorItem) Error() string {
	return "mock_error"
}

func (m *mockSafeErrorItem) Resolve(errData map[string]any) any {
	errData["mock_resolve"] = 123
	return errData
}
