package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

func init() {
	Providers[NameWechatMobile] = func() Provider {
		return NewWechatProvider("mobile")
	}
	Providers[NameWechatWeb] = func() Provider {
		return NewWechatProvider("web")
	}
}

var _ Provider = (*Wechat)(nil)

// NameWechat is the unique name of the Wechat provider.
const NameWechatMobile string = "wechat_mobile"
const NameWechatWeb string = "wechat_web"

// Wechat allows authentication via Wechat OAuth2.
type Wechat struct {
	BaseProvider
}

// NewWechatProvider creates a new Wechat provider instance with some defaults.
func NewWechatProvider(authType string) *Wechat {
	var scopes []string
	if authType == "mobile" {
		scopes = []string{"snsapi_userinfo"}
	} else {
		scopes = []string{"snsapi_login"}
	}
	return &Wechat{
		BaseProvider: BaseProvider{
			ctx:         context.Background(),
			displayName: fmt.Sprintf("Wechat (%s)", authType),
			scopes:      scopes,
			authURL:     "https://open.weixin.qq.com/connect/qrconnect",
			tokenURL:    "https://api.weixin.qq.com/sns/oauth2/access_token",
			userInfoURL: "https://api.weixin.qq.com/sns/userinfo",
		},
	}
}

// Wechat requires the `appid` and `lang` parameters to be passed in the request URL.
func (p *Wechat) BuildAuthURL(state string, opts ...oauth2.AuthCodeOption) string {
	opts = append(opts, oauth2.SetAuthURLParam("appid", p.clientId))
	opts = append(opts, oauth2.SetAuthURLParam("lang", "cn")) // Optional, default is "cn"
	return p.BaseProvider.BuildAuthURL(state, opts...)
}

// FetchAuthUser returns an AuthUser instance based on the provided token.
//
// API reference: https://developers.weixin.qq.com/doc/oplatform/en/Mobile_App/WeChat_Login/Development_Guide.html
func (p *Wechat) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	openid, _ := token.Extra("openid").(string)
	unionid, _ := token.Extra("unionid").(string)
	authorizedScopes, _ := token.Extra("scope").(string)
	scopes := strings.Split(authorizedScopes, ",")

	user := &AuthUser{
		Id:   openid,
		Name: "",
		RawUser: map[string]any{
			"openid":  openid,
			"unionid": unionid, // saving for future use
		},
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	// Scope "snsapi_userinfo" is always included in mobile native auth.
	// But for web oauth flow, it is not added by default.
	if slices.Contains(scopes, "snsapi_userinfo") {
		data, err := p.FetchRawUserInfo(token)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(data, &user.RawUser); err != nil {
			return nil, err
		}

		extracted := struct {
			OpenId     string   `json:"openid"`
			UnionId    string   `json:"unionid"`
			Nickname   string   `json:"nickname"`
			Headimgurl string   `json:"headimgurl"`
			Priviledge []string `json:"privilege"`
		}{}
		if err := json.Unmarshal(data, &extracted); err != nil {
			return nil, err
		}

		// Set nickname to [AuthUser.Name] instead of [AuthUser.Username] because
		// the real username is not accessible via Wechat OAuth2 API.
		user.Name = extracted.Nickname
		user.AvatarURL = extracted.Headimgurl
	}

	return user, nil
}

// FetchToken implements Provider.FetchToken() interface method.
//
// Other than the traditinal "code" and "grant_type" parameters, Wechat requires
// the "appid" and "secret" parameters to be passed in the request body.
// For example:
//
//	GET https://api.weixin.qq.com/sns/oauth2/access_token?appid=APPID&secret=SECRET&code=CODE&grant_type=authorization_code
//
// API reference: https://developers.weixin.qq.com/doc/oplatform/en/Mobile_App/WeChat_Login/Development_Guide.html
func (p *Wechat) FetchToken(code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	opts = append(opts, oauth2.SetAuthURLParam("appid", p.clientId))
	opts = append(opts, oauth2.SetAuthURLParam("secret", p.clientSecret))

	return p.BaseProvider.FetchToken(code, opts...)
}

// FetchRawUserInfo implements [Provider.FetchRawUserInfo] interface.
//
// Wechat user info endpoint does not use normal `Bearer` authentication header.
// Instead, it requires `access_token` and `openid` to be passed as query parameters.
// For example:
//
//	GET https://api.weixin.qq.com/sns/userinfo?access_token=ACCESS_TOKEN&openid=OPENID
//
// API reference: https://developers.weixin.qq.com/doc/oplatform/en/Website_App/WeChat_Login/Authorized_Interface_Calling_UnionID.html
func (p *Wechat) FetchRawUserInfo(token *oauth2.Token) ([]byte, error) {
	req, err := http.NewRequestWithContext(p.ctx, "GET", p.userInfoURL, nil)
	if err != nil {
		return nil, err
	}

	// Adding `access_token` and `openid` as query parameters
	q := req.URL.Query()
	q.Add("access_token", token.AccessToken)
	q.Add("openid", token.Extra("openid").(string))
	req.URL.RawQuery = q.Encode()

	return p.sendRawUserInfoRequest(req, token)
}
