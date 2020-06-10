package hopper

import (
	"net/url"
	//"net/http"
	//"errors"
	//"github.com/gorilla/schema"
	hopperApi "github.com/hopperteam/hopper-api/golang"
	"sp-slack/config"
	"sp-slack/logger"
)

type Receivers map[string]string

func SendNotifications(heading string, content string, receivers Receivers, messageId string) []string {
	var notIds []string

	for id, subscription := range receivers {
		callback := getCallback(id, messageId)
		not := hopperApi.DefaultNotification(heading, content).
		IsDone(false).
		Action(
			hopperApi.TextAction("Reply", callback),
		)

		notId, err := Api.PostNotification(subscription, not)
		if err != nil {
			logger.Error(err)
			continue
		}
		notIds = append(notIds, notId)
	}

	return notIds
}

func UpdateNotifications(content string, notIds []string) {
	update := &hopperApi.NotificationUpdate{
		Content: hopperApi.StrPtr(content),
	}
	for _, id := range notIds {
		err := Api.UpdateNotification(id, update)
		if err != nil {
			logger.Error(err)
		}
	}
}

func DeleteNotifications(notIds []string) {
	for _, id := range notIds {
		err := Api.DeleteNotification(id)
		if err != nil {
			logger.Error(err)
		}
	}
}

var baseCallback = config.BaseUrl + "/reply?"

func getCallback(userId string, messageId string) string {
	params := url.Values{
		"userId": {userId},
		"messageId": {messageId},
	}
	return baseCallback + params.Encode()
}
