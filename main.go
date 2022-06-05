package main

import (
	"log"
	"net/http"

	paymentApi "github.com/Neiva07/queue-management-service/api/payments"
	userApi "github.com/Neiva07/queue-management-service/api/users"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/users/{userId}/payment-requests", paymentApi.CreatePaymentRequest).Methods("POST")
	router.HandleFunc("/api/v1/users", userApi.CreateStudent).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))
}
