package templates

// Available variables:
//
// ```
// Record        *models.Record
// AppName       string
// AppUrl        string
// ProviderNames []string
// ```
const PasswordLoginAlertBody = `
{{define "content"}}
	<p>Hello,</p>
	<p>
		Just to let you know that someone has logged in to your {{.AppName}} account using a password while you already have
		OAuth2
		{{range $index, $provider := .ProviderNames }}
			{{if $index}}|{{end}}
		    {{ $provider }}
		{{ end }}
		auth linked.
	</p>
	<p>If you have recently signed in with a password, you may disregard this email.</p>
	<p><strong>If you don't recognize the above action, you should immediately change your {{.AppName}} account password.</strong></p>
	<p>
	  Thanks,<br/>
	  {{.AppName}} team
	</p>
{{end}}
`
