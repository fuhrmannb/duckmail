package duckmail

import "gobot.io/x/gobot/drivers/gpio"

type DuckLed struct {
	*gpio.LedDriver
}

func (l *DuckLed) Halt() error {
	return l.Off()
}
