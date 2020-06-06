package handler

import (
	"net/http"
	"github.com/slack-go/slack"
	"sp-slack/logger"
)

func Subscribe(w http.ResponseWriter, r *http.Request) {
	logger.Info("subscribe endpoint hit")
	s, err := parseCommand(r)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Info(s.UserID)
}

func Unsubscribe(w http.ResponseWriter, r *http.Request) {
	logger.Info("unsubscribe endpoint hit")
	s, err := parseCommand(r)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Info(s.UserID)
}

func parseCommand(r *http.Request) (slack.SlashCommand, error) {
	if !isFromSlack(r) {
		return slack.SlashCommand{}, notFromSlack
	}

	return slack.SlashCommandParse(r)
}
