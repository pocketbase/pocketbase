package auth

import (
	"fmt"

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
		scopes:     []string{"identify", "email"},
		authUrl:    "https://discord.com/api/oauth2/authorize",
		tokenUrl:   "https://discord.com/api/oauth2/token",
		userApiUrl: "https://discord.com/api/users/@me",
	}}
}

// FetchAuthUser returns an AuthUser instance from Discord's user api.
func (p *Discord) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	// https://discord.com/developers/docs/resources/user#user-object
	rawData := struct {
		Id            string `json:"id"`
		Username      string `json:"username"`
		Discriminator string `json:"discriminator"`
		Email         string `json:"email"`
		Avatar        string `json:"avatar"`
	}{}

	if err := p.FetchRawUserData(token, &rawData); err != nil {
		return nil, err
	}

	// Build a full avatar URL using the avatar hash provided in the API response
	// https://discord.com/developers/docs/reference#image-formatting
	avatarUrl := fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", rawData.Id, rawData.Avatar)

	// Concatenate the user's username and discriminator into a single username string
	username := fmt.Sprintf("%s#%s", rawData.Username, rawData.Discriminator)

	user := &AuthUser{
		Id:        rawData.Id,
		Name:      username,
		Username:  rawData.Username,
		Email:     rawData.Email,
		AvatarUrl: avatarUrl,
	}

	return user, nil
}
