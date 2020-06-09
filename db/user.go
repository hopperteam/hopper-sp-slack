package db

import (
	"errors"
	"context"
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/slack-go/slack"
	"sp-slack/logger"
)

var empty = ""
var noSubscription = "-"
var unauthedUser = errors.New("user has no access")

type dbUser struct {
	Id primitive.ObjectID `bson:"_id,omitempty"`
	SlackId string `bson:"slackId,omitempty"`
	Subscription string `bson:"subscription,omitempty"`
	Name string `bson:"name,omitempty"`
	Token string `bson:"token,omitempty"`
}

type User struct {
	SlackId string 
	Subscription string
	Name string
}

func (user *dbUser) key() bson.M {
	return bson.M{ "slackId": user.SlackId }
}

func (user *dbUser) update() bson.M {
	update := bson.M{}
	if user.Subscription != "" {
		update["subscription"] = user.Subscription
	}
	if user.Name != "" {
		update["name"] = user.Name
	}
	if user.Token != "" {
		update["token"] = user.Token
	}
	return update
}

func (user *User) HasSubscription() bool {
	return user.Subscription != noSubscription && user.Subscription != empty
}

func newUser(slackUser *slack.User) *dbUser {
	return &dbUser{
		SlackId: slackUser.ID,
		Name: slackUser.RealName,
	}
}

func (user *dbUser) toPrimitiveUser() (*User) {
	return &User{
		SlackId : user.SlackId,
		Subscription: user.Subscription,
		Name: user.Name,
	}
}

func upsertUser(user *dbUser) error {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	_, err := userCollection.UpdateOne(
		ctx,
		user.key(),
		bson.D{
			{"$set", user.update()},
		},
		upsertOpt,
	)

	return err
}

func SelectUser(slackId string) (*User, error) {
	user, err := selectUser(slackId)
	return user.toPrimitiveUser(), err
}

func selectUser(slackId string) (*dbUser, error) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	user := &dbUser{SlackId: slackId}
	err := userCollection.FindOne(
		ctx,
		user.key(),
		).Decode(user)

	return user, err
}

func selectUsers() (*[]dbUser, error) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	var users = &[]dbUser{}
	cursor, err := userCollection.Find(
		ctx,
		bson.M{},
	)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, users)
	if err != nil {
		return nil, err
	}

	return users, err
}

func PersistUserAccess(slackId string, token string) bool {
	var ok = true

	err := updateUserAuth(slackId, token)
	if err != nil {
		logger.Error(err)
		return false
	}
	err = updateKnownChannels(slackId)
	if err != nil {
		logger.Error(err)
		ok = false
	}

	return ok
}

func updateUserAuth(slackId string, token string) error {
	return upsertUser(&dbUser{
		SlackId: slackId,
		Token: token,
	})
}

func updateKnownChannels(slackId string) error {
	api := getUserApi(slackId)
	if api == nil {
		return unauthedUser
	}
	convs, _, err := api.GetConversations(&slack.GetConversationsParameters{
		ExcludeArchived: "true",
		Types: []string{
			"public_channel",
			"private_channel",
			"mpim",
			"im",
		},
	})
	if err != nil {
		return err
	}
	for _, conv := range convs {
		err = addChannelMembers(&conv, slackId)
		if err != nil {
			return err
		}
		err = upsertChannel(newChannel(&conv))
		if err != nil {
			return err
		}
	}
	return nil
}

func AddSubscriptionToUser(slackId string, subscription string) error {
	return upsertUser(&dbUser{
		SlackId: slackId,
		Subscription: subscription,
	})
}

func RemoveSubscriptionFromUser(slackId string) error {
	return AddSubscriptionToUser(slackId, noSubscription)
}
