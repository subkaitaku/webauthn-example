package main

import (
	"log"
	"net/http"
	"strings"
)

func BeginLogin(w http.ResponseWriter, r *http.Request) {
	// get username
	if r.Method != "GET" {
		http.NotFound(w, r)
		return
	}

	username := strings.TrimPrefix(r.URL.Path, "/login/begin/")
	user, err := usersDB.GetUser(username)
	if err != nil {
		log.Println(err)
		jsonResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// generate PublicKeyCredentialRequestOptions, session data
	options, sessionData, err := Wc.BeginLogin(user)
	if err != nil {
		log.Println(err)
		jsonResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "authentication",
		Value: sessionDb.StartSession(sessionData),
		Path:  "/",
	})

	jsonResponse(w, options, http.StatusOK)
}

func FinishLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}

	username := strings.TrimPrefix(r.URL.Path, "/login/finish/")
	user, err := usersDB.GetUser(username)
	if err != nil {
		log.Println(err)
		jsonResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// load the session data
	cookie, err := r.Cookie("authentication")
	if err != nil {
		log.Println("cookie:", err)
		jsonResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	sessionData, err := sessionDb.GetSession(cookie.Value)
	if err != nil {
		log.Println("session:", err)
		jsonResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// in an actual implementation, we should perform additional checks on
	// the returned 'credential', i.e. check 'credential.Authenticator.CloneWarning'
	// and then increment the credentials counter
	c, err := Wc.FinishLogin(user, *sessionData, r)
	if err != nil {
		log.Println(err)
		jsonResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if c.Authenticator.CloneWarning {
		log.Println("cloned key detected")
		jsonResponse(w, "cloned key detected", http.StatusBadRequest)
		return
	}

	sessionDb.DeleteSession(cookie.Value)
	log.Println("user: ", username, "finishing logging in")

	jsonResponse(w, "Login Success", http.StatusOK)
}
