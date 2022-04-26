package services

import (
	"log"

	uuid "github.com/nu7hatch/gouuid"
)

func GeneratePixKey() {

	key, err := uuid.NewV4()

	if err != nil {
		log.Fatalf("Something went wrong creating pix key")
	}

}
