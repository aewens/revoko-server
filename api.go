package main

import (
	"fmt"
	"log"
	"flag"
	"net/http"
	"encoding/json"

	"github.com/gorilla/mux"
)

type Entry struct {
	ID			int		`json:"ID"`
	ParentID	int		`json:"ParentID"`	
	Value		string	`json:"Value"`
}

type Entries []Entry

var entries = Entries{
	{
		ID:			1,	
		ParentID:	0,
		Value:		"Test",
	},
};

func welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

func getEntries(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(entries)
}

func main() {
	portFlag := flag.Int("port", 3030, "Port to serve HTTP server")
	flag.Parse()

	httpPort := fmt.Sprintf(":%d", *portFlag)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api", welcome)
	router.HandleFunc("/api/entries", getEntries).Methods("GET")
	log.Fatal(http.ListenAndServe(httpPort, router))
}
