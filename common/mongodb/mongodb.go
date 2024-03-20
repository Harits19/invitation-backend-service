package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func InitConnection() error {
	fmt.Println("init mongodb connection")
	context, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	uri := "mongodb://root:example@38.47.180.46:27017/"
	dbName := "invitation"
	client, err := mongo.Connect(context, options.Client().ApplyURI(uri))
	if err != nil {

		return err
	}

	DB = client.Database(dbName)
	return nil
}
