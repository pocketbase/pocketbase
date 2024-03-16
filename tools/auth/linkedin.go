package auth

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

var _ Provider = (*LinkedIn)(nil)

// NameLinkedIn is the unique name of the LinkedIn provider.
const NameLinkedIn string = "linkedin"

// LinkedIn allows authentication via LinkedIn OAuth2.
type LinkedIn struct {
	*baseProvider
}

// NewLinkedInProvider creates new LinkedIn provider instance with some defaults.
func NewLinkedInProvider() *LinkedIn {
	return &LinkedIn{&baseProvider{
		ctx:         context.Background(),
		displayName: "LinkedIn",
		pkce:        false,
		scopes: 		 []string{"openid", "profile", "email" },
		authUrl:     "https://www.linkedin.com/oauth/v2/authorization",
		tokenUrl:    "https://www.linkedin.com/oauth/v2/accessToken",
		userApiUrl:  "https://api.linkedin.com/v2/userinfo",
	}}
}

// FetchAuthUser returns an AuthUser instance based on the LinkedIn's user api.
//
// API reference: https://developers.kakao.com/docs/latest/en/kakaologin/rest-api#req-user-info-response
func (p *LinkedIn) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserData(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Id       string `json:"id"`
		Username string `json:"username"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           extracted.Id,
		Username:     extracted.Username,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	return user, nil
}
