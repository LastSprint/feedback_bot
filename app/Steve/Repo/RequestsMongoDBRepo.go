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

	year, week := time.Now().ISOWeek()

	filter := bson.M{
		"channel_id": channelID,
		"year":       year,
		"week":       week,
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

// GetCountForThisWeek will return reports count which which was made for current time.Now().ISOWeek() week
func (rep *RequestsMongoDBRepo) GetCountForThisWeek(channelID string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(rep.ConnectionString))

	if err != nil {
		return 0, err
	}

	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Println("[WARN] Couldn't close DB connection with error", err.Error())
		}
	}()

	collection := client.Database("ops").Collection("requests")

	year, week := time.Now().ISOWeek()

	filter := bson.M{
		"channel_id": channelID,
		"year":       year,
		"week":       week,
	}

	res := collection.FindOne(context.Background(), filter)

	if err := res.Err(); err != nil {
		return 0, err
	}

	var result_val struct {
		RequestsCount int `bson:"requests_count"`
	}

	if err = res.Decode(&result_val); err != nil {
		return 0, err
	}

	return result_val.RequestsCount, nil
}
