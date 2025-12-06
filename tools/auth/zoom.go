package auth

import (
	"context"
	"encoding/json"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

func init() {
	Providers[NameZoom] = wrapFactory(NewZoomProvider)
}

var _ Provider = (*Zoom)(nil)

// NameZoom is the unique name of the Zoom provider.
const NameZoom string = "zoom"

// Zoom allows authentication via Zoom OAuth2.
type Zoom struct {
	BaseProvider
}

// NewZoomProvider creates a new Zoom provider instance with some defaults.
func NewZoomProvider() *Zoom {
	// https://developers.zoom.us/docs/integrations/oauth/
	return &Zoom{BaseProvider{
		ctx:         context.Background(),
		displayName: "Zoom",
		pkce:        true,
		scopes:      []string{"user:read:user"},
		authURL:     "https://zoom.us/oauth/authorize",
		tokenURL:    "https://zoom.us/oauth/token",
		userInfoURL: "https://api.zoom.us/v2/users/me",
	}}
}

// FetchAuthUser returns an AuthUser instance based on Zoom's user api.
//
// API reference: https://developers.zoom.us/docs/api/rest/reference/user/methods/#operation/user
func (p *Zoom) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserInfo(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Id        string `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		PicURL    string `json:"pic_url"`
		Verified  int    `json:"verified"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           extracted.Id,
		Name:         extracted.FirstName + " " + extracted.LastName,
		AvatarURL:    extracted.PicURL,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	if extracted.Verified == 1 {
		user.Email = extracted.Email
	}

	return user, nil
}
