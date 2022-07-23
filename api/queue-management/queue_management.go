package api

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/Neiva07/queue-management-service/api"
	"github.com/Neiva07/queue-management-service/persistence"
	"github.com/Neiva07/queue-management-service/services"
	"github.com/gorilla/mux"
)

type DequeueUserRequest struct {
	Quantity uint64 `json:"quantity"`
}

func EnqueueUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	params := mux.Vars(r)

	userId := params["userId"]

	u := persistence.GetUserById(ctx, userId)

	if u == nil {
		log.Println("User not found")
		m := api.Message(http.StatusNotFound, errors.New("User not found"))
		api.Response(w, m)
		return
	}

	if u.TicketQuantity <= 0 {
		log.Printf("User with email %v do not have any ticket to redeem\n", u.Email)
		m := api.Message(http.StatusPaymentRequired, errors.New("No tickets to redeem"))
		api.Response(w, m)
		return
	}

	u.TicketQuantity -= 1

	err := services.EnqueueUser(u)

	if err != nil {
		log.Println(err)
		m := api.Message(http.StatusInternalServerError, err)
		api.Response(w, m)
		return
	}

	err = persistence.RedeemTicket(userId)

	if err != nil {
		log.Println(err)
		m := api.Message(http.StatusUnprocessableEntity, err)
		api.Response(w, m)
		return
	}

	ticket, err := persistence.CreateTicket(u)

	if err != nil {
		log.Println(err)
		m := api.Message(http.StatusUnprocessableEntity, err)
		api.Response(w, m)
		return
	}

	log.Printf("User enqueued successfully %v", &u)

	m := api.Message(http.StatusOK, ticket)
	w.WriteHeader(http.StatusOK)
	api.Response(w, m)

}

func DequeueUser(w http.ResponseWriter, r *http.Request) {

	request := &DequeueUserRequest{}

	err := json.NewDecoder(r.Body).Decode(request)

	if err != nil {
		log.Println(err)
		m := api.Message(http.StatusBadRequest, err)
		api.Response(w, m)
		return
	}

	users := services.PopUsers(request.Quantity)

	//notify users

	if len(users) == 0 {
		log.Println("No users to dequeue")
		m := api.Message(http.StatusNotFound, users)
		w.WriteHeader(http.StatusNotFound)
		api.Response(w, m)
		return
	}

	log.Printf("User dequeued successfully %v", users)
	m := api.Message(http.StatusOK, users)
	w.WriteHeader(http.StatusOK)
	api.Response(w, m)
}

func ValidateTicket(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	ticketKey := params["ticketKey"]

	ticket, err := persistence.ValidateTicket(ticketKey)

	if err != nil {
		m := api.Message(http.StatusNotFound, nil)
		w.WriteHeader(http.StatusNotFound)
		api.Response(w, m)
		return
	}

	//notify user

	m := api.Message(http.StatusOK, ticket)
	w.WriteHeader(http.StatusOK)
	api.Response(w, m)
	return

}
