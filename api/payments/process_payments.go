package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Neiva07/queue-management-service/api"
	"github.com/Neiva07/queue-management-service/persistence"
	"github.com/gorilla/mux"
)

// type PIX struct {
// 	Key         string `json:"chave"`
// 	Type        string `json:"tipo"`
// 	TotalValue  string `json:"valor"`
// 	ProcessedAt string `json:"horario"`
// 	Status      string `json:"status"`
// }

// type ProcessPaymentRequest struct {
// 	PIX []PIX `json:"pix"`
// }

// func ProcessPayment(w http.ResponseWriter, r *http.Request) {
// 	log.Printf("Processing Payment from PIX")

// 	request := &ProcessPaymentRequest{}

// 	json.NewDecoder(r.Body).Decode(request)

// 	log.Println(request)

// 	m := map[string]interface{}{}

// 	json.NewDecoder(r.Body).Decode(&m)

// 	log.Println(m)

// }
type ProcessPaymentRequest struct {
	Status string `json:"status"`
}

func ProcessPayment(w http.ResponseWriter, r *http.Request) {
	log.Printf("Processing Payment from PIX")

	request := &ProcessPaymentRequest{}

	params := mux.Vars(r)

	paymentKey := params["paymentKey"]
	userId := params["userId"]

	err := json.NewDecoder(r.Body).Decode(request)

	if err != nil {
		m := api.Message(http.StatusBadRequest, fmt.Sprint("Error trying to decode process payment request body :%s", err.Error()))
		api.Response(w, m)
		return
	}

	payment, err := persistence.ProcessPayment(paymentKey, request.Status)

	if err != nil {
		m := api.Message(http.StatusUnprocessableEntity, fmt.Sprint("Error trying to process payment for key %s with error :%s", paymentKey, err.Error()))
		api.Response(w, m)
		return
	}

	err = persistence.BuyTickets(userId, payment.Quantity)

	if err != nil {
		m := api.Message(http.StatusUnprocessableEntity, fmt.Sprint("Error trying to update tickets for user %s with error %s", userId, err.Error()))
		api.Response(w, m)
		return
	}

	log.Println("Payment processed successfully for key %s and status %s", paymentKey, request.Status)

	w.WriteHeader(http.StatusOK)
	api.Response(w, nil)

}
