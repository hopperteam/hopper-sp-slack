package db

import (
	"context"
	"time"
	"strings"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sp-slack/logger"
)

var appId string

type app struct {
	Id primitive.ObjectID `bson:"_id"`
	StrRep string `bson:"strRep"`
}

func CreateApp(strRep string) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	a := app{
		StrRep: strRep,
	}
	_, err := appCollection.InsertOne(ctx, a)
	if err != nil {
		logger.Fatalf("could not persist app: %s", err.Error())
	}
}

func GetApp() (string, error) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	a := &app{}
	err := appCollection.FindOne(ctx, bson.M{}).Decode(a)
	if err != nil && !strings.Contains(err.Error(), "mongo: no documents in result") {
		logger.Fatal(err)
	}
	return a.StrRep, err
}
