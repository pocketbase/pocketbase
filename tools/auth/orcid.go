package auth

import (
	"context"
	"encoding/xml"
	"fmt"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

func init() {
	Providers[NameORCID] = wrapFactory(NewORCIDProvider)
}

var _ Provider = (*ORCID)(nil)

// NameORCID is the unique name of the ORCID provider.
const NameORCID string = "ORCID"

// ORCID allows authentication via ORCID OAuth2.
type ORCID struct {
	BaseProvider
}

// NewORCIDProvider creates new ORCID provider instance with some defaults.
func NewORCIDProvider() *ORCID {
	return &ORCID{BaseProvider{
		ctx:         context.Background(),
		displayName: "ORCID",
		pkce:        true,
		scopes: []string{
			"/authenticate",
		},
		authURL:     "https://orcid.org/oauth/authorize",
		tokenURL:    "https://orcid.org/oauth/token",
		userInfoURL: "", // this is set later as it must be derived from the returned token
	}}
}

// FetchAuthUser returns an AuthUser instance based on the ORCID's user api.
//
// API reference: https://info.orcid.org/documentation/integration-guide/
func (p *ORCID) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {

	// deriving userInfoURL from the iD (i.e. username) returned in the token
	iD, ok := token.Extra("orcid").(string)
	if !ok || iD == "" {
		return nil, fmt.Errorf("Failed to get ORCID iD from OAuth2 token")
	}
	p.userInfoURL = `https://pub.orcid.org/v3.0/` + iD + `/person`

	data, err := p.FetchRawUserInfo(token)
	if err != nil {
		return nil, err
	}

	extracted := struct {
		Name struct {
			GivenNames string `xml:"given-names"`
			FamilyName string `xml:"family-name"`
			CreditName string `xml:"credit-name"`
		} `xml:"name"`
		Emails struct {
			Email struct {
				Email string `xml:"email"`
			} `xml:"email"`
		} `xml:"emails"`
	}{}

	if err := xml.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	name := extracted.Name.CreditName
	if name == "" {
		// GivenNames is a required field on ORCID, so it will always be set
		name = extracted.Name.GivenNames
		if extracted.Name.FamilyName != "" {
			name += " " + extracted.Name.FamilyName
		}
	}

	user := &AuthUser{
		Name:         extracted.Name.CreditName,
		Username:     iD,
		Email:        extracted.Emails.Email.Email,
		RawUser:      map[string]any{"rawXML": data},
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Id:           iD,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	return user, nil
}
