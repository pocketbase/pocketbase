package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/vk"
)

func init() {
	Providers[NameVK] = wrapFactory(NewVKProvider)
}

var _ Provider = (*VK)(nil)

// NameVK is the unique name of the VK provider.
const NameVK string = "vk"

// @todo mark as deprecated
//
// VK allows authentication via VK OAuth2.
type VK struct {
	BaseProvider
}

// NewVKProvider creates new VK provider instance with some defaults.
//
// Docs: https://dev.vk.com/api/oauth-parameters
func NewVKProvider() *VK {
	return &VK{BaseProvider{
		ctx:         context.Background(),
		order:       15,
		logo:        `<svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" fill="none"><path fill="#07f" d="M0 23C0 12.2 0 6.7 3.4 3.4 6.7 0 12.2 0 23 0h2c10.8 0 16.3 0 19.6 3.4C48 6.7 48 12.2 48 23v2c0 10.8 0 16.3-3.4 19.6C41.3 48 35.8 48 25 48h-2c-10.8 0-16.3 0-19.6-3.4C0 41.3 0 35.8 0 25z"/><path fill="#fff" d="M25.5 34.6c-10.9 0-17.1-7.5-17.4-20h5.5c.2 9.2 4.2 13 7.4 13.8V14.6h5.2v7.9c3.1-.3 6.4-4 7.6-7.9h5.1a15 15 0 0 1-7 10c2.6 1.2 6.7 4.3 8.2 10h-5.7a10 10 0 0 0-8.2-7.2v7.2z"/></svg>`,
		displayName: "ВКонтакте",
		pkce:        false, // VK currently doesn't support PKCE and throws an error if PKCE params are send
		scopes:      []string{"email"},
		authURL:     vk.Endpoint.AuthURL,
		tokenURL:    vk.Endpoint.TokenURL,
		userInfoURL: "https://api.vk.com/method/users.get?fields=photo_max,screen_name&v=5.131",
	}}
}

// FetchAuthUser returns an AuthUser instance based on VK's user api.
//
// API reference: https://dev.vk.com/method/users.get
func (p *VK) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserInfo(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Response []struct {
			Id        int64  `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Username  string `json:"screen_name"`
			AvatarURL string `json:"photo_max"`
		} `json:"response"`
	}{}

	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	if len(extracted.Response) == 0 {
		return nil, errors.New("missing response entry")
	}

	user := &AuthUser{
		Id:           strconv.FormatInt(extracted.Response[0].Id, 10),
		Name:         strings.TrimSpace(extracted.Response[0].FirstName + " " + extracted.Response[0].LastName),
		Username:     extracted.Response[0].Username,
		AvatarURL:    extracted.Response[0].AvatarURL,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	if email := token.Extra("email"); email != nil {
		user.Email = fmt.Sprint(email)
	}

	return user, nil
}
