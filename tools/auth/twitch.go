package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/twitch"
)

func init() {
	Providers[NameTwitch] = wrapFactory(NewTwitchProvider)
}

var _ Provider = (*Twitch)(nil)

// NameTwitch is the unique name of the Twitch provider.
const NameTwitch string = "twitch"

// Twitch allows authentication via Twitch OAuth2.
type Twitch struct {
	BaseProvider
}

// NewTwitchProvider creates new Twitch provider instance with some defaults.
func NewTwitchProvider() *Twitch {
	return &Twitch{BaseProvider{
		ctx:         context.Background(),
		order:       23,
		logo:        `<svg xmlns="http://www.w3.org/2000/svg" width="256" height="268" preserveAspectRatio="xMidYMid"><path fill="#5a3e85" d="M17 0 0 47v186h64v35h35l35-35h52l70-70V0zm24 23h192v128l-41 41h-64l-35 35v-35H41zm64 117h23V70h-23zm64 0h23V70h-23z"/></svg>`,
		displayName: "Twitch",
		pkce:        true,
		scopes:      []string{"user:read:email"},
		authURL:     twitch.Endpoint.AuthURL,
		tokenURL:    twitch.Endpoint.TokenURL,
		userInfoURL: "https://api.twitch.tv/helix/users",
	}}
}

// FetchAuthUser returns an AuthUser instance based the Twitch's user api.
//
// API reference: https://dev.twitch.tv/docs/api/reference#get-users
func (p *Twitch) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserInfo(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Data []struct {
			Id              string `json:"id"`
			Login           string `json:"login"`
			DisplayName     string `json:"display_name"`
			Email           string `json:"email"`
			ProfileImageURL string `json:"profile_image_url"`
		} `json:"data"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	if len(extracted.Data) == 0 {
		return nil, errors.New("failed to fetch AuthUser data")
	}

	user := &AuthUser{
		Id:           extracted.Data[0].Id,
		Name:         extracted.Data[0].DisplayName,
		Username:     extracted.Data[0].Login,
		Email:        extracted.Data[0].Email,
		AvatarURL:    extracted.Data[0].ProfileImageURL,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	return user, nil
}

// FetchRawUserInfo implements Provider.FetchRawUserInfo interface method.
//
// This differ from BaseProvider because Twitch requires the Client-Id header.
func (p *Twitch) FetchRawUserInfo(token *oauth2.Token) ([]byte, error) {
	req, err := http.NewRequestWithContext(p.ctx, "GET", p.userInfoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Client-Id", p.clientId)

	return p.sendRawUserInfoRequest(req, token)
}
