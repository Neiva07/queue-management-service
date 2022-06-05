package persistence

import (
	"context"
	"log"
	"time"

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

func CreateUser(email, userType, cpf, name string) (*User, error) {

	newUser := &User{
		CreatedAt: time.Now(),
		Email:     email,
		Name:      name,
		UserType:  userType,
	}

	_, err := usersCollection.InsertOne(context.TODO(), newUser)

	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func GetUser(ctx context.Context, userId string) *User {
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
