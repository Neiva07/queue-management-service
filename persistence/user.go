package persistence

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var usersCollection *mongo.Collection

func init() {
	log.Println("Initiating the users collection...")

	isUnique := true

	usersEmailIndexModel := mongo.IndexModel{
		Keys: bson.M{
			"email": 1, // query by email
		},
		Options: &options.IndexOptions{
			Unique: &isUnique,
		},
	}

	usersCollection = DB.Collection("users")

	_, err := usersCollection.Indexes().CreateOne(context.TODO(), usersEmailIndexModel)

	if err != nil {
		log.Println("Error creating the user-email index", err)
	}

}

func CreateUser(newUser *User) (*User, error) {

	_, err := usersCollection.InsertOne(context.TODO(), newUser)

	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func GetUserByEmail(ctx context.Context, email string) *User {
	var user User

	userFilterParamters := bson.M{"email": email}

	err := usersCollection.FindOne(ctx, userFilterParamters).Decode(&user)

	if err != nil {
		log.Printf("Error while querying email %s \n", email)
		log.Println(err.Error())
		return nil
	}

	return &user

}

func GetUserById(ctx context.Context, userId string) *User {
	var user User

	s, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		log.Printf("Invalid user id %s \n", userId)
		log.Println(err.Error())
		return nil
	}
	log.Println(s)

	userFilterParamters := bson.M{"_id": s}

	err = usersCollection.FindOne(ctx, userFilterParamters).Decode(&user)

	if err != nil {
		log.Printf("Error while querying userId %s \n", userId)
		log.Println(err.Error())
		return nil
	}

	return &user
}

func BuyTickets(userId string, quantity uint64) error {

	s, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		log.Printf("Invalid user id %s \n", userId)
		return err
	}

	filter := bson.D{{"_id", s}}

	ticketChange := bson.D{{"ticketQuantity", quantity}}

	update := bson.D{{"$inc", ticketChange}}

	result := usersCollection.FindOneAndUpdate(context.TODO(), filter, update)

	return result.Err()

}

func RedeemTicket(userId string) error {

	s, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		log.Printf("Invalid user id %s \n", userId)
		return err
	}

	filter := bson.D{{"_id", s}}

	ticketChange := bson.D{{"ticketQuantity", -1}}

	update := bson.D{{"$inc", ticketChange}}

	result := usersCollection.FindOneAndUpdate(context.TODO(), filter, update)

	return result.Err()

}
