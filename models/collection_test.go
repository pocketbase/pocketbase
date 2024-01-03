package models_test

import (
	"encoding/json"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/types"
)

func TestCollectionTableName(t *testing.T) {
	t.Parallel()

	m := models.Collection{}
	if m.TableName() != "_collections" {
		t.Fatalf("Unexpected table name, got %q", m.TableName())
	}
}

func TestCollectionBaseFilesPath(t *testing.T) {
	t.Parallel()

	m := models.Collection{}

	m.RefreshId()

	expected := m.Id
	if m.BaseFilesPath() != expected {
		t.Fatalf("Expected path %s, got %s", expected, m.BaseFilesPath())
	}
}

func TestCollectionIsBase(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		collection models.Collection
		expected   bool
	}{
		{models.Collection{}, false},
		{models.Collection{Type: "unknown"}, false},
		{models.Collection{Type: models.CollectionTypeBase}, true},
		{models.Collection{Type: models.CollectionTypeAuth}, false},
	}

	for i, s := range scenarios {
		result := s.collection.IsBase()
		if result != s.expected {
			t.Errorf("(%d) Expected %v, got %v", i, s.expected, result)
		}
	}
}

func TestCollectionIsAuth(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		collection models.Collection
		expected   bool
	}{
		{models.Collection{}, false},
		{models.Collection{Type: "unknown"}, false},
		{models.Collection{Type: models.CollectionTypeBase}, false},
		{models.Collection{Type: models.CollectionTypeAuth}, true},
	}

	for i, s := range scenarios {
		result := s.collection.IsAuth()
		if result != s.expected {
			t.Errorf("(%d) Expected %v, got %v", i, s.expected, result)
		}
	}
}

func TestCollectionMarshalJSON(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		name       string
		collection models.Collection
		expected   string
	}{
		{
			"no type",
			models.Collection{Name: "test"},
			`{"id":"","created":"","updated":"","name":"test","type":"","system":false,"schema":[],"indexes":[],"listRule":null,"viewRule":null,"createRule":null,"updateRule":null,"deleteRule":null,"options":{}}`,
		},
		{
			"unknown type + non empty options",
			models.Collection{Name: "test", Type: "unknown", ListRule: types.Pointer("test_list"), Options: types.JsonMap{"test": 123}, Indexes: types.JsonArray[string]{"idx_test"}},
			`{"id":"","created":"","updated":"","name":"test","type":"unknown","system":false,"schema":[],"indexes":["idx_test"],"listRule":"test_list","viewRule":null,"createRule":null,"updateRule":null,"deleteRule":null,"options":{}}`,
		},
		{
			"base type + non empty options",
			models.Collection{Name: "test", Type: models.CollectionTypeBase, ListRule: types.Pointer("test_list"), Options: types.JsonMap{"test": 123}},
			`{"id":"","created":"","updated":"","name":"test","type":"base","system":false,"schema":[],"indexes":[],"listRule":"test_list","viewRule":null,"createRule":null,"updateRule":null,"deleteRule":null,"options":{}}`,
		},
		{
			"auth type + non empty options",
			models.Collection{BaseModel: models.BaseModel{Id: "test"}, Type: models.CollectionTypeAuth, Options: types.JsonMap{"test": 123, "allowOAuth2Auth": true, "minPasswordLength": 4, "onlyVerified": true}},
			`{"id":"test","created":"","updated":"","name":"","type":"auth","system":false,"schema":[],"indexes":[],"listRule":null,"viewRule":null,"createRule":null,"updateRule":null,"deleteRule":null,"options":{"allowEmailAuth":false,"allowOAuth2Auth":true,"allowUsernameAuth":false,"exceptEmailDomains":null,"manageRule":null,"minPasswordLength":4,"onlyEmailDomains":null,"onlyVerified":true,"requireEmail":false}}`,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			result, err := s.collection.MarshalJSON()
			if err != nil {
				t.Fatalf("Unexpected error %v", err)
			}

			if string(result) != s.expected {
				t.Fatalf("Expected\n%v\ngot\n%v", s.expected, string(result))
			}
		})
	}
}

func TestCollectionBaseOptions(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		name       string
		collection models.Collection
		expected   string
	}{
		{
			"no type",
			models.Collection{Options: types.JsonMap{"test": 123}},
			"{}",
		},
		{
			"unknown type",
			models.Collection{Type: "anything", Options: types.JsonMap{"test": 123}},
			"{}",
		},
		{
			"different type",
			models.Collection{Type: models.CollectionTypeAuth, Options: types.JsonMap{"test": 123, "minPasswordLength": 4}},
			"{}",
		},
		{
			"base type",
			models.Collection{Type: models.CollectionTypeBase, Options: types.JsonMap{"test": 123}},
			"{}",
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			result := s.collection.BaseOptions()

			encoded, err := json.Marshal(result)
			if err != nil {
				t.Fatal(err)
			}

			if strEncoded := string(encoded); strEncoded != s.expected {
				t.Fatalf("Expected \n%v \ngot \n%v", s.expected, strEncoded)
			}
		})
	}
}

