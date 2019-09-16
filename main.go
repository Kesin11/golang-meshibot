package main

// 参考 https://github.com/tcnksm/go-slack-interactive

import (
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/nlopes/slack"
)

func main() {
	os.Exit(_main(os.Args[1:]))
}

func _main(args []string) int {
	var env EnvConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		return 1
	}

	hotPepper := NewHotPepperClient(env.HotpepperKey)

	client := slack.New(
		env.BotToken,
		// slack.OptionDebug(true),
		// slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)
	slackListener := &SlackListener{
		client:           client,
		botUserID:        env.BotID,
		restaurantClient: hotPepper,
	}

	slackListener.ListenAndResponse()
	return 0
}
