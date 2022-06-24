package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Neiva07/queue-management-service/api"
	"github.com/gorilla/mux"
)

func test(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	m := map[string]interface{}{"status": 200}
	api.Response(w, m)
}

func test2(res http.ResponseWriter, req *http.Request) {
	fmt.Fprint(res, "Hello World!")
}

var (
	certFile = flag.String("cert", "localhost.crt", "A PEM eoncoded certificate file.")
	keyFile  = flag.String("key", "localhost.key", "A PEM encoded private key file.")
	caFile   = flag.String("CA", "chain-pix-prod.crt", "A PEM eoncoded CA's certificate file.")
)

func main() {
	flag.Parse()

	// Load client cert
	cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
	if err != nil {
		log.Fatal(err)
	}

	// Load CA cert
	caCert, err := ioutil.ReadFile(*caFile)
	if err != nil {
		log.Fatal(err)
	}

	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(caCert)

	if !ok {
		log.Fatal("not ok")
	}

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	router := mux.NewRouter()

	router.HandleFunc("/", test2).Methods("GET")
	router.HandleFunc("/pix", test2).Methods("POST", "PUT")
	router.HandleFunc("/webhook/webhook", test2).Methods("POST", "PUT")
	router.HandleFunc("/webhook", test).Methods("POST", "PUT")
	router.HandleFunc("/webhook/pix", test2).Methods("POST", "PUT")
	log.Println("Running Webhook!")

	s := &http.Server{
		Addr:      ":443",
		Handler:   router, // use `http.DefaultServeMux`
		TLSConfig: tlsConfig,
	}

	err = s.ListenAndServeTLS("", "")
	// err := http.ListenAndServeTLS(":443", "localhost.crt", "localhost.key", router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
