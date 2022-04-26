package persistence

import "time"

type User struct {
	UserID    string    `bson:"_id"`
	Email     string    `bson:"email"`
	Status    string    `bson:"status,omitempty"`
	createdAt time.Time `bson:"createdAt"`
	Name      string    `bson:"name"`
}
