package db

import (
	"context"
	"time"
	"strings"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sp-slack/logger"
)

type State struct {
	Id primitive.ObjectID `bson:"_id"`
	Key string `bson:"key"`
	Value string `bson:"value"`
}

func CreateState(key string, value string) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	state := State{
		Key: key,
		Value: value,
	}
	_, err := stateCollection.InsertOne(
		ctx,
		state,
	)

	if err != nil {
		logger.Fatalf("could not persist state %s: %s (%s)", key, value, err.Error())
	}
}

func UpdateState(key string, value string) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	_, err := stateCollection.UpdateOne(
		ctx,
		bson.M{ "key": key },
		bson.D{ 
			{ "$set", bson.D{{ "value", value }}},
		 },
	)

	if err != nil {
		logger.Fatalf("could not update state %s: %s (%s)", key, value, err.Error())
	}
}

func SelectState(key string) string {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	state := &State{}
	err := stateCollection.FindOne(
		ctx,
		bson.M{ "key": key },
		).Decode(state)

	if err != nil {
		if !strings.Contains(err.Error(), "mongo: no documents in result") {
			logger.Fatal(err)
		} else {
			CreateState(key, "")
		}
	}
	return state.Value
}
