package auth

import (
	"context"
	"encoding/json"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

func init() {
	Providers[NameInstagram] = wrapFactory(NewInstagramProvider)
}

var _ Provider = (*Instagram)(nil)

// NameInstagram is the unique name of the Instagram provider.
const NameInstagram string = "instagram2" // "2" suffix to avoid conflicts with the old deprecated version

// Instagram allows authentication via Instagram Login OAuth2.
type Instagram struct {
	BaseProvider
}

// NewInstagramProvider creates new Instagram provider instance with some defaults.
func NewInstagramProvider() *Instagram {
	return &Instagram{BaseProvider{
		ctx:         context.Background(),
		displayName: "Instagram",
		pkce:        true,
		scopes:      []string{"instagram_business_basic"},
		authURL:     "https://www.instagram.com/oauth/authorize",
		tokenURL:    "https://api.instagram.com/oauth/access_token",
		userInfoURL: "https://graph.instagram.com/me?fields=id,username,account_type,user_id,name,profile_picture_url,followers_count,follows_count,media_count",
	}}
}

// FetchAuthUser returns an AuthUser instance based on the Instagram Login user api response.
//
// API reference: https://developers.facebook.com/docs/instagram-platform/instagram-api-with-instagram-login/get-started#fields
func (p *Instagram) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserInfo(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	// include list of granted permissions to RawUser's payload
	if _, ok := rawUser["permissions"]; !ok {
		if permissions := token.Extra("permissions"); permissions != nil {
			rawUser["permissions"] = permissions
		}
	}

	extracted := struct {
		Id        string `json:"user_id"`
		Name      string `json:"name"`
		Username  string `json:"username"`
		AvatarURL string `json:"profile_picture_url"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           extracted.Id,
		Username:     extracted.Username,
		Name:         extracted.Name,
		AvatarURL:    extracted.AvatarURL,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	return user, nil
}
