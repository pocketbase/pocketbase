package mails

import (
	"html/template"
	"net/mail"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/mails/templates"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tokens"
)

// SendRecordPasswordReset sends a password reset request email to the specified user.
func SendRecordPasswordReset(app core.App, authRecord *models.Record) error {
	token, tokenErr := tokens.NewRecordResetPasswordToken(app, authRecord)
	if tokenErr != nil {
		return tokenErr
	}

	mailClient := app.NewMailClient()

	event := &core.MailerRecordEvent{
		MailClient: mailClient,
		Record:     authRecord,
		Meta:       map[string]any{"token": token},
	}

	sendErr := app.OnMailerBeforeRecordResetPasswordSend().Trigger(event, func(e *core.MailerRecordEvent) error {
		settings := app.Settings()

		subject, body, err := resolveEmailTemplate(app, token, settings.Meta.ResetPasswordTemplate)
		if err != nil {
			return err
		}

		return e.MailClient.Send(
			mail.Address{
				Name:    settings.Meta.SenderName,
				Address: settings.Meta.SenderAddress,
			},
			mail.Address{Address: e.Record.GetString(schema.FieldNameEmail)},
			subject,
			body,
			nil,
		)
	})

	if sendErr == nil {
		app.OnMailerAfterRecordResetPasswordSend().Trigger(event)
	}

	return sendErr
}

// SendRecordVerification sends a verification request email to the specified user.
func SendRecordVerification(app core.App, authRecord *models.Record) error {
	token, tokenErr := tokens.NewRecordVerifyToken(app, authRecord)
	if tokenErr != nil {
		return tokenErr
	}

	mailClient := app.NewMailClient()

	event := &core.MailerRecordEvent{
		MailClient: mailClient,
		Record:     authRecord,
		Meta:       map[string]any{"token": token},
	}

	sendErr := app.OnMailerBeforeRecordVerificationSend().Trigger(event, func(e *core.MailerRecordEvent) error {
		settings := app.Settings()

		subject, body, err := resolveEmailTemplate(app, token, settings.Meta.VerificationTemplate)
		if err != nil {
			return err
		}

		return e.MailClient.Send(
			mail.Address{
				Name:    settings.Meta.SenderName,
				Address: settings.Meta.SenderAddress,
			},
			mail.Address{Address: e.Record.GetString(schema.FieldNameEmail)},
			subject,
			body,
			nil,
		)
	})

	if sendErr == nil {
		app.OnMailerAfterRecordVerificationSend().Trigger(event)
	}

	return sendErr
}

// SendUserChangeEmail sends a change email confirmation email to the specified user.
func SendRecordChangeEmail(app core.App, record *models.Record, newEmail string) error {
	token, tokenErr := tokens.NewRecordChangeEmailToken(app, record, newEmail)
	if tokenErr != nil {
		return tokenErr
	}

	mailClient := app.NewMailClient()

	event := &core.MailerRecordEvent{
		MailClient: mailClient,
		Record:     record,
		Meta: map[string]any{
			"token":    token,
			"newEmail": newEmail,
		},
	}

	sendErr := app.OnMailerBeforeRecordChangeEmailSend().Trigger(event, func(e *core.MailerRecordEvent) error {
		settings := app.Settings()

		subject, body, err := resolveEmailTemplate(app, token, settings.Meta.ConfirmEmailChangeTemplate)
		if err != nil {
			return err
		}

		return e.MailClient.Send(
			mail.Address{
				Name:    settings.Meta.SenderName,
				Address: settings.Meta.SenderAddress,
			},
			mail.Address{Address: newEmail},
			subject,
			body,
			nil,
		)
	})

	if sendErr == nil {
		app.OnMailerAfterRecordChangeEmailSend().Trigger(event)
	}

	return sendErr
}

func resolveEmailTemplate(
	app core.App,
	token string,
	emailTemplate core.EmailTemplate,
) (subject string, body string, err error) {
	settings := app.Settings()

	subject, rawBody, _ := emailTemplate.Resolve(
		settings.Meta.AppName,
		settings.Meta.AppUrl,
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
