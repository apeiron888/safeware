package email

import (
	"bytes"
	"fmt"
	"html/template"
)

type EmailService struct {
	smtpHost    string
	smtpPort    int
	smtpUser    string
	smtpPass    string
	fromAddress string
	frontendURL string
}

func NewEmailService(host string, port int, user, pass, from, frontendURL string) *EmailService {
	return &EmailService{
		smtpHost:    host,
		smtpPort:    port,
		smtpUser:    user,
		smtpPass:    pass,
		fromAddress: from,
		frontendURL: frontendURL,
	}
}

func (e *EmailService) SendEmail(to, subject, body string) error {
	// For MVP, we'll just log the email instead of actually sending
	// In production, uncomment the actual SMTP code below
	fmt.Printf("\n=== EMAIL ===\nTo: %s\nSubject: %s\nBody:\n%s\n=============\n", to, subject, body)
	return nil

	/* PRODUCTION CODE - Uncomment when ready:
	auth := smtp.PlainAuth("", e.smtpUser, e.smtpPass, e.smtpHost)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" + body + "\r\n")

	addr := fmt.Sprintf("%s:%d", e.smtpHost, e.smtpPort)
	return smtp.SendMail(addr, auth, e.fromAddress, []string{to}, msg)
	*/
}

func (e *EmailService) SendVerificationEmail(to, token string) error {
	verifyURL := fmt.Sprintf("%s/verify-email?token=%s", e.frontendURL, token)

	tmpl := `
	<html>
	<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
		<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
			<h2 style="color: #0ea5e9;">Welcome to SafeWare!</h2>
			<p>Thank you for registering. Please verify your email address by clicking the button below:</p>
			<div style="margin: 30px 0;">
				<a href="{{.VerifyURL}}" style="background-color: #0ea5e9; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; display: inline-block;">Verify Email</a>
			</div>
			<p style="color: #666; font-size: 14px;">If the button doesn't work, copy and paste this link into your browser:</p>
			<p style="color: #0ea5e9; font-size: 14px;">{{.VerifyURL}}</p>
			<p style="color: #999; font-size: 12px; margin-top: 30px;">This link will expire in 24 hours.</p>
		</div>
	</body>
	</html>
	`

	t := template.Must(template.New("verify").Parse(tmpl))
	var body bytes.Buffer
	if err := t.Execute(&body, map[string]string{"VerifyURL": verifyURL}); err != nil {
		return err
	}

	return e.SendEmail(to, "Verify your SafeWare account", body.String())
}

func (e *EmailService) SendPasswordResetEmail(to, token string) error {
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", e.frontendURL, token)

	tmpl := `
	<html>
	<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
		<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
			<h2 style="color: #0ea5e9;">Password Reset Request</h2>
			<p>We received a request to reset your password. Click the button below to create a new password:</p>
			<div style="margin: 30px 0;">
				<a href="{{.ResetURL}}" style="background-color: #0ea5e9; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; display: inline-block;">Reset Password</a>
			</div>
			<p style="color: #666; font-size: 14px;">If the button doesn't work, copy and paste this link into your browser:</p>
			<p style="color: #0ea5e9; font-size: 14px;">{{.ResetURL}}</p>
			<p style="color: #999; font-size: 12px; margin-top: 30px;">This link will expire in 1 hour. If you didn't request this, please ignore this email.</p>
		</div>
	</body>
	</html>
	`

	t := template.Must(template.New("reset").Parse(tmpl))
	var body bytes.Buffer
	if err := t.Execute(&body, map[string]string{"ResetURL": resetURL}); err != nil {
		return err
	}

	return e.SendEmail(to, "Reset your SafeWare password", body.String())
}
