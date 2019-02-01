package duckmail

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
	mailgun "gopkg.in/mailgun/mailgun-go.v1"
)

func StartController(cfg *RootCfg) error {
	// Robot cfg
	robot := gobot.NewRobot("duckmail")

	// Arduino firmata cfg
	firmataAdaptor := firmata.NewAdaptor(cfg.Arduino.Path)
	firmataAdaptor.SetName("Arduino-Firmata")
	robot.AddConnection(firmataAdaptor)

	// Notifiers cfg
	// Mailgun
	mailNotif := &MailgunNotification{
		Mailgun:       mailgun.NewMailgun(cfg.Mailgun.Domain, cfg.Mailgun.PrivateKey, cfg.Mailgun.PublicKey),
		SenderAddress: cfg.Mailgun.SenderAddress,
		SendTimeout:   cfg.Mailgun.SendTimeout,
	}
	// Discord
	discordSession, err := discordgo.New("Bot " + cfg.Discord.Token)
	if err != nil {
		return err
	}
	err = discordSession.Open()
	if err != nil {
		return err
	}
	discordNotif := &DiscordNotification{
		Session:     discordSession,
		SendTimeout: cfg.Discord.SendTimeout,
	}
	// All in one slice!
	notifiers := []Notifier{mailNotif, discordNotif}

	// Mailbox creation
	var boxes []*Mailbox
	for _, boxCfg := range cfg.Mailboxes {
		box := &Mailbox{
			LED:                     &DuckLed{gpio.NewLedDriver(firmataAdaptor, boxCfg.Arduino.LedPin)},
			LEDNotifDuration:        cfg.Arduino.LEDNotifDuration,
			LEDNotifPushingInterval: cfg.Arduino.LEDNotifPushingInterval,
			LDR:                     aio.NewAnalogSensorDriver(firmataAdaptor, boxCfg.Arduino.LDRPin, cfg.Arduino.LDRPollingInterval),
			LDRTrigger:              boxCfg.Arduino.LDRTrigger,
			LDRWindowSize:           cfg.Arduino.LDRWindowSize,
			Notifiers:               notifiers,
			Person:                  boxCfg.Person,
		}
		// Configure LED
		box.LED.SetName(fmt.Sprintf("LED-%v-pin%v", boxCfg.Person.Name, box.LED.Pin()))
		robot.AddDevice(box.LED)

		// Configure LDR
		box.LDR.SetName(fmt.Sprintf("LDR-%v-pin%v", boxCfg.Person.Name, box.LDR.Pin()))
		robot.AddDevice(box.LDR)

		boxes = append(boxes, box)
	}

	// Start box when robot starts
	robot.Work = func() {
		for _, b := range boxes {
			b.Start()
		}
	}

	// Let's go friends!
	if err := robot.Start(); err != nil {
		return err
	}
	return nil
}
