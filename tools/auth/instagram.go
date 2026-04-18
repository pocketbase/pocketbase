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
		order:       6,
		logo:        `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 132 132"><defs><radialGradient xlink:href="#a" id="c" cx="158.4" cy="578.1" r="65" fx="158.4" fy="578.1" gradientTransform="matrix(0 -1.98198 1.8439 0 -1031.4 454)" gradientUnits="userSpaceOnUse"/><radialGradient xlink:href="#b" id="d" cx="147.7" cy="473.5" r="65" fx="147.7" fy="473.5" gradientTransform="rotate(78.7 1103.9 776.3)scale(.88596 3.6529)" gradientUnits="userSpaceOnUse"/><linearGradient id="b"><stop offset="0" stop-color="#3771c8"/><stop offset=".1" stop-color="#3771c8"/><stop offset="1" stop-color="#60f" stop-opacity="0"/></linearGradient><linearGradient id="a"><stop offset="0" stop-color="#fd5"/><stop offset=".1" stop-color="#fd5"/><stop offset=".5" stop-color="#ff543e"/><stop offset="1" stop-color="#c837ab"/></linearGradient></defs><path fill="url(#c)" d="M65 0C38 0 30 0 28.4.2c-5.6.4-9 1.3-12.8 3.2a28 28 0 0 0-15 21.3c-.4 3-.6 3.6-.6 19.1V65c0 27 0 35 .2 36.6.4 5.4 1.3 8.8 3 12.5 3.5 7.2 10 12.5 17.8 14.5q4 1 9.5 1.3a2913 2913 0 0 0 68.8 0c4.4-.2 7-.6 9.8-1.3 7.8-2 14.2-7.3 17.7-14.5 1.8-3.7 2.7-7.2 3.1-12.3a2759 2759 0 0 0 0-73.6c-.4-5.2-1.3-8.8-3.1-12.5q-2.1-4.3-5.6-7.6A28 28 0 0 0 105.4.6c-3-.4-3.7-.6-19.2-.6z" transform="translate(1 1)"/><path fill="url(#d)" d="M65 0C38 0 30 0 28.4.2c-5.6.4-9 1.3-12.8 3.2a28 28 0 0 0-15 21.3c-.4 3-.6 3.6-.6 19.1V65c0 27 0 35 .2 36.6.4 5.4 1.3 8.8 3 12.5 3.5 7.2 10 12.5 17.8 14.5q4 1 9.5 1.3a2913 2913 0 0 0 68.8 0c4.4-.2 7-.6 9.8-1.3 7.8-2 14.2-7.3 17.7-14.5 1.8-3.7 2.7-7.2 3.1-12.3a2759 2759 0 0 0 0-73.6c-.4-5.2-1.3-8.8-3.1-12.5q-2.1-4.3-5.6-7.6A28 28 0 0 0 105.4.6c-3-.4-3.7-.6-19.2-.6z" transform="translate(1 1)"/><path fill="#fff" d="M66 18c-13 0-14.7 0-19.8.3s-8.6 1-11.6 2.2a24 24 0 0 0-8.5 5.6 24 24 0 0 0-5.6 8.5 35 35 0 0 0-2.2 11.6C18 51.3 18 53 18 66s0 14.7.3 19.8 1 8.6 2.2 11.6q1.7 4.7 5.6 8.5 3.8 4 8.5 5.6c3 1.2 6.5 2 11.6 2.2s6.8.3 19.8.3 14.7 0 19.8-.3 8.6-1 11.6-2.2q4.7-1.7 8.5-5.6t5.6-8.5c1.2-3 2-6.5 2.2-11.6s.3-6.8.3-19.8 0-14.7-.3-19.8-1-8.6-2.2-11.6a24 24 0 0 0-5.6-8.5 24 24 0 0 0-8.5-5.6c-3-1.2-6.5-2-11.6-2.2S79 18 66 18m-4.3 8.7H66c12.8 0 14.3 0 19.4.2 4.7.2 7.2 1 9 1.7 2.2.8 3.7 1.9 5.4 3.6q2.4 2.2 3.6 5.5c.7 1.7 1.5 4.2 1.7 8.9.2 5 .3 6.6.3 19.4s-.1 14.3-.3 19.4c-.2 4.7-1 7.2-1.7 8.9q-1.1 3.1-3.6 5.5a15 15 0 0 1-5.5 3.6c-1.7.7-4.2 1.4-8.9 1.6-5 .3-6.6.3-19.4.3a334 334 0 0 1-19.4-.3 28 28 0 0 1-8.9-1.6q-3.1-1.3-5.5-3.6a15 15 0 0 1-3.6-5.5 25 25 0 0 1-1.7-9c-.2-5-.2-6.5-.2-19.3s0-14.4.2-19.4a26 26 0 0 1 1.7-9q1.2-3.1 3.6-5.4 2.4-2.5 5.5-3.6a25 25 0 0 1 9-1.7c4.3-.2 6-.3 15-.3zm30 8a5.8 5.8 0 1 0 0 11.4 5.8 5.8 0 0 0 0-11.5M66 41.2a24.7 24.7 0 1 0 0 49.4 24.7 24.7 0 0 0 0-49.3m0 8.7a16 16 0 1 1 0 32 16 16 0 0 1 0-32"/></svg>`,
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
