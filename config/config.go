package config

import (
	"errors"
	"os"
)

var Port string

var ClientId string
var ClientSecret string
var Token string
var BaseUrl string

var SlackApi string

func Init() error {
	Port = os.Getenv("PORT")

	ClientId = os.Getenv("CLIENT_ID")
	ClientSecret = os.Getenv("CLIENT_SECRET")
	Token = os.Getenv("VERIFICATION_TOKEN")
	BaseUrl = os.Getenv("BASE_URL")

	SlackApi = os.Getenv("SLACK_API")

	if Port == "" || ClientId == "" || ClientSecret == "" || Token == "" || BaseUrl == "" || SlackApi == "" {
		return errors.New("config incomplete")
	}

	return nil
}
