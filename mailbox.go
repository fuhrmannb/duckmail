package duckmail

import (
	"log"
	"time"

	movingaverage "github.com/RobinUS2/golang-moving-average"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
)

type Mailbox struct {
	LED                     *DuckLed
	LEDNotifDuration        time.Duration
	LEDNotifPushingInterval time.Duration

	LDR        *aio.AnalogSensorDriver
	LDRTrigger struct {
		Min int
		Max int
	}
	LDRWindowSize int

	Notifiers []Notifier
	Person    Person

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
		m.ldrValue = movingaverage.New(m.LDRWindowSize)
	}
	m.ldrValue.Add(float64(s.(int)))

	if m.ldrValue.Avg() < float64(m.LDRTrigger.Min) {
		if m.ledTicker == nil {
			log.Printf("Mail received for %v (LDR value: %v)\n", m.Person.Name, m.ldrValue.Avg())

			// Blink LED
			m.ledTicker = m.blinkLed()

			// Notify Mailgun/Discord/...
			for _, not := range m.Notifiers {
				err := not.Send(m.Person)
				if err != nil {
					log.Printf("Error during %v notification: %v", not.Name(), err)
				}
			}
		}
	} else if m.ldrValue.Avg() > float64(m.LDRTrigger.Max) {
		if m.ledTicker != nil {
			log.Printf("Mailbox %v now empty (LDR value: %v)\n", m.Person.Name, m.ldrValue.Avg())
			m.ledTicker.Stop()
			m.ledTicker = nil
			m.LED.Off()
		}
	}
}

func (m *Mailbox) blinkLed() *time.Ticker {
	ledRange := int(512 * m.LEDNotifPushingInterval / m.LEDNotifDuration)
	m.ledSum = 0
	t := gobot.Every(m.LEDNotifPushingInterval, func() {
		m.ledSum = (m.ledSum + ledRange) % 512
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
