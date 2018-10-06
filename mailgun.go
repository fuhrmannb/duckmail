package duckmail

import (
	"gopkg.in/mailgun/mailgun-go.v1"
	"log"
	"time"
)

const (
	MailgunSubject = "DuckMail - You have a new mail"
	MailgunBody    = `Hey!
Check your mail inbox to retrieve your (IRL) mail!
Maybe a good news... :)

Duckmail, a duck who likes to check mails
`
)

type MailgunNotification struct {
	Mailgun       mailgun.Mailgun
	SenderAddress string
	SendTimeout   time.Duration

	nextSend time.Time
}

func (m *MailgunNotification) Send(p Person) error {
	// Do not resend a email before timeout expired
	if time.Now().Before(m.nextSend) {
		return nil
	}

	msg := m.Mailgun.NewMessage(m.SenderAddress, MailgunSubject, MailgunBody, p.Email)
	resp, id, err := m.Mailgun.Send(msg)
	if err != nil {
		return err
	}

	// Update next send time
	m.nextSend = time.Now().Add(m.SendTimeout)

	log.Printf("Email sent with Mailgun (ID: %s Resp: %s)", id, resp)
	return nil

}
