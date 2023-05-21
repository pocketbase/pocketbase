package auth

import (
	"encoding/json"
	"fmt"
	"strconv"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/vk"
)

var _ Provider = (*Vk)(nil)

// NameVk is the unique name of the Vk provider.
const NameVk string = "vk"

// Vk allows authentication via Vk OAuth2.
type Vk struct {
	*baseProvider
}

// NewVkProvider creates new Vk provider instance with some defaults.
//
// Docs: https://dev.vk.com/api/oauth-parameters
func NewVkProvider() *Vk {
	return &Vk{&baseProvider{
		scopes:     []string{"email"},
		authUrl:    vk.Endpoint.AuthURL,
		tokenUrl:   vk.Endpoint.TokenURL,
		userApiUrl: "https://api.vk.com/method/users.get?fields=photo_max,screen_name&v=5.131",
	}}
}

// FetchAuthUser returns an AuthUser instance based on Vk's user api.
//
// API reference: https://dev.vk.com/method/users.get
func (p *Vk) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
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
		return nil, err
	}

	user := &AuthUser{
		Id:           strconv.Itoa(extracted.Response[0].Id),
		Name:         extracted.Response[0].FirstName + " " + extracted.Response[0].LastName,
		Username:     extracted.Response[0].Username,
		AvatarUrl:    extracted.Response[0].AvatarUrl,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	email := token.Extra("email")
	if email != nil {
		user.Email = fmt.Sprint(email)
	}

	return user, nil
}
