package smtp

import (
	"log"
	"os"
	"testing"

	"github.com/SaeedAlian/econest/api/config"
)

func TestAuthHandler(t *testing.T) {
	if config.Env.Env != "test" {
		log.Panic("environment is not on test!!")
		os.Exit(1)
	}

	smtpServer := NewSMTPServer(
		config.Env.SMTPHost,
		config.Env.SMTPPort,
		config.Env.SMTPEmail,
		config.Env.SMTPPassword,
	)

	t.Run("should send email successfully", func(t *testing.T) {
		err := smtpServer.SendMail(
			config.Env.SMTPTestRecipientAddress,
			"test mail",
			"This is a test email",
		)
		if err != nil {
			t.Fatal(err)
		}
	})
}
