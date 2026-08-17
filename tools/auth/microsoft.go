package auth

import (
	"context"
	"encoding/json/v2"
	"errors"
	"slices"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/spf13/cast"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

func init() {
	Providers[NameMicrosoft] = wrapFactory(NewMicrosoftProvider)
}

var _ Provider = (*Microsoft)(nil)

// NameMicrosoft is the unique name of the Microsoft provider.
const NameMicrosoft string = "microsoft"

// extraIdTokenEmailClaim is the name of the extra map entry that
// specifies which email extraction method to use
const extraIdTokenEmailClaim string = "idTokenEmailClaim"

// Microsoft allows authentication via Azure AD/Entra ID OAuth2.
type Microsoft struct {
	BaseProvider
}

// NewMicrosoftProvider creates new Microsoft AD provider instance with some defaults.
func NewMicrosoftProvider() *Microsoft {
	endpoints := microsoft.AzureADEndpoint("")
	return &Microsoft{BaseProvider{
		ctx:         context.Background(),
		order:       3,
		logo:        `<svg xmlns="http://www.w3.org/2000/svg" width="256" height="256" preserveAspectRatio="xMidYMid"><path fill="#f1511b" d="M122 122H0V0h122z"/><path fill="#80cc28" d="M256 122H134V0h122z"/><path fill="#00adef" d="M122 256H0V134h122z"/><path fill="#fbbc09" d="M256 256H134V134h122z"/></svg>`,
		displayName: "Microsoft",
		pkce:        true,
		scopes:      []string{"User.Read"},
		authURL:     endpoints.AuthURL,
		tokenURL:    endpoints.TokenURL,
		userInfoURL: "https://graph.microsoft.com/v1.0/me",
	}}
}

// SetExtra implements Provider.SetExtra() interface method.
//
// If the [extraIdTokenEmailClaim] data is set it will also add "openid"
// to the list of default scopes in order to be able to get an id_token.
func (p *Microsoft) SetExtra(data map[string]any) {
	p.extra = data

	if cast.ToString(p.extra[extraIdTokenEmailClaim]) != "" {
		scopes := p.Scopes()
		if !slices.Contains(scopes, "openid") {
			scopes = append(scopes, "openid")
			p.SetScopes(scopes)
		}
	}
}

// FetchAuthUser returns an AuthUser instance based on the Microsoft's user api.
//
// Graph explorer:  https://developer.microsoft.com/en-us/graph/graph-explorer
// API reference:   https://learn.microsoft.com/en-us/graph/api/user-get
// Optional claims: https://learn.microsoft.com/en-us/entra/identity-platform/optional-claims-reference
func (p *Microsoft) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	// @todo with the future update to v2 endpoint consider skipping the request
	// if id_token is available (we need to make sure that the graph's id is the same as id_token's sub!)
	data, err := p.FetchRawUserInfo(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Id   string `json:"id"`
		Name string `json:"displayName"`
		Mail string `json:"mail"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           extracted.Id,
		Name:         extracted.Name,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	// decide which email to trust and assign
	switch p.extra[extraIdTokenEmailClaim] {
	case "any_verified":
		user.Email = p.extractIdTokenVerifiedPrimaryEmail(token)
		if user.Email == "" {
			user.Email = p.extractIdTokenVerifiedXmsEdovEmail(token)
		}
	case "verified_primary_email":
		user.Email = p.extractIdTokenVerifiedPrimaryEmail(token)
	case "email_and_xms_edov":
		user.Email = p.extractIdTokenVerifiedXmsEdovEmail(token)
	case "email":
		user.Email = p.extractIdTokenEmail(token)
	default:
		// This is kept to avoid introducing breaking changes and generally
		// it is considered safe because the provider was originally created
		// for single-tenants apps. Furthermore the value is expected to be
		// synced with the id_token's `email` claim which since 2023
		// by *default* would be empty if it is unverified.
		user.Email = extracted.Mail
	}

	return user, nil
}

func (p *Microsoft) extractIdTokenClaims(trustedIdToken *oauth2.Token) (jwt.MapClaims, error) {
	idToken, _ := trustedIdToken.Extra("id_token").(string)
	if idToken == "" {
		return nil, errors.New("empty id_token")
	}

	claims := jwt.MapClaims{}
	_, _, err := jwt.NewParser().ParseUnverified(idToken, claims)
	if err != nil {
		return nil, err
	}

	// loosely validate common claims
	// (we don't check the signature because the token is expected to be
	// from a trusted source, usually from a direct TLS communication with the provider)
	jwtValidator := jwt.NewValidator(
		jwt.WithIssuedAt(),
		jwt.WithLeeway(idTokenLeeway),
	)
	err = jwtValidator.Validate(claims)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func (p *Microsoft) extractIdTokenVerifiedPrimaryEmail(trustedIdToken *oauth2.Token) string {
	claims, _ := p.extractIdTokenClaims(trustedIdToken)

	email, _ := claims["verified_primary_email"].(string)

	return email
}

func (p *Microsoft) extractIdTokenVerifiedXmsEdovEmail(trustedIdToken *oauth2.Token) string {
	claims, _ := p.extractIdTokenClaims(trustedIdToken)

	if !cast.ToBool(claims["xms_edov"]) {
		return ""
	}

	email, _ := claims["email"].(string)

	return email
}

func (p *Microsoft) extractIdTokenEmail(trustedIdToken *oauth2.Token) string {
	claims, _ := p.extractIdTokenClaims(trustedIdToken)

	email, _ := claims["email"].(string)

	return email
}
