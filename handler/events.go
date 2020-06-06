package handler

import (
	"encoding/json"
	"github.com/slack-go/slack/slackevents"
	"net/http"
	"sp-slack/logger"
	"sp-slack/utils"
)

//var api = slack.New(os.Getenv("BU_OAUTH_ACCESS"), slack.OptionDebug(true))

func HandleEvents(w http.ResponseWriter, r *http.Request) {
	logger.Info("event triggered")
	body, eventsAPIEvent, err := parseEvent(r)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if eventsAPIEvent.Type == slackevents.URLVerification {
		handleUrlVerification(body, w)
		return
	}

	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		handleInnerEvent(eventsAPIEvent.InnerEvent)
		return
	}
}

func parseEvent(r *http.Request) (string, slackevents.EventsAPIEvent, error) {
	if !isFromSlack(r) {
		return "", slackevents.EventsAPIEvent{}, notFromSlack
	}

	body, err := utils.MvReqBodyToStr(r)
	if err != nil {
		return body, slackevents.EventsAPIEvent{}, err
	}
	event, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	return body, event, err
}

func handleUrlVerification(body string, w http.ResponseWriter) {
	logger.Info("verifying url")
	var res *slackevents.ChallengeResponse
	err := json.Unmarshal([]byte(body), &res)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "text")
	w.Write([]byte(res.Challenge))
}

func handleInnerEvent(event slackevents.EventsAPIInnerEvent) {
	switch event.Data.(type) {
	case *slackevents.MessageEvent:
		logger.Infof("%+v", event.Data)
	}
}
