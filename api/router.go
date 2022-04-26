package api

import (
	"net/http"

	api "github.com/Neiva07/queue-management-service/api/payments"
	"github.com/gorilla/mux"
)

func init() {
	router := mux.NewRouter()
	router.HandleFunc("/users/{userId}/payment-requests", api.CreatePaymentRequest).Methods("POST")
	http.Handle("/api/v1/", router)

	http.ListenAndServe(":8000", router)
}
