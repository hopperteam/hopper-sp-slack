package handler

import (
	"net/http"
	"sp-slack/logger"
	"sp-slack/oauth"
	"sp-slack/db"
	"sp-slack/utils"
)

func AddToSlack(w http.ResponseWriter, r *http.Request) {
	url := oauth.GenerateButtonUrl()
	utils.SendPlainText(url, w)
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	// TODO check if state matches
	code := r.URL.Query().Get("code")
	response, err := oauth.GetOAuthV2AccessToken(code)
	if err != nil {
		logger.Error(err)
		utils.SendPlainText("Could not get permanent access to workspace", w)
		return
	}

	ok := db.PersistTeam(response.Team.ID, response.AccessToken)
	if !ok {
		utils.SendPlainText("Could not gather workspace information", w)
		return
	}

	utils.SendPlainText("Success", w)
}
