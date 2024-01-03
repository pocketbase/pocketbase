package apis_test

import (
	"encoding/json"
	"errors"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/apis"
)

func TestNewApiErrorWithRawData(t *testing.T) {
	t.Parallel()

	e := apis.NewApiError(
		300,
		"message_test",
		"rawData_test",
	)

	result, _ := json.Marshal(e)
	expected := `{"code":300,"message":"Message_test.","data":{}}`

	if string(result) != expected {
		t.Errorf("Expected %v, got %v", expected, string(result))
	}

	if e.Error() != "Message_test." {
		t.Errorf("Expected %q, got %q", "Message_test.", e.Error())
	}

	if e.RawData() != "rawData_test" {
		t.Errorf("Expected rawData %v, got %v", "rawData_test", e.RawData())
	}
}

func TestNewApiErrorWithValidationData(t *testing.T) {
	t.Parallel()

	e := apis.NewApiError(
		300,
		"message_test",
		validation.Errors{
			"err1": errors.New("test error"), // should be normalized
			"err2": validation.ErrRequired,
			"err3": validation.Errors{
				"sub1": errors.New("test error"), // should be normalized
				"sub2": validation.ErrRequired,
				"sub3": validation.Errors{
					"sub11": validation.ErrRequired,
				},
			},
		},
	)

	result, _ := json.Marshal(e)
	expected := `{"code":300,"message":"Message_test.","data":{"err1":{"code":"validation_invalid_value","message":"Invalid value."},"err2":{"code":"validation_required","message":"Cannot be blank."},"err3":{"sub1":{"code":"validation_invalid_value","message":"Invalid value."},"sub2":{"code":"validation_required","message":"Cannot be blank."},"sub3":{"sub11":{"code":"validation_required","message":"Cannot be blank."}}}}}`

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
		{"", nil, `{"code":404,"message":"The requested resource wasn't found.","data":{}}`},
		{"demo", "rawData_test", `{"code":404,"message":"Demo.","data":{}}`},
		{"demo", validation.Errors{"err1": validation.NewError("test_code", "test_message")}, `{"code":404,"message":"Demo.","data":{"err1":{"code":"test_code","message":"Test_message."}}}`},
	}

	for i, scenario := range scenarios {
		e := apis.NewNotFoundError(scenario.message, scenario.data)
		result, _ := json.Marshal(e)

		if string(result) != scenario.expected {
			t.Errorf("(%d) Expected \n%v, \ngot \n%v", i, scenario.expected, string(result))
		}
	}
}

func TestNewBadRequestError(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		message  string
		data     any
		expected string
	}{
		{"", nil, `{"code":400,"message":"Something went wrong while processing your request.","data":{}}`},
		{"demo", "rawData_test", `{"code":400,"message":"Demo.","data":{}}`},
		{"demo", validation.Errors{"err1": validation.NewError("test_code", "test_message")}, `{"code":400,"message":"Demo.","data":{"err1":{"code":"test_code","message":"Test_message."}}}`},
	}

	for i, scenario := range scenarios {
		e := apis.NewBadRequestError(scenario.message, scenario.data)
		result, _ := json.Marshal(e)

		if string(result) != scenario.expected {
			t.Errorf("(%d) Expected \n%v, \ngot \n%v", i, scenario.expected, string(result))
		}
	}
}

func TestNewForbiddenError(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		message  string
		data     any
		expected string
	}{
		{"", nil, `{"code":403,"message":"You are not allowed to perform this request.","data":{}}`},
		{"demo", "rawData_test", `{"code":403,"message":"Demo.","data":{}}`},
		{"demo", validation.Errors{"err1": validation.NewError("test_code", "test_message")}, `{"code":403,"message":"Demo.","data":{"err1":{"code":"test_code","message":"Test_message."}}}`},
	}

	for i, scenario := range scenarios {
		e := apis.NewForbiddenError(scenario.message, scenario.data)
		result, _ := json.Marshal(e)

		if string(result) != scenario.expected {
			t.Errorf("(%d) Expected \n%v, \ngot \n%v", i, scenario.expected, string(result))
		}
	}
}

func TestNewUnauthorizedError(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		message  string
		data     any
		expected string
	}{
		{"", nil, `{"code":401,"message":"Missing or invalid authentication token.","data":{}}`},
		{"demo", "rawData_test", `{"code":401,"message":"Demo.","data":{}}`},
		{"demo", validation.Errors{"err1": validation.NewError("test_code", "test_message")}, `{"code":401,"message":"Demo.","data":{"err1":{"code":"test_code","message":"Test_message."}}}`},
	}

	for i, scenario := range scenarios {
		e := apis.NewUnauthorizedError(scenario.message, scenario.data)
		result, _ := json.Marshal(e)

		if string(result) != scenario.expected {
			t.Errorf("(%d) Expected \n%v, \ngot \n%v", i, scenario.expected, string(result))
		}
	}
}
