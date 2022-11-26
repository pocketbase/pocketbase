package templates

// Available variables:
//
// ```
// Admin     *models.Admin
// AppName   string
// AppUrl    string
// Token     string
// ActionUrl string
// ```
const AdminPasswordResetBody = `
{{define "content"}}
	<p>Hello,</p>
	<p>Follow this link to reset your admin password for {{.AppName}}.</p>
	<p>
		<a class="btn" href="{{.ActionUrl}}" target="_blank" rel="noopener">Reset password</a>
	</p>
	<p><i>If you did not request to reset your password, please ignore this email and the link will expire on its own.</i></p>
{{end}}
`
