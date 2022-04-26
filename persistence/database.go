package persistence

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database = initiateDatabase()

func initiateDatabase() *mongo.Database {
	opt := options.Client().ApplyURI("mongodb://localhost:27017")
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("Initiating the database...")

	return client.Database("QueueManagement")
}
