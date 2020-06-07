package handler

import (
	"net/http"
	"sp-slack/logger"
	"sp-slack/oauth"
)

func AddToSlack(w http.ResponseWriter, r *http.Request) {
	url := oauth.GenerateButtonUrl()
	w.Write([]byte(url))
	logger.Info("button url returned")
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	// TODO check if state matches
	code := r.URL.Query().Get("code")
	logger.Infof("temp code: %+v", code)
	response, err := oauth.GetOAuthV2AccessToken(code)
	if err != nil {
		logger.Error(err)
		w.Write([]byte(err.Error()))
		return
	}
	logger.Infof("access token for team %s with scopes %s : %s", response.Team.ID, response.Scope, response.AccessToken)
	w.Write([]byte("Top"))
}
