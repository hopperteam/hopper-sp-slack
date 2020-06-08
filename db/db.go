package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"time"
	"sp-slack/logger"
	"sp-slack/config"
)

var database *mongo.Database
var stateCollection *mongo.Collection

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

	database = client.Database(config.DbName)
	stateCollection = database.Collection("state")

	createIndex(stateCollection, bson.M{ "key": 1 })
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