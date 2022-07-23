package persistence

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var ticketsCollection *mongo.Collection

func init() {
	log.Println("Initiating the tickets collection...")

	userIDIndexModel := mongo.IndexModel{
		Keys: bson.M{
			"user._id": 1, // query by user id
		},
	}

	usersCollection = DB.Collection("tickets")

	_, err := usersCollection.Indexes().CreateOne(context.TODO(), userIDIndexModel)

	if err != nil {
		log.Println("Error creating the userID index", err)
	}
}

func CreateTicket(user *User) (*Ticket, error) {

	ticket, err := generateTicket()

	ticket.User = *user

	if err != nil {
		log.Println("Error generating ticket")
		return nil, err
	}

	_, err = ticketsCollection.InsertOne(context.TODO(), ticket)

	if err != nil {
		log.Println("Error saving the ticket in the database")
		return nil, err
	}

	return ticket, nil
}

func ValidateTicket(ticketKey string) (*Ticket, error) {

	filter := bson.D{{"_id", ticketKey}, {"status", bson.D{{"$ne", "VALID"}}}}

	ticketChange := bson.D{{"status", "PROCESSED"}}

	update := bson.D{{"$set", ticketChange}}

	var ticket *Ticket

	result := paymentsCollection.FindOneAndUpdate(context.TODO(), filter, update)

	if result.Err() != nil {
		log.Printf("Error while validating the ticket with key: %s", ticketKey)
		return nil, result.Err()
	}

	result.Decode(&ticket)

	log.Printf("Successfully validated ticket %s", ticketKey)

	return ticket, nil
}

func generateTicket() (*Ticket, error) {

	id, err := uuid.NewUUID()

	if err != nil {
		log.Println("Error creating the uuid")
		return nil, err
	}

	image, err := qrcode.Encode(id.String(), qrcode.Medium, 256)

	log.Println(string(image))

	if err != nil {
		log.Println("Error creating the qr code")
		return nil, err
	}

	return &Ticket{
		TicketQRCodeKey:   id.String(),
		TicketQRCodeImage: string(image),
		Status:            "VALID",
	}, nil
}
