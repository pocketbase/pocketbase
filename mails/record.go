package mails

import (
	"html/template"
	"net/mail"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/mails/templates"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/settings"
	"github.com/pocketbase/pocketbase/tokens"
	"github.com/pocketbase/pocketbase/tools/mailer"
)

// @todo remove after the refactoring
//
// SendRecordPasswordLoginAlert sends a OAuth2 password login alert to the specified auth record.
func SendRecordPasswordLoginAlert(app core.App, authRecord *models.Record, providerNames ...string) error {
	params := struct {
		AppName       string
		AppUrl        string
		Record        *models.Record
		ProviderNames []string
	}{
		AppName:       app.Settings().Meta.AppName,
		AppUrl:        app.Settings().Meta.AppUrl,
		Record:        authRecord,
		ProviderNames: providerNames,
	}

	mailClient := app.NewMailClient()

	// resolve body template
	body, renderErr := resolveTemplateContent(params, templates.Layout, templates.PasswordLoginAlertBody)
	if renderErr != nil {
		return renderErr
	}

	message := &mailer.Message{
		From: mail.Address{
			Name:    app.Settings().Meta.SenderName,
			Address: app.Settings().Meta.SenderAddress,
		},
		To:      []mail.Address{{Address: authRecord.Email()}},
		Subject: "Password login alert",
		HTML:    body,
	}

	return mailClient.Send(message)
}

// SendRecordPasswordReset sends a password reset request email to the specified user.
func SendRecordPasswordReset(app core.App, authRecord *models.Record) error {
	token, tokenErr := tokens.NewRecordResetPasswordToken(app, authRecord)
	if tokenErr != nil {
		return tokenErr
	}

	mailClient := app.NewMailClient()

	subject, body, err := resolveEmailTemplate(app, token, app.Settings().Meta.ResetPasswordTemplate)
	if err != nil {
		return err
	}

	message := &mailer.Message{
		From: mail.Address{
			Name:    app.Settings().Meta.SenderName,
			Address: app.Settings().Meta.SenderAddress,
		},
		To:      []mail.Address{{Address: authRecord.Email()}},
		Subject: subject,
		HTML:    body,
	}

	event := new(core.MailerRecordEvent)
	event.MailClient = mailClient
	event.Message = message
	event.Collection = authRecord.Collection()
	event.Record = authRecord
	event.Meta = map[string]any{"token": token}

	return app.OnMailerBeforeRecordResetPasswordSend().Trigger(event, func(e *core.MailerRecordEvent) error {
		if err := e.MailClient.Send(e.Message); err != nil {
			return err
		}

		return app.OnMailerAfterRecordResetPasswordSend().Trigger(e)
	})
}

// SendRecordVerification sends a verification request email to the specified user.
func SendRecordVerification(app core.App, authRecord *models.Record) error {
	token, tokenErr := tokens.NewRecordVerifyToken(app, authRecord)
	if tokenErr != nil {
		return tokenErr
	}

	mailClient := app.NewMailClient()

	subject, body, err := resolveEmailTemplate(app, token, app.Settings().Meta.VerificationTemplate)
	if err != nil {
		return err
	}

	message := &mailer.Message{
		From: mail.Address{
			Name:    app.Settings().Meta.SenderName,
			Address: app.Settings().Meta.SenderAddress,
		},
		To:      []mail.Address{{Address: authRecord.Email()}},
		Subject: subject,
		HTML:    body,
	}

	event := new(core.MailerRecordEvent)
	event.MailClient = mailClient
	event.Message = message
	event.Collection = authRecord.Collection()
	event.Record = authRecord
	event.Meta = map[string]any{"token": token}

	return app.OnMailerBeforeRecordVerificationSend().Trigger(event, func(e *core.MailerRecordEvent) error {
		if err := e.MailClient.Send(e.Message); err != nil {
			return err
		}

		return app.OnMailerAfterRecordVerificationSend().Trigger(e)
	})
}

// SendRecordChangeEmail sends a change email confirmation email to the specified user.
func SendRecordChangeEmail(app core.App, record *models.Record, newEmail string) error {
	token, tokenErr := tokens.NewRecordChangeEmailToken(app, record, newEmail)
	if tokenErr != nil {
		return tokenErr
	}

	mailClient := app.NewMailClient()

	subject, body, err := resolveEmailTemplate(app, token, app.Settings().Meta.ConfirmEmailChangeTemplate)
	if err != nil {
		return err
	}

	message := &mailer.Message{
		From: mail.Address{
			Name:    app.Settings().Meta.SenderName,
			Address: app.Settings().Meta.SenderAddress,
		},
		To:      []mail.Address{{Address: newEmail}},
		Subject: subject,
		HTML:    body,
	}

	event := new(core.MailerRecordEvent)
	event.MailClient = mailClient
	event.Message = message
	event.Collection = record.Collection()
	event.Record = record
	event.Meta = map[string]any{
		"token":    token,
		"newEmail": newEmail,
	}

	return app.OnMailerBeforeRecordChangeEmailSend().Trigger(event, func(e *core.MailerRecordEvent) error {
		if err := e.MailClient.Send(e.Message); err != nil {
			return err
		}

		return app.OnMailerAfterRecordChangeEmailSend().Trigger(e)
	})
}

func resolveEmailTemplate(
	app core.App,
	token string,
	emailTemplate settings.EmailTemplate,
) (subject string, body string, err error) {
	subject, rawBody, _ := emailTemplate.Resolve(
		app.Settings().Meta.AppName,
		app.Settings().Meta.AppUrl,
		token,
	)

	params := struct {
		HtmlContent template.HTML
	}{
		HtmlContent: template.HTML(rawBody),
	}

	body, err = resolveTemplateContent(params, templates.Layout, templates.HtmlBody)
	if err != nil {
		return "", "", err
	}

	return subject, body, nil
}
