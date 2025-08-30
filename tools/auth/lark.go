package auth

import (
	"context"
	"encoding/json"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

func init() {
	Providers[NameLark] = wrapFactory(NewLarkProvider)
}

var _ Provider = (*Lark)(nil)

// NameLark is the unique name of the Lark provider.
const NameLark string = "lark"

// Lark allows authentication via Lark OAuth2.
type Lark struct {
	BaseProvider
}

// NewLarkProvider creates new Lark provider instance with some defaults.
func NewLarkProvider() *Lark {
	return &Lark{BaseProvider{
		ctx:         context.Background(),
		displayName: "Lark",
		pkce:        true,
		// Lark has two domains with the same API: feishu.cn and larksuite.com.
		// The former is used in China and the latter is used in the other regions.
		// We choose feishu.cn as a default, matching the behavior of Lark's official SDK.
		// Endpoint URLs can be overridden from the frontend if needed.
		// SDK Reference: https://github.com/larksuite/oapi-sdk-go
		authURL:     "https://accounts.feishu.cn/open-apis/authen/v1/authorize",
		tokenURL:    "https://open.feishu.cn/open-apis/authen/v2/oauth/token",
		userInfoURL: "https://open.feishu.cn/open-apis/authen/v1/user_info",
	}}
}

// FetchAuthUser returns an AuthUser instance based the Lark's user api.
//
// API reference: https://open.feishu.cn/document/server-docs/authentication-management/login-state-management/get
func (p *Lark) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
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
			// https://open.feishu.cn/document/platform-overveiw/basic-concepts/user-identity-introduction/introduction#3f2d4b63
			UnionId   string `json:"union_id"`
			Name      string `json:"name"`
			AvatarURL string `json:"avatar_url"`
		} `json:"data"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           extracted.Data.UnionId,
		Name:         extracted.Data.Name,
		AvatarURL:    extracted.Data.AvatarURL,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	return user, nil
}
