package handler

import (
	"net/http"
	"github.com/slack-go/slack"
	"sp-slack/logger"
	"sp-slack/hopper"
	"sp-slack/utils"
	"sp-slack/db"
)

var genericError string = "an error occured"

func HandleCommand(w http.ResponseWriter, r *http.Request) {
	s, err := parseCommand(r)
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Infof("received command %s", s.Command)

	user, err := db.SelectUser(s.UserID)
	if err != nil {
		logger.Error(err)
		return
	}
	
	switch s.Command {
	case "/subscribe":
		subscribe(user, w)
		break
	case "/unsubscribe":
		subscribe(user, w)
		break
	default:
		logger.Warnf("received unsupported command %s", s.Command)
	}
}

func Callback(w http.ResponseWriter, r *http.Request) {
	sub, err := hopper.ParseSubscribeResponse(r)
	if err != nil {
		logger.Error(err)
		utils.SendPlainText("Did not create subscription", w)
		return
	}

	if hasUserActiveSubscription(sub.UserId) {
		utils.SendPlainText("Can't have more than one Slack subscription", w)
		return
	}

	err = db.AddSubscriptionToUser(sub.UserId, sub.SubscriptionId)
	if err != nil {
		logger.Error(err)
		utils.SendPlainText(genericError, w)
		return
	}

	logger.Infof("added subscription for user %s", sub.UserId)

	utils.SendPlainText("Success", w)
}

func subscribe(user *db.User, w http.ResponseWriter) {
	url, err := hopper.CreateSubscribeRequest(user.SlackId, user.Name)
	if err != nil {
		logger.Error(err)
		url = genericError
	}

	utils.SendEquemeral(url, w)
}

func unsubscribe(user *db.User, w http.ResponseWriter) {
	msg := "successfully unsubscribed"
	err := db.RemoveSubscriptionFromUser(user.SlackId)
	if err != nil {
		logger.Error(err)
		msg = genericError
	}

	logger.Infof("removed subscription for user %s", user.SlackId)

	utils.SendEquemeral(msg, w)
}

func parseCommand(r *http.Request) (slack.SlashCommand, error) {
	if !isFromSlack(r) {
		return slack.SlashCommand{}, notFromSlack
	}

	return slack.SlashCommandParse(r)
}

func hasUserActiveSubscription(slackId string) bool {
	user, err := db.SelectUser(slackId)
	if err != nil {
		logger.Error(err)
		return true
	}
	if user.HasSubscription() {
		logger.Warn("has active subscription")
		return true
	}
	return false
}
