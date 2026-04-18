package auth

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

func init() {
	Providers[NameDiscord] = wrapFactory(NewDiscordProvider)
}

var _ Provider = (*Discord)(nil)

// NameDiscord is the unique name of the Discord provider.
const NameDiscord string = "discord"

// Discord allows authentication via Discord OAuth2.
type Discord struct {
	BaseProvider
}

// NewDiscordProvider creates a new Discord provider instance with some defaults.
func NewDiscordProvider() *Discord {
	// https://discord.com/developers/docs/topics/oauth2
	// https://discord.com/developers/docs/resources/user#get-current-user
	return &Discord{BaseProvider{
		ctx:         context.Background(),
		order:       12,
		logo:        `<svg xmlns="http://www.w3.org/2000/svg" width="256" height="199" preserveAspectRatio="xMidYMid"><path fill="#5865f2" d="M216.9 16.6A209 209 0 0 0 164 0c-2.2 4.1-4.9 9.6-6.7 14a194 194 0 0 0-58.6 0C97 9.6 94.2 4.1 92 0a208 208 0 0 0-53 16.6A222 222 0 0 0 1 165a211 211 0 0 0 65 33 161 161 0 0 0 13.8-22.8q-11.5-4.4-21.8-10.6l5.3-4.3a149 149 0 0 0 129.6 0q2.6 2.3 5.3 4.3a136 136 0 0 1-21.9 10.6q6 12 13.9 22.9a211 211 0 0 0 64.8-33.2c5.3-56.3-9-105.1-38-148.4M85.5 135.1c-12.7 0-23-11.8-23-26.2s10.1-26.2 23-26.2 23.2 11.8 23 26.2c0 14.4-10.2 26.2-23 26.2m85 0c-12.6 0-23-11.8-23-26.2s10.2-26.2 23-26.2 23.3 11.8 23 26.2c0 14.4-10.1 26.2-23 26.2"/></svg>`,
		displayName: "Discord",
		pkce:        true,
		scopes:      []string{"identify", "email"},
		authURL:     "https://discord.com/api/oauth2/authorize",
		tokenURL:    "https://discord.com/api/oauth2/token",
		userInfoURL: "https://discord.com/api/users/@me",
	}}
}

// FetchAuthUser returns an AuthUser instance from Discord's user api.
//
// API reference:  https://discord.com/developers/docs/resources/user#user-object
func (p *Discord) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
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
		GlobalName    string `json:"global_name"`
		Username      string `json:"username"`
		Discriminator string `json:"discriminator"`
		Avatar        string `json:"avatar"`
		Email         string `json:"email"`
		Verified      bool   `json:"verified"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	// build a full avatar URL using the avatar hash provided in the API response
	// https://discord.com/developers/docs/reference#image-formatting
	var avatarURL string
	if extracted.Avatar != "" {
		avatarURL = fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", extracted.Id, extracted.Avatar)
	}

	name := extracted.GlobalName
	if name == "" {
		// fallback to username+discriminator
		//
		// Note: Discord migrated to unique usernames without discriminators.
		// Legacy accounts still have a non-zero discriminator (e.g. "1234").
		// See https://support.discord.com/hc/en-us/articles/12620128861463-New-Usernames-Display-Names.
		name = extracted.Username
		if extracted.Discriminator != "" && extracted.Discriminator != "0" {
			name += "#" + extracted.Discriminator
		}
	}

	user := &AuthUser{
		Id:           extracted.Id,
		Name:         name,
		Username:     extracted.Username,
		AvatarURL:    avatarURL,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	if extracted.Verified {
		user.Email = extracted.Email
	}

	return user, nil
}
