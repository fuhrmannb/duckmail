package duckmail

import (
	"fmt"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
	"gopkg.in/mailgun/mailgun-go.v1"
	"time"
)

const (
	LDRPollingInterval = 10 * time.Millisecond
)

func StartController(cfg *RootCfg) {
	// Robot cfg
	robot := gobot.NewRobot("duckmail")

	// Arduino firmata cfg
	firmataAdaptor := firmata.NewAdaptor(cfg.Arduino.Path)
	firmataAdaptor.SetName("Arduino-Firmata")
	robot.AddConnection(firmataAdaptor)

	// Mailgun cfg
	mailNotif := &MailgunNotification{
		Mailgun: mailgun.NewMailgun(cfg.Mailgun.Domain, cfg.Mailgun.PrivateKey, cfg.Mailgun.PublicKey),
	}

	// Mailbox creation
	var boxes []*Mailbox
	for _, boxCfg := range cfg.Mailboxes {
		box := &Mailbox{
			LED:       &DuckLed{gpio.NewLedDriver(firmataAdaptor, boxCfg.Arduino.LedPin)},
			LDR:       aio.NewAnalogSensorDriver(firmataAdaptor, boxCfg.Arduino.LDRPin, LDRPollingInterval),
			MailNotif: mailNotif,
			Person:    boxCfg.Person,
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
	robot.Start()
}
