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
		order:       22,
		logo:        `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 48 48"><defs><radialGradient id="a" cx="48.5" cy="-.9" r="64.8" fx="48.5" fy="-.9" gradientUnits="userSpaceOnUse"><stop offset="0" stop-color="#9f42c6"/><stop offset=".3" stop-color="#a041c3"/><stop offset=".4" stop-color="#a43ebb"/><stop offset=".5" stop-color="#aa39ad"/><stop offset=".6" stop-color="#b4339a"/><stop offset=".7" stop-color="#c02b81"/><stop offset=".8" stop-color="#cf2061"/><stop offset=".9" stop-color="#e1143c"/><stop offset="1" stop-color="#f50613"/><stop offset="1" stop-color="red"/></radialGradient></defs><path d="M48 11.3v25.4C48 43 43 48 36.7 48H11.3C5 48 0 43 0 36.7V11.3C0 5 5 0 11.3 0h25.4A11 11 0 0 1 48 11.3" style="fill:url(#a)"/><path d="m13.6 18 8 7.9 1.4-1.5-8-7.9zM28 32.4l1.5-1.5-2.2-2.2L47.6 8.4q-.2-1-.8-2.1L24.5 28.7zM13 18.7 11.4 20l14.4 14.4 1.4-1.4-4.3-4.3L46.4 5.4 45 3.7 21.5 27.3zm34.9-9.1L28.7 28.8l1.5 1.4L48 12.4v-1.1zM25.2 22.3l-8-8-1.4 1.5 7.9 8zM41.3 35c0 3.4-2.8 6.2-6.2 6.2H13A6 6 0 0 1 6.8 35V13c0-3.4 2.8-6.2 6.2-6.2h20.8V4.6H12.9a8.3 8.3 0 0 0-8.3 8.3V35c0 4.6 3.7 8.3 8.3 8.3H35c4.6 0 8.3-3.7 8.3-8.3v-3.5h-2z" style="fill:#fff"/></svg>`,
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
