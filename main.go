package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-webauthn/webauthn/webauthn"
)

var (
	Wc *webauthn.WebAuthn
	err error
)

func main() {
	Wc, err = webauthn.New(&webauthn.Config{
		RPDisplayName: "Example.",          // Display Name for your site
		RPID:          "localhost",             // Generally the domain name for your site
		RPOrigin:      "http://localhost:8080", // The origin URL for WebAuthn requests
	})

	if err != nil {
		log.Fatal("failed to create WebAuthn from config:", err)
	}

	http.HandleFunc("/register/begin/", BeginRegistration)
	http.HandleFunc("/register/finish/", FinishRegistration)

	http.HandleFunc("/login/begin/", BeginLogin)
	http.HandleFunc("/login/finish/", FinishLogin)

	// index.html
	http.Handle("/", http.FileServer(http.Dir("./templates/")))

	serverAddress := "localhost:8080"

	log.Println("starting server at", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, nil))
}

func jsonResponse(w http.ResponseWriter, d interface{}, c int) {
	dj, err := json.Marshal(d)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(c)
	fmt.Fprintf(w, "%s", dj)
}
