package config

import (
	"os"
	"sp-slack/logger"
)

var Port string

var ClientId string
var ClientSecret string
var Secret string
var BaseUrl string

var SlackApi string

func Init() {
	Port = getStr("PORT")

	ClientId = getStr("CLIENT_ID")
	ClientSecret = getStr("CLIENT_SECRET")
	Secret = getStr("SIGNING_SECRET")
	BaseUrl = getStr("BASE_URL")

	SlackApi = getStr("SLACK_API")
}

func getStr(key string) string {
	val := os.Getenv(key)
	if val == "" {
		logger.Fatalf("required env %s not specified", key)
	}
	return val
}
