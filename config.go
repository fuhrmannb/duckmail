package duckmail

type RootCfg struct {
	Arduino   ArduinoCfg
	Mailgun   MailgunCfg
	Mailboxes []MailboxCfg
}

type ArduinoCfg struct {
	Path string
}

type MailgunCfg struct {
	Domain     string
	PrivateKey string `mapstructure:"private_key"`
	PublicKey  string `mapstructure:"public_key"`
}

type MailboxArduinoCfg struct {
	LedPin string `mapstructure:"led_pin"`
	LDRPin string `mapstructure:"ldr_pin"`
}

type MailboxCfg struct {
	Person  Person
	Arduino MailboxArduinoCfg
}
