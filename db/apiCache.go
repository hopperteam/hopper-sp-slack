package db

import (
	"github.com/slack-go/slack"
	"sp-slack/logger"
)

var apiCache = make(map[string]*slack.Client)

func initApiCache() {
	teams, err := selectTeams()
	if err == nil {
		for _, team := range *teams {
			apiCache[team.TeamId] = createApi(team.Token)
		}
	}
}

func getApi(teamId string) *slack.Client {
	api, ok := apiCache[teamId]
	if !ok {
		api = _addApi(teamId)		
	}
	return api
}

func _addApi(teamId string) *slack.Client {
	team, err := selectTeam(teamId)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return createApi(team.Token)
}

func createApi(token string) *slack.Client {
	return slack.New(token, slack.OptionDebug(true))
}
