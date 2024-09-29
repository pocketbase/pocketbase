package mailer

import (
	"net/smtp"
	"testing"
)

func TestLoginAuthStart(t *testing.T) {
	auth := smtpLoginAuth{username: "test", password: "123456"}

	scenarios := []struct {
		name        string
		serverInfo  *smtp.ServerInfo
		expectError bool
	}{
		{
			"localhost without tls",
			&smtp.ServerInfo{TLS: false, Name: "localhost"},
			false,
		},
		{
			"localhost with tls",
			&smtp.ServerInfo{TLS: true, Name: "localhost"},
			false,
		},
		{
			"127.0.0.1 without tls",
			&smtp.ServerInfo{TLS: false, Name: "127.0.0.1"},
			false,
		},
		{
			"127.0.0.1 with tls",
			&smtp.ServerInfo{TLS: false, Name: "127.0.0.1"},
			false,
		},
		{
			"::1 without tls",
			&smtp.ServerInfo{TLS: false, Name: "::1"},
			false,
		},
		{
			"::1 with tls",
			&smtp.ServerInfo{TLS: false, Name: "::1"},
			false,
		},
		{
			"non-localhost without tls",
			&smtp.ServerInfo{TLS: false, Name: "example.com"},
			true,
		},
		{
			"non-localhost with tls",
			&smtp.ServerInfo{TLS: true, Name: "example.com"},
			false,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			method, resp, err := auth.Start(s.serverInfo)

			hasErr := err != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr %v, got %v", s.expectError, hasErr)
			}

			if hasErr {
				return
			}

			if len(resp) != 0 {
				t.Fatalf("Expected empty data response, got %v", resp)
			}

			if method != "LOGIN" {
				t.Fatalf("Expected LOGIN, got %v", method)
			}
		})
	}
}

func TestLoginAuthNext(t *testing.T) {
	auth := smtpLoginAuth{username: "test", password: "123456"}

	{
		// example|false
		r1, err := auth.Next([]byte("example:"), false)
		if err != nil {
			t.Fatalf("[example|false] Unexpected error %v", err)
		}
		if len(r1) != 0 {
			t.Fatalf("[example|false] Expected empty part, got %v", r1)
		}

		// example|true
		r2, err := auth.Next([]byte("example:"), true)
		if err != nil {
			t.Fatalf("[example|true] Unexpected error %v", err)
		}
		if len(r2) != 0 {
			t.Fatalf("[example|true] Expected empty part, got %v", r2)
		}
	}

	// ---------------------------------------------------------------

	{
		// username:|false
		r1, err := auth.Next([]byte("username:"), false)
		if err != nil {
			t.Fatalf("[username|false] Unexpected error %v", err)
		}
		if len(r1) != 0 {
			t.Fatalf("[username|false] Expected empty part, got %v", r1)
		}

		// username:|true
		r2, err := auth.Next([]byte("username:"), true)
		if err != nil {
			t.Fatalf("[username|true] Unexpected error %v", err)
		}
		if str := string(r2); str != auth.username {
			t.Fatalf("[username|true] Expected %s, got %s", auth.username, str)
		}

		// uSeRnAmE:|true
		r3, err := auth.Next([]byte("uSeRnAmE:"), true)
		if err != nil {
			t.Fatalf("[uSeRnAmE|true] Unexpected error %v", err)
		}
		if str := string(r3); str != auth.username {
			t.Fatalf("[uSeRnAmE|true] Expected %s, got %s", auth.username, str)
		}
	}

	// ---------------------------------------------------------------

	{
		// password:|false
		r1, err := auth.Next([]byte("password:"), false)
		if err != nil {
			t.Fatalf("[password|false] Unexpected error %v", err)
		}
		if len(r1) != 0 {
			t.Fatalf("[password|false] Expected empty part, got %v", r1)
		}

		// password:|true
		r2, err := auth.Next([]byte("password:"), true)
		if err != nil {
			t.Fatalf("[password|true] Unexpected error %v", err)
		}
		if str := string(r2); str != auth.password {
			t.Fatalf("[password|true] Expected %s, got %s", auth.password, str)
		}

		// pAsSwOrD:|true
		r3, err := auth.Next([]byte("pAsSwOrD:"), true)
		if err != nil {
			t.Fatalf("[pAsSwOrD|true] Unexpected error %v", err)
		}
		if str := string(r3); str != auth.password {
			t.Fatalf("[pAsSwOrD|true] Expected %s, got %s", auth.password, str)
		}
	}
}
