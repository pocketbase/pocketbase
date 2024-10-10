package validators_test

import (
	"errors"
	"fmt"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core/validators"
)

func TestJoinValidationErrors(t *testing.T) {
	scenarios := []struct {
		errA     error
		errB     error
		expected string
	}{
		{nil, nil, "<nil>"},
		{errors.New("abc"), nil, "abc"},
		{nil, errors.New("abc"), "abc"},
		{errors.New("abc"), errors.New("456"), "abc\n456"},
		{validation.Errors{"test1": errors.New("test1_err")}, nil, "test1: test1_err."},
		{nil, validation.Errors{"test2": errors.New("test2_err")}, "test2: test2_err."},
		{validation.Errors{}, errors.New("456"), "\n456"},
		{errors.New("456"), validation.Errors{}, "456\n"},
		{validation.Errors{"test1": errors.New("test1_err")}, errors.New("456"), "test1: test1_err."},
		{errors.New("456"), validation.Errors{"test2": errors.New("test2_err")}, "test2: test2_err."},
		{validation.Errors{"test1": errors.New("test1_err")}, validation.Errors{"test2": errors.New("test2_err")}, "test1: test1_err; test2: test2_err."},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%#T_%T", i, s.errA, s.errB), func(t *testing.T) {
			result := fmt.Sprintf("%v", validators.JoinValidationErrors(s.errA, s.errB))
			if result != s.expected {
				t.Fatalf("Expected\n%v\ngot\n%v", s.expected, result)
			}
		})
	}
}
