package Repo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"log"
	"time"
)

type RequestsMongoDBRepo struct {
	ConnectionString string
}

func (rep *RequestsMongoDBRepo) IncrementRequestCount(channelID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(rep.ConnectionString))

	if err != nil {
		return err
	}

	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Println("[WARN] Couldn't close DB connection with error", err.Error())
		}
	}()

	collection := client.Database("ops").Collection("requests")

	filter := bson.M{
		"channel_id": channelID,
	}

	update := bson.M{
		"$set": bson.M{"channel_id": channelID},
		"$inc": bson.M{"requests_count": 1},
	}

	_, err = collection.UpdateOne(context.TODO(), filter, update, options.Update().SetUpsert(true))

	if err != nil {
		return err
	}

	return nil
}
