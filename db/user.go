package db

import (
	"context"
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/slack-go/slack"
)

var empty = ""
var noSubscription = "-"

type dbUser struct {
	Id primitive.ObjectID `bson:"_id,omitempty"`
	SlackId string `bson:"slackId,omitempty"`
	Subscription string `bson:"subscription,omitempty"`
	Name string `bson:"name,omitempty"`
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
	return update
}

func (user *User) HasSubscription() bool {
	return user.Subscription != noSubscription && user.Subscription != empty
}

func newUser(slackUser *slack.User) *User {
	return &User{
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

func (user *User) toDBUser() (*dbUser) {
	return &dbUser{
		SlackId : user.SlackId,
		Subscription: user.Subscription,
		Name: user.Name,
	}
}

func upsertUser(user *User) error {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	userEntity := user.toDBUser()
	_, err := userCollection.UpdateOne(
		ctx,
		userEntity.key(),
		bson.D{
			{"$set", userEntity.update()},
		},
		upsertOpt,
	)

	return err
}

func SelectUser(slackId string) (*User, error) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	user := &dbUser{SlackId: slackId}
	err := userCollection.FindOne(
		ctx,
		user.key(),
		).Decode(user)

	return user.toPrimitiveUser(), err
}

func AddSubscriptionToUser(slackId string, subscription string) error {
	return upsertUser(&User{
		SlackId: slackId,
		Subscription: subscription,
	})
}

func RemoveSubscriptionFromUser(slackId string) error {
	return AddSubscriptionToUser(slackId, noSubscription)
}
