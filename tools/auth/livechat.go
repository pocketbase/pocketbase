package auth

import (
	"encoding/json"

	"golang.org/x/oauth2"
)

var _ Provider = (*Livechat)(nil)

const NameLivechat = "livechat"

type Livechat struct {
	*baseProvider
}

func NewLivechatProvider() *Livechat {
	return &Livechat{
		&baseProvider{
			authUrl:    "https://accounts.livechat.com/",
			tokenUrl:   "https://accounts.livechat.com/token",
			userApiUrl: "https://accounts.livechat.com/v2/accounts/me",
		},
	}
}

func (p *Livechat) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserData(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Id        string `json:"account_id"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarUrl string `json:"avatar_url"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           extracted.Id,
		Name:         extracted.Name,
		Username:     extracted.Email,
		Email:        extracted.Email,
		AvatarUrl:    extracted.AvatarUrl,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	return user, nil
}
