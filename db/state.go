package db

import (
	"context"
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sp-slack/logger"
)

type dbState struct {
	Id primitive.ObjectID `bson:"_id"`
	Key string `bson:"key"`
	Value string `bson:"value"`
}

func UpsertState(key string, value string) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	_, err := stateCollection.UpdateOne(
		ctx,
		bson.M{ "key": key },
		bson.D{ 
			{ "$set", bson.D{{ "value", value }}},
		},
		upsertOpt,
	)

	if err != nil {
		logger.Fatalf("could not upsert state %s: %s (%s)", key, value, err.Error())
	}
}

func SelectState(key string) string {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	state := &dbState{}
	err := stateCollection.FindOne(
		ctx,
		bson.M{ "key": key },
		).Decode(state)

	if err != nil && !emptyResult(err) {
		logger.Fatal(err)
	}
	return state.Value
}
