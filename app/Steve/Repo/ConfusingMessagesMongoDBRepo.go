package Repo

import (
	"context"
	"github.com/LastSprint/feedback_bot/Steve/Models/Entry"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type ConfusingMessagesMongoDBRepo struct {
	ConnectionString string
}

func (rep *ConfusingMessagesMongoDBRepo) Save(message Entry.ConfusingMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(rep.ConnectionString))

	if err != nil {
		return err
	}

	collection := client.Database("ops").Collection("confusing_messages")

	filter := Entry.ConfusingMessage{
		MessageId: message.MessageId,
	}

	update := bson.M{
		"$set": message,
	}

	_, err = collection.UpdateOne(context.TODO(), filter, update, options.Update().SetUpsert(true))

	if err != nil {
		return err
	}

	return nil
}
