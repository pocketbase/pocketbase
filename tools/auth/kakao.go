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
		order:       14,
		logo:        `<svg xmlns="http://www.w3.org/2000/svg" width="2500" height="2500" viewBox="0 0 256 256"><path fill="#ffe812" d="M256 236q-2 18-20 20H20q-18-2-20-20V20Q2 2 20 0h216q18 2 20 20z"/><path d="M128 36C71 36 24 73 24 118c0 29 19 55 49 69l-11 38 1 3h3l44-29 18 1c57 0 104-37 104-82s-47-82-104-82"/><path fill="#ffe812" d="M71 147q-6-1-6-6v-36H55a6 6 0 0 1 0-11h31a6 6 0 0 1 0 11h-9v36q-1 5-6 6m52 0q-3 0-5-3l-3-8H97l-3 8q-2 3-5 3l-4-1q-3-1-1-9l14-38q2-4 8-5 6 1 8 5l14 38q2 8-1 9zm-11-22-6-17-6 17zm26 21q-5-1-6-6v-40q1-6 6-6 7 0 7 6v35h12q5 0 6 5 0 7-6 6zm33 1q-5-1-6-6v-41a6 6 0 0 1 12 0v12l17-16 3-2 5 2 1 4-1 4-14 13 15 20 1 4-2 4-4 1-5-2-14-19-2 2v14a6 6 0 0 1-6 6"/></svg>`,
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
