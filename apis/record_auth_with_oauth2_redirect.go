package apis

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/subscriptions"
	"github.com/pocketbase/pocketbase/ui"
)

const (
	oauth2SubscriptionTopic               string = "@oauth2"
	oauth2RedirectFailurePath             string = "../_/#/auth/oauth2-redirect-failure"
	oauth2RedirectSuccessPath             string = "../_/#/auth/oauth2-redirect-success"
	oauth2RedirectAppleNameStoreKeyPrefix string = "@redirect_name_"
)

type oauth2RedirectData struct {
	State string `form:"state" json:"state"`
	Code  string `form:"code" json:"code"`
	Error string `form:"error" json:"error,omitempty"`

	// returned by Apple only
	AppleUser string `form:"user" json:"-"`
}

func oauth2SubscriptionRedirect(e *core.RequestEvent) error {
	data := oauth2RedirectData{}

	if e.Request.Method == http.MethodPost {
		if err := e.BindBody(&data); err != nil {
			e.App.Logger().Debug("Failed to read OAuth2 redirect data", "error", err)
			return failureRedirect(e)
		}
	} else {
		query := e.Request.URL.Query()
		data.State = query.Get("state")
		data.Code = query.Get("code")
		data.Error = query.Get("error")
	}

	if data.State == "" {
		e.App.Logger().Debug("Missing OAuth2 state parameter")
		return failureRedirect(e)
	}

	client, err := e.App.SubscriptionsBroker().ClientById(data.State)
	if err != nil || client.IsDiscarded() || !client.HasSubscription(oauth2SubscriptionTopic) {
		e.App.Logger().Debug("Missing or invalid OAuth2 subscription client", "error", err, "clientId", data.State)
		return failureRedirect(e)
	}
	defer client.Unsubscribe(oauth2SubscriptionTopic)

	// temporary store the Apple user's name so that it can be later retrieved with the authWithOAuth2 call
	// (see https://github.com/pocketbase/pocketbase/issues/7090)
	if data.AppleUser != "" && data.Error == "" && data.Code != "" {
		nameErr := parseAndStoreAppleRedirectName(
			e.App,
			oauth2RedirectAppleNameStoreKeyPrefix+data.Code,
			data.AppleUser,
		)
		if nameErr != nil {
			// non-critical error
			e.App.Logger().Debug("Failed to parse and load Apple Redirect name data", "error", nameErr)
		}
	}

	encodedData, err := json.Marshal(data)
	if err != nil {
		e.App.Logger().Debug("Failed to marshalize OAuth2 redirect data", "error", err)
		return failureRedirect(e)
	}

	msg := subscriptions.Message{
		Name: oauth2SubscriptionTopic,
		Data: encodedData,
	}

	client.Send(msg)

	if data.Error != "" || data.Code == "" {
		e.App.Logger().Debug("Failed OAuth2 redirect due to an error or missing code parameter", "error", data.Error, "clientId", data.State)
		return failureRedirect(e)
	}

	return successRedirect(e)
}

func redirectStatusCode(e *core.RequestEvent) int {
	if e.Request.Method != http.MethodGet {
		return http.StatusSeeOther
	}

	return http.StatusTemporaryRedirect
}

func failureRedirect(e *core.RequestEvent) error {
	// fallback if UI is not bundled
	if ui.DistDirFS == nil {
		return e.String(http.StatusOK, "Failed to authenticate. You can close this window and go back to the app to try again.")
	}

	return e.Redirect(redirectStatusCode(e), oauth2RedirectFailurePath)
}

func successRedirect(e *core.RequestEvent) error {
	// fallback if UI is not bundled
	if ui.DistDirFS == nil {
		return e.HTML(http.StatusOK, "Auth completed. You can close this window and go back to the app.")
	}

	return e.Redirect(redirectStatusCode(e), oauth2RedirectSuccessPath)
}

// parseAndStoreAppleRedirectName extracts the first and last name
// from serializedNameData and temporary store them in the app.Store.
//
// This is hacky workaround to forward safely and seamlessly the Apple
// redirect user's name back to the OAuth2 auth handler.
//
// Note that currently Apple is the only provider that behaves like this and
// for now it is unnecessary to check whether the redirect is coming from Apple or not.
//
// Ideally this shouldn't be needed and will be removed in the future
// once Apple adds a dedicated userinfo endpoint.
func parseAndStoreAppleRedirectName(app core.App, nameKey string, serializedNameData string) error {
	if serializedNameData == "" {
		return nil
	}

	// just in case to prevent storing large strings in memory
	if len(nameKey) > 1000 {
		return errors.New("nameKey is too large")
	}

	// https://developer.apple.com/documentation/signinwithapple/incorporating-sign-in-with-apple-into-other-platforms#Handle-the-response
	extracted := struct {
		Name struct {
			FirstName string `json:"firstName"`
			LastName  string `json:"lastName"`
		} `json:"name"`
	}{}
	if err := json.Unmarshal([]byte(serializedNameData), &extracted); err != nil {
		return err
	}

	fullName := extracted.Name.FirstName + " " + extracted.Name.LastName

	// truncate just in case to prevent storing large strings in memory
	if len(fullName) > 150 {
		fullName = fullName[:150]
	}

	fullName = strings.TrimSpace(fullName)
	if fullName == "" {
		return nil
	}

	// store (and remove)
	app.Store().Set(nameKey, fullName)
	time.AfterFunc(1*time.Minute, func() {
		app.Store().Remove(nameKey)
	})

	return nil
}
