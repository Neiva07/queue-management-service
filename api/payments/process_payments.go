package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Neiva07/queue-management-service/api"
	"github.com/Neiva07/queue-management-service/persistence"
	"github.com/Neiva07/queue-management-service/services"
	"github.com/gorilla/mux"
)

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

func CheckUserPayments(w http.ResponseWriter, r *http.Request) {
	log.Println("Checking payments...")

	params := mux.Vars(r)

	userId := params["userId"]

	payments, err := persistence.GetUnprocessedPaymentsByID(userId)

	log.Println(payments)

	if err != nil {
		m := api.Message(http.StatusNotFound, fmt.Sprint("Error trying to get payments for user %s", userId))
		log.Println(m)
		api.Response(w, m)
		return
	}

	for _, payment := range payments {

		chargeDetails := services.CheckCharge(payment.PaymentKey)

		log.Println(chargeDetails)

		if payment.Status != chargeDetails.Status && chargeDetails.Status == "CONCLUIDA" {
			payment, err := persistence.ProcessPayment(payment.PaymentKey, chargeDetails.Status)

			if err != nil {
				log.Println("Error trying to process payment for key %s with error :%s", payment.PaymentKey, err.Error())
				continue
			}
			err = persistence.BuyTickets(userId, payment.Quantity)

			if err != nil {
				log.Println("Error trying to update tickets for user %s with error %s", userId, err.Error())
			}
		}
	}

	user := persistence.GetUserById(context.Background(), userId)

	if user == nil {
		m := api.Message(http.StatusNotFound, fmt.Sprint("Error trying to get user tickets for user %s", userId))
		log.Println(m)
		api.Response(w, m)
		return
	}

	m := api.Message(200, &user)
	api.Response(w, m)
}
