package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var cancelConnect context.CancelFunc
var loadoutCollection *mongo.Collection

func initMongoDB(config *Config) error {
	var ctx context.Context
	ctx, cancelConnect = context.WithCancel(context.Background())
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Database.ConnectURI))
	if err != nil {
		return err
	}

	loadoutCollection = client.Database(config.Database.DBName).Collection("loadout")
	return nil
}

func addLoadout(loadout map[string]interface{}) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := loadoutCollection.InsertOne(ctx, loadout)
	if err != nil {
		return "", err
	}

	objectID := res.InsertedID.(primitive.ObjectID)
	return objectID.Hex(), nil
}

func getLoadout(id string) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result bson.M
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = loadoutCollection.FindOne(ctx, bson.D{{"_id", objID}}).Decode(&result)
	if err == nil {
		return result, nil
	}

	return nil, err
}

func getOldLoadout(id string) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result bson.M
	err := loadoutCollection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func closeMongoDB() {
	if cancelConnect != nil {
		cancelConnect()
	}
}
