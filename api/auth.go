package api

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func CreateToken(userID string) (string, error) {

	jwtSecret, isValid := os.LookupEnv("JWT_SECRET")

	if !isValid {
		err := errors.New("Error getting new jwt secret")
		return "", err
	}

	claim := jwt.StandardClaims{
		Issuer:    userID,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		IssuedAt:  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	return token.SignedString([]byte(jwtSecret))
}

func isTokenValid(token string) bool {

	jwtSecret, isValid := os.LookupEnv("JWT_SECRET")

	if !isValid {
		err := errors.New("Error getting new jwt secret")
		log.Println(err)
		return false
	}

	_, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {

		if _, isValid := t.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, fmt.Errorf("invalid token :%v", token)
		}

		return []byte(jwtSecret), nil
	})

	return err == nil
}
