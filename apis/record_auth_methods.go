package apis

import (
	"log/slog"
	"net/http"
	"slices"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/auth"
	"github.com/pocketbase/pocketbase/tools/security"
	"golang.org/x/oauth2"
)

type otpResponse struct {
	Enabled  bool  `json:"enabled"`
	Duration int64 `json:"duration"` // in seconds
}

type mfaResponse struct {
	Enabled  bool  `json:"enabled"`
	Duration int64 `json:"duration"` // in seconds
}

type passwordResponse struct {
	IdentityFields []string `json:"identityFields"`
	Enabled        bool     `json:"enabled"`
}

type oauth2Response struct {
	Providers []providerInfo `json:"providers"`
	Enabled   bool           `json:"enabled"`
}

type providerInfo struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	State       string `json:"state"`
	AuthURL     string `json:"authURL"`

	// @todo
	// deprecated: use AuthURL instead
	// AuthUrl will be removed after dropping v0.22 support
	AuthUrl string `json:"authUrl"`

	// technically could be omitted if the provider doesn't support PKCE,
	// but to avoid breaking existing typed clients we'll return them as empty string
	CodeVerifier        string `json:"codeVerifier"`
	CodeChallenge       string `json:"codeChallenge"`
	CodeChallengeMethod string `json:"codeChallengeMethod"`
}

type authMethodsResponse struct {
	Password passwordResponse `json:"password"`
	OAuth2   oauth2Response   `json:"oauth2"`
	MFA      mfaResponse      `json:"mfa"`
	OTP      otpResponse      `json:"otp"`

	// legacy fields
	// @todo remove after dropping v0.22 support
	AuthProviders    []providerInfo `json:"authProviders"`
	UsernamePassword bool           `json:"usernamePassword"`
	EmailPassword    bool           `json:"emailPassword"`
}

func (amr *authMethodsResponse) fillLegacyFields() {
	amr.EmailPassword = amr.Password.Enabled && slices.Contains(amr.Password.IdentityFields, "email")

	amr.UsernamePassword = amr.Password.Enabled && slices.Contains(amr.Password.IdentityFields, "username")

	if amr.OAuth2.Enabled {
		amr.AuthProviders = amr.OAuth2.Providers
	}
}

func recordAuthMethods(e *core.RequestEvent) error {
	collection, err := findAuthCollection(e)
	if err != nil {
		return err
	}

	result := authMethodsResponse{
		Password: passwordResponse{
			IdentityFields: make([]string, 0, len(collection.PasswordAuth.IdentityFields)),
		},
		OAuth2: oauth2Response{
			Providers: make([]providerInfo, 0, len(collection.OAuth2.Providers)),
		},
		OTP: otpResponse{
			Enabled: collection.OTP.Enabled,
		},
		MFA: mfaResponse{
			Enabled: collection.MFA.Enabled,
		},
	}

	if collection.PasswordAuth.Enabled {
		result.Password.Enabled = true
		result.Password.IdentityFields = collection.PasswordAuth.IdentityFields
	}

	if collection.OTP.Enabled {
		result.OTP.Duration = collection.OTP.Duration
	}

	if collection.MFA.Enabled {
		result.MFA.Duration = collection.MFA.Duration
	}

	if !collection.OAuth2.Enabled {
		result.fillLegacyFields()

		return e.JSON(http.StatusOK, result)
	}

	result.OAuth2.Enabled = true

	for _, config := range collection.OAuth2.Providers {
		provider, err := config.InitProvider()
		if err != nil {
			e.App.Logger().Debug(
				"Failed to setup OAuth2 provider",
				slog.String("name", config.Name),
				slog.String("error", err.Error()),
			)
			continue // skip provider
		}

		info := providerInfo{
			Name:        config.Name,
			DisplayName: provider.DisplayName(),
			State:       security.RandomString(30),
		}

		if info.DisplayName == "" {
			info.DisplayName = config.Name
		}

		urlOpts := []oauth2.AuthCodeOption{}

		// custom providers url options
		switch config.Name {
		case auth.NameApple:
			// see https://developer.apple.com/documentation/sign_in_with_apple/sign_in_with_apple_js/incorporating_sign_in_with_apple_into_other_platforms#3332113
			urlOpts = append(urlOpts, oauth2.SetAuthURLParam("response_mode", "form_post"))
		}

		if provider.PKCE() {
			info.CodeVerifier = security.RandomString(43)
			info.CodeChallenge = security.S256Challenge(info.CodeVerifier)
			info.CodeChallengeMethod = "S256"
			urlOpts = append(urlOpts,
				oauth2.SetAuthURLParam("code_challenge", info.CodeChallenge),
				oauth2.SetAuthURLParam("code_challenge_method", info.CodeChallengeMethod),
			)
		}

		info.AuthURL = provider.BuildAuthURL(
			info.State,
			urlOpts...,
		) + "&redirect_uri=" // empty redirect_uri so that users can append their redirect url

		info.AuthUrl = info.AuthURL

		result.OAuth2.Providers = append(result.OAuth2.Providers, info)
	}

	result.fillLegacyFields()

	return e.JSON(http.StatusOK, result)
}
