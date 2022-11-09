package auth

import (
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

// NewSpotifyProvider creates a new Spotify provider instance with some defaults.
func NewKakaoProvider() *Kakao {
	return &Kakao{&baseProvider{
		scopes:     []string{"profile_nickname", "profile_image", "account_email", "name", "gender", "age_range", "birthyear"},
		authUrl:    kakao.Endpoint.AuthURL,
		tokenUrl:   kakao.Endpoint.TokenURL,
		userApiUrl: "https://kapi.kakao.com/v2/user/me",
	}}
}

// FetchAuthUser returns an AuthUser instance based on the Spotify's user api.
func (p *Kakao) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	// https://developer.spotify.com/documentation/web-api/reference/#/operations/get-current-users-profile
	rawData := struct {
		Id         string `json:"id"`
		Name       string `json:"kakao_account.name"`
		Email      string `json:"kakao_account.email"`
		ProfileImg string `json:"kakao_account.profile.thumbnail_image_url"`
	}{}

	if err := p.FetchRawUserData(token, &rawData); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:        rawData.Id,
		Name:      rawData.Name,
		Email:     rawData.Email,
		AvatarUrl: rawData.ProfileImg,
	}

	return user, nil
}
