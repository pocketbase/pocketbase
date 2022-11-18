package auth

import (
	"strconv"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/kakao"
)

var _ Provider = (*Kakao)(nil)

// NameKakao is the unique name of the Kakao provider.
const NameKakao string = "kakao"

// Kakao allows authentication via Kakao OAuth2.
type Kakao struct {
	*baseProvider
}

// NewKakaoProvider creates a new Kakao provider instance with some defaults.
func NewKakaoProvider() *Kakao {
	return &Kakao{&baseProvider{
		scopes:     []string{"account_email", "profile_nickname", "profile_image"},
		authUrl:    kakao.Endpoint.AuthURL,
		tokenUrl:   kakao.Endpoint.TokenURL,
		userApiUrl: "https://kapi.kakao.com/v2/user/me",
	}}
}

// FetchAuthUser returns an AuthUser instance based on the Kakao's user api.
func (p *Kakao) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	// https://developers.kakao.com/docs/latest/en/kakaologin/rest-api#req-user-info-response
	rawData := struct {
		Id      int `json:"id"`
		Profile struct {
			Nickname string `json:"nickname"`
			ImageUrl string `json:"profile_image"`
		} `json:"properties"`
		KakaoAccount struct {
			Email           string `json:"email"`
			IsEmailVerified bool   `json:"is_email_verified"`
			IsEmailValid    bool   `json:"is_email_valid"`
		} `json:"kakao_account"`
	}{}

	if err := p.FetchRawUserData(token, &rawData); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:        strconv.Itoa(rawData.Id),
		Username:  rawData.Profile.Nickname,
		AvatarUrl: rawData.Profile.ImageUrl,
	}
	if rawData.KakaoAccount.IsEmailValid && rawData.KakaoAccount.IsEmailVerified {
		user.Email = rawData.KakaoAccount.Email
	}

	return user, nil
}
