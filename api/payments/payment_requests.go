package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Neiva07/queue-management-service/persistence"
	"github.com/gorilla/mux"
)

type CreatePaymentRequestRequest struct {
	Quantity      uint64 `json:"quantity"`
	PaymentMethod uint8  `json:"paymentMethod"`
}

func CreatePaymentRequest(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	userId := params["userId"]

	log.Printf("Generate Payment Request")

	request := &CreatePaymentRequestRequest{
		Quantity:      1,
		PaymentMethod: 0,
	}

	json.NewDecoder(r.Body).Decode(request)

	user := persistence.GetUser(userId)

}
