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
const UserPasswordResetBody = `
{{define "content"}}
	<p>Hello,</p>
	<p>Click on the button below to reset your password.</p>
	<p>
		<a class="btn" href="{{.ActionUrl}}">Reset password</a>
		<a class="fallback-link" href="{{.ActionUrl}}">{{.ActionUrl}}</a>
	</p>
	<p><i>If you didn’t ask to reset your password, you can ignore this email.</i></p>
	<p>
		Thanks,<br/>
		{{.AppName}} team
	</p>
{{end}}
`
