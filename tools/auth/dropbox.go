package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

func init() {
	Providers[NameDropbox] = wrapFactory(NewDropboxProvider)
}

var _ Provider = (*Dropbox)(nil)

// NameDropbox is the unique name of the Dropbox provider.
const NameDropbox string = "dropbox"

// Dropbox allows authentication via Dropbox OAuth2.
type Dropbox struct {
	BaseProvider
}

// NewDropboxProvider creates a new Dropbox provider instance with some defaults.
func NewDropboxProvider() *Dropbox {
	// https://developers.dropbox.com/oauth-guide
	// https://www.dropbox.com/developers/documentation/http/documentation#users-get_current_account
	return &Dropbox{BaseProvider{
		ctx:         context.Background(),
		displayName: "Dropbox",
		pkce:        true,
		scopes:      []string{"account_info.read"},
		authURL:     "https://www.dropbox.com/oauth2/authorize",
		tokenURL:    "https://api.dropboxapi.com/oauth2/token",
		userInfoURL: "https://api.dropboxapi.com/2/users/get_current_account",
	}}
}

// FetchAuthUser returns an AuthUser instance based on Dropbox's user api.
//
// API reference: https://www.dropbox.com/developers/documentation/http/documentation#users-get_current_account
func (p *Dropbox) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserInfo(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		AccountId     string `json:"account_id"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Name          struct {
			DisplayName string `json:"display_name"`
		} `json:"name"`
		ProfilePhotoURL string `json:"profile_photo_url"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           extracted.AccountId,
		Name:         extracted.Name.DisplayName,
		AvatarURL:    extracted.ProfilePhotoURL,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	if extracted.EmailVerified {
		user.Email = extracted.Email
	}

	return user, nil
}

// FetchRawUserInfo implements Provider.FetchRawUserInfo interface method.
//
// Dropbox's get_current_account endpoint requires a POST request with "null" body.
func (p *Dropbox) FetchRawUserInfo(token *oauth2.Token) ([]byte, error) {
	req, err := http.NewRequestWithContext(p.ctx, "POST", p.userInfoURL, strings.NewReader("null"))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return p.sendRawUserInfoRequest(req, token)
}
