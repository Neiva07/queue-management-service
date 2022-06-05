package persistence

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var paymentsCollection *mongo.Collection

func init() {
	log.Println("Initiating the payments collection...")

	usersEmailIndexModel := mongo.IndexModel{
		Keys: bson.M{
			"user.email": 1,
		},
	}

	paymentsCollection = DB.Collection("payments")

	_, err := paymentsCollection.Indexes().CreateOne(context.TODO(), usersEmailIndexModel)

	if err != nil {
		log.Println("Error creating the user-email index", err)
	}
}

func CreatePayment(payment *Payment) error {

	_, err := paymentsCollection.InsertOne(context.TODO(), payment)

	return err
}
