package Repo

import (
	"github.com/LastSprint/feedback_bot/Steve/Models/Entry"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"log"
	"time"
)

type ReactionsMongoDBRepo struct {
	ConnectionString string
}

func (r *ReactionsMongoDBRepo) AddReactionIfNotAddedPreviously(reaction, channelId, messageId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(r.ConnectionString))

	if err != nil {
		return err
	}

	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Println("[WARN] Couldn't close DB connection with error", err.Error())
		}
	}()

	collection := client.Database("ops").Collection("reactions")

	filter := Entry.MsgReaction{
		ChannelId: channelId,
		MessageId: messageId,
	}

	update := bson.M{
		"$addToSet": bson.M{
			"reactions": reaction,
		},
	}

	_, err = collection.UpdateOne(context.TODO(), filter, update, options.Update().SetUpsert(true))

	if err != nil {
		return err
	}

	return nil
}

func (r *ReactionsMongoDBRepo) RemoveReactionIfPossible(reaction, channelId, messageId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(r.ConnectionString))

	if err != nil {
		return err
	}

	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Println("[WARN] Couldn't close DB connection with error", err.Error())
		}
	}()

	collection := client.Database("ops").Collection("reactions")

	filter := Entry.MsgReaction{
		ChannelId: channelId,
		MessageId: messageId,
	}

	update := bson.M{
		"$pull": bson.M{
			"reactions": reaction,
		},
	}

	_, err = collection.UpdateOne(context.TODO(), filter, update, options.Update().SetUpsert(false))

	if err != nil {
		return err
	}

	return nil
}
