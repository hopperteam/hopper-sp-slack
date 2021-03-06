package handler

import (
	"net/http"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"sp-slack/logger"
	"sp-slack/db"
	"sp-slack/hopper"
)

func processMessage(messageEvent *slackevents.MessageEvent, teamId string) {
	if messageEvent.IsEdited() {
		updateNotifications(messageEvent.Message)
		return
	}
	if messageEvent.SubType == "message_deleted" {
		deleteNotifications(messageEvent.PreviousMessage)
		return
	}
	sendNotifications(messageEvent, teamId)
}

func HandleReply(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "content-type")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	reply, err := hopper.ParseReply(r)
	if err != nil {
		logger.Error(err)
		return
	}
	api := db.GetUserApi(reply.UserId)
	if api == nil {
		logger.Error("user did not authorize bot")
		return
	}
	_, _, err = api.PostMessage(reply.ChannelId, slack.MsgOptionText(reply.Text, true), slack.MsgOptionAsUser(true))
	if err != nil {
		logger.Error(err)
	}
}

func sendNotifications(messageEvent *slackevents.MessageEvent, teamId string) {
	heading := createHeading(messageEvent.User, messageEvent.Channel, teamId)
	receivers := getReceivers(messageEvent.User, messageEvent.Channel)
	ids := hopper.SendNotifications(heading, messageEvent.Text, receivers, messageEvent.Channel)
	err := db.InsertNotifications(ids, messageEvent.ClientMsgID)
	if err != nil {
		logger.Error(err)
	}
}

func updateNotifications(messageEvent *slackevents.MessageEvent) {
	ids := getHopperNotificationIds(messageEvent.ClientMsgID)
	hopper.UpdateNotifications(messageEvent.Text, ids)
}

func deleteNotifications(messageEvent *slackevents.MessageEvent) {
	ids := getHopperNotificationIds(messageEvent.ClientMsgID)
	hopper.DeleteNotifications(ids)
	err := db.DeleteNotifications(ids)
	if err != nil {
		logger.Error(err)
	}
}

func createHeading(userId string, channelId string, teamId string) string {
	return getUserName(userId) + " (" + getChannelName(channelId) + ": " + getTeamName(teamId) + ")"
}

func getReceivers(senderId string, channelId string) hopper.Receivers {
	receivers := make(hopper.Receivers)
	channel, err := db.SelectChannel(channelId)
	if err != nil {
		logger.Error(err)
		return receivers
	}
	for _, userId := range channel.Members {
		if userId == senderId {
			continue
		}
		user, err := db.SelectUser(userId)
		if err != nil {
			logger.Error(err)
			continue
		}
		if user.Subscription != "" {
			receivers[userId] = user.Subscription
		}
	}
	return receivers
}

func getHopperNotificationIds(messageId string) []string {
	ids, err := db.GetNotificationIds(messageId)
	if err != nil {
		logger.Error(err)
	}
	return ids
}

func getUserName(userId string) string {
	user, err := db.SelectUser(userId)
	if err != nil {
		logger.Error(err)
		return "unknown"
	}
	return user.Name
}

func getChannelName(channelId string) string {
	channel, err := db.SelectChannel(channelId)
	if err != nil {
		logger.Error(err)
		return "unknown"
	}
	return channel.Name
}

func getTeamName(teamId string) string {
	team, err := db.SelectTeam(teamId)
	if err != nil {
		logger.Error(err)
		return "unknown"
	}
	return team.Name
}
