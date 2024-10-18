package forms_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/forms"
)

func TestRealtimeSubscribeValidate(t *testing.T) {
	t.Parallel()

	validSubscriptionsLimit := make([]string, 1000)
	for i := 0; i < len(validSubscriptionsLimit); i++ {
		validSubscriptionsLimit[i] = fmt.Sprintf(`"%d"`, i)
	}
	invalidSubscriptionsLimit := make([]string, 1001)
	for i := 0; i < len(invalidSubscriptionsLimit); i++ {
		invalidSubscriptionsLimit[i] = fmt.Sprintf(`"%d"`, i)
	}

	scenarios := []struct {
		name           string
		data           string
		expectedErrors []string
	}{
		{
			"empty data",
			`{}`,
			[]string{"clientId"},
		},
		{
			"clientId > max chars limit",
			`{"clientId":"` + strings.Repeat("a", 256) + `"}`,
			[]string{"clientId"},
		},
		{
			"clientId <= max chars limit",
			`{"clientId":"` + strings.Repeat("a", 255) + `"}`,
			[]string{},
		},
		{
			"subscriptions > max limit",
			`{"clientId":"test", "subscriptions":[` + strings.Join(invalidSubscriptionsLimit, ",") + `]}`,
			[]string{"subscriptions"},
		},
		{
			"subscriptions <= max limit",
			`{"clientId":"test", "subscriptions":[` + strings.Join(validSubscriptionsLimit, ",") + `]}`,
			[]string{},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			form := forms.NewRealtimeSubscribe()

			err := json.Unmarshal([]byte(s.data), &form)
			if err != nil {
				t.Fatal(err)
			}

			result := form.Validate()

			// parse errors
			errs, ok := result.(validation.Errors)
			if !ok && result != nil {
				t.Fatalf("Failed to parse errors %v", result)
				return
			}

			// check errors
			if len(errs) > len(s.expectedErrors) {
				t.Fatalf("Expected error keys %v, got %v", s.expectedErrors, errs)
			}
			for _, k := range s.expectedErrors {
				if _, ok := errs[k]; !ok {
					t.Fatalf("Missing expected error key %q in %v", k, errs)
				}
			}
		})
	}
}