func TestCollectionAuthOptions(t *testing.T) {
	t.Parallel()

	options := types.JsonMap{"test": 123, "minPasswordLength": 4}
	expectedSerialization := `{"manageRule":null,"allowOAuth2Auth":false,"allowUsernameAuth":false,"allowEmailAuth":false,"requireEmail":false,"exceptEmailDomains":null,"onlyVerified":false,"onlyEmailDomains":null,"minPasswordLength":4}`

	scenarios := []struct {
		name       string
		collection models.Collection
		expected   string
	}{
		{
			"no type",
			models.Collection{Options: options},
			expectedSerialization,
		},
		{
			"unknown type",
			models.Collection{Type: "anything", Options: options},
			expectedSerialization,
		},
		{
			"different type",
			models.Collection{Type: models.CollectionTypeBase, Options: options},
			expectedSerialization,
		},
		{
			"auth type",
			models.Collection{Type: models.CollectionTypeAuth, Options: options},
			expectedSerialization,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			result := s.collection.AuthOptions()

			encoded, err := json.Marshal(result)
			if err != nil {
				t.Fatal(err)
			}

			if strEncoded := string(encoded); strEncoded != s.expected {
				t.Fatalf("Expected \n%v \ngot \n%v", s.expected, strEncoded)
			}
		})
	}
}

func TestCollectionViewOptions(t *testing.T) {
	t.Parallel()

	options := types.JsonMap{"query": "select id from demo1", "minPasswordLength": 4}
	expectedSerialization := `{"query":"select id from demo1"}`

	scenarios := []struct {
		name       string
		collection models.Collection
		expected   string
	}{
		{
			"no type",
			models.Collection{Options: options},
			expectedSerialization,
		},
		{
			"unknown type",
			models.Collection{Type: "anything", Options: options},
			expectedSerialization,
		},
		{
			"different type",
			models.Collection{Type: models.CollectionTypeBase, Options: options},
			expectedSerialization,
		},
		{
			"view type",
			models.Collection{Type: models.CollectionTypeView, Options: options},
			expectedSerialization,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			result := s.collection.ViewOptions()

			encoded, err := json.Marshal(result)
			if err != nil {
				t.Fatal(err)
			}

			if strEncoded := string(encoded); strEncoded != s.expected {
				t.Fatalf("Expected \n%v \ngot \n%v", s.expected, strEncoded)
			}
		})
	}
}

func TestNormalizeOptions(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		name       string
		collection models.Collection
		expected   string // serialized options
	}{
		{
			"unknown type",
			models.Collection{Type: "unknown", Options: types.JsonMap{"test": 123, "minPasswordLength": 4}},
			"{}",
		},
		{
			"base type",
			models.Collection{Type: models.CollectionTypeBase, Options: types.JsonMap{"test": 123, "minPasswordLength": 4}},
			"{}",
		},
		{
			"auth type",
			models.Collection{Type: models.CollectionTypeAuth, Options: types.JsonMap{"test": 123, "minPasswordLength": 4}},
			`{"allowEmailAuth":false,"allowOAuth2Auth":false,"allowUsernameAuth":false,"exceptEmailDomains":null,"manageRule":null,"minPasswordLength":4,"onlyEmailDomains":null,"onlyVerified":false,"requireEmail":false}`,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			if err := s.collection.NormalizeOptions(); err != nil {
				t.Fatalf("Unexpected error %v", err)
			}

			encoded, err := json.Marshal(s.collection.Options)
			if err != nil {
				t.Fatal(err)
			}

			if strEncoded := string(encoded); strEncoded != s.expected {
				t.Fatalf("Expected \n%v \ngot \n%v", s.expected, strEncoded)
			}
		})
	}
}

func TestDecodeOptions(t *testing.T) {
	t.Parallel()

	m := models.Collection{
		Options: types.JsonMap{"test": 123},
	}

	result := struct {
		Test int
	}{}

	if err := m.DecodeOptions(&result); err != nil {
		t.Fatal(err)
	}

	if result.Test != 123 {
		t.Fatalf("Expected %v, got %v", 123, result.Test)
	}
}

