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

	userIDStatusIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{"user._id", 1},
			{"status", 1},
		}, Options: nil,
	}

	paymentsCollection = DB.Collection("payments")

	_, err := paymentsCollection.Indexes().CreateOne(context.TODO(), userIDStatusIndexModel)

	if err != nil {
		log.Fatal("Error creating the userID-status index", err)
	}
}

func CreatePayment(payment *Payment) error {

	_, err := paymentsCollection.InsertOne(context.TODO(), payment)

	return err
}

func GetUnprocessedPaymentsByID(userID string) ([]*Payment, error) {

	s, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Printf("Invalid user id %s \n", userID)
		return nil, err
	}

	filter := bson.D{{"user._id", s}, {"status", bson.D{{"$ne", "CONCLUIDA"}}}}

	cursor, err := paymentsCollection.Find(context.Background(), filter)

	if err != nil {
		log.Printf("Error finding payments for user %v", userID)
		return nil, err
	}

	payments := []*Payment{}

	for cursor.Next(context.Background()) {
		payment := &Payment{}

		err := cursor.Decode(payment)

		if err != nil {
			log.Printf("Error while parsing payments for user %s", userID)
			return nil, err
		}

		payments = append(payments, payment)
	}

	return payments, nil

}

func ProcessPayment(paymentKey string, status string) (*Payment, error) {

	filter := bson.D{{"_id", paymentKey}, {"status", bson.D{{"$ne", "PROCESSED"}}}}

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
