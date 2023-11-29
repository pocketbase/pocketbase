package auth

import (
	"context"
	"testing"

	"golang.org/x/oauth2"
)

func TestContext(t *testing.T) {
	b := baseProvider{}

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
	b := baseProvider{}

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
	b := baseProvider{}

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
	b := baseProvider{}

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
	b := baseProvider{}

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
	b := baseProvider{}

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

func TestRedirectUrl(t *testing.T) {
	b := baseProvider{}

	before := b.RedirectUrl()
	if before != "" {
		t.Fatalf("Expected RedirectUrl to be empty, got %v", before)
	}

	b.SetRedirectUrl("test")

	after := b.RedirectUrl()
	if after != "test" {
		t.Fatalf("Expected RedirectUrl to be 'test', got %v", after)
	}
}

func TestAuthUrl(t *testing.T) {
	b := baseProvider{}

	before := b.AuthUrl()
	if before != "" {
		t.Fatalf("Expected authUrl to be empty, got %v", before)
	}

	b.SetAuthUrl("test")

	after := b.AuthUrl()
	if after != "test" {
		t.Fatalf("Expected authUrl to be 'test', got %v", after)
	}
}

func TestTokenUrl(t *testing.T) {
	b := baseProvider{}

	before := b.TokenUrl()
	if before != "" {
		t.Fatalf("Expected tokenUrl to be empty, got %v", before)
	}

	b.SetTokenUrl("test")

	after := b.TokenUrl()
	if after != "test" {
		t.Fatalf("Expected tokenUrl to be 'test', got %v", after)
	}
}

func TestUserApiUrl(t *testing.T) {
	b := baseProvider{}

	before := b.UserApiUrl()
	if before != "" {
		t.Fatalf("Expected userApiUrl to be empty, got %v", before)
	}

	b.SetUserApiUrl("test")

	after := b.UserApiUrl()
	if after != "test" {
		t.Fatalf("Expected userApiUrl to be 'test', got %v", after)
	}
}

func TestBuildAuthUrl(t *testing.T) {
	b := baseProvider{
		authUrl:      "authUrl_test",
		tokenUrl:     "tokenUrl_test",
		redirectUrl:  "redirectUrl_test",
		clientId:     "clientId_test",
		clientSecret: "clientSecret_test",
		scopes:       []string{"test_scope"},
	}

	expected := "authUrl_test?access_type=offline&client_id=clientId_test&prompt=consent&redirect_uri=redirectUrl_test&response_type=code&scope=test_scope&state=state_test"
	result := b.BuildAuthUrl("state_test", oauth2.AccessTypeOffline, oauth2.ApprovalForce)

	if result != expected {
		t.Errorf("Expected auth url %q, got %q", expected, result)
	}
}

func TestClient(t *testing.T) {
	b := baseProvider{}

	result := b.Client(&oauth2.Token{})
	if result == nil {
		t.Error("Expected *http.Client instance, got nil")
	}
}

func TestOauth2Config(t *testing.T) {
	b := baseProvider{
		authUrl:      "authUrl_test",
		tokenUrl:     "tokenUrl_test",
		redirectUrl:  "redirectUrl_test",
		clientId:     "clientId_test",
		clientSecret: "clientSecret_test",
		scopes:       []string{"test"},
	}

	result := b.oauth2Config()

	if result.RedirectURL != b.RedirectUrl() {
		t.Errorf("Expected redirectUrl %s, got %s", b.RedirectUrl(), result.RedirectURL)
	}

	if result.ClientID != b.ClientId() {
		t.Errorf("Expected clientId %s, got %s", b.ClientId(), result.ClientID)
	}

	if result.ClientSecret != b.ClientSecret() {
		t.Errorf("Expected clientSecret %s, got %s", b.ClientSecret(), result.ClientSecret)
	}

	if result.Endpoint.AuthURL != b.AuthUrl() {
		t.Errorf("Expected authUrl %s, got %s", b.AuthUrl(), result.Endpoint.AuthURL)
	}

	if result.Endpoint.TokenURL != b.TokenUrl() {
		t.Errorf("Expected authUrl %s, got %s", b.TokenUrl(), result.Endpoint.TokenURL)
	}

	if len(result.Scopes) != len(b.Scopes()) || result.Scopes[0] != b.Scopes()[0] {
		t.Errorf("Expected scopes %s, got %s", b.Scopes(), result.Scopes)
	}
}
