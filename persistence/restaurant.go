package persistence

import "time"

type Restaurant struct {
	RestaurantID string    `bson:"_id"`
	Name         string    `bson:"name"`
	Email        string    `bson:"email"`
	Status       string    `bson:"status,omitempty"`
	createdAt    time.Time `bson:"createdAt"`
}
