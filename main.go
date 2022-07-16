package main

import (
	"log"
	"net/http"

	"github.com/Neiva07/queue-management-service/api"
	paymentApi "github.com/Neiva07/queue-management-service/api/payments"
	queueApi "github.com/Neiva07/queue-management-service/api/queue-management"
	userApi "github.com/Neiva07/queue-management-service/api/users"
	"github.com/joho/godotenv"

	"github.com/gorilla/mux"
)

func main() {

	godotenv.Load()

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/users/{userId}/payments", paymentApi.CreatePaymentRequest).Methods("POST")
	router.HandleFunc("/api/v1/users/{userId}/payments", paymentApi.CheckUserPayments).Methods("PATCH")
	router.HandleFunc("/api/v1/users/{userId}/payments/{paymentKey}", paymentApi.ProcessPayment).Methods("PATCH")
	router.HandleFunc("/api/v1/queue-management/users/{userId}", queueApi.EnqueueUser).Methods("POST")
	router.HandleFunc("/api/v1/queue-management/users", queueApi.DequeueUser).Methods("POST")
	router.Use(api.AuthMiddleware)

	router.HandleFunc("/api/v1/signup", userApi.RegisterUser).Methods("POST")
	router.HandleFunc("/api/v1/signin", userApi.Login).Methods("POST")
	router.HandleFunc("/api/v1/webhook", api.CreateWebhook).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", router))
}
