package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Starting webhook server on :8080")

	http.HandleFunc("/webhook", webhookHandler)

	http.ListenAndServe(":8080", nil)
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received webhook request")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Webhook received"))
}
