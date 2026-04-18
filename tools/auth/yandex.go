package auth

import (
	"context"
	"encoding/json"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/yandex"
)

func init() {
	Providers[NameYandex] = wrapFactory(NewYandexProvider)
}

var _ Provider = (*Yandex)(nil)

// NameYandex is the unique name of the Yandex provider.
const NameYandex string = "yandex"

// Yandex allows authentication via Yandex OAuth2.
type Yandex struct {
	BaseProvider
}

// NewYandexProvider creates new Yandex provider instance with some defaults.
//
// Docs: https://yandex.ru/dev/id/doc/en/
func NewYandexProvider() *Yandex {
	return &Yandex{BaseProvider{
		ctx:         context.Background(),
		order:       4,
		logo:        `<svg xmlns="http://www.w3.org/2000/svg" width="44" height="44" fill="none"><rect width="44" height="44" fill="#fc3f1d" rx="22"/><path fill="#fff" d="M25.2 12.3H23c-3.8 0-5.7 2-5.7 4.8 0 3.2 1.3 4.8 4.1 6.7l2.3 1.6-6.4 9.8h-5.1l6-8.9c-3.5-2.5-5.4-4.8-5.4-8.9 0-5 3.5-8.6 10.2-8.6h6.7v26.4h-4.5z"/></svg>`,
		displayName: "Yandex",
		pkce:        true,
		scopes:      []string{"login:email", "login:avatar", "login:info"},
		authURL:     yandex.Endpoint.AuthURL,
		tokenURL:    yandex.Endpoint.TokenURL,
		userInfoURL: "https://login.yandex.ru/info",
	}}
}

// FetchAuthUser returns an AuthUser instance based on Yandex's user api.
//
// API reference: https://yandex.ru/dev/id/doc/en/user-information#response-format
func (p *Yandex) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserInfo(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Id            string `json:"id"`
		Name          string `json:"real_name"`
		Username      string `json:"login"`
		Email         string `json:"default_email"`
		IsAvatarEmpty bool   `json:"is_avatar_empty"`
		AvatarId      string `json:"default_avatar_id"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           extracted.Id,
		Name:         extracted.Name,
		Username:     extracted.Username,
		Email:        extracted.Email,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	if !extracted.IsAvatarEmpty {
		user.AvatarURL = "https://avatars.yandex.net/get-yapic/" + extracted.AvatarId + "/islands-200"
	}

	return user, nil
}
