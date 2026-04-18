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
		order:       26,
		logo:        `<svg xmlns="http://www.w3.org/2000/svg" width="340" height="340" fill="none"><path stroke="#000" stroke-width="40" d="M170 20a150 150 0 1 0 0 300 150 150 0 0 0 0-300Z" clip-rule="evenodd"/><path fill="#000" stroke="#000" stroke-width="10" d="M190.2 213.5a8 8 0 0 1-7.6 3l-2-.8a9 9 0 0 1-2.3-2l-.9-1.3-8.8-14.2-8.8 14.2a8 8 0 0 1-6.8 4.3 8 8 0 0 1-6.8-4.4l-38.6-56.2a9 9 0 0 1-2-5.8c0-4.8 3.4-8.6 7.7-8.6 2.8 0 5.3 1.6 6.6 4l32.7 48.2 9.1-15a8 8 0 0 1 6.9-4.4q4.2.2 6.4 3.8l9.5 15.5 51.2-73.2q2.3-3.8 6.5-4c4.3 0 7.8 4 7.8 8.7q0 3.2-1.8 5.4z"/></svg>`,
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
