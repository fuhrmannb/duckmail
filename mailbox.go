package duckmail

import (
	"fmt"
	"github.com/RobinUS2/golang-moving-average"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"log"
	"time"
)

const (
	LDRWindowSize           = 10
	LEDNotifDuration        = 2 * time.Second
	LEDNotifPushingInterval = 10 * time.Millisecond
	LEDRange                = int(512 * LEDNotifPushingInterval / LEDNotifDuration)
)

type Mailbox struct {
	LED        *DuckLed
	LDR        *aio.AnalogSensorDriver
	LDRTrigger int
	MailNotif  *MailgunNotification
	Person     Person

	ledTicker *time.Ticker
	ledSum    int
	ldrValue  *movingaverage.MovingAverage
}

func (m *Mailbox) Start() {
	// Check sensor
	m.LDR.On(aio.Data, m.onLDRValue)
}

func (m *Mailbox) onLDRValue(s interface{}) {
	if m.ldrValue == nil {
		m.ldrValue = movingaverage.New(LDRWindowSize)
	}
	m.ldrValue.Add(float64(s.(int)))

	if m.ldrValue.Avg() < float64(m.LDRTrigger) {
		if m.ledTicker == nil {
			fmt.Printf("Mail received! (LDR value: %v)\n", m.ldrValue.Avg())
			m.ledTicker = m.blinkLed()
			err := m.MailNotif.Send(m.Person)
			if err != nil {
				log.Printf("Error during Mailgun email sending: %v", err)
			}
		}
	} else {
		if m.ledTicker != nil {
			fmt.Printf("Email empty (LDR value: %v)\n", m.ldrValue.Avg())
			m.ledTicker.Stop()
			m.ledTicker = nil
			m.LED.Off()
		}
	}
}

func (m *Mailbox) blinkLed() *time.Ticker {
	m.ledSum = 0
	t := gobot.Every(LEDNotifPushingInterval, func() {
		m.ledSum = (m.ledSum + LEDRange) % 512
		var brightness uint8
		if m.ledSum > 255 {
			brightness = uint8(511 - m.ledSum)
		} else {
			brightness = uint8(m.ledSum)
		}
		m.LED.Brightness(brightness)
	})
	return t
}
