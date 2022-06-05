package persistence

type Payment struct {
	Quantity      uint64 `bson:"quantity"`
	User          User
	PaymentKey    string
	PaymentMethod uint8
	Status        string
	UnitPrice     uint64 //cents
}
