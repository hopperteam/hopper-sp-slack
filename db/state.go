package db

import (
	"context"
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sp-slack/logger"
)

type dbState struct {
	Id primitive.ObjectID `bson:"_id,omitempty"`
	Key string `bson:"key,omitempty"`
	Value string `bson:"value,omitempty"`
}

func (state *dbState) key() bson.M {
	return bson.M{ "key": state.Key }
}

func (state *dbState) update() bson.M {
	update := bson.M{}
	if state.Value != "" {
		update["value"] = state.Value
	}
	return update
}

func UpsertState(key string, value string) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	state := &dbState{
		Key: key,
		Value: value,
	}
	_, err := stateCollection.UpdateOne(
		ctx,
		state.key(),
		bson.D{ 
			{ "$set", state.update()},
		},
		upsertOpt,
	)

	if err != nil {
		logger.Fatalf("could not upsert state %s: %s (%s)", key, value, err.Error())
	}
}

func SelectState(key string) string {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	state := &dbState{Key: key}
	err := stateCollection.FindOne(
		ctx,
		state.key(),
		).Decode(state)

	if err != nil && !emptyResult(err) {
		logger.Fatal(err)
	}
	return state.Value
}
