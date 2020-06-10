package db

import (
	"errors"
	"context"
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sp-slack/logger"
)

var unknownTeam = errors.New("unknown team")

type dbTeam struct {
	Id primitive.ObjectID `bson:"_id,omitempty"`
	TeamId string `bson:"teamId,omitempty"`
	Token string `bson:"token,omitempty"`
	Name string `bson:"name,omitempty"`
}

type Team struct {
	TeamId string
	Name string
}

func (team *dbTeam) key() bson.M {
	return bson.M{ "teamId": team.TeamId }
}

func (team *dbTeam) update() bson.M {
	update := bson.M{}
	if team.Token != "" {
		update["token"] = team.Token
	}
	if team.Name != "" {
		update["name"] = team.Name
	}
	return update
}

func (team *dbTeam) toPrimitiveTeam() (*Team) {
	return &Team{
		TeamId : team.TeamId,
		Name: team.Name,
	}
}

func SelectTeam(teamId string) (*Team, error) {
	team, err := selectTeam(teamId)
	return team.toPrimitiveTeam(), err
}

func selectTeam(teamId string) (*dbTeam, error) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	team := &dbTeam{TeamId: teamId}
	err := teamCollection.FindOne(
		ctx,
		team.key(),
		).Decode(team)

	return team, err
}

func selectTeams() (*[]dbTeam, error) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	var teams = &[]dbTeam{}
	cursor, err := teamCollection.Find(
		ctx,
		bson.M{},
	)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, teams)
	if err != nil {
		return nil, err
	}

	return teams, err
}

// forfeits when it can't upsert auth details
func PersistTeam(teamId string, token string) bool {
	var ok = true

	err := updateTeamAuth(teamId, token)
	if err != nil {
		logger.Error(err)
		return false
	}
	err = updateTeamDetails(teamId)
	if err != nil {
		logger.Error(err)
		ok = false
	}
	err = updateTeamMembers(teamId)
	if err != nil {
		logger.Error(err)
		ok = false
	}

	return ok
}

func upsertTeam(team *dbTeam) error {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	_, err := teamCollection.UpdateOne(
		ctx,
		team.key(),
		bson.D{ 
			{ "$set", team.update()},
		},
		upsertOpt,
	)

	return err
}

func updateTeamAuth(teamId string, token string) error {
	return upsertTeam(&dbTeam{
		TeamId: teamId,
		Token: token,
	})
}

func updateTeamDetails(teamId string) error {
	api := GetTeamApi(teamId)
	if api == nil {
		return unknownTeam
	}
	slackTeam, err := api.GetTeamInfo()
	if err != nil {
		return err
	}

	team := &dbTeam{
		TeamId: teamId,
		Name: slackTeam.Name,
	}

	return upsertTeam(team)
}

func updateTeamMembers(teamId string) error {
	api := GetTeamApi(teamId)
	if api == nil {
		return unknownTeam
	}
	users, err := api.GetUsers()
	if err != nil {
		return err
	}
	for _, user := range users {
		err = upsertUser(newUser(&user))
		if err != nil {
			return err
		}
	}
	return nil
}
