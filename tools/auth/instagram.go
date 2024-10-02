package auth

import (
	"context"
	"encoding/json"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

var _ Provider = (*Instagram)(nil)

// NameInstagram is the unique name of the Instagram provider.
const NameInstagram string = "instagram"

// Instagram allows authentication via Instagram OAuth2.
type Instagram struct {
	*baseProvider
}

// NewInstagramProvider creates new Instagram provider instance with some defaults.
func NewInstagramProvider() *Instagram {
	return &Instagram{&baseProvider{
		ctx:         context.Background(),
		displayName: "Instagram",
		pkce:        true,
		scopes:      []string{"instagram_business_basic"},
		authUrl:     "https://www.instagram.com/oauth/authorize",
		tokenUrl:    "https://api.instagram.com/oauth/access_token",
		userApiUrl:  "https://graph.instagram.com/me?fields=id,username,account_type,user_id,name,profile_picture_url,followers_count,follows_count,media_count",
	}}
}

// FetchAuthUser returns an AuthUser instance based on the Instagram's login fields.
//
// API reference: https://developers.facebook.com/docs/instagram-platform/instagram-api-with-instagram-login/get-started#fields
func (p *Instagram) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserData(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	// @note the extracted "id" is a app scoped id, to get the actual IG ID use the RawUser's map key "user_id"
	extracted := struct {
		Id       string `json:"id"`
		Username string `json:"username"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           extracted.Id,
		Username:     extracted.Username,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	return user, nil
}
