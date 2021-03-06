package duckmail

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	DiscordDMBody    = "Hey! T'as un ptit courrier sympa dans la boîte aux lettres, peut-être une bonne nouvelle :)"
	DiscordNotifName = "Discord"
)

type DiscordNotification struct {
	Session     *discordgo.Session
	SendTimeout time.Duration

	nextSend time.Time
}

func (m *DiscordNotification) Name() string {
	return DiscordNotifName
}

func (d *DiscordNotification) Send(p Person) error {
	// No notification if no discord ID specified
	if p.DiscordID == "" {
		return nil
	}

	// Do not resend a notification before timeout expired
	if time.Now().Before(d.nextSend) {
		return nil
	}

	user, err := d.Session.User(p.DiscordID)
	if err != nil {
		return err
	}
	userCh, err := d.Session.UserChannelCreate(user.ID)
	if err != nil {
		return err
	}
	_, err = d.Session.ChannelMessageSend(userCh.ID, DiscordDMBody)
	if err != nil {
		return err
	}

	// Update next send time
	d.nextSend = time.Now().Add(d.SendTimeout)

	log.Printf("Notification sent to Discord (User: %v)", user.Username)
	return nil
}
