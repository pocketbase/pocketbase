package mails

import (
	"net/mail"
	"strings"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/mails/templates"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tokens"
)

func prepareUserEmailBody(
	app core.App,
	user *models.User,
	token string,
	actionUrl string,
	bodyTemplate string,
) (string, error) {
	settings := app.Settings()

	// replace action url placeholder params (if any)
	actionUrlParams := map[string]string{
		core.EmailPlaceholderAppUrl: settings.Meta.AppUrl,
		core.EmailPlaceholderToken:  token,
	}
	for k, v := range actionUrlParams {
		actionUrl = strings.ReplaceAll(actionUrl, k, v)
	}
	var urlErr error
	actionUrl, urlErr = normalizeUrl(actionUrl)
	if urlErr != nil {
		return "", urlErr
	}

	params := struct {
		AppName   string
		AppUrl    string
		User      *models.User
		Token     string
		ActionUrl string
	}{
		AppName:   settings.Meta.AppName,
		AppUrl:    settings.Meta.AppUrl,
		User:      user,
		Token:     token,
		ActionUrl: actionUrl,
	}

	return resolveTemplateContent(params, templates.Layout, bodyTemplate)
}

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
		body, err := prepareUserEmailBody(
			app,
			user,
			token,
			app.Settings().Meta.UserResetPasswordUrl,
			templates.UserPasswordResetBody,
		)
		if err != nil {
			return err
		}

		return e.MailClient.Send(
			mail.Address{
				Name:    app.Settings().Meta.SenderName,
				Address: app.Settings().Meta.SenderAddress,
			},
			mail.Address{Address: e.User.Email},
			("Reset your " + app.Settings().Meta.AppName + " password"),
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
		body, err := prepareUserEmailBody(
			app,
			user,
			token,
			app.Settings().Meta.UserVerificationUrl,
			templates.UserVerificationBody,
		)
		if err != nil {
			return err
		}

		return e.MailClient.Send(
			mail.Address{
				Name:    app.Settings().Meta.SenderName,
				Address: app.Settings().Meta.SenderAddress,
			},
			mail.Address{Address: e.User.Email},
			("Verify your " + app.Settings().Meta.AppName + " email"),
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
		body, err := prepareUserEmailBody(
			app,
			user,
			token,
			app.Settings().Meta.UserConfirmEmailChangeUrl,
			templates.UserConfirmEmailChangeBody,
		)
		if err != nil {
			return err
		}

		return e.MailClient.Send(
			mail.Address{
				Name:    app.Settings().Meta.SenderName,
				Address: app.Settings().Meta.SenderAddress,
			},
			mail.Address{Address: newEmail},
			("Confirm your " + app.Settings().Meta.AppName + " new email address"),
			body,
			nil,
		)
	})

	if sendErr == nil {
		app.OnMailerAfterUserChangeEmailSend().Trigger(event)
	}

	return sendErr
}