func TestSetOptions(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		name       string
		collection models.Collection
		options    any
		expected   string // serialized options
	}{
		{
			"no type",
			models.Collection{},
			map[string]any{},
			"{}",
		},
		{
			"unknown type + non empty options",
			models.Collection{Type: "unknown", Options: types.JsonMap{"test": 123}},
			map[string]any{"test": 456, "minPasswordLength": 4},
			"{}",
		},
		{
			"base type",
			models.Collection{Type: models.CollectionTypeBase, Options: types.JsonMap{"test": 123}},
			map[string]any{"test": 456, "minPasswordLength": 4},
			"{}",
		},
		{
			"auth type",
			models.Collection{Type: models.CollectionTypeAuth, Options: types.JsonMap{"test": 123}},
			map[string]any{"test": 456, "minPasswordLength": 4},
			`{"allowEmailAuth":false,"allowOAuth2Auth":false,"allowUsernameAuth":false,"exceptEmailDomains":null,"manageRule":null,"minPasswordLength":4,"onlyEmailDomains":null,"onlyVerified":false,"requireEmail":false}`,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			if err := s.collection.SetOptions(s.options); err != nil {
				t.Fatalf("Unexpected error %v", err)
			}

			encoded, err := json.Marshal(s.collection.Options)
			if err != nil {
				t.Fatal(err)
			}

			if strEncoded := string(encoded); strEncoded != s.expected {
				t.Fatalf("Expected\n%v\ngot\n%v", s.expected, strEncoded)
			}
		})
	}
}

func TestCollectionBaseOptionsValidate(t *testing.T) {
	t.Parallel()

	opt := models.CollectionBaseOptions{}
	if err := opt.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestCollectionAuthOptionsValidate(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		name           string
		options        models.CollectionAuthOptions
		expectedErrors []string
	}{
		{
			"empty",
			models.CollectionAuthOptions{},
			nil,
		},
		{
			"empty string ManageRule",
			models.CollectionAuthOptions{ManageRule: types.Pointer("")},
			[]string{"manageRule"},
		},
		{
			"minPasswordLength < 5",
			models.CollectionAuthOptions{MinPasswordLength: 3},
			[]string{"minPasswordLength"},
		},
		{
			"minPasswordLength > 72",
			models.CollectionAuthOptions{MinPasswordLength: 73},
			[]string{"minPasswordLength"},
		},
		{
			"both OnlyDomains and ExceptDomains set",
			models.CollectionAuthOptions{
				OnlyEmailDomains:   []string{"example.com", "test.com"},
				ExceptEmailDomains: []string{"example.com", "test.com"},
			},
			[]string{"onlyEmailDomains", "exceptEmailDomains"},
		},
		{
			"only OnlyDomains set",
			models.CollectionAuthOptions{
				OnlyEmailDomains: []string{"example.com", "test.com"},
			},
			[]string{},
		},
		{
			"only ExceptEmailDomains set",
			models.CollectionAuthOptions{
				ExceptEmailDomains: []string{"example.com", "test.com"},
			},
			[]string{},
		},
		{
			"all fields with valid data",
			models.CollectionAuthOptions{
				ManageRule:         types.Pointer("test"),
				AllowOAuth2Auth:    true,
				AllowUsernameAuth:  true,
				AllowEmailAuth:     true,
				RequireEmail:       true,
				ExceptEmailDomains: []string{"example.com", "test.com"},
				OnlyEmailDomains:   nil,
				MinPasswordLength:  5,
			},
			[]string{},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			result := s.options.Validate()

			// parse errors
			errs, ok := result.(validation.Errors)
			if !ok && result != nil {
				t.Fatalf("Failed to parse errors %v", result)
			}

			if len(errs) != len(s.expectedErrors) {
				t.Fatalf("Expected error keys %v, got errors \n%v", s.expectedErrors, result)
			}

			for key := range errs {
				if !list.ExistInSlice(key, s.expectedErrors) {
					t.Fatalf("Unexpected error key %q in \n%v", key, errs)
				}
			}
		})
	}
}

func TestCollectionViewOptionsValidate(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		name           string
		options        models.CollectionViewOptions
		expectedErrors []string
	}{
		{
			"empty",
			models.CollectionViewOptions{},
			[]string{"query"},
		},
		{
			"valid data",
			models.CollectionViewOptions{
				Query: "test123",
			},
			[]string{},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			result := s.options.Validate()

			// parse errors
			errs, ok := result.(validation.Errors)
			if !ok && result != nil {
				t.Fatalf("Failed to parse errors %v", result)
			}

			if len(errs) != len(s.expectedErrors) {
				t.Fatalf("Expected error keys %v, got errors \n%v", s.expectedErrors, result)
			}

			for key := range errs {
				if !list.ExistInSlice(key, s.expectedErrors) {
					t.Fatalf("Unexpected error key %q in \n%v", key, errs)
				}
			}
		})
	}
}
