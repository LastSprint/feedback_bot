package Repo

import (
	"context"
	"github.com/LastSprint/feedback_bot/Steve/Models/Entry"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
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

	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Println("[WARN] Couldn't close DB connection with error", err.Error())
		}
	}()

	collection := client.Database("ops").Collection("confusing_messages")

	filter := Entry.ConfusingMessage{
		MessageId: message.MessageId,
		ReportType: message.ReportType,
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

func (rep *ConfusingMessagesMongoDBRepo) GetCountForThisWeek(channelID string) (map[string]int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(rep.ConnectionString))

	if err != nil {
		return nil, err
	}

	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Println("[WARN] Couldn't close DB connection with error", err.Error())
		}
	}()

	collection := client.Database("ops").Collection("confusing_messages")

	timeDiff := time.Hour * 24 * 6

	filter := bson.M{
		"$and": bson.A{
			bson.M{
				"reportDate": bson.M{
					"$lte": time.Now().Add(time.Duration(timeDiff)),
				},
			}, bson.M{
				"reportDate": bson.M{
					"$gte": time.Now().Add(time.Duration(timeDiff) * -1),
				},
			},
		},
	}

	cursor, err := collection.Find(context.TODO(), filter)

	if err != nil {
		return nil, err
	}

	messages := []Entry.ConfusingMessage{}

	if err = cursor.All(context.Background(), &messages); err != nil {
		return nil, err
	}

	result := map[string]int{}

	for _, val := range messages {
		count, ok := result[val.ReportType]

		if !ok {
			result[val.ReportType] = 1
			continue
		}

		result[val.ReportType] = count + 1
	}

	return result, err
}
