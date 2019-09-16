package main

// EnvConfig 環境変数
type EnvConfig struct {
	// BotToken is bot user token to access to slack API.
	BotToken string `envconfig:"BOT_TOKEN" required:"true"`

	// BotID is bot user ID.
	BotID string `envconfig:"BOT_ID" required:"true"`

	// HotpepperKey is HotPepper API KEY
	HotpepperKey string `envconfig:"HOTPEPPER_KEY" required:"true"`
}
