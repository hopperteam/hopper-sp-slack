package hopper

import (
	"net/http"
	"errors"
	"github.com/gorilla/schema"
	hopperApi "github.com/hopperteam/hopper-api/golang"
	"sp-slack/config"
)

func CreateSubscribeRequest(userId string, userName string) (string, error) {
	callback := config.BaseUrl + "/callback?userId=" + userId
	url, err := App.CreateSubscribeRequest(callback, hopperApi.StrPtr(userName))
	if err != nil {
		return "", err
	}
	return url, nil
}

type callbackResponse struct {
	Status string `schema:"status"`
	Error string `schema:"error"`
	SubscriptionId string `schema:"id"`
	UserId string `schema:"userId"`
}

type Subscription struct {
	SubscriptionId string
	UserId string
}

var schemaDecoder = schema.NewDecoder()

func ParseSubscribeResponse(r *http.Request) (*Subscription, error) {
	var res callbackResponse
	err := schemaDecoder.Decode(&res, r.URL.Query())
	if err != nil {
		return nil, err
	} 
	if res.Status == "error" {
		return nil, errors.New(res.Error)
	}
	sub := &Subscription{
		SubscriptionId: res.SubscriptionId,
		UserId: res.UserId,
	}
	return sub, nil
}