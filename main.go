package main

// 参考 https://github.com/tcnksm/go-slack-interactive

import (
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/nlopes/slack"
)

type envConfig struct {
	// BotToken is bot user token to access to slack API.
	BotToken string `envconfig:"BOT_TOKEN" required:"true"`

	// BotID is bot user ID.
	BotID string `envconfig:"BOT_ID" required:"true"`
}

func main() {
	os.Exit(_main(os.Args[1:]))
}

func _main(args []string) int {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		return 1
	}

	client := slack.New(
		env.BotToken,
		// slack.OptionDebug(true),
		// slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)
	slackListener := &SlackListener{
		Client:    client,
		BotUserID: env.BotID,
	}

	slackListener.ListenAndResponse()
	return 0
}
