package email

import (
	"bytes"
	"fmt"
	"net/smtp"
	"os"
	"text/template"
	"time"
)

func SendMail(subject string, email string, body string) error {
	password := os.Getenv("SMTP_PASSWORD")
	emailAcc := os.Getenv("EMAILACC")
	start := time.Now()
	auth := smtp.PlainAuth(
		"",
		emailAcc,
		password,
		"smtp.gmail.com",
	)

	msg := "Subject: " + subject + "\n" + body
	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		emailAcc,
		[]string{email},
		[]byte(msg),
	)
	if err != nil {
		return fmt.Errorf("SendMail failed: %w", err)
	}

	fmt.Printf("SendMail done in %s\n", time.Since(start))
	return nil
}
func SendHTML(subject string, email string, name string) error {
	password := os.Getenv("SMTP_PASSWORD")
	emailAcc := os.Getenv("EMAILACC")
	start := time.Now()

	templatePath := "./utils/email/email.html"

	var body bytes.Buffer
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	if err := t.Execute(&body, struct{ Name string }{Name: name}); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	auth := smtp.PlainAuth(
		"",
		emailAcc,
		password, // Make sure to use an application-specific password here
		"smtp.gmail.com",
	)

	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"
	msg := "Subject: " + subject + "\n" + headers + "\n\n" + body.String()
	err = smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		emailAcc,
		[]string{email},
		[]byte(msg),
	)
	if err != nil {
		return fmt.Errorf("SendMail failed: %w", err)
	}

	fmt.Printf("SendHTML done in %s\n", time.Since(start))
	return nil
}
