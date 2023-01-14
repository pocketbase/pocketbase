package auth

/*
 *	https://goauthentik.io/
 *  authentik is an open-source Identity Provider focused on flexibility and versatility
 */

import (
	"encoding/json"
	"strconv"

	"golang.org/x/oauth2"
)

var _ Provider = (*Authentik)(nil)

// Unique name of this authentification provider
const NameAuthentik string = "authentik"

// Authentik is a self hosted authentification provider
type Authentik struct {
	*baseProvider
}

// Create a new Authentik provider instance with some defaults
func NewAuthentikProvider() *Authentik {
	return &Authentik{&baseProvider{
		// Default scopes provided to any oauth2 providers
		scopes: []string{
			"email",
			"profile",
			"openid",
		},
		authUrl:    "",
		tokenUrl:   "",
		userApiUrl: "",
	}}
}

// FetchAuthUser returns an AuthUser instance based on the Authentik user api
// Reference:
func (p *Authentik) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserData(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Id       int    `json:"id"`
		Name     string `json:"name"`
		Username string `json:"preferred_username"`
		Email    string `json:"email"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	// Username value in authentik is not restricted with min / max length
	if len(extracted.Username) <= 3 || len(extracted.Username) > 100 {
		// Update username to be empty, and use only e-mail address
		extracted.Username = ""
	}

	user := &AuthUser{
		Id:           strconv.Itoa(extracted.Id),
		Name:         extracted.Name,
		Username:     extracted.Username,
		Email:        extracted.Email,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	return user, nil
}
