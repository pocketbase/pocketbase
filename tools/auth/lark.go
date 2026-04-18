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
		order:       19,
		logo:        `<svg xmlns="http://www.w3.org/2000/svg" width="256" height="256" fill="none"><path fill="#00d6b9" d="M137 126v-1l2-1v-1l3-2 3-3 3-3 2-3 3-2 2-3 4-3 2-2 4-3 9-7 5-3 9-3 2-1a135 135 0 0 0-26-51 12 12 0 0 0-9-4H56l-2 1 1 2a288 288 0 0 1 82 93"/><path fill="#3370ff" d="M98 212c51 0 95-28 118-69l3-5-6 9-4 4-6 5-2 2-4 3-11 4-6 2h-15l-4-1h-3l-2-1-5-1-3-1-4-1-3-1-3-1-1-1-3-1h-2l-3-2h-2l-2-1-3-1-2-1-2-1-2-1h-2l-1-1-2-1h-1l-1-1-2-1-2-1-2-1-2-1a285 285 0 0 1-81-60l-3 1v94a12 12 0 0 0 5 11 135 135 0 0 0 76 22"/><path fill="#133c9a" d="M248 90a78 78 0 0 0-58-5l-2 1-14 6-6 4-7 6-2 2-4 3-2 3-3 2-2 3-3 3-3 3-3 2v1l-2 1-1 1-1 1a137 137 0 0 1-28 20l2 1 1 1h2l1 1 1 1h2l2 1 2 1 2 1 3 1 2 1h2l3 2h2l3 1 2 1 2 1 3 1 4 1 8 2 2 1 7 1h15l6-2 7-2 6-4 2-1 2-2 3-1 7-8 6-9 1-2 12-24a77 77 0 0 1 16-22"/></svg>`,
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
