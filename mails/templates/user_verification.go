package templates

// Available variables:
//
// ```
// User      *models.User
// AppName   string
// AppUrl    string
// Token     string
// ActionUrl string
// ```
const UserVerificationBody = `
{{define "content"}}
	<p>Hello,</p>
	<p>Thank you for joining us at {{.AppName}}.</p>
	<p>Click on the button below to verify your email address.</p>
	<p>
		<a class="btn" href="{{.ActionUrl}}" target="_blank" rel="noopener">Verify</a>
		<a class="fallback-link" href="{{.ActionUrl}}" target="_blank" rel="noopener">{{.ActionUrl}}</a>
	</p>
	<p>
		Thanks,<br/>
		{{.AppName}} team
	</p>
{{end}}
`
