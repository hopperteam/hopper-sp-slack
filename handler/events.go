package handler

import (
	"encoding/json"
	"github.com/slack-go/slack/slackevents"
	"net/http"
	"sp-slack/logger"
	"sp-slack/utils"
)

func HandleEvents(w http.ResponseWriter, r *http.Request) {
	logger.Info("event triggered")
	eventsAPIEvent, err := parseEvent(r)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if eventsAPIEvent.Type == slackevents.URLVerification {
		handleUrlVerification(eventsAPIEvent, w)
		return
	}

	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		handleInnerEvent(eventsAPIEvent.InnerEvent)
		return
	}
}

func parseEvent(r *http.Request) (slackevents.EventsAPIEvent, error) {
	var event = slackevents.EventsAPIEvent{}
	if !isFromSlack(r) {
		return event, notFromSlack
	}

	body, err := utils.MvReqBodyToStr(r)
	if err != nil {
		return event, err
	}

	event, err = slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	return event, err
}

func handleUrlVerification(event slackevents.EventsAPIEvent, w http.ResponseWriter) {
	logger.Info("verifying url")
	uvEvent, ok := event.Data.(*slackevents.EventsAPIURLVerificationEvent)
	if !ok {
		logger.Error("did not contain the expected event")
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "text")
	utils.SendPlainText(uvEvent.Challenge, w)
}

func handleInnerEvent(event slackevents.EventsAPIInnerEvent) {
	switch event.Data.(type) {
	case *slackevents.MessageEvent:
		logger.Infof("%+v", event.Data)
	}
}
