package apis

import (
	"encoding/json"
	"net/http"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/subscriptions"
)

const (
	oauth2SubscriptionTopic   string = "@oauth2"
	oauth2RedirectFailurePath string = "../_/#/auth/oauth2-redirect-failure"
	oauth2RedirectSuccessPath string = "../_/#/auth/oauth2-redirect-success"
)

type oauth2RedirectData struct {
	State string `form:"state" json:"state"`
	Code  string `form:"code" json:"code"`
	Error string `form:"error" json:"error,omitempty"`
}

func oauth2SubscriptionRedirect(e *core.RequestEvent) error {
	redirectStatusCode := http.StatusTemporaryRedirect
	if e.Request.Method != http.MethodGet {
		redirectStatusCode = http.StatusSeeOther
	}

	data := oauth2RedirectData{}

	if e.Request.Method == http.MethodPost {
		if err := e.BindBody(&data); err != nil {
			e.App.Logger().Debug("Failed to read OAuth2 redirect data", "error", err)
			return e.Redirect(redirectStatusCode, oauth2RedirectFailurePath)
		}
	} else {
		query := e.Request.URL.Query()
		data.State = query.Get("state")
		data.Code = query.Get("code")
		data.Error = query.Get("error")
	}

	if data.State == "" {
		e.App.Logger().Debug("Missing OAuth2 state parameter")
		return e.Redirect(redirectStatusCode, oauth2RedirectFailurePath)
	}

	client, err := e.App.SubscriptionsBroker().ClientById(data.State)
	if err != nil || client.IsDiscarded() || !client.HasSubscription(oauth2SubscriptionTopic) {
		e.App.Logger().Debug("Missing or invalid OAuth2 subscription client", "error", err, "clientId", data.State)
		return e.Redirect(redirectStatusCode, oauth2RedirectFailurePath)
	}
	defer client.Unsubscribe(oauth2SubscriptionTopic)

	encodedData, err := json.Marshal(data)
	if err != nil {
		e.App.Logger().Debug("Failed to marshalize OAuth2 redirect data", "error", err)
		return e.Redirect(redirectStatusCode, oauth2RedirectFailurePath)
	}

	msg := subscriptions.Message{
		Name: oauth2SubscriptionTopic,
		Data: encodedData,
	}

	client.Send(msg)

	if data.Error != "" || data.Code == "" {
		e.App.Logger().Debug("Failed OAuth2 redirect due to an error or missing code parameter", "error", data.Error, "clientId", data.State)
		return e.Redirect(redirectStatusCode, oauth2RedirectFailurePath)
	}

	return e.Redirect(redirectStatusCode, oauth2RedirectSuccessPath)
}
