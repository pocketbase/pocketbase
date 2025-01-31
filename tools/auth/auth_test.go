package auth_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/tools/auth"
)

func TestProvidersCount(t *testing.T) {
	expected := 30

	if total := len(auth.Providers); total != expected {
		t.Fatalf("Expected %d providers, got %d", expected, total)
	}
}

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

	// gitea
	p, err = auth.NewProviderByName(auth.NameGitea)
	if err != nil {
		t.Errorf("Expected nil, got error %v", err)
	}
	if _, ok := p.(*auth.Gitea); !ok {
		t.Error("Expected to be instance of *auth.Gitea")
	}

	// oidc
	p, err = auth.NewProviderByName(auth.NameOIDC)
	if err != nil {
		t.Errorf("Expected nil, got error %v", err)
	}
	if _, ok := p.(*auth.OIDC); !ok {
		t.Error("Expected to be instance of *auth.OIDC")
	}

	// oidc2
	p, err = auth.NewProviderByName(auth.NameOIDC + "2")
	if err != nil {
		t.Errorf("Expected nil, got error %v", err)
	}
	if _, ok := p.(*auth.OIDC); !ok {
		t.Error("Expected to be instance of *auth.OIDC")
	}

	// oidc3
	p, err = auth.NewProviderByName(auth.NameOIDC + "3")
	if err != nil {
		t.Errorf("Expected nil, got error %v", err)
	}
	if _, ok := p.(*auth.OIDC); !ok {
		t.Error("Expected to be instance of *auth.OIDC")
	}

	// apple
	p, err = auth.NewProviderByName(auth.NameApple)
	if err != nil {
		t.Errorf("Expected nil, got error %v", err)
	}
	if _, ok := p.(*auth.Apple); !ok {
		t.Error("Expected to be instance of *auth.Apple")
	}

	// instagram
	p, err = auth.NewProviderByName(auth.NameInstagram)
	if err != nil {
		t.Errorf("Expected nil, got error %v", err)
	}
	if _, ok := p.(*auth.Instagram); !ok {
		t.Error("Expected to be instance of *auth.Instagram")
	}

	// vk
	p, err = auth.NewProviderByName(auth.NameVK)
	if err != nil {
		t.Errorf("Expected nil, got error %v", err)
	}
	if _, ok := p.(*auth.VK); !ok {
		t.Error("Expected to be instance of *auth.VK")
	}

	// yandex
	p, err = auth.NewProviderByName(auth.NameYandex)
	if err != nil {
		t.Errorf("Expected nil, got error %v", err)
	}
	if _, ok := p.(*auth.Yandex); !ok {
		t.Error("Expected to be instance of *auth.Yandex")
	}

	// patreon
	p, err = auth.NewProviderByName(auth.NamePatreon)
	if err != nil {
		t.Errorf("Expected nil, got error %v", err)
	}
	if _, ok := p.(*auth.Patreon); !ok {
		t.Error("Expected to be instance of *auth.Patreon")
	}

	// mailcow
	p, err = auth.NewProviderByName(auth.NameMailcow)
	if err != nil {
		t.Errorf("Expected nil, got error %v", err)
	}
	if _, ok := p.(*auth.Mailcow); !ok {
		t.Error("Expected to be instance of *auth.Mailcow")
	}

	// bitbucket
	p, err = auth.NewProviderByName(auth.NameBitbucket)
	if err != nil {
		t.Errorf("Expected nil, got error %v", err)
	}
	if _, ok := p.(*auth.Bitbucket); !ok {
		t.Error("Expected to be instance of *auth.Bitbucket")
	}

	// planningcenter
	p, err = auth.NewProviderByName(auth.NamePlanningcenter)
	if err != nil {
		t.Errorf("Expected nil, got error %v", err)
	}
	if _, ok := p.(*auth.Planningcenter); !ok {
		t.Error("Expected to be instance of *auth.Planningcenter")
	}

	// notion
	p, err = auth.NewProviderByName(auth.NameNotion)
	if err != nil {
		t.Errorf("Expected nil, got error %v", err)
	}
	if _, ok := p.(*auth.Notion); !ok {
		t.Error("Expected to be instance of *auth.Notion")
	}

	// monday
	p, err = auth.NewProviderByName(auth.NameMonday)
	if err != nil {
		t.Errorf("Expected nil, got error %v", err)
	}
	if _, ok := p.(*auth.Monday); !ok {
		t.Error("Expected to be instance of *auth.Monday")
	}

	// wakatime
	p, err = auth.NewProviderByName(auth.NameWakatime)
	if err != nil {
		t.Errorf("Expected nil, got error %v", err)
	}
	if _, ok := p.(*auth.Wakatime); !ok {
		t.Error("Expected to be instance of *auth.Wakatime")
	}

	// linear
	p, err = auth.NewProviderByName(auth.NameLinear)
	if err != nil {
		t.Errorf("Expected nil, got error %v", err)
	}
	if _, ok := p.(*auth.Linear); !ok {
		t.Error("Expected to be instance of *auth.Linear")
	}

	// trakt
	p, err = auth.NewProviderByName(auth.NameTrakt)
	if err != nil {
		t.Errorf("Expected nil, got error %v", err)
	}
	if _, ok := p.(*auth.Trakt); !ok {
		t.Error("Expected to be instance of *auth.Trakt")
	}
}
