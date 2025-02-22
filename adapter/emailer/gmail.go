package emailer

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/jordan-wright/email"
)

const (
	_smtpAuthAddress   = "smtp.gmail.com"
	_smtpServerAddress = "smtp.gmail.com:587"
)

// Gmail struct (Actual Implementation)
type Gmail struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

// NewGmail creates a new instance of Gmail service
func NewGmail(name, email, password string) *Gmail {
	return &Gmail{
		name:              name,
		fromEmailAddress:  email,
		fromEmailPassword: password,
	}
}

// Send sends an email with the specified details
func (sender *Gmail) Send(
	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attachFiles []string,
) error {

	mail := email.NewEmail()

	mail.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddress)
	mail.Subject = subject
	mail.HTML = []byte(content)
	mail.To = to
	mail.Cc = cc
	mail.Bcc = bcc

	for _, file := range attachFiles {
		fileReader, err := os.Open(file)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", file, err)
		}
		defer fileReader.Close()

		_, err = mail.Attach(fileReader, file, "application/octet-stream")
		if err != nil {
			return fmt.Errorf("failed to attach file: %w", err)
		}
	}

	smtpAuth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, _smtpAuthAddress)

	return mail.Send(_smtpServerAddress, smtpAuth)
}

// SenderEmail returns the sender's email address
func (sender *Gmail) SenderEmail() string {
	return sender.fromEmailAddress
}
