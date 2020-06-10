package db

import (
	"github.com/slack-go/slack"
	"sp-slack/logger"
)

var teamApiCache = make(map[string]*slack.Client)
var userApiCache = make(map[string]*slack.Client)

func initApiCache() {
	teams, err := selectTeams()
	if err == nil {
		for _, team := range *teams {
			teamApiCache[team.TeamId] = _createApi(team.Token)
		}
	}
	users, err := selectUsers()
	if err == nil {
		for _, user := range *users {
			if user.Token != "" {
				userApiCache[user.SlackId] = _createApi(user.Token)
			}
		}
	}
}

func GetTeamApi(teamId string) *slack.Client {
	api, ok := teamApiCache[teamId]
	if !ok {
		api = _addTeamApi(teamId)		
	}
	return api
}

func GetUserApi(slackId string) *slack.Client {
	api, ok := userApiCache[slackId]
	if !ok {
		api = _addUserApi(slackId)		
	}
	return api
}

func _addTeamApi(teamId string) *slack.Client {
	team, err := selectTeam(teamId)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return _createApi(team.Token)
}

func _addUserApi(slackId string) *slack.Client {
	user, err := selectUser(slackId)
	if err != nil {
		logger.Error(err)
		return nil
	}
	if user.Token == "" {
		return nil
	}
	return _createApi(user.Token)
}

func _createApi(token string) *slack.Client {
	return slack.New(token, slack.OptionDebug(true))
}
