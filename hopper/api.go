package hopper

import (
	hopperApi "github.com/hopperteam/hopper-api/golang"
	"github.com/gorilla/schema"
	"sp-slack/logger"
	"sp-slack/state"
	"sp-slack/config"
)

var Api *hopperApi.HopperApi
var App *hopperApi.App

var schemaDecoder = schema.NewDecoder()

func InitApi() {
	Api = hopperApi.CreateHopperApi(hopperApi.HopperProd)

	ok, err := Api.CheckConnectivity()
	if !ok {
		logger.Fatalf("could not connect to hopper: %s", err.Error())
	}

	strRep := state.GetAppStr()
	if strRep == "" {
		createApp()
	} else {
		parseApp(strRep)
	}	
}

func createApp() {
	var err error
	logger.Info("registering hopper app")
	App, err = Api.CreateApp("Slack", config.BaseUrl, "https://production-cdn.bonus.ly/assets/integration_logos/slack-logo-square-17b7d0d31e59a2aa5a44986849d45d2fc1f9565f47dc4781ab3b218182e7e505.png", "https://hoppercloud.net", "team@hoppercloud.net")
	if err != nil {
		logger.Fatal(err)
	}
	strRep, err := App.Serialize()
	if err != nil {
		logger.Fatal(err)
	}
	state.SetAppStr(strRep)
}

func parseApp(strRep string) {
	var err error
	logger.Info("found existing app entry")
	App, err = Api.DeserializeApp(strRep)
	if err != nil {
		logger.Fatal(err)
	}
}
