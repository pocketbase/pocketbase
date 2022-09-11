package mails

import (
	"html/template"
	"net/mail"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/mails/templates"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tokens"
)

// SendUserPasswordReset sends a password reset request email to the specified user.
func SendUserPasswordReset(app core.App, user *models.User) error {
	token, tokenErr := tokens.NewUserResetPasswordToken(app, user)
	if tokenErr != nil {
		return tokenErr
	}

	mailClient := app.NewMailClient()

	event := &core.MailerUserEvent{
		MailClient: mailClient,
		User:       user,
		Meta:       map[string]any{"token": token},
	}

	sendErr := app.OnMailerBeforeUserResetPasswordSend().Trigger(event, func(e *core.MailerUserEvent) error {
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
			mail.Address{Address: e.User.Email},
			subject,
			body,
			nil,
		)
	})

	if sendErr == nil {
		app.OnMailerAfterUserResetPasswordSend().Trigger(event)
	}

	return sendErr
}

// SendUserVerification sends a verification request email to the specified user.
func SendUserVerification(app core.App, user *models.User) error {
	token, tokenErr := tokens.NewUserVerifyToken(app, user)
	if tokenErr != nil {
		return tokenErr
	}

	mailClient := app.NewMailClient()

	event := &core.MailerUserEvent{
		MailClient: mailClient,
		User:       user,
		Meta:       map[string]any{"token": token},
	}

	sendErr := app.OnMailerBeforeUserVerificationSend().Trigger(event, func(e *core.MailerUserEvent) error {
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
			mail.Address{Address: e.User.Email},
			subject,
			body,
			nil,
		)
	})

	if sendErr == nil {
		app.OnMailerAfterUserVerificationSend().Trigger(event)
	}

	return sendErr
}

// SendUserChangeEmail sends a change email confirmation email to the specified user.
func SendUserChangeEmail(app core.App, user *models.User, newEmail string) error {
	token, tokenErr := tokens.NewUserChangeEmailToken(app, user, newEmail)
	if tokenErr != nil {
		return tokenErr
	}

	mailClient := app.NewMailClient()

	event := &core.MailerUserEvent{
		MailClient: mailClient,
		User:       user,
		Meta: map[string]any{
			"token":    token,
			"newEmail": newEmail,
		},
	}

	sendErr := app.OnMailerBeforeUserChangeEmailSend().Trigger(event, func(e *core.MailerUserEvent) error {
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
		app.OnMailerAfterUserChangeEmailSend().Trigger(event)
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
