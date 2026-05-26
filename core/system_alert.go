package core

import (
	"bytes"
	"errors"
	"html"
	"html/template"
	"net/mail"

	"github.com/pocketbase/pocketbase/tools/mailer"
)

const systemAlertHTML = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    <meta name="viewport" content="width=device-width,initial-scale=1" />
    <style>
        body, html {
            padding: 0;
            margin: 0;
            border: 0;
            color: #16161a;
            background: #fff;
            font-size: 14px;
            line-height: 20px;
            font-weight: normal;
            font-family: Source Sans Pro, sans-serif, emoji;
        }
        body {
            padding: 20px 30px;
        }
        p {
            display: block;
            margin: 10px 0;
            font-family: inherit;
        }
        small {
            font-size: 12px;
            line-height: 16px;
        }
        strong {
            font-weight: bold;
        }
        em, i {
            font-style: italic;
        }
        a {
            color: inherit;
        }
        .alert {
        	padding: 15px;
        	background: #e4e8ec;
        	border-radius: 5px;
        	white-space: pre-wrap;
        }
    </style>
</head>
<body>
	<p>{{.AppName}} system alert occurred:</p>
	<p class="alert"><strong>{{.AlertDetails}}</strong></p>
	<p>For more information you could explore the logs in the dashboard of your application.</p>
</body>
</html>`

// sendSystemAlertToAllSuperusers sends a system error alert to all superusers.
//
// note: unexported for now until there is clarity around the planned log level alerts.
func sendSystemAlertToAllSuperusers(app App, subject string, details string) error {
	superusers, err := app.FindAllRecords(CollectionNameSuperusers)
	if err != nil {
		return err
	}

	var alertErrors []error
	for _, superuser := range superusers {
		err := sendSystemAlert(app, superuser, subject, details)
		if err != nil {
			alertErrors = append(alertErrors, err)
		}
	}

	return errors.Join(alertErrors...)
}

// sendSystemAlert sends a system error alert to a single superuser.
//
// note: unexported for now until there is clarity around the planned log level alerts.
func sendSystemAlert(app App, superuser *Record, subject string, details string) error {
	if !superuser.IsSuperuser() {
		return errors.New("system alerts can be sent only to superusers")
	}

	if subject == "" || details == "" {
		return errors.New("system alerts subject and details are required")
	}

	data := struct {
		AppName      string
		AlertDetails string
	}{
		AppName:      app.Settings().Meta.AppName,
		AlertDetails: details,
	}

	tpl := template.New("system_alert")

	var parseErr error
	tpl, parseErr = tpl.Parse(systemAlertHTML)
	if parseErr != nil {
		return parseErr
	}

	var buff bytes.Buffer
	executeErr := tpl.Execute(&buff, data)
	if executeErr != nil {
		return executeErr
	}

	message := &mailer.Message{
		From: mail.Address{
			Name:    app.Settings().Meta.SenderName,
			Address: app.Settings().Meta.SenderAddress,
		},
		To:      []mail.Address{{Address: superuser.Email()}},
		Subject: "[" + app.Settings().Meta.AppName + " system alert] " + html.EscapeString(subject),
		HTML:    buff.String(),
	}

	return app.NewMailClient().Send(message)
}
