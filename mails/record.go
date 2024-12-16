package mails

import (
	"html"
	"html/template"
	"net/mail"
	"slices"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/mails/templates"
	"github.com/pocketbase/pocketbase/tools/mailer"
)

// SendRecordAuthAlert sends a new device login alert to the specified auth record.
func SendRecordAuthAlert(app core.App, authRecord *core.Record) error {
	mailClient := app.NewMailClient()

	subject, body, err := resolveEmailTemplate(app, authRecord, authRecord.Collection().AuthAlert.EmailTemplate, nil)
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
	event.App = app
	event.Mailer = mailClient
	event.Message = message
	event.Record = authRecord

	return app.OnMailerRecordAuthAlertSend().Trigger(event, func(e *core.MailerRecordEvent) error {
		return e.Mailer.Send(e.Message)
	})
}

// SendRecordOTP sends OTP email to the specified auth record.
//
// This method will also update the "sentTo" field of the related OTP record to the mail sent To address (if the OTP exists and not already assigned).
func SendRecordOTP(app core.App, authRecord *core.Record, otpId string, pass string) error {
	mailClient := app.NewMailClient()

	subject, body, err := resolveEmailTemplate(app, authRecord, authRecord.Collection().OTP.EmailTemplate, map[string]any{
		core.EmailPlaceholderOTPId: otpId,
		core.EmailPlaceholderOTP:   pass,
	})
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
	event.App = app
	event.Mailer = mailClient
	event.Message = message
	event.Record = authRecord
	event.Meta = map[string]any{
		"otpId":    otpId,
		"password": pass,
	}

	return app.OnMailerRecordOTPSend().Trigger(event, func(e *core.MailerRecordEvent) error {
		err := e.Mailer.Send(e.Message)
		if err != nil {
			return err
		}

		var toAddress string
		if len(e.Message.To) > 0 {
			toAddress = e.Message.To[0].Address
		}
		if toAddress == "" {
			return nil
		}

		otp, err := e.App.FindOTPById(otpId)
		if err != nil {
			e.App.Logger().Warn(
				"Unable to find OTP to update its sentTo field (either it was already deleted or the id is nonexisting)",
				"error", err,
				"otpId", otpId,
			)
			return nil
		}

		if otp.SentTo() != "" {
			return nil // was already sent to another target
		}

		otp.SetSentTo(toAddress)
		if err = e.App.Save(otp); err != nil {
			e.App.Logger().Error(
				"Failed to update OTP sentTo field",
				"error", err,
				"otpId", otpId,
				"to", toAddress,
			)
		}

		return nil
	})
}

// SendRecordPasswordReset sends a password reset request email to the specified auth record.
func SendRecordPasswordReset(app core.App, authRecord *core.Record) error {
	token, tokenErr := authRecord.NewPasswordResetToken()
	if tokenErr != nil {
		return tokenErr
	}

	mailClient := app.NewMailClient()

	subject, body, err := resolveEmailTemplate(app, authRecord, authRecord.Collection().ResetPasswordTemplate, map[string]any{
		core.EmailPlaceholderToken: token,
	})
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
	event.App = app
	event.Mailer = mailClient
	event.Message = message
	event.Record = authRecord
	event.Meta = map[string]any{"token": token}

	return app.OnMailerRecordPasswordResetSend().Trigger(event, func(e *core.MailerRecordEvent) error {
		return e.Mailer.Send(e.Message)
	})
}

// SendRecordVerification sends a verification request email to the specified auth record.
func SendRecordVerification(app core.App, authRecord *core.Record) error {
	token, tokenErr := authRecord.NewVerificationToken()
	if tokenErr != nil {
		return tokenErr
	}

	mailClient := app.NewMailClient()

	subject, body, err := resolveEmailTemplate(app, authRecord, authRecord.Collection().VerificationTemplate, map[string]any{
		core.EmailPlaceholderToken: token,
	})
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
	event.App = app
	event.Mailer = mailClient
	event.Message = message
	event.Record = authRecord
	event.Meta = map[string]any{"token": token}

	return app.OnMailerRecordVerificationSend().Trigger(event, func(e *core.MailerRecordEvent) error {
		return e.Mailer.Send(e.Message)
	})
}

// SendRecordChangeEmail sends a change email confirmation email to the specified auth record.
func SendRecordChangeEmail(app core.App, authRecord *core.Record, newEmail string) error {
	token, tokenErr := authRecord.NewEmailChangeToken(newEmail)
	if tokenErr != nil {
		return tokenErr
	}

	mailClient := app.NewMailClient()

	subject, body, err := resolveEmailTemplate(app, authRecord, authRecord.Collection().ConfirmEmailChangeTemplate, map[string]any{
		core.EmailPlaceholderToken: token,
	})
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
	event.App = app
	event.Mailer = mailClient
	event.Message = message
	event.Record = authRecord
	event.Meta = map[string]any{
		"token":    token,
		"newEmail": newEmail,
	}

	return app.OnMailerRecordEmailChangeSend().Trigger(event, func(e *core.MailerRecordEvent) error {
		return e.Mailer.Send(e.Message)
	})
}

var nonescapeTypes = []string{
	core.FieldTypeAutodate,
	core.FieldTypeDate,
	core.FieldTypeBool,
	core.FieldTypeNumber,
}

func resolveEmailTemplate(
	app core.App,
	authRecord *core.Record,
	emailTemplate core.EmailTemplate,
	placeholders map[string]any,
) (subject string, body string, err error) {
	if placeholders == nil {
		placeholders = map[string]any{}
	}

	// register default system placeholders
	if _, ok := placeholders[core.EmailPlaceholderAppName]; !ok {
		placeholders[core.EmailPlaceholderAppName] = app.Settings().Meta.AppName
	}
	if _, ok := placeholders[core.EmailPlaceholderAppURL]; !ok {
		placeholders[core.EmailPlaceholderAppURL] = app.Settings().Meta.AppURL
	}

	// register default auth record placeholders
	for _, field := range authRecord.Collection().Fields {
		if field.GetHidden() {
			continue
		}

		fieldPlacehodler := "{RECORD:" + field.GetName() + "}"
		if _, ok := placeholders[fieldPlacehodler]; !ok {
			val := authRecord.GetString(field.GetName())

			// note: the escaping is not strictly necessary but for just in case
			// the user decide to store and render the email as plain html
			if !slices.Contains(nonescapeTypes, field.Type()) {
				val = html.EscapeString(val)
			}

			placeholders[fieldPlacehodler] = val
		}
	}

	subject, rawBody := emailTemplate.Resolve(placeholders)

	params := struct {
		HTMLContent template.HTML
	}{
		HTMLContent: template.HTML(rawBody),
	}

	body, err = resolveTemplateContent(params, templates.Layout, templates.HTMLBody)
	if err != nil {
		return "", "", err
	}

	return subject, body, nil
}
