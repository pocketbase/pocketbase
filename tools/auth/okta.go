package auth

import (
	"encoding/json"
	"fmt"

	"golang.org/x/oauth2"
)

var _ Provider = (*Okta)(nil)

// NameOkta is the unique name of the Okta provider.
const NameOkta string = "okta"

// Okta allows authentication via Okta OAuth2.
type Okta struct {
	*baseProvider
}

// NewOktaProvider creates new Okta provider instance with some defaults.
func NewOktaProvider() *Okta {
	return &Okta{&baseProvider{
		scopes: []string{
			"openid", // minimal requirement to return the id
			"email",
			"profile",
		},
	}}
}

// FetchAuthUser returns an AuthUser instance based on https://developer.okta.com/docs/reference/api/oidc/#userinfo.
func (p *Okta) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserData(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}
	fmt.Println(rawUser)
	extracted := struct {
		Id    string `json:"sub"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}
	fmt.Println(extracted)
	user := &AuthUser{
		Id:           extracted.Id,
		Name:         extracted.Name,
		Email:        extracted.Email,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
	fmt.Println(user)
	return user, nil
}
