package auth

import (
	"context"
	"encoding/json"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

func init() {
	Providers[NameBox] = wrapFactory(NewBoxProvider)
}

var _ Provider = (*Box)(nil)

// NameBox is the unique name of the Box provider.
const NameBox = "box"

// Box is an auth provider for Box.
type Box struct {
	BaseProvider
}

// NewBoxProvider creates a new Box provider instance with some defaults.
func NewBoxProvider() *Box {
	return &Box{BaseProvider{
		ctx:         context.Background(),
		displayName: "Box",
		pkce:        true,
		scopes:      []string{"root_readwrite"},
		authURL:     "https://account.box.com/api/oauth2/authorize",
		tokenURL:    "https://api.box.com/oauth2/token",
		userInfoURL: "https://api.box.com/2.0/users/me",
	}}
}

// FetchAuthUser returns an AuthUser instance based on the Box's user API.
//
// API reference: https://developer.box.com/reference/get-users-me/
func (p *Box) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserInfo(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		ID                            string `json:"id"`
		Name                          string `json:"name"`
		Login                         string `json:"login"`
		AvatarURL                     string `json:"avatar_url"`
		Language                      string `json:"language"`
		Timezone                      string `json:"timezone"`
		SpaceAmount                   int64  `json:"space_amount"`
		SpaceUsed                     int64  `json:"space_used"`
		MaxUploadSize                 int64  `json:"max_upload_size"`
		Status                        string `json:"status"`
		JobTitle                      string `json:"job_title"`
		Phone                         string `json:"phone"`
		Address                       string `json:"address"`
		Role                          string `json:"role"`
		TrackingCodes                 []any  `json:"tracking_codes"`
		CanSeeManagedUsers            bool   `json:"can_see_managed_users"`
		IsSyncEnabled                 bool   `json:"is_sync_enabled"`
		IsExternalCollabRestricted    bool   `json:"is_external_collab_restricted"`
		IsExemptFromDeviceLimits      bool   `json:"is_exempt_from_device_limits"`
		IsExemptFromLoginVerification bool   `json:"is_exempt_from_login_verification"`
		Enterprise                    struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"enterprise"`
		MyTags               []string `json:"my_tags"`
		Hostname             string   `json:"hostname"`
		IsPlatformAccessOnly bool     `json:"is_platform_access_only"`
		ExternalAppUserId    string   `json:"external_app_user_id"`
	}{}

	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           extracted.ID,
		Name:         extracted.Name,
		Username:     extracted.Login,
		Email:        extracted.Login, // Box uses login as email identifier
		AvatarURL:    extracted.AvatarURL,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	return user, nil
}
