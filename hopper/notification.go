package hopper

import (
	"net/url"
	"net/http"
	hopperApi "github.com/hopperteam/hopper-api/golang"
	"sp-slack/config"
	"sp-slack/logger"
	"sp-slack/utils"
)

type Receivers map[string]string

func SendNotifications(heading string, content string, receivers Receivers, channelId string) []string {
	var notIds []string

	for id, subscription := range receivers {
		callback := getCallback(id, channelId)
		not := hopperApi.DefaultNotification(heading, content).
		IsDone(false).
		Action(
			hopperApi.TextAction("Reply", callback).MarkAsDone(true),
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

func getCallback(userId string, channelId string) string {
	baseCallback := config.BaseUrl + "/reply?"
	params := url.Values{
		"userId": {userId},
		"channelId": {channelId},
	}
	return baseCallback + params.Encode()
}

type Reply struct {
	UserId string `schema:"userId"`
	ChannelId string `schema:"channelId"`
	Text string `json:"text"`
}

func ParseReply(r *http.Request) (*Reply, error) {
	var reply Reply
	err := schemaDecoder.Decode(&reply, r.URL.Query())
	if err != nil {
		return nil, err
	}
	err = utils.ParseJSONReq(r, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}
