package settings

// Common settings placeholder tokens
const (
	EmailPlaceholderAppName   string = "{APP_NAME}"
	EmailPlaceholderAppUrl    string = "{APP_URL}"
	EmailPlaceholderToken     string = "{TOKEN}"
	EmailPlaceholderActionUrl string = "{ACTION_URL}"
)

var defaultVerificationTemplate = EmailTemplate{
	Subject: "Verify your " + EmailPlaceholderAppName + " email",
	Body: `<p>Hello,</p>
<p>Thank you for joining us at ` + EmailPlaceholderAppName + `.</p>
<p>Click on the button below to verify your email address.</p>
<p>
  <a class="btn" href="` + EmailPlaceholderActionUrl + `" target="_blank" rel="noopener">Verify</a>
</p>
<p>
  Thanks,<br/>
  ` + EmailPlaceholderAppName + ` team
</p>`,
	ActionUrl: EmailPlaceholderAppUrl + "/_/#/auth/confirm-verification/" + EmailPlaceholderToken,
}

var defaultResetPasswordTemplate = EmailTemplate{
	Subject: "Reset your " + EmailPlaceholderAppName + " password",
	Body: `<p>Hello,</p>
<p>Click on the button below to reset your password.</p>
<p>
  <a class="btn" href="` + EmailPlaceholderActionUrl + `" target="_blank" rel="noopener">Reset password</a>
</p>
<p><i>If you didn't ask to reset your password, you can ignore this email.</i></p>
<p>
  Thanks,<br/>
  ` + EmailPlaceholderAppName + ` team
</p>`,
	ActionUrl: EmailPlaceholderAppUrl + "/_/#/auth/confirm-password-reset/" + EmailPlaceholderToken,
}

var defaultConfirmEmailChangeTemplate = EmailTemplate{
	Subject: "Confirm your " + EmailPlaceholderAppName + " new email address",
	Body: `<p>Hello,</p>
<p>Click on the button below to confirm your new email address.</p>
<p>
  <a class="btn" href="` + EmailPlaceholderActionUrl + `" target="_blank" rel="noopener">Confirm new email</a>
</p>
<p><i>If you didn't ask to change your email address, you can ignore this email.</i></p>
<p>
  Thanks,<br/>
  ` + EmailPlaceholderAppName + ` team
</p>`,
	ActionUrl: EmailPlaceholderAppUrl + "/_/#/auth/confirm-email-change/" + EmailPlaceholderToken,
}
