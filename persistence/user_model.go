package persistence

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	UserID    primitive.ObjectID `json:"userId" bson:"_id"`
	Email     string             `json:"email" bson:"email"`
	Status    *string            `json:"status,omitempty" bson:"status,omitempty"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	Name      string             `json:"name" bson:"name"`
	UserType  string             `json:"userType" bson:"userType"`
	CPF       string             `json:"cpf" bson:"cpf"`
}

var USER_TYPE_PRICE_MAP = map[string]uint64{
	"student":   100,
	"professor": 300,
}

func (u *User) GetTotalPayment(quantity uint64) uint64 {

	return quantity * USER_TYPE_PRICE_MAP[u.UserType]
}
