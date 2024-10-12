package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"golang.org/x/oauth2"
)

func TestContext(t *testing.T) {
	b := BaseProvider{}

	before := b.Scopes()
	if before != nil {
		t.Errorf("Expected nil context, got %v", before)
	}

	b.SetContext(context.Background())

	after := b.Scopes()
	if after != nil {
		t.Error("Expected non-nil context")
	}
}

func TestDisplayName(t *testing.T) {
	b := BaseProvider{}

	before := b.DisplayName()
	if before != "" {
		t.Fatalf("Expected displayName to be empty, got %v", before)
	}

	b.SetDisplayName("test")

	after := b.DisplayName()
	if after != "test" {
		t.Fatalf("Expected displayName to be 'test', got %v", after)
	}
}

func TestPKCE(t *testing.T) {
	b := BaseProvider{}

	before := b.PKCE()
	if before != false {
		t.Fatalf("Expected pkce to be %v, got %v", false, before)
	}

	b.SetPKCE(true)

	after := b.PKCE()
	if after != true {
		t.Fatalf("Expected pkce to be %v, got %v", true, after)
	}
}

func TestScopes(t *testing.T) {
	b := BaseProvider{}

	before := b.Scopes()
	if len(before) != 0 {
		t.Fatalf("Expected 0 scopes, got %v", before)
	}

	b.SetScopes([]string{"test1", "test2"})

	after := b.Scopes()
	if len(after) != 2 {
		t.Fatalf("Expected 2 scopes, got %v", after)
	}
}

func TestClientId(t *testing.T) {
	b := BaseProvider{}

	before := b.ClientId()
	if before != "" {
		t.Fatalf("Expected clientId to be empty, got %v", before)
	}

	b.SetClientId("test")

	after := b.ClientId()
	if after != "test" {
		t.Fatalf("Expected clientId to be 'test', got %v", after)
	}
}

func TestClientSecret(t *testing.T) {
	b := BaseProvider{}

	before := b.ClientSecret()
	if before != "" {
		t.Fatalf("Expected clientSecret to be empty, got %v", before)
	}

	b.SetClientSecret("test")

	after := b.ClientSecret()
	if after != "test" {
		t.Fatalf("Expected clientSecret to be 'test', got %v", after)
	}
}

func TestRedirectURL(t *testing.T) {
	b := BaseProvider{}

	before := b.RedirectURL()
	if before != "" {
		t.Fatalf("Expected RedirectURL to be empty, got %v", before)
	}

	b.SetRedirectURL("test")

	after := b.RedirectURL()
	if after != "test" {
		t.Fatalf("Expected RedirectURL to be 'test', got %v", after)
	}
}

func TestAuthURL(t *testing.T) {
	b := BaseProvider{}

	before := b.AuthURL()
	if before != "" {
		t.Fatalf("Expected authURL to be empty, got %v", before)
	}

	b.SetAuthURL("test")

	after := b.AuthURL()
	if after != "test" {
		t.Fatalf("Expected authURL to be 'test', got %v", after)
	}
}

func TestTokenURL(t *testing.T) {
	b := BaseProvider{}

	before := b.TokenURL()
	if before != "" {
		t.Fatalf("Expected tokenURL to be empty, got %v", before)
	}

	b.SetTokenURL("test")

	after := b.TokenURL()
	if after != "test" {
		t.Fatalf("Expected tokenURL to be 'test', got %v", after)
	}
}

func TestUserInfoURL(t *testing.T) {
	b := BaseProvider{}

	before := b.UserInfoURL()
	if before != "" {
		t.Fatalf("Expected userInfoURL to be empty, got %v", before)
	}

	b.SetUserInfoURL("test")

	after := b.UserInfoURL()
	if after != "test" {
		t.Fatalf("Expected userInfoURL to be 'test', got %v", after)
	}
}

func TestExtra(t *testing.T) {
	b := BaseProvider{}

	before := b.Extra()
	if before != nil {
		t.Fatalf("Expected extra to be empty, got %v", before)
	}

	extra := map[string]any{"a": 1, "b": 2}

	b.SetExtra(extra)

	after := b.Extra()

	rawExtra, err := json.Marshal(extra)
	if err != nil {
		t.Fatal(err)
	}

	rawAfter, err := json.Marshal(after)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(rawExtra, rawAfter) {
		t.Fatalf("Expected extra to be\n%s\ngot\n%s", rawExtra, rawAfter)
	}

	// ensure that it was shallow copied
	after["b"] = 3
	if d := b.Extra(); d["b"] != 2 {
		t.Fatalf("Expected extra to remain unchanged, got\n%v", d)
	}
}

func TestBuildAuthURL(t *testing.T) {
	b := BaseProvider{
		authURL:      "authURL_test",
		tokenURL:     "tokenURL_test",
		redirectURL:  "redirectURL_test",
		clientId:     "clientId_test",
		clientSecret: "clientSecret_test",
		scopes:       []string{"test_scope"},
	}

	expected := "authURL_test?access_type=offline&client_id=clientId_test&prompt=consent&redirect_uri=redirectURL_test&response_type=code&scope=test_scope&state=state_test"
	result := b.BuildAuthURL("state_test", oauth2.AccessTypeOffline, oauth2.ApprovalForce)

	if result != expected {
		t.Errorf("Expected auth url %q, got %q", expected, result)
	}
}

func TestClient(t *testing.T) {
	b := BaseProvider{}

	result := b.Client(&oauth2.Token{})
	if result == nil {
		t.Error("Expected *http.Client instance, got nil")
	}
}

func TestOauth2Config(t *testing.T) {
	b := BaseProvider{
		authURL:      "authURL_test",
		tokenURL:     "tokenURL_test",
		redirectURL:  "redirectURL_test",
		clientId:     "clientId_test",
		clientSecret: "clientSecret_test",
		scopes:       []string{"test"},
	}

	result := b.oauth2Config()

	if result.RedirectURL != b.RedirectURL() {
		t.Errorf("Expected redirectURL %s, got %s", b.RedirectURL(), result.RedirectURL)
	}

	if result.ClientID != b.ClientId() {
		t.Errorf("Expected clientId %s, got %s", b.ClientId(), result.ClientID)
	}

	if result.ClientSecret != b.ClientSecret() {
		t.Errorf("Expected clientSecret %s, got %s", b.ClientSecret(), result.ClientSecret)
	}

	if result.Endpoint.AuthURL != b.AuthURL() {
		t.Errorf("Expected authURL %s, got %s", b.AuthURL(), result.Endpoint.AuthURL)
	}

	if result.Endpoint.TokenURL != b.TokenURL() {
		t.Errorf("Expected authURL %s, got %s", b.TokenURL(), result.Endpoint.TokenURL)
	}

	if len(result.Scopes) != len(b.Scopes()) || result.Scopes[0] != b.Scopes()[0] {
		t.Errorf("Expected scopes %s, got %s", b.Scopes(), result.Scopes)
	}
}
