package persistence

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payment struct {
	PaymentKey    string `json:"paymentKey" bson:"_id"`
	Quantity      uint64 `json:"quantity" bson:"quantity"`
	User          User
	PaymentMethod uint8  `json:"paymentMethod" bson:"paymentMethod"`
	Status        string `json:"status" bson:"status"`
	UnitPrice     uint64 `json:"unitPrice" bson:"unitPrice"` //cents
	TotalCost     uint64 `json:"totalCost" bson:"totalCost"`
	QRCodeKey     string `json:"qrCodeKey" bson:"qrCodeKey"`
}

type Restaurant struct {
	RestaurantID string    `bson:"_id"`
	Name         string    `bson:"name"`
	Email        string    `bson:"email"`
	Status       string    `bson:"status,omitempty"`
	createdAt    time.Time `bson:"createdAt"`
}

type Ticket struct {
	TicketQRCodeKey   string `bson:"_id" json:"ticketQRCodeKey"`
	TicketQRCodeImage string `bson:"ticketQRCodeImage" json:"ticketQRCodeImage"`
	Status            string `bson:"status" json:"status"`
	User              User   `bson:"user" json:"user"`
}

type User struct {
	UserID         primitive.ObjectID `json:"userId" bson:"_id"`
	Email          string             `json:"email" bson:"email"`
	Status         *string            `json:"status,omitempty" bson:"status,omitempty"`
	CreatedAt      time.Time          `json:"createdAt" bson:"createdAt"`
	Name           string             `json:"name" bson:"name"`
	UserType       string             `json:"userType" bson:"userType"`
	CPF            string             `json:"cpf" bson:"cpf"`
	TicketQuantity uint64             `json:"ticketQuantity" bson:"ticketQuantity"`
	Password       []byte             `json:"password" bson:"password"`
}

var USER_TYPE_PRICE_MAP = map[string]uint64{
	"student":   1,
	"professor": 3,
}

func (u *User) GetTotalPayment(quantity uint64) uint64 {

	return quantity * USER_TYPE_PRICE_MAP[u.UserType]
}
