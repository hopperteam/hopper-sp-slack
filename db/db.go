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

	createIndex(stateCollection, bson.M{ "key": 1 })
	createIndex(userCollection, bson.M{ "slackId": 1 })
	createIndex(teamCollection, bson.M{ "teamId": 1 })
}

func createIndex(col *mongo.Collection, keys bson.M) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Options: options.Index().SetUnique(true),
		Keys: keys,
	})

	if err != nil {
		logger.Fatal(err)
	}
}

func emptyResult(err error) bool {
	return strings.Contains(err.Error(), "mongo: no documents in result")
}