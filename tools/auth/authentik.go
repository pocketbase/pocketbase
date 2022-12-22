package auth

import (
	"fmt"

	"golang.org/x/oauth2"
)

var _ Provider = (*Authentik)(nil)

const NameAuthentik string = "authentik"

type Authentik struct {
	*baseProvider
}

func NewAuthentikProvider() *Authentik {
	return &Authentik{&baseProvider{
		scopes:   []string{},
		authUrl:  "",
		tokenUrl: "",
	}}
}

func (p *Authentik) BuildAuthUrl(state string, opts ...oauth2.AuthCodeOption) string {
	return p.oauth2Config().AuthCodeURL(state, opts...)
}

func (p *Authentik) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	// TODO: Implement auth user function
	return nil, fmt.Errorf("not implemented")
}
