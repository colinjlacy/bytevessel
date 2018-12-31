package main

import (
	"byteVessel/bundler"
	"byteVessel/emailer"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type EmailRequest struct {
	Foldername   string `json: foldername`
	EmailAddress string `json: emailAddress`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/email", emailFile).Methods("POST")
	log.Fatal(http.ListenAndServe(":9000", router))
}

func emailFile(w http.ResponseWriter, r *http.Request) {
	// parse req body
	var params EmailRequest
	_ = json.NewDecoder(r.Body).Decode(&params)
	filepath, err := bundler.Bundle(params.Foldername)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonData := map[string]string{"error": err.Error()}
		_ = json.NewEncoder(w).Encode(jsonData)
		return
	}
	if err := emailer.EmailFile(filepath, params.EmailAddress); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonData := map[string]string{"error": err.Error()}
		_ = json.NewEncoder(w).Encode(jsonData)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte{})
}

func storeFile() {

}

func retrieveFile() {

}
