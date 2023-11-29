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

var _ Provider = (*VK)(nil)

// NameVK is the unique name of the VK provider.
const NameVK string = "vk"

// VK allows authentication via VK OAuth2.
type VK struct {
	*baseProvider
}

// NewVKProvider creates new VK provider instance with some defaults.
//
// Docs: https://dev.vk.com/api/oauth-parameters
func NewVKProvider() *VK {
	return &VK{&baseProvider{
		ctx:         context.Background(),
		displayName: "ВКонтакте",
		pkce:        false, // VK currently doesn't support PKCE and throws an error if PKCE params are send
		scopes:      []string{"email"},
		authUrl:     vk.Endpoint.AuthURL,
		tokenUrl:    vk.Endpoint.TokenURL,
		userApiUrl:  "https://api.vk.com/method/users.get?fields=photo_max,screen_name&v=5.131",
	}}
}

// FetchAuthUser returns an AuthUser instance based on VK's user api.
//
// API reference: https://dev.vk.com/method/users.get
func (p *VK) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserData(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Response []struct {
			Id        int    `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Username  string `json:"screen_name"`
			AvatarUrl string `json:"photo_max"`
		} `json:"response"`
	}{}

	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	if len(extracted.Response) == 0 {
		return nil, errors.New("missing response entry")
	}

	user := &AuthUser{
		Id:           strconv.Itoa(extracted.Response[0].Id),
		Name:         strings.TrimSpace(extracted.Response[0].FirstName + " " + extracted.Response[0].LastName),
		Username:     extracted.Response[0].Username,
		AvatarUrl:    extracted.Response[0].AvatarUrl,
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
