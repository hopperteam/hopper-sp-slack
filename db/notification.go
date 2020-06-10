package db

import (
	"context"
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type dbNotification struct {
	Id primitive.ObjectID `bson:"_id,omitempty"`
	MsgId string `bson:"msgId,omitempty"`
	HopperId string `bson:"hopperId,omitempty"`
}

func InsertNotifications(ids []string, msgId string) error {
	for _, id := range ids {
		err := insertNotification(&dbNotification{MsgId: msgId, HopperId: id})
		if err != nil {
			return err
		}
	}
	return nil
}

func GetNotificationIds(msgId string) ([]string, error) {
	nots, err := selectNotifications(msgId)
	if err != nil {
		return []string{}, err
	}
	var ids []string
	for _, not := range *nots {
		ids = append(ids, not.HopperId)
	}
	return ids, nil
}

func DeleteNotifications(ids []string) error {
	for _, id := range ids {
		err := deleteNotification(id)
		if err != nil {
			return err
		}
	}
	return nil
}

func insertNotification(notification *dbNotification) error {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	_, err := notificationCollection.InsertOne(
		ctx,
		notification,
	)

	return err
}

func selectNotifications(msgId string) (*[]dbNotification, error) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	var notifications = &[]dbNotification{}
	cursor, err := notificationCollection.Find(
		ctx,
		bson.M{
			"msgId": msgId,
		},
	)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, notifications)
	if err != nil {
		return nil, err
	}

	return notifications, err
}

func deleteNotification(hopperId string) error {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	_, err := notificationCollection.DeleteOne(
		ctx,
		bson.M{
			"hopperId": hopperId,
		},
	)

	return err
}
