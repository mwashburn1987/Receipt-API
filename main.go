package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// store receipt responses in memory with map
var ReceiptResponses = make(map[string]ReceiptPointsResponse)

func main() {
	// create router
	r := mux.NewRouter().StrictSlash(true)
	// list handlers
	r.HandleFunc("/receipts/process", ProcessReceipt).Methods("POST")
	r.HandleFunc("/receipts/{id}/point", GetPoints).Methods("GET")

	//listen for incoming requests
	log.Fatal(http.ListenAndServe(":8080", r))
}
