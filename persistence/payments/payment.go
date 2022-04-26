package payments

type Payment struct {
	Quantity      uint64
	UserID        string
	PaymentKey    string
	PaymentMethod uint8
	Status        uint8
	UnitPrice     uint64 //centavos
}
