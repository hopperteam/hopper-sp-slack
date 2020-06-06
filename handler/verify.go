package handler

import (
	"errors"
	"net/http"
	"github.com/slack-go/slack"
	"sp-slack/config"
	"sp-slack/logger"
	"sp-slack/utils"
)

var notFromSlack = errors.New("request did not originate from slack")

func isFromSlack(r *http.Request) bool {
	body, err := utils.CpReqBodyToStr(r)
	if err != nil {
		logger.Error(err)
		return false
	}

	verifier, err := slack.NewSecretsVerifier(r.Header, config.Secret)
	if err != nil {
		logger.Error(err)
		return false
	}
	_, err = verifier.Write([]byte(body))
	if err != nil {
		logger.Error(err)
		return false
	}
	err = verifier.Ensure()
	if err != nil {
		logger.Error(err)
		return false
	}

	return true
}