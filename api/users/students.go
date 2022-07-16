package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Neiva07/queue-management-service/api"
	"github.com/Neiva07/queue-management-service/persistence"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUserRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	CPF      string `json:"cpf"`
	Password string `json:"password"`
	UserType string `json:"userType"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterUserResponse struct {
	User  *persistence.User `json:"user"`
	Token string            `json:"token"`
}

type LoginResponse struct {
	User  *persistence.User `json:"user"`
	Token string            `json:"token"`
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {

	log.Println("Receiving a registration request")

	request := &RegisterUserRequest{}

	json.NewDecoder(r.Body).Decode(request)

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), 7)

	newUser := &persistence.User{
		CreatedAt:      time.Now(),
		Email:          request.Email,
		Name:           request.Name,
		UserType:       request.UserType,
		CPF:            request.CPF,
		TicketQuantity: 0,
		Password:       encryptedPassword,
		UserID:         primitive.NewObjectID(),
	}

	user, err := persistence.CreateUser(newUser)

	if err != nil {
		messageString := fmt.Sprint("Error creating user with email %s", request.Email)
		m := api.Message(http.StatusBadRequest, messageString)
		log.Printf(messageString)
		log.Println(err.Error())
		api.Response(w, m)
		return
	}

	log.Printf("successfully created user with email %s", user.Email)

	token, err := api.CreateToken(user.Email)

	if err != nil {
		messageString := fmt.Sprint("Failed creating token for user %s", request.Email)
		m := api.Message(http.StatusBadRequest, messageString)
		log.Printf(messageString)
		log.Println(err.Error())
		api.Response(w, m)
		return
	}

	response := RegisterUserResponse{
		User:  user,
		Token: token,
	}

	m := api.Message(http.StatusCreated, response)
	api.Response(w, m)

}

func Login(w http.ResponseWriter, r *http.Request) {
	log.Println("Receiving a login request")

	request := &LoginRequest{}

	json.NewDecoder(r.Body).Decode(request)

	user := persistence.GetUserByEmail(context.Background(), request.Email)

	if user == nil {
		messageString := fmt.Sprintf("User with %s not found", request.Email)
		m := api.Message(http.StatusNotFound, messageString)
		log.Printf(messageString)
		api.Response(w, m)
		return
	}

	log.Println("Checking password for user ", user.Email)

	err := bcrypt.CompareHashAndPassword(user.Password, []byte(request.Password))

	log.Println("Correct Password for user ", user.Email)

	if err != nil {
		messageString := fmt.Sprint("Incorrect password for email %s", request.Email)
		log.Printf(messageString)
		m := api.Message(http.StatusUnauthorized, messageString)
		log.Println(err.Error())
		api.Response(w, m)
		return
	}

	token, err := api.CreateToken(user.Email)

	if err != nil {
		messageString := fmt.Sprint("Failed creating token for user %s", request.Email)
		m := api.Message(http.StatusBadRequest, messageString)
		log.Printf(messageString)
		log.Println(err.Error())
		api.Response(w, m)
		return
	}

	response := &LoginResponse{
		User:  user,
		Token: token,
	}

	log.Println("Successfully log in user ", user.Email)

	m := api.Message(http.StatusOK, response)
	api.Response(w, m)
	return
}
