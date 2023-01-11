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

	// twitter
	p, err = auth.NewProviderByName(auth.NameTwitter)
	if err != nil {
		t.Errorf("Expected nil, got error %v", err)
	}
	if _, ok := p.(*auth.Twitter); !ok {
		t.Error("Expected to be instance of *auth.Twitter")
	}

	// discord
	p, err = auth.NewProviderByName(auth.NameDiscord)
	if err != nil {
		t.Errorf("Expected nil, got error %v", err)
	}
	if _, ok := p.(*auth.Discord); !ok {
		t.Error("Expected to be instance of *auth.Discord")
	}

	// microsoft
	p, err = auth.NewProviderByName(auth.NameMicrosoft)
	if err != nil {
		t.Errorf("Expected nil, got error %v", err)
	}
	if _, ok := p.(*auth.Microsoft); !ok {
		t.Error("Expected to be instance of *auth.Microsoft")
	}

	// spotify
	p, err = auth.NewProviderByName(auth.NameSpotify)
	if err != nil {
		t.Errorf("Expected nil, got error %v", err)
	}
	if _, ok := p.(*auth.Spotify); !ok {
		t.Error("Expected to be instance of *auth.Spotify")
	}

	// kakao
	p, err = auth.NewProviderByName(auth.NameKakao)
	if err != nil {
		t.Errorf("Expected nil, got error %v", err)
	}
	if _, ok := p.(*auth.Kakao); !ok {
		t.Error("Expected to be instance of *auth.Kakao")
	}

	// twitch
	p, err = auth.NewProviderByName(auth.NameTwitch)
	if err != nil {
		t.Errorf("Expected nil, got error %v", err)
	}
	if _, ok := p.(*auth.Twitch); !ok {
		t.Error("Expected to be instance of *auth.Twitch")
	}

	// strava
	p, err = auth.NewProviderByName(auth.NameStrava)
	if err != nil {
		t.Errorf("Expected nil, got error %v", err)
	}
	if _, ok := p.(*auth.Strava); !ok {
		t.Error("Expected to be instance of *auth.Strava")
	}

	// gitee
	p, err = auth.NewProviderByName(auth.NameGitee)
	if err != nil {
		t.Errorf("Expected nil, got error %v", err)
	}
	if _, ok := p.(*auth.Gitee); !ok {
		t.Error("Expected to be instance of *auth.Gitee")
	}

	// livechat
	p, err = auth.NewProviderByName(auth.NameLivechat)
	if err != nil {
		t.Errorf("Expected nil, got error %v", err)
	}
	if _, ok := p.(*auth.Livechat); !ok {
		t.Error("Expected to be instance of *auth.Livechat")
	}
}
