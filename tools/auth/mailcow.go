package auth

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

func init() {
	Providers[NameMailcow] = wrapFactory(NewMailcowProvider)
}

var _ Provider = (*Mailcow)(nil)

// NameMailcow is the unique name of the mailcow provider.
const NameMailcow string = "mailcow"

// Mailcow allows authentication via mailcow OAuth2.
type Mailcow struct {
	BaseProvider
}

// NewMailcowProvider creates a new mailcow provider instance with some defaults.
func NewMailcowProvider() *Mailcow {
	return &Mailcow{BaseProvider{
		ctx:         context.Background(),
		displayName: "mailcow",
		pkce:        true,
		scopes:      []string{"profile"},
	}}
}

// FetchAuthUser returns an AuthUser instance based on mailcow's user api.
//
// API reference: https://github.com/mailcow/mailcow-dockerized/blob/master/data/web/oauth/profile.php
func (p *Mailcow) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserInfo(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Id       string `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		FullName string `json:"full_name"`
		Active   int    `json:"active"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	if extracted.Active != 1 {
		return nil, errors.New("the mailcow user is not active")
	}

	user := &AuthUser{
		Id:           extracted.Id,
		Name:         extracted.FullName,
		Username:     extracted.Username,
		Email:        extracted.Email,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	// mailcow usernames are usually just the email adresses, so we just take the part in front of the @
	if strings.Contains(user.Username, "@") {
		user.Username = strings.Split(user.Username, "@")[0]
	}

	return user, nil
}
