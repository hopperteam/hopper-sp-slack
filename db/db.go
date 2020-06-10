package db

import (
	"strings"
	"context"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"sp-slack/logger"
	"sp-slack/config"
)

var database *mongo.Database
var stateCollection *mongo.Collection
var userCollection *mongo.Collection
var teamCollection *mongo.Collection
var channelCollection *mongo.Collection
var notificationCollection *mongo.Collection

type entity interface {
	key() bson.M
	update() bson.M
}

var upsertOpt = options.Update().SetUpsert(true)

func ConnectDB() {
	dbOptions := options.Client().ApplyURI(config.DbConStr)
	client, err := mongo.NewClient(dbOptions)
	if err != nil {
		logger.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		logger.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		logger.Fatalf("could not ping db: %s", err)
	}
	database = client.Database(config.DbName)

	initCollections()

	initApiCache()
}

func initCollections() {
	stateCollection = database.Collection("state")
	userCollection = database.Collection("user")
	teamCollection = database.Collection("team")
	channelCollection = database.Collection("channel")
	notificationCollection = database.Collection("notification")

	createIndex(stateCollection, bson.M{ "key": 1 }, true)
	createIndex(userCollection, bson.M{ "slackId": 1 }, true)
	createIndex(teamCollection, bson.M{ "teamId": 1 }, true)
	createIndex(channelCollection, bson.M{ "channelId": 1 }, true)
	createIndex(notificationCollection, bson.M{ "msgId": 1 }, false)
}

func createIndex(col *mongo.Collection, keys bson.M, unique bool) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Options: options.Index().SetUnique(unique),
		Keys: keys,
	})

	if err != nil {
		logger.Fatal(err)
	}
}

func emptyResult(err error) bool {
	return strings.Contains(err.Error(), "mongo: no documents in result")
}