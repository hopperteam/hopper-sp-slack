package handler

import (
	"bytes"
	"encoding/json"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"net/http"
	"os"
	"sp-slack/config"
	"sp-slack/logger"
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

func parseEvent(r *http.Request) (body string, event slackevents.EventsAPIEvent, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	body = buf.String()
	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: config.Token}))
	return body, eventsAPIEvent, err
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
