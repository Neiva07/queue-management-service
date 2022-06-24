package persistence

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func ProcessPayment(paymentKey string, status string) (*Payment, error) {

	s, err := primitive.ObjectIDFromHex(paymentKey)

	if err != nil {
		log.Printf("Invalid payment key %s \n", paymentKey)
		return nil, err
	}

	filter := bson.D{{"_id", s}, {"status", bson.D{{"$ne", "PROCESSED"}}}}

	paymentChange := bson.D{{"status", status}}

	update := bson.D{{"$set", paymentChange}}

	var payment *Payment

	result := paymentsCollection.FindOneAndUpdate(context.TODO(), filter, update)

	if result.Err() != nil {
		log.Printf("Error while updating payment process with key: %s into database", paymentKey)
		return nil, result.Err()
	}

	result.Decode(&payment)

	log.Printf("Successfully updated payment process %s into the database", paymentKey)

	return payment, nil
}
