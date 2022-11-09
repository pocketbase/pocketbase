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
	// https://developers.kakao.com/docs/latest/en/kakaologin/prerequisite#personal-information
	rawData := struct {
		Id      int `json:"id"`
		Profile struct {
			Name     string `json:"nickname"`
			ImageUrl string `json:"profile_image"`
		} `json:"properties"`
		KakaoAccount struct {
			Email string `json:"email"`
		} `json:"kakao_account"`
	}{}

	if err := p.FetchRawUserData(token, &rawData); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:        strconv.Itoa(rawData.Id),
		Name:      rawData.Profile.Name,
		Email:     rawData.KakaoAccount.Email,
		AvatarUrl: rawData.Profile.ImageUrl,
	}

	return user, nil
}
