package slackbot

import (
	"github.com/lestrrat/go-pdebug"
	"github.com/nlopes/slack"
)

var slackClient *slack.Client

func init() {
	var token string
	if err := readEnvConfigFile("Slack API token", "SLACKBOT_API_TOKEN_FILE", &token); err != nil {
		panic(err)
	}
	if token == "" {
		panic("token is empty")
	}

	pdebug.Printf("token = '%s'", token)
	slackClient = slack.New(token)
}

// Dummy for now
func Run(_ string) error {
	return Watch()
}
