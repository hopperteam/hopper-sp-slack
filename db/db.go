package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"sp-slack/logger"
	"sp-slack/config"
)

var database *mongo.Database
var appCollection *mongo.Collection

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
	appCollection = database.Collection("app")
}