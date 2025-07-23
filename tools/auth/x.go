package auth

import (
	"context"
	"encoding/json"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

type X struct {
	BaseProvider
}

const NameX string = "X"

func init() {
	Providers[NameX] = wrapFactory(NewXProvider)
}

func NewXProvider() *X {
	return &X{BaseProvider{
		ctx:         context.Background(),
		displayName: "X",
		pkce:        true,
		scopes: []string{
			"users.read",
			"tweet.read",
		},
		authURL:     "https://x.com/i/oauth2/authorize",
		tokenURL:    "https://api.x.com/2/oauth2/token",
		userInfoURL: "https://api.x.com/2/users/me?user.fields=id,name,username,profile_image_url",
	}}
}

func (p *X) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserInfo(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Data struct {
			Id              string `json:"id"`
			Name            string `json:"name"`
			Username        string `json:"username"`
			ProfileImageURL string `json:"profile_image_url"`
		} `json:"data"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           extracted.Data.Id,
		Name:         extracted.Data.Name,
		Username:     extracted.Data.Username,
		AvatarURL:    extracted.Data.ProfileImageURL,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	return user, nil
}
