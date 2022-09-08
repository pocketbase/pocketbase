package mails

import (
	"net/mail"

	"github.com/pocketbase/pocketbase/core"
)

// SendCustomEmail sends an email to the email address with the template string (html/template) as html body with the PocketBase mail client.
// The hooks OnMailerBeforeCustomEmailSend and OnMailerAfterCustomEmailSend are triggered.
func SendCustomEmail(app core.App, email string, title string, template string, data any) error {
	mailClient := app.NewMailClient()

	event := &core.MailerCustomEvent{
		MailClient: mailClient,
		Email:      email,
		Title:      title,
		Meta:       data,
	}

	sendErr := app.OnMailerBeforeCustomEmailSend().Trigger(event, func(e *core.MailerCustomEvent) error {
		// resolve body template
		body, renderErr := resolveTemplateContent(data, template)
		if renderErr != nil {
			return renderErr
		}

		return e.MailClient.Send(
			mail.Address{
				Name:    app.Settings().Meta.SenderName,
				Address: app.Settings().Meta.SenderAddress,
			},
			mail.Address{Address: email},
			title,
			body,
			nil,
		)
	})

	if sendErr == nil {
		app.OnMailerAfterCustomEmailSend().Trigger(event)
	}

	return sendErr
}
