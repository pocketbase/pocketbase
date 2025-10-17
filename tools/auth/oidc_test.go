package auth

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TestValidateIdTokenSignature_RSA tests the existing RSA support
func TestValidateIdTokenSignature_RSA(t *testing.T) {
	// Generate RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA key: %v", err)
	}

	// Create a JWT token signed with RSA
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": "1234567890",
		"aud": "test-client",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	token.Header["kid"] = "test-rsa-key"

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		t.Fatalf("Failed to sign token: %v", err)
	}

	// Create mock JWKS server with RSA key
	jwksServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		publicKey := &privateKey.PublicKey

		// Encode RSA public key to JWK format
		n := base64.RawURLEncoding.EncodeToString(publicKey.N.Bytes())
		e := base64.RawURLEncoding.EncodeToString([]byte{byte(publicKey.E >> 16), byte(publicKey.E >> 8), byte(publicKey.E)})

		jwks := map[string]any{
			"keys": []map[string]any{
				{
					"kty": "RSA",
					"kid": "test-rsa-key",
					"use": "sig",
					"alg": "RS256",
					"n":   n,
					"e":   e,
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jwks)
	}))
	defer jwksServer.Close()

	// Test validation
	err = validateIdTokenSignature(context.Background(), tokenString, jwksServer.URL, "test-rsa-key")
	if err != nil {
		t.Errorf("Expected no error for valid RSA token, got: %v", err)
	}
}

// TestValidateIdTokenSignature_Ed25519 tests the new Ed25519 support
func TestValidateIdTokenSignature_Ed25519(t *testing.T) {
	// Generate Ed25519 key pair
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate Ed25519 key: %v", err)
	}

	// Create a JWT token signed with Ed25519
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims{
		"sub": "1234567890",
		"aud": "test-client",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	token.Header["kid"] = "test-ed25519-key"

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		t.Fatalf("Failed to sign token: %v", err)
	}

	// Create mock JWKS server with Ed25519 key
	jwksServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Encode Ed25519 public key to JWK format
		x := base64.RawURLEncoding.EncodeToString(publicKey)

		jwks := map[string]any{
			"keys": []map[string]any{
				{
					"kty": "OKP",
					"kid": "test-ed25519-key",
					"use": "sig",
					"alg": "EdDSA",
					"crv": "Ed25519",
					"x":   x,
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jwks)
	}))
	defer jwksServer.Close()

	// Test validation
	err = validateIdTokenSignature(context.Background(), tokenString, jwksServer.URL, "test-ed25519-key")
	if err != nil {
		t.Errorf("Expected no error for valid Ed25519 token, got: %v", err)
	}
}

// TestValidateIdTokenSignature_UnsupportedKeyType tests error handling for unsupported key types
func TestValidateIdTokenSignature_UnsupportedKeyType(t *testing.T) {
	// Create a dummy token (signature doesn't matter here as we're testing key type validation)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "1234567890",
		"aud": "test-client",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	token.Header["kid"] = "test-ec-key"
	tokenString, _ := token.SignedString([]byte("secret"))

	// Create mock JWKS server with unsupported key type
	jwksServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwks := map[string]any{
			"keys": []map[string]any{
				{
					"kty": "EC",
					"kid": "test-ec-key",
					"use": "sig",
					"alg": "ES256",
					"crv": "P-256",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jwks)
	}))
	defer jwksServer.Close()

	// Test validation
	err := validateIdTokenSignature(context.Background(), tokenString, jwksServer.URL, "test-ec-key")
	if err == nil {
		t.Error("Expected error for unsupported key type, got nil")
	}
	if err != nil && err.Error() != "unsupported key type: EC" {
		t.Errorf("Expected 'unsupported key type: EC' error, got: %v", err)
	}
}

// TestValidateIdTokenSignature_UnsupportedOKPCurve tests error handling for unsupported OKP curves
func TestValidateIdTokenSignature_UnsupportedOKPCurve(t *testing.T) {
	// Create a dummy token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "1234567890",
		"aud": "test-client",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	token.Header["kid"] = "test-x25519-key"
	tokenString, _ := token.SignedString([]byte("secret"))

	// Create mock JWKS server with unsupported OKP curve
	jwksServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwks := map[string]any{
			"keys": []map[string]any{
				{
					"kty": "OKP",
					"kid": "test-x25519-key",
					"use": "sig",
					"alg": "EdDSA",
					"crv": "X25519", // Unsupported curve
					"x":   "test",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jwks)
	}))
	defer jwksServer.Close()

	// Test validation
	err := validateIdTokenSignature(context.Background(), tokenString, jwksServer.URL, "test-x25519-key")
	if err == nil {
		t.Error("Expected error for unsupported OKP curve, got nil")
	}
	if err != nil && err.Error() != "unsupported OKP curve: X25519" {
		t.Errorf("Expected 'unsupported OKP curve: X25519' error, got: %v", err)
	}
}

// TestValidateIdTokenSignature_InvalidEd25519KeyLength tests error handling for invalid Ed25519 key length
func TestValidateIdTokenSignature_InvalidEd25519KeyLength(t *testing.T) {
	// Create a dummy token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "1234567890",
		"aud": "test-client",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	token.Header["kid"] = "test-invalid-key"
	tokenString, _ := token.SignedString([]byte("secret"))

	// Create mock JWKS server with invalid Ed25519 key length
	jwksServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Use an invalid length (Ed25519 requires 32 bytes)
		invalidKey := base64.RawURLEncoding.EncodeToString([]byte("short"))

		jwks := map[string]any{
			"keys": []map[string]any{
				{
					"kty": "OKP",
					"kid": "test-invalid-key",
					"use": "sig",
					"alg": "EdDSA",
					"crv": "Ed25519",
					"x":   invalidKey,
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jwks)
	}))
	defer jwksServer.Close()

	// Test validation
	err := validateIdTokenSignature(context.Background(), tokenString, jwksServer.URL, "test-invalid-key")
	if err == nil {
		t.Error("Expected error for invalid Ed25519 key length, got nil")
	}
	expectedError := fmt.Sprintf("invalid Ed25519 public key length: %d", len("short"))
	if err != nil && err.Error() != expectedError {
		t.Errorf("Expected '%s' error, got: %v", expectedError, err)
	}
}
