package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Neiva07/queue-management-service/services"
)

type CreateWebhookRequest struct {
	URL string `json:"url"`
}

func CreateWebhook(w http.ResponseWriter, r *http.Request) {

	log.Println("Creating Webhook")

	request := &CreateWebhookRequest{}
	json.NewDecoder(r.Body).Decode(request)

	log.Printf("Creating webhook for url %v", request.URL)

	res, err := services.AddPIXWebhook(request.URL)

	if err != nil {
		m := Message(http.StatusBadRequest, fmt.Sprint("Error registering new webhook", err.Error()))
		Response(w, m)
		return
	}

	m := Message(http.StatusOK, res)

	Response(w, m)
}
