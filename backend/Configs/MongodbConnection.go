package Configs

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongodbClient *mongo.Client

func IntializeMongodb() {
	mongouri := os.Getenv("MONGODB_URI")

	var err error

	MongodbClient, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(mongouri))

	if err != nil {
		panic(err)
	}
}
