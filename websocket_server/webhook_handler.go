package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
)

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	fmt.Println("Received webhook:", string(body))
	w.WriteHeader(http.StatusOK)
}
