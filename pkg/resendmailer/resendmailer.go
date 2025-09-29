package resendmailer

import (
	"fmt"

	"github.com/resend/resend-go/v2"
)

type ResendMailer interface {
	Send(to, subject, body string)
}

type resendMailer struct {
	ApiKey string
	From   string
}

func NewResendMailer(apiKey, from string) ResendMailer {
	return &resendMailer{
		ApiKey: apiKey,
		From:   from,
	}
}

func (r *resendMailer) Send(to, subject, body string) {
	client := resend.NewClient(r.ApiKey)

	fmt.Println("Sending email to:", to)
	fmt.Println("Email subject:", subject)
	fmt.Println("Email body:", body)
	fmt.Println("from:", r.From)

	params := &resend.SendEmailRequest{
		From:    r.From,
		To:      []string{to},
		Html:    body,
		Subject: subject,
	}

	sent, err := client.Emails.Send(params)

	if err != nil {
		fmt.Println("Error sending email:", err)
		return
	}
	fmt.Println("Email sent successfully:", sent.Id)

}
