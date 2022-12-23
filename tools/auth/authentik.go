package auth

/*
 *	https://goauthentik.io/
 *  authentik is an open-source Identity Provider focused on flexibility and versatility
 */

import (
	"fmt"

	"golang.org/x/oauth2"
)

var _ Provider = (*Authentik)(nil)

// Name of the oauth provider
const NameAuthentik string = "authentik"

type Authentik struct {
	*baseProvider
}

func NewAuthentikProvider() *Authentik {
	return &Authentik{&baseProvider{
		scopes:   []string{"profile"},
		authUrl:  "",
		tokenUrl: "",
	}}
}

func (p *Authentik) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	// TODO: Implement auth user function
	return nil, fmt.Errorf("not implemented")
}
