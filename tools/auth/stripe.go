package auth

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

func init() {
	Providers[NameStripe] = wrapFactory(NewStripeProvider)
}

var _ Provider = (*Stripe)(nil)

// NameStripe is the unique name of the Stripe provider.
const NameStripe string = "stripe"

// Stripe allows authentication via Stripe Connect OAuth2.
type Stripe struct {
	BaseProvider
}

// NewStripeProvider creates a new Stripe provider instance with some defaults.
func NewStripeProvider() *Stripe {
	// https://docs.stripe.com/connect/oauth-reference
	return &Stripe{BaseProvider{
		ctx:         context.Background(),
		displayName: "Stripe",
		pkce:        false,
		scopes:      []string{"read_write"},
		authURL:     "https://connect.stripe.com/oauth/authorize",
		tokenURL:    "https://connect.stripe.com/oauth/token",
		userInfoURL: "https://api.stripe.com/v1/accounts",
	}}
}

// FetchAuthUser returns an AuthUser instance based on Stripe's account api.
//
// API reference: https://docs.stripe.com/api/accounts/retrieve
func (p *Stripe) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserInfo(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Id              string `json:"id"`
		Email           string `json:"email"`
		BusinessProfile struct {
			Name string `json:"name"`
		} `json:"business_profile"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           extracted.Id,
		Name:         extracted.BusinessProfile.Name,
		Email:        extracted.Email,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	return user, nil
}

// FetchRawUserInfo implements Provider.FetchRawUserInfo interface method.
//
// Stripe returns the account ID in the token response as stripe_user_id,
// which we use to fetch the account details.
func (p *Stripe) FetchRawUserInfo(token *oauth2.Token) ([]byte, error) {
	stripeUserId := token.Extra("stripe_user_id")
	if stripeUserId == nil {
		return nil, errors.New("missing stripe_user_id in token response")
	}

	accountId, ok := stripeUserId.(string)
	if !ok || accountId == "" {
		return nil, errors.New("invalid stripe_user_id in token response")
	}

	req, err := http.NewRequestWithContext(p.ctx, "GET", p.userInfoURL+"/"+accountId, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	result, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		return nil, errors.New("failed to fetch Stripe account info: " + string(result))
	}

	return result, nil
}
