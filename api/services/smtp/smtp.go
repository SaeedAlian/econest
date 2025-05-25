package smtp

import (
	"fmt"
	"log"
	"net/smtp"
)

type SMTPServer struct {
	mail string
	pass string
	host string
	port string
}

func NewSMTPServer(host string, port string, mail string, pass string) *SMTPServer {
	return &SMTPServer{mail: mail, pass: pass, host: host, port: port}
}

func (s *SMTPServer) SendMail(to string, subject string, body string) error {
	log.Printf("Sending mail at %s:%s ...", s.host, s.port)

	headers := make(map[string]string)
	headers["From"] = s.mail
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"UTF-8\""

	messageString := ""
	for k, v := range headers {
		messageString += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	messageString += "\r\n" + body

	message := []byte(messageString)

	auth := smtp.PlainAuth("", s.mail, s.pass, s.host)

	err := smtp.SendMail(s.host+":"+s.port, auth, s.mail, []string{to}, message)
	if err != nil {
		return err
	}

	return nil
}

func (s *SMTPServer) SendEmailVerificationRequestMail(
	userFullName string,
	userEmail string,
	verificationLink string,
	websiteName string,
	websiteUrl string,
	expirationInMinutes int,
) error {
	return s.SendMail(
		userEmail,
		fmt.Sprintf("%s: Email Verification Request", websiteName),
		fmt.Sprintf(`<p>Hi %s</p>

<p>Thank you for signing up at %s</p>

<p>To complete your registration and activate your account, please verify your email address by clicking the button below:</p>

<a href="%s" style="display:inline-block;padding:10px 20px;background-color:#007BFF;color:white;text-decoration:none;border-radius:5px;">👉 Verify Email Address</a>

<p>If the button doesn’t work, copy and paste this link into your browser:</p>
<p>%s</p>

<p>This link will expire in %d minutes, so please verify your email as soon as possible.</p>

<p>If you didn’t create an account with us, you can safely ignore this message.</p>

<p>Thanks, The %s Team %s</p>
	`, userFullName, websiteName, verificationLink, verificationLink, expirationInMinutes, websiteName, websiteUrl),
	)
}

func (s *SMTPServer) SendPasswordResetRequestMail(
	userFullName string,
	userEmail string,
	resetLink string,
	websiteName string,
	websiteUrl string,
	expirationInMinutes int,
) error {
	return s.SendMail(
		userEmail,
		fmt.Sprintf("%s: Reset Password Request", websiteName),
		fmt.Sprintf(`
<p>Hi %s,</p>

<p>We received a request to reset your password for your %s account.</p>

<p>Click the button below to reset your password:</p>

<a href="%s" style="display:inline-block;padding:10px 20px;background-color:#007BFF;color:white;text-decoration:none;border-radius:5px;">🔐 Reset Password</a>

<p>If the button doesn’t work, copy and paste this link into your browser:</p>
<p>%s</p>

<p>This link will expire in %d minutes for your security.</p>

<p>If you didn’t request a password reset, please ignore this email or contact our support.</p>

<p>Thanks,<br>The %s Team</p>
	`, userFullName, websiteName, resetLink, resetLink, expirationInMinutes, websiteName),
	)
}
