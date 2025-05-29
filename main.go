package main

import (
	"appsec/accountholder"
	"appsec/apiUtil"
	"appsec/db"
	"appsec/destination"
	"appsec/incomingdata"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	// Open DB Connection
	lErr := db.OpenConn()
	if lErr != nil {
		log.Println("Error while connecting db -", lErr)
	} else {
		// Close DB Connection
		defer db.CloseConn()

		// Initiate http.Client
		apiUtil.Init()

		log.Print("Server Started")

		r := mux.NewRouter()
		// API to store or retrieve the account details
		r.HandleFunc("/account", accountholder.AccountHolder)
		// API to store or retrieve the destination details
		r.HandleFunc("/destination", destination.Designation)
		// API to retrieve the destinations of the given account
		r.HandleFunc("/getaccdest", destination.GetAccDestination)
		// API to identify the destinations of the given secret and pass the jsondata to the destination.
		r.HandleFunc("/server/incoming_data", incomingdata.PassIncomingData)

		// Serving running on the port.
		http.ListenAndServe(":8888", r)
	}
}
