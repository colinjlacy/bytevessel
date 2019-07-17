package main

import (
	"byteVessel/bundler"
	"byteVessel/cloudreach/dropbox"
	"byteVessel/emailer"
	"encoding/json"
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type EmailRequest struct {
	FolderName   string `json: folderName`
	EmailAddress string `json: emailAddress`
}

type CloudRequest struct {
	Foldername  string `json: folderName`
	Destination string `json: destination`
}

var email, password, token *string

func main() {
	// TODO should be passed via env vars
	// TODO authorize the app to use the service email account properly
	email = flag.String("e", "", "the email address you'd like to use to send scan files")
	password = flag.String("p", "", "the password to access the specified email address")
	token = flag.String("t", "", "the cloud storage auth token to use with this application")
	flag.Parse()

	// TODO: add user auth so that other users can add their dropbox account
	if *token == "" {
		log.Fatal("You'll need to provide a cloud token in order to use this application")
	}
	if *email == "" {
		log.Fatal("You must provide an email address in order to use this application")
	}
	if *password == "" {
		log.Fatal("You must provide an email password in order to use this application")
	}

	router := mux.NewRouter()
	router.HandleFunc("/email", emailFile).Methods("POST")
	router.HandleFunc("/store", storeFile).Methods("POST")
	log.Fatal(http.ListenAndServe(":9000", router))
}

func emailFile(w http.ResponseWriter, r *http.Request) {
	// parse req body
	var params EmailRequest
	_ = json.NewDecoder(r.Body).Decode(&params)
	filepath, err := bundler.Bundle(params.FolderName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonData := map[string]string{"error": err.Error()}
		_ = json.NewEncoder(w).Encode(jsonData)
		return
	}
	if err := emailer.EmailFile(filepath, params.EmailAddress, *email, *password); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonData := map[string]string{"error": err.Error()}
		_ = json.NewEncoder(w).Encode(jsonData)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte{})
}

func storeFile(w http.ResponseWriter, r *http.Request) {
	var params CloudRequest
	_ = json.NewDecoder(r.Body).Decode(&params)
	filepath, err := bundler.Bundle(params.Foldername)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonData := map[string]string{"error": err.Error()}
		_ = json.NewEncoder(w).Encode(jsonData)
		return
	}
	if err := dropbox.AscendFile(filepath, params.Destination, *token); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonData := map[string]string{"error": err.Error()}
		_ = json.NewEncoder(w).Encode(jsonData)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte{})
}

func retrieveFile() {

}
