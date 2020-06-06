package config

import (
	"errors"
	"os"
)

var Port string

var ClientId string
var ClientSecret string
var Secret string
var BaseUrl string

var SlackApi string

func Init() error {
	Port = os.Getenv("PORT")

	ClientId = os.Getenv("CLIENT_ID")
	ClientSecret = os.Getenv("CLIENT_SECRET")
	Secret = os.Getenv("SIGNING_SECRET")
	BaseUrl = os.Getenv("BASE_URL")

	SlackApi = os.Getenv("SLACK_API")

	if Port == "" || ClientId == "" || ClientSecret == "" || Secret == "" || BaseUrl == "" || SlackApi == "" {
		return errors.New("config incomplete")
	}

	return nil
}
