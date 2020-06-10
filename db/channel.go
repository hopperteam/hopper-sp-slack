package db

import (
	"context"
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/slack-go/slack"
	"sp-slack/logger"
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
		Name: getChannelName(slackChannel),
		Members: slackChannel.Members,
	}
}

func getChannelName(slackChannel *slack.Channel) string {
	if slackChannel.IsChannel {
		return slackChannel.Name
	}
	if slackChannel.IsIM {
		return "direct message"
	}
	if slackChannel.IsMpIM {
		return "group chat"
	}
	// check this last as MpIM also have IsGroup set to true
	if slackChannel.IsGroup {
		return slackChannel.Name
	}
	return "unknown"
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

func PersistPublicChannels(teamId string) bool {
	api := getTeamApi(teamId)
	if api == nil {
		return false
	}
	err := persistChannels(api)
	if err != nil {
		logger.Error(err)
		return false
	}
	return true
}

func PersistAllChannels(slackId string) bool {
	api := getUserApi(slackId)
	if api == nil {
		return false
	}
	err := persistChannels(api)
	if err != nil {
		logger.Error(err)
		return false
	}
	return true
}

func persistChannels(api *slack.Client) error {
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
		err = addChannelMembers(&conv, api)
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

func addChannelMembers(channel *slack.Channel, api *slack.Client) error {
	members, _, err := api.GetUsersInConversation(&slack.GetUsersInConversationParameters {
		ChannelID: channel.ID,
	})
	if err != nil {
		return err
	}
	channel.Members = members
	return nil
}