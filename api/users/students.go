package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Neiva07/queue-management-service/api"
	"github.com/Neiva07/queue-management-service/persistence"
)

type CreateStudentRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	CPF   string `json:"cpf"`
}

func CreateStudent(w http.ResponseWriter, r *http.Request) {

	log.Println("Receiving an create student request")

	request := &CreateStudentRequest{}

	json.NewDecoder(r.Body).Decode(request)

	user, err := persistence.CreateUser(request.Email, "student", request.CPF, request.Name)

	if err != nil {
		log.Printf("Error creating user with email %s", request.Email)
		log.Println(err.Error())
	}

	log.Printf("successfully created student with email %s", user.Email)

	m := api.Message(http.StatusCreated, user)
	api.Response(w, m)

}
