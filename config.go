package duckmail

import "time"

type RootCfg struct {
	Arduino   ArduinoCfg
	Mailgun   MailgunCfg
	Mailboxes []MailboxCfg
}

type ArduinoCfg struct {
	Path                    string
	LDRPollingInterval      time.Duration `mapstructure:"ldr_polling_interval"`
	LDRWindowSize           int           `mapstructure:"ldr_window_size"`
	LEDNotifDuration        time.Duration `mapstructure:"led_notif_duration"`
	LEDNotifPushingInterval time.Duration `mapstructure:"led_notif_pushing_interval"`
}

type MailgunCfg struct {
	Domain        string
	PrivateKey    string        `mapstructure:"private_key"`
	PublicKey     string        `mapstructure:"public_key"`
	SenderAddress string        `mapstructure:"sender_address"`
	SendTimeout   time.Duration `mapstructure:"send_timeout"`
}

type MailboxArduinoCfg struct {
	LedPin     string `mapstructure:"led_pin"`
	LDRPin     string `mapstructure:"ldr_pin"`
	LDRTrigger int    `mapstructure:"ldr_trigger"`
}

type MailboxCfg struct {
	Person  Person
	Arduino MailboxArduinoCfg
}
