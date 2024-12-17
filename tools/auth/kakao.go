package auth

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/kakao"
)

func init() {
	Providers[NameKakao] = wrapFactory(NewKakaoProvider)
}

var _ Provider = (*Kakao)(nil)

// NameKakao is the unique name of the Kakao provider.
const NameKakao string = "kakao"

// Kakao allows authentication via Kakao OAuth2.
type Kakao struct {
	BaseProvider
}

// NewKakaoProvider creates a new Kakao provider instance with some defaults.
func NewKakaoProvider() *Kakao {
	return &Kakao{BaseProvider{
		ctx:         context.Background(),
		displayName: "Kakao",
		pkce:        true,
		scopes:      []string{"account_email", "profile_nickname", "profile_image"},
		authURL:     kakao.Endpoint.AuthURL,
		tokenURL:    kakao.Endpoint.TokenURL,
		userInfoURL: "https://kapi.kakao.com/v2/user/me",
	}}
}

// FetchAuthUser returns an AuthUser instance based on the Kakao's user api.
//
// API reference: https://developers.kakao.com/docs/latest/en/kakaologin/rest-api#req-user-info-response
func (p *Kakao) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserInfo(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Profile struct {
			Nickname string `json:"nickname"`
			ImageURL string `json:"profile_image"`
		} `json:"properties"`
		KakaoAccount struct {
			Email           string `json:"email"`
			IsEmailVerified bool   `json:"is_email_verified"`
			IsEmailValid    bool   `json:"is_email_valid"`
		} `json:"kakao_account"`
		Id int64 `json:"id"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           strconv.FormatInt(extracted.Id, 10),
		Username:     extracted.Profile.Nickname,
		AvatarURL:    extracted.Profile.ImageURL,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	if extracted.KakaoAccount.IsEmailValid && extracted.KakaoAccount.IsEmailVerified {
		user.Email = extracted.KakaoAccount.Email
	}

	return user, nil
}
