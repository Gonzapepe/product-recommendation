package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongoDB(uri string) (*mongo.Client, error) {

	clientOptions := options.Client().ApplyURI(uri)
	

	mongoClient, err := mongo.Connect(context.TODO(), clientOptions)

	err = mongoClient.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	log.Println("Conectado a mongo")
	return mongoClient, nil
}
