package mails

import (
	"fmt"
	"net/mail"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/mails/templates"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tokens"
	"github.com/pocketbase/pocketbase/tools/mailer"
	"github.com/pocketbase/pocketbase/tools/rest"
)

// SendAdminPasswordReset sends a password reset request email to the specified admin.
func SendAdminPasswordReset(app core.App, admin *models.Admin) error {
	token, tokenErr := tokens.NewAdminResetPasswordToken(app, admin)
	if tokenErr != nil {
		return tokenErr
	}

	actionUrl, urlErr := rest.NormalizeUrl(fmt.Sprintf(
		"%s/_/#/confirm-password-reset/%s",
		app.Settings().Meta.AppUrl,
		token,
	))
	if urlErr != nil {
		return urlErr
	}

	params := struct {
		AppName   string
		AppUrl    string
		Admin     *models.Admin
		Token     string
		ActionUrl string
	}{
		AppName:   app.Settings().Meta.AppName,
		AppUrl:    app.Settings().Meta.AppUrl,
		Admin:     admin,
		Token:     token,
		ActionUrl: actionUrl,
	}

	mailClient := app.NewMailClient()

	// resolve body template
	body, renderErr := resolveTemplateContent(params, templates.Layout, templates.AdminPasswordResetBody)
	if renderErr != nil {
		return renderErr
	}

	message := &mailer.Message{
		From: mail.Address{
			Name:    app.Settings().Meta.SenderName,
			Address: app.Settings().Meta.SenderAddress,
		},
		To:      []mail.Address{{Address: admin.Email}},
		Subject: "Reset admin password",
		HTML:    body,
	}

	event := new(core.MailerAdminEvent)
	event.MailClient = mailClient
	event.Message = message
	event.Admin = admin
	event.Meta = map[string]any{"token": token}

	return app.OnMailerBeforeAdminResetPasswordSend().Trigger(event, func(e *core.MailerAdminEvent) error {
		if err := e.MailClient.Send(e.Message); err != nil {
			return err
		}

		return app.OnMailerAfterAdminResetPasswordSend().Trigger(e)
	})
}
