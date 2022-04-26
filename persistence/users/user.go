package persistence

import (
	"context"

	"github.com/Neiva07/queue-management-service/persistence"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var usersCollection *mongo.Collection

func init() {

	usersCollection = persistence.DB.Collection("users")
}

func GetUser(userId string) *User {
	var user *User

	userFilterParamters := bson.M{"_id": userId}

	usersCollection.FindOne(context.TODO(), userFilterParamters).Decode(user)

	return user
}

// func CreateUser(email string, name string) *User {

// 	newUser := &User{
// 		UserID: ?

// 	}

// 	usersCollection.InsertOne(context.TODO(), newUser)
// }
