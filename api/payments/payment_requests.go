package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	api "github.com/Neiva07/queue-management-service/api"
	"github.com/Neiva07/queue-management-service/persistence"
	"github.com/Neiva07/queue-management-service/services"
	"github.com/gorilla/mux"
)

type CreatePaymentRequestRequest struct {
	Quantity      uint64 `json:"quantity"`
	PaymentMethod uint8  `json:"paymentMethod"`
}

type CreatePaymentRequestResponse struct {
	QRCodeKey    string `json:"qrCodeKey"`
	QRCodeImage  string `json:"qrCodeImage"`
	PaymentValue uint64 `json:"paymentValue"` //centavos
	Status       string `json:"status"`
}

func CreatePaymentRequest(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	params := mux.Vars(r)

	userId := params["userId"]

	log.Printf("Generating Payment Request")

	request := &CreatePaymentRequestRequest{
		Quantity:      1,
		PaymentMethod: 0,
	}

	json.NewDecoder(r.Body).Decode(request)

	user := persistence.GetUserById(ctx, userId)

	if user == nil {
		m := api.Message(http.StatusNotFound, "User Not Found")
		api.Response(w, m)
		return
	}

	totalPayment := user.GetTotalPayment(request.Quantity)

	generatePixResponse, err := services.GeneratePixKey(user, totalPayment)

	if err != nil {
		m := api.Message(http.StatusBadRequest, fmt.Sprint("Something went wrong generating PIX", err.Error()))
		api.Response(w, m)
		return
	}

	payment := &persistence.Payment{
		Quantity:      request.Quantity,
		User:          *user,
		QRCodeKey:     generatePixResponse.QRCode.Key,
		PaymentKey:    generatePixResponse.TXID,
		PaymentMethod: request.PaymentMethod,
		Status:        generatePixResponse.Status,
		UnitPrice:     persistence.USER_TYPE_PRICE_MAP[user.UserType],
	}

	err = persistence.CreatePayment(payment)

	if err != nil {
		message := fmt.Sprint("Error recording payment in the database", err.Error())
		log.Println(message)
		m := api.Message(http.StatusBadRequest, message)
		api.Response(w, m)
		return
	}

	log.Printf("Created payment sucessfully for user %v", user.Email)

	paymentRequestResponse := CreatePaymentRequestResponse{
		QRCodeKey:    generatePixResponse.QRCode.Key,
		QRCodeImage:  generatePixResponse.QRCode.Image,
		PaymentValue: totalPayment,
		Status:       generatePixResponse.Status,
	}

	responseBody := api.Message(http.StatusCreated, paymentRequestResponse)

	log.Println(responseBody)

	api.Response(w, responseBody)
}
