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
const UserConfirmEmailChangeBody = `
{{define "content"}}
	<p>Hello,</p>
	<p>Click on the button below to confirm your new email address.</p>
	<p>
		<a class="btn" href="{{.ActionUrl}}" target="_blank" rel="noopener">Confirm new email</a>
		<a class="fallback-link" href="{{.ActionUrl}}" target="_blank" rel="noopener">{{.ActionUrl}}</a>
	</p>
	<p><i>If you didn’t ask to change your email address, you can ignore this email.</i></p>
	<p>
		Thanks,<br/>
		{{.AppName}} team
	</p>
{{end}}
`
