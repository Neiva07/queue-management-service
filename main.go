package main

import (
	"log"
	"net/http"

	"github.com/Neiva07/queue-management-service/api"
	paymentApi "github.com/Neiva07/queue-management-service/api/payments"
	userApi "github.com/Neiva07/queue-management-service/api/users"
	"github.com/joho/godotenv"

	"github.com/gorilla/mux"
)

func main() {

	godotenv.Load()

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/users/{userId}/payments", paymentApi.CreatePaymentRequest).Methods("POST")
	router.HandleFunc("/api/v1/users/{userId}/payments/{paymentKey}", paymentApi.ProcessPayment).Methods("PATCH")
	router.HandleFunc("/api/v1/users", userApi.CreateStudent).Methods("POST")
	router.HandleFunc("/api/v1/webhook", api.CreateWebhook).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", router))
}
