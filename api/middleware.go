package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

var UNAUTHENTICATED_ROUTES map[string]bool

func init() {
	UNAUTHENTICATED_ROUTES = make(map[string]bool)
	UNAUTHENTICATED_ROUTES["/api/v1/signin"] = true
	UNAUTHENTICATED_ROUTES["/api/v1/signup"] = true
}

func AuthMiddleware(next http.Handler) http.Handler {

	log.Println("Checking authentication...")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if ok := UNAUTHENTICATED_ROUTES[r.URL.Path]; ok {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			fmt.Println("Malformed token")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Malformed Token"))
			return
		}
		token := authHeader[1]

		if !isTokenValid(token) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid Token"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
