package duckmail

type Person struct {
	Name  string
	Email string
	// Note: This is NOT the username (such as Azerty#1234) but the ID
	// More info to retrieve ID: https://support.discordapp.com/hc/en-us/articles/206346498-Where-can-I-find-my-User-Server-Message-ID-
	DiscordID string `mapstructure:"discord_id"`
}
