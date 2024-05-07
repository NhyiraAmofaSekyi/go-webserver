package utils

import (
	"bytes"
	"fmt"
	"net/smtp"
	"os"
	"text/template"
	"time"
)

func SendMail() error {
	password := os.Getenv("SMTP_PASSWORD")
	start := time.Now()
	auth := smtp.PlainAuth(
		"",
		"nhyiraamofasekyi@gmail.com",
		password,
		"smtp.gmail.com",
	)

	msg := "Subject: My special subject\nbody"
	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"nhyiraamofasekyi@gmail.com",
		[]string{"nhyiraamofasekyi@gmail.com"},
		[]byte(msg),
	)
	if err != nil {
		return fmt.Errorf("SendMail failed: %w", err)
	}

	fmt.Printf("SendMail done in %s\n", time.Since(start))
	return nil
}
func SendHTML(subject string) error {
	password := os.Getenv("SMTP_PASSWORD")
	start := time.Now()

	templatePath := "./email.html"

	var body bytes.Buffer
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	if err := t.Execute(&body, struct{ Name string }{Name: "robby"}); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	auth := smtp.PlainAuth(
		"",
		"nhyiraamofasekyi@gmail.com",
		password, // Make sure to use an application-specific password here
		"smtp.gmail.com",
	)

	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"
	msg := "Subject: " + subject + "\n" + headers + "\n\n" + body.String()
	err = smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"nhyiraamofasekyi@gmail.com",
		[]string{"nhyiraamofasekyi@gmail.com"},
		[]byte(msg),
	)
	if err != nil {
		return fmt.Errorf("SendMail failed: %w", err)
	}

	fmt.Printf("SendHTML done in %s\n", time.Since(start))
	return nil
}
