package auth_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/tools/auth"
)

func TestNewProviderByName(t *testing.T) {
	var err error
	var p auth.Provider

	// invalid
	p, err = auth.NewProviderByName("invalid")
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if p != nil {
		t.Errorf("Expected provider to be nil, got %v", p)
	}

	// google
	p, err = auth.NewProviderByName(auth.NameGoogle)
	if err != nil {
		t.Errorf("Expected nil, got error %v", err)
	}
	if _, ok := p.(*auth.Google); !ok {
		t.Error("Expected to be instance of *auth.Google")
	}

	// facebook
	p, err = auth.NewProviderByName(auth.NameFacebook)
	if err != nil {
		t.Errorf("Expected nil, got error %v", err)
	}
	if _, ok := p.(*auth.Facebook); !ok {
		t.Error("Expected to be instance of *auth.Facebook")
	}

	// github
	p, err = auth.NewProviderByName(auth.NameGithub)
	if err != nil {
		t.Errorf("Expected nil, got error %v", err)
	}
	if _, ok := p.(*auth.Github); !ok {
		t.Error("Expected to be instance of *auth.Github")
	}

	// gitlab
	p, err = auth.NewProviderByName(auth.NameGitlab)
	if err != nil {
		t.Errorf("Expected nil, got error %v", err)
	}
	if _, ok := p.(*auth.Gitlab); !ok {
		t.Error("Expected to be instance of *auth.Gitlab")
	}
}
