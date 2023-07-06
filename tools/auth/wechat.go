package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

var _ Provider = (*WeChat)(nil)

// NameWeChat is the unique name of the WeChat provider.
const NameWeChat string = "wechat"

// Github allows authentication via Github OAuth2.
type WeChat struct {
	*baseProvider
}

// NewWeChatProvider creates new WeChat provider instance with some defaults.
func NewWeChatProvider() *WeChat {
	return &WeChat{&baseProvider{
		ctx:        context.Background(),
		scopes:     []string{"snsapi_login"},
		authUrl:    "https://open.weixin.qq.com/connect/qrconnect",
		tokenUrl:   "https://api.weixin.qq.com/sns/oauth2/access_token",
		userApiUrl: "https://api.weixin.qq.com/sns/userinfo",
	}}
}

type expirationTime int32

// tokenJSON is the struct representing the HTTP response from WeChat OAuth2
// providers returning a token or error in JSON form.
type tokenJSON struct {
	AccessToken  string         `json:"access_token"`
	RefreshToken string         `json:"refresh_token"`
	ExpiresIn    expirationTime `json:"expires_in"` // in seconds
	// error fields
	// https://developers.weixin.qq.com/doc/oplatform/Website_App/WeChat_Login/Wechat_Login.html
	// {"errcode":40030,"errmsg":"invalid refresh_token"}
	ErrorCode    string `json:"errcode"`
	ErrorMessage string `json:"errmsg"`
}

func (e *tokenJSON) expiry() (t time.Time) {
	if v := e.ExpiresIn; v != 0 {
		return time.Now().Add(time.Duration(v) * time.Second)
	}
	return
}
func (p *WeChat) FetchToken(code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	// WeChat OAuth2's token response will be "plain/text" even if it is a json
	// Therefore, we have to fetch and parse manually
	v := url.Values{
		"grant_type": {"authorization_code"},
		"code":       {code},
		"appid":      {p.clientId},
		"secret":     {p.clientSecret},
	}
	req, err := http.NewRequest("POST", p.tokenUrl, strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}
	client := oauth2.NewClient(p.ctx, nil)

	r, err := client.Do(req.WithContext(p.ctx))
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1<<20))
	r.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("oauth2: cannot fetch token: %v", err)
	}

	failureStatus := r.StatusCode < 200 || r.StatusCode > 299
	retrieveError := &oauth2.RetrieveError{
		Response: r,
		Body:     body,
		// attempt to populate error detail below
	}

	var token *oauth2.Token
	var tj tokenJSON
	if err = json.Unmarshal(body, &tj); err != nil {
		if failureStatus {
			return nil, retrieveError
		}
		return nil, fmt.Errorf("oauth2: cannot parse json: %v", err)
	}
	retrieveError.ErrorCode = tj.ErrorCode
	retrieveError.ErrorDescription = tj.ErrorMessage
	token = &oauth2.Token{
		AccessToken:  tj.AccessToken,
		RefreshToken: tj.RefreshToken,
		Expiry:       tj.expiry(),
	}

	var raw interface{}
	json.Unmarshal(body, &raw) // no error checks for optional fields

	if failureStatus || retrieveError.ErrorCode != "" {
		return nil, retrieveError
	}
	if token.AccessToken == "" {
		return nil, errors.New("oauth2: server response missing access_token")
	}
	return token.WithExtra(raw), nil
}

// Refer to: https://developers.weixin.qq.com/doc/oplatform/Website_App/WeChat_Login/Authorized_Interface_Calling_UnionID.html
func (p *WeChat) BuildUserApiUrl(token *oauth2.Token) string {
	var buf bytes.Buffer
	buf.WriteString(p.userApiUrl)
	v := url.Values{
		"access_token": {token.AccessToken},
		"openid":       {token.Extra("openid").(string)},
	}
	if strings.Contains(p.userApiUrl, "?") {
		buf.WriteByte('&')
	} else {
		buf.WriteByte('?')
	}
	buf.WriteString(v.Encode())
	return buf.String()
}

// FetchAuthUser returns an AuthUser instance based the WeChat's user api.
//
// API reference: https://developers.weixin.qq.com/doc/oplatform/en/Website_App/WeChat_Login/Wechat_Login.html
func (p *WeChat) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	req, err := http.NewRequest("GET", p.BuildUserApiUrl(token), nil)
	if err != nil {
		return nil, err
	}
	// access token is in the query parameter, and no need to set here.
	client := oauth2.NewClient(p.ctx, nil)

	r, err := client.Do(req.WithContext(p.ctx))
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	// http.Client.Get doesn't treat non 2xx responses as error
	if r.StatusCode >= 400 {
		return nil, fmt.Errorf(
			"failed to fetch OAuth2 user profile via %s (%d):\n%s",
			p.userApiUrl,
			r.StatusCode,
			string(data),
		)
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		OpenId     string `json:"openid"`
		UnionId    string `json:"unionid"`
		NickName   string `json:"nickname"`
		HeadImgUrl string `json:"headimgurl"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           extracted.UnionId,
		Name:         extracted.NickName,
		AvatarUrl:    extracted.HeadImgUrl,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	return user, nil
}
