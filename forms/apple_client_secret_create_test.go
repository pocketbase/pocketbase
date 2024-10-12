package forms_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/tests"
)

func TestAppleClientSecretCreateValidateAndSubmit(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	encodedKey, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatal(err)
	}

	privatePem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: encodedKey,
		},
	)

	scenarios := []struct {
		name        string
		formData    map[string]any
		expectError bool
	}{
		{
			"empty data",
			map[string]any{},
			true,
		},
		{
			"invalid data",
			map[string]any{
				"clientId":   "",
				"teamId":     "123456789",
				"keyId":      "123456789",
				"privateKey": "-----BEGIN PRIVATE KEY----- invalid -----END PRIVATE KEY-----",
				"duration":   -1,
			},
			true,
		},
		{
			"valid data",
			map[string]any{
				"clientId":   "123",
				"teamId":     "1234567890",
				"keyId":      "1234567891",
				"privateKey": string(privatePem),
				"duration":   1,
			},
			false,
		},
	}

	for _, s := range scenarios {
		form := forms.NewAppleClientSecretCreate(app)

		rawData, marshalErr := json.Marshal(s.formData)
		if marshalErr != nil {
			t.Errorf("[%s] Failed to marshalize the scenario data: %v", s.name, marshalErr)
			continue
		}

		// load data
		loadErr := json.Unmarshal(rawData, form)
		if loadErr != nil {
			t.Errorf("[%s] Failed to load form data: %v", s.name, loadErr)
			continue
		}

		secret, err := form.Submit()

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("[%s] Expected hasErr %v, got %v (%v)", s.name, s.expectError, hasErr, err)
		}

		if hasErr {
			continue
		}

		if secret == "" {
			t.Errorf("[%s] Expected non-empty secret", s.name)
		}

		claims := jwt.MapClaims{}
		token, _, err := jwt.NewParser().ParseUnverified(secret, claims)
		if err != nil {
			t.Errorf("[%s] Failed to parse token: %v", s.name, err)
		}

		if alg := token.Header["alg"]; alg != "ES256" {
			t.Errorf("[%s] Expected %q alg header, got %q", s.name, "ES256", alg)
		}

		if kid := token.Header["kid"]; kid != form.KeyId {
			t.Errorf("[%s] Expected %q kid header, got %q", s.name, form.KeyId, kid)
		}
	}
}
