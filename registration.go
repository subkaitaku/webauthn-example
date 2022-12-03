package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/go-webauthn/webauthn/protocol"
)

func BeginRegistration(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.NotFound(w, r)
		return
	}

	username := strings.TrimPrefix(r.URL.Path, "/register/begin/")
	user, err := usersDB.GetUser(username)
	if err != nil {
		displayName := strings.Split(username, "@")[0]
		user = NewUser(username, displayName)
		usersDB.PutUser(user)
	}

	registerOptions := func(credCreationOpts *protocol.PublicKeyCredentialCreationOptions) {
		credCreationOpts.CredentialExcludeList = user.CredentialExcludeList()
	}

	options, sessionData, err := Wc.BeginRegistration(
		user,
		registerOptions,
	)

	if err != nil {
		log.Println(err)
		jsonResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "registration",
		Value: sessionDb.StartSession(sessionData),
		Path:  "/",
	})

	jsonResponse(w, options, http.StatusOK)
}

func FinishRegistration(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}

	username := strings.TrimPrefix(r.URL.Path, "/register/finish/")
	user, err := usersDB.GetUser(username)
	if err != nil {
		log.Println(err)
		jsonResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("registration")
	if err != nil {
		log.Println("cookie:", err)
		jsonResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	sessionData, err := sessionDb.GetSession(cookie.Value)
	if err != nil {
		log.Println("cookie:", err)
		jsonResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	credential, err := Wc.FinishRegistration(user, *sessionData, r)
	if err != nil {
		log.Println("finalising: ", err)
		jsonResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.AddCredential(*credential)

	sessionDb.DeleteSession(cookie.Value)
	log.Println("user: ", username, "finishing registration")

	jsonResponse(w, "Registration Success", http.StatusOK)
}
