package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

func init() {
	Providers[NameWechatNative] = func() Provider {
		return NewWechatProvider(NameWechatNative)
	}
	Providers[NameWechatWeb] = func() Provider {
		return NewWechatProvider(NameWechatWeb)
	}
	EquivalentProviders[NameWechatNative] = []string{NameWechatWeb}
	EquivalentProviders[NameWechatWeb] = []string{NameWechatNative}
}

var _ Provider = (*Wechat)(nil)

// NameWechat is the unique name of the Wechat provider.
const NameWechatNative string = "wechat_native"
const NameWechatWeb string = "wechat_web"

// Wechat allows authentication via Wechat OAuth2.
type Wechat struct {
	BaseProvider
}

// NewWechatProvider creates a new Wechat provider instance with some defaults.
func NewWechatProvider(providerName string) *Wechat {
	var scopes []string
	var displayName string
	if providerName == NameWechatNative {
		scopes = []string{"snsapi_userinfo"}
		displayName = "Wechat (Native)"
	} else {
		scopes = []string{"snsapi_login"}
		displayName = "Wechat (Web)"
	}
	return &Wechat{
		BaseProvider: BaseProvider{
			ctx:         context.Background(),
			displayName: displayName,
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

	// Replace the underlying http client to fix WeChat's weird OAuth2 API issue.
	if p.Context().Value(oauth2.HTTPClient) == nil {
		p.SetContext(context.WithValue(p.Context(), oauth2.HTTPClient, NewJsonHttpClient()))
	}

	return p.BaseProvider.BuildAuthURL(state, opts...)
}

// FetchAuthUser returns an AuthUser instance based on the provided token.
//
// API reference: https://developers.weixin.qq.com/doc/oplatform/en/Mobile_App/WeChat_Login/Development_Guide.html
func (p *Wechat) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	unionid, _ := token.Extra("unionid").(string)

	// TODO: Support silent implict login flow for Mini Apps and In-Wechat Browsers.
	// 难点在于：
	// 当用户用户关注公众号前，在微信内访问h5时，静默登录的用户只有openid，没有unionid。此时记录在数据库中的用户unionid是空的。
	// 当用户关注公众号后，在微信内访问h5网页时，静默登录的用户有openid和unionid。
	//
	// 小程序同理：
	// 用户首次访问当前主体开发的小程序时，如果采用静默登录，用户只有openid，没有unionid。
	// 如果用户此前访问过当前主体开发的小程序，并授权过userinfo, 下次静默登录时，用户有openid和unionid。
	//
	// 我们标志用户身份时，是采用 openid 还是 unionid 呢？
	// 如果采用 openid，则用户在app端/网页端/小程序端登录时，会产生多个不同的身份。
	//
	// 权衡之后，为了确保多端访问身份的统一，我们决定默认采用 unionid 作为用户身份标志。
	// 这样就只能拒绝静默登录且无unionid的用户登录。
	//
	// 一个完美的解决方案是：默认采用unionid作为用户身份标志，但是当获取不到unionid时，
	// 允许用户使用openid作为身份标志登录，同时下次登录如果获取到unionid时，更新用户的身份标志为unionid。
	// 但是这个方案实现太复杂，因此暂时只能拒绝无unionid的登录了。
	//
	// PS: 禁止无unionid登录只会影响到小程序和公众号的*无授权*静默登录功能。网页端扫码登录和移动端登录都不会受影响。
	if unionid == "" {
		return nil, fmt.Errorf("暂不支持静默授权且无unionid的用户登录。请使用snsapi_userinfo重新登录。")
	}

	// Scope "snsapi_userinfo" is always included in mobile native auth.
	// But for web oauth flow, it is not added by default.
	data, err := p.FetchRawUserInfo(token)
	if err != nil {
		return nil, err
	}

	// Update user.RawUser
	var rawUser map[string]any
	if err := json.Unmarshal(data, &rawUser); err != nil {
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

	user := &AuthUser{
		Id:           unionid,
		Name:         extracted.Nickname,
		AvatarURL:    extracted.Headimgurl,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
	user.Expiry, _ = types.ParseDateTime(token.Expiry)

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

// JsonTransport is a custom [http.RoundTripper] that forces
// the response Content-Type to "application/json" in WeChat Token API.
type JsonTransport struct {
	Base http.RoundTripper
}

// RoundTrip executes a single HTTP transaction and modifies the response header.
func (t *JsonTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Use the default transport if none is provided
	transport := t.Base
	if transport == nil {
		transport = http.DefaultTransport
	}

	resp, err := transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// Fix a weird Wechat OAuth2 API issue:
	// The access token API response with a JSON, but in its head, the content-type is "text/plain".
	// This will make golang.org/x/oauth2 fail to parse the response.
	// By manually setting the content-type to "application/json", we can avoid this issue.
	if strings.HasPrefix(req.URL.Path, "/sns/oauth2/access_token") {
		resp.Header.Set("Content-Type", "application/json")
	}

	return resp, nil
}

func NewJsonHttpClient() *http.Client {
	return &http.Client{
		Transport: &JsonTransport{},
	}
}
