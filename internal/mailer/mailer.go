// Filename: internal/mailer/mailer.go
// Description: Mailer to send static email templates

package mailer

import (
	"bytes"
	"embed"
	"html/template"
	"time"

	"gopkg.in/mail.v2"
)

//go:embed "templates"
var templateFS embed.FS

type Mailer struct {
	dialer *mail.Dialer
	sender string
}

// Configuring a SMTP connection instance with credentials from mailtrap
func New(host string, port int, username, password, sender string) Mailer {
	dialer := mail.NewDialer(host, port, username, password)
	dialer.Timeout = 5 * time.Second
	return Mailer{
		dialer: dialer,
		sender: sender,
	}
}

// Sending the email to the user, using the data parameter to inject dynamic content into the template
func (m Mailer) Send(recipient, templateFile string, data any) error {
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)

	if err != nil {
		return err
	}

	// Filling in subject data
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)

	if err != nil {
		return err
	}

	// Filling in plaintext data
	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)

	if err != nil {
		return err
	}

	// Filling in HTML data
	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)

	if err != nil {
		return err
	}

	// Crafting the message from the parsed templates
	msg := mail.NewMessage()
	msg.SetHeader("From", m.sender)
	msg.SetHeader("To", recipient)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	// Sending the email message, 3 times before giving up
	for i := 0; i < 3; i++ {
		err = m.dialer.DialAndSend(msg)
		if err == nil {
			return nil // Successfully sent the email
		}

		// Wait 500 milliseconds before trying again
		time.Sleep(500 * time.Millisecond)
	}

	return err // Return the last error if all attempts fail
}
