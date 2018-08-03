package duckmail

import (
	"fmt"
	"gopkg.in/mailgun/mailgun-go.v1"
)

const (
	MailgunSender  = "noreply@duckmail.com"
	MailgunSubject = "DuckMail - You have a new mail"
	MailgunBody    = `Hey!
Check your mail inbox to retrieve your (IRL) mail!
Maybe a good news... :)

Duckmail, a duck who likes to check mails
`
)

type MailgunNotification struct {
	Mailgun mailgun.Mailgun
}

func (m *MailgunNotification) Send(p Person) error {
	msg := m.Mailgun.NewMessage(MailgunSender, MailgunSubject, MailgunBody, p.Email)
	resp, id, err := m.Mailgun.Send(msg)
	if err != nil {
		return err
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)
	return nil

}
