package auth

import (
	"context"
	"encoding/json"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

func init() {
	Providers[NameWakatime] = wrapFactory(NewWakatimeProvider)
}

var _ Provider = (*Wakatime)(nil)

// NameWakatime is the unique name of the Wakatime provider.
const NameWakatime = "wakatime"

// Wakatime is an auth provider for Wakatime.
type Wakatime struct {
	BaseProvider
}

// NewWakatimeProvider creates a new Wakatime provider instance with some defaults.
func NewWakatimeProvider() *Wakatime {
	return &Wakatime{BaseProvider{
		ctx:         context.Background(),
		displayName: "WakaTime",
		pkce:        true,
		scopes:      []string{"email"},
		authURL:     "https://wakatime.com/oauth/authorize",
		tokenURL:    "https://wakatime.com/oauth/token",
		userInfoURL: "https://wakatime.com/api/v1/users/current",
	}}
}

// FetchAuthUser returns an AuthUser instance based on the Wakatime's user API.
//
// API reference: https://wakatime.com/developers#users
func (p *Wakatime) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
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
			Id               string `json:"id"`
			DisplayName      string `json:"display_name"`
			Username         string `json:"username"`
			Email            string `json:"email"`
			Photo            string `json:"photo"`
			IsPhotoPublic    bool   `json:"photo_public"`
			IsEmailConfirmed bool   `json:"is_email_confirmed"`
		} `json:"data"`
	}{}

	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           extracted.Data.Id,
		Name:         extracted.Data.DisplayName,
		Username:     extracted.Data.Username,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	// note: we don't check for is_email_public field because PocketBase
	// has its own emailVisibility flag which is false by default
	if extracted.Data.IsEmailConfirmed {
		user.Email = extracted.Data.Email
	}

	if extracted.Data.IsPhotoPublic {
		user.AvatarURL = extracted.Data.Photo
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	return user, nil
}
