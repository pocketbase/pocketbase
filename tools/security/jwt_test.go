package security_test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pocketbase/pocketbase/tools/security"
)

func TestParseUnverifiedJWT(t *testing.T) {
	// invalid formatted JWT
	result1, err1 := security.ParseUnverifiedJWT("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoidGVzdCJ9")
	if err1 == nil {
		t.Error("Expected error got nil")
	}
	if len(result1) > 0 {
		t.Error("Expected no parsed claims, got", result1)
	}

	// properly formatted JWT with INVALID claims
	// {"name": "test", "exp":1516239022}
	result2, err2 := security.ParseUnverifiedJWT("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoidGVzdCIsImV4cCI6MTUxNjIzOTAyMn0.xYHirwESfSEW3Cq2BL47CEASvD_p_ps3QCA54XtNktU")
	if err2 == nil {
		t.Error("Expected error got nil")
	}
	if len(result2) != 2 || result2["name"] != "test" {
		t.Errorf("Expected to have 2 claims, got %v", result2)
	}

	// properly formatted JWT with VALID claims (missing exp)
	// {"name": "test"}
	result3, err3 := security.ParseUnverifiedJWT("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoidGVzdCJ9.ml0QsTms3K9wMygTu41ZhKlTyjmW9zHQtoS8FUsCCjU")
	if err3 != nil {
		t.Error("Expected nil, got", err3)
	}
	if len(result3) != 1 || result3["name"] != "test" {
		t.Errorf("Expected to have 1 claim, got %v", result3)
	}

	// properly formatted JWT with VALID claims (valid exp)
	// {"name": "test", "exp": 2208985261}
	result4, err4 := security.ParseUnverifiedJWT("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoidGVzdCIsImV4cCI6MjIwODk4NTI2MX0._0KQu60hYNx5wkBIpEaoX35shXRicb0X_0VdWKWb-3k")
	if err4 != nil {
		t.Error("Expected nil, got", err4)
	}
	if len(result4) != 2 || result4["name"] != "test" {
		t.Errorf("Expected to have 2 claims, got %v", result4)
	}
}

func TestParseJWT(t *testing.T) {
	scenarios := []struct {
		token        string
		secret       string
		expectError  bool
		expectClaims jwt.MapClaims
	}{
		// invalid formatted JWT
		{
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoidGVzdCJ9",
			"test",
			true,
			nil,
		},
		// properly formatted JWT with INVALID claims and INVALID secret
		// {"name": "test", "exp": 1516239022}
		{
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoidGVzdCIsImV4cCI6MTUxNjIzOTAyMn0.xYHirwESfSEW3Cq2BL47CEASvD_p_ps3QCA54XtNktU",
			"invalid",
			true,
			nil,
		},
		// properly formatted JWT with INVALID claims and VALID secret
		// {"name": "test", "exp": 1516239022}
		{
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoidGVzdCIsImV4cCI6MTUxNjIzOTAyMn0.xYHirwESfSEW3Cq2BL47CEASvD_p_ps3QCA54XtNktU",
			"test",
			true,
			nil,
		},
		// properly formatted JWT with VALID claims and INVALID secret
		// {"name": "test", "exp": 1898636137}
		{
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoidGVzdCIsImV4cCI6MTg5ODYzNjEzN30.gqRkHjpK5s1PxxBn9qPaWEWxTbpc1PPSD-an83TsXRY",
			"invalid",
			true,
			nil,
		},
		// properly formatted EXPIRED JWT with VALID secret
		// {"name": "test", "exp": 1652097610}
		{
			"eyJhbGciOiJIUzI1NiJ9.eyJuYW1lIjoidGVzdCIsImV4cCI6OTU3ODczMzc0fQ.0oUUKUnsQHs4nZO1pnxQHahKtcHspHu4_AplN2sGC4A",
			"test",
			true,
			nil,
		},
		// properly formatted JWT with VALID claims and VALID secret
		// {"name": "test", "exp": 1898636137}
		{
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoidGVzdCIsImV4cCI6MTg5ODYzNjEzN30.gqRkHjpK5s1PxxBn9qPaWEWxTbpc1PPSD-an83TsXRY",
			"test",
			false,
			jwt.MapClaims{"name": "test", "exp": 1898636137.0},
		},
		// properly formatted JWT with VALID claims (without exp) and VALID secret
		// {"name": "test"}
		{
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoidGVzdCJ9.ml0QsTms3K9wMygTu41ZhKlTyjmW9zHQtoS8FUsCCjU",
			"test",
			false,
			jwt.MapClaims{"name": "test"},
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%s", i, s.token), func(t *testing.T) {
			result, err := security.ParseJWT(s.token, s.secret)

			hasErr := err != nil

			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr %v, got %v (%v)", s.expectError, hasErr, err)
			}

			if len(result) != len(s.expectClaims) {
				t.Fatalf("Expected %v claims got %v", s.expectClaims, result)
			}

			for k, v := range s.expectClaims {
				v2, ok := result[k]
				if !ok {
					t.Fatalf("Missing expected claim %q", k)
				}
				if v != v2 {
					t.Fatalf("Expected %v for %q claim, got %v", v, k, v2)
				}
			}
		})
	}
}

func TestNewJWT(t *testing.T) {
	scenarios := []struct {
		claims      jwt.MapClaims
		key         string
		duration    time.Duration
		expectError bool
	}{
		// empty, zero duration
		{jwt.MapClaims{}, "", 0, true},
		// empty, 10 seconds duration
		{jwt.MapClaims{}, "", 10 * time.Second, false},
		// non-empty, 10 seconds duration
		{jwt.MapClaims{"name": "test"}, "test", 10 * time.Second, false},
	}

	for i, scenario := range scenarios {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			token, tokenErr := security.NewJWT(scenario.claims, scenario.key, scenario.duration)
			if tokenErr != nil {
				t.Fatalf("Expected NewJWT to succeed, got error %v", tokenErr)
			}

			claims, parseErr := security.ParseJWT(token, scenario.key)

			hasParseErr := parseErr != nil
			if hasParseErr != scenario.expectError {
				t.Fatalf("Expected hasParseErr to be %v, got %v (%v)", scenario.expectError, hasParseErr, parseErr)
			}

			if scenario.expectError {
				return
			}

			if _, ok := claims["exp"]; !ok {
				t.Fatalf("Missing required claim exp, got %v", claims)
			}

			// clear exp claim to match with the scenario ones
			delete(claims, "exp")

			if len(claims) != len(scenario.claims) {
				t.Fatalf("Expected %v claims, got %v", scenario.claims, claims)
			}

			for j, k := range claims {
				if claims[j] != scenario.claims[j] {
					t.Fatalf("Expected %v for %q claim, got %v", claims[j], k, scenario.claims[j])
				}
			}
		})
	}
}
