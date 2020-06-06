package oauth

import (
	"errors"
	"net/http"
	"net/url"
	"sp-slack/config"
	"sp-slack/logger"
	"sp-slack/utils"
	"strings"
)

var client *http.Client = &http.Client{}

func GenerateButtonUrl() string {
	var scopes = []string{
		"channels:history",
		"groups:history",
		"im:history",
		"mpim:history",
		"commands",
	}

	var params = url.Values{
		"client_id":    {config.ClientId},
		"scope":        {strings.Join(scopes, " ")},
		"redirect_uri": {config.BaseUrl + "/redirect"},
		"state":        {"test"},
	}

	return "https://slack.com/oauth/v2/authorize?" + params.Encode()
}

func GetOAuthV2AccessToken(code string) (res *OAuthV2Response, err error) {
	vals := url.Values{
		"client_id":     {config.ClientId},
		"client_secret": {config.ClientSecret},
		"code":          {code},
		"redirect_uri":  {config.BaseUrl + "/redirect"},
	}
	response, err := postUrlEncoded(config.SlackApi+"/oauth.v2.access", vals)
	if err != nil {
		logger.Error(err)
		return nil, errors.New("did not get access from slack")
	}

	auth := &OAuthV2Response{}
	err = utils.ParseJSONRes(response, auth)
	if err != nil {
		logger.Error(err)
		return nil, errors.New("received invalid data from slack")
	}
	return auth, err
}

func postUrlEncoded(url string, content url.Values) (res *http.Response, err error) {
	reqBody := strings.NewReader(content.Encode())
	req, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("response code not okay")
	}

	return response, nil
}
