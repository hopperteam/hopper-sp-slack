package db

import (
	"context"
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/slack-go/slack"
)

type dbChannel struct {
	Id primitive.ObjectID `bson:"_id,omitempty"`
	ChannelId string `bson:"channelId,omitempty"`
	Name string `bson:"name,omitempty"`
	Members []string `bson:"members,omitempty"`
}

type Channel struct {
	ChannelId string 
	Name string
	Members []string
}

func (channel *dbChannel) key() bson.M {
	return bson.M{ "channelId": channel.ChannelId }
}

func (channel *dbChannel) update() bson.M {
	update := bson.M{}
	if channel.Name != "" {
		update["name"] = channel.Name
	}
	if len(channel.Members) > 0 {
		update["members"] = channel.Members
	}
	return update
}

func newChannel(slackChannel *slack.Channel) *dbChannel {
	return &dbChannel{
		ChannelId: slackChannel.ID,
		Name: slackChannel.Name,
		Members: slackChannel.Members,
	}
}

func (channel *dbChannel) toPrimitiveChannel() (*Channel) {
	return &Channel{
		ChannelId : channel.ChannelId,
		Name: channel.Name,
		Members: channel.Members,
	}
}

func upsertChannel(channel *dbChannel) error {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	_, err := channelCollection.UpdateOne(
		ctx,
		channel.key(),
		bson.D{
			{"$set", channel.update()},
		},
		upsertOpt,
	)

	return err
}

func SelectChannel(channelId string) (*Channel, error) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	channel := &dbChannel{ChannelId: channelId}
	err := channelCollection.FindOne(
		ctx,
		channel.key(),
		).Decode(channel)

	return channel.toPrimitiveChannel(), err
}

func addChannelMembers(channel *slack.Channel, channelId string) error {
	api := getUserApi(channelId)
	if api == nil {
		return unauthedUser
	}
	members, _, err := api.GetUsersInConversation(&slack.GetUsersInConversationParameters {
		ChannelID: channel.ID,
	})
	if err != nil {
		return err
	}
	channel.Members = members
	return nil
}