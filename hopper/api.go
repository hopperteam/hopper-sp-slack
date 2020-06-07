package hopper

import (
	hopperApi "github.com/hopperteam/hopper-api/golang"
	"sp-slack/logger"
	"sp-slack/db"
	"sp-slack/config"
)

var Api *hopperApi.HopperApi
var App *hopperApi.App

func InitApi() {
	Api = hopperApi.CreateHopperApi(hopperApi.HopperDev)

	ok, err := Api.CheckConnectivity()
	if !ok {
		logger.Fatalf("could not connect to hopper: %s", err.Error())
	}

	strRep, err := db.GetApp()
	if err == nil {
		logger.Info("found existing app entry")
		App, err = Api.DeserializeApp(strRep)
		if err != nil {
			logger.Fatal(err)
		}
		return
	}
	logger.Info("registering hopper app")
	App, err := Api.CreateApp("Slack", config.BaseUrl, "https://production-cdn.bonus.ly/assets/integration_logos/slack-logo-square-17b7d0d31e59a2aa5a44986849d45d2fc1f9565f47dc4781ab3b218182e7e505.png", "https://hoppercloud.net", "team@hoppercloud.net")
	if err != nil {
		logger.Fatal(err)
	}
	strRep, err = App.Serialize()
	if err != nil {
		logger.Fatal(err)
	}
	db.CreateApp(strRep)
}
