package auth

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

var _ Provider = (*Discord)(nil)

// NameDiscord is the unique name of the Discord provider.
const NameDiscord string = "discord"

// Discord allows authentication via Discord OAuth2.
type Discord struct {
	*baseProvider
}

// NewDiscordProvider creates a new Discord provider instance with some defaults.
func NewDiscordProvider() *Discord {
	// https://discord.com/developers/docs/topics/oauth2
	// https://discord.com/developers/docs/resources/user#get-current-user
	return &Discord{&baseProvider{
		ctx:         context.Background(),
		displayName: "Discord",
		pkce:        true,
		scopes:      []string{"identify", "email"},
		authUrl:     "https://discord.com/api/oauth2/authorize",
		tokenUrl:    "https://discord.com/api/oauth2/token",
		userApiUrl:  "https://discord.com/api/users/@me",
	}}
}

// FetchAuthUser returns an AuthUser instance from Discord's user api.
//
// API reference:  https://discord.com/developers/docs/resources/user#user-object
func (p *Discord) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserData(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Id            string `json:"id"`
		Username      string `json:"username"`
		Discriminator string `json:"discriminator"`
		Avatar        string `json:"avatar"`
		Email         string `json:"email"`
		Verified      bool   `json:"verified"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	// Build a full avatar URL using the avatar hash provided in the API response
	// https://discord.com/developers/docs/reference#image-formatting
	avatarUrl := fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", extracted.Id, extracted.Avatar)

	// Concatenate the user's username and discriminator into a single username string
	username := fmt.Sprintf("%s#%s", extracted.Username, extracted.Discriminator)

	user := &AuthUser{
		Id:           extracted.Id,
		Name:         username,
		Username:     extracted.Username,
		AvatarUrl:    avatarUrl,
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
