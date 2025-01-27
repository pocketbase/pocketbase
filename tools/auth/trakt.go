package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

func init() {
	Providers[NameTrakt] = wrapFactory(NewTraktProvider)
}

var _ Provider = (*Trakt)(nil)

// NameTrakt is the unique name of the Trakt provider.
const NameTrakt string = "trakt"

// Trakt allows authentication via Trakt OAuth2.
type Trakt struct {
	BaseProvider
}

// NewTraktProvider creates new Trakt provider instance with some defaults.
func NewTraktProvider() *Trakt {
	return &Trakt{BaseProvider{
		ctx:         context.Background(),
		displayName: "Trakt",
		pkce:        true,
		authURL:     "https://trakt.tv/oauth/authorize",
		tokenURL:    "https://api.trakt.tv/oauth/token",
		userInfoURL: "https://api.trakt.tv/users/settings",
	}}
}

// FetchAuthUser returns an AuthUser instance based on Trakt's user settings API.
// API reference: https://trakt.docs.apiary.io/#reference/users/settings/retrieve-settings
func (p *Trakt) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserInfo(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		User struct {
			Username string `json:"username"`
			Name     string `json:"name"`
			Ids      struct {
				Slug string `json:"slug"`
				UUID string `json:"uuid"`
			} `json:"ids"`
			Images struct {
				Avatar struct {
					Full string `json:"full"`
				} `json:"avatar"`
			} `json:"images"`
		} `json:"user"`
	}{}

	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           extracted.User.Ids.UUID,
		Username:     extracted.User.Username,
		Name:         extracted.User.Name,
		AvatarURL:    extracted.User.Images.Avatar.Full,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	return user, nil
}

// FetchRawUserInfo implements Provider.FetchRawUserInfo interface method.
//
// This differ from BaseProvider because Trakt requires a number of
// mandatory headers for all requests
// (https://trakt.docs.apiary.io/#introduction/required-headers).
func (p *Trakt) FetchRawUserInfo(token *oauth2.Token) ([]byte, error) {
	req, err := http.NewRequestWithContext(p.ctx, "GET", p.userInfoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-type", "application/json")
	req.Header.Set("trakt-api-key", p.clientId)
	req.Header.Set("trakt-api-version", "2")

	return p.sendRawUserInfoRequest(req, token)
}
