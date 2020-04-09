package main

import (
	"os"
	"fmt"
	"log"
	"flag"
	"syscall"
	"net/http"
	//"io/ioutil"
	"os/signal"
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

func cleanup() {
	fmt.Println("Closing server")
}

func main() {
	// Handle SIGTERM
	sigterm := make(chan os.Signal)
	signal.Notify(sigterm, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigterm
		cleanup()
		os.Exit(1)
	}()

	// Parse flags
	portFlag := flag.Int("port", 3030, "Port to serve HTTP server")
	flag.Parse()

	// Setup routes for API
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api", welcome)
	router.HandleFunc("/api/entries", getEntries).Methods("GET")

	// Start HTTP server
	httpPort := fmt.Sprintf(":%d", *portFlag)
	log.Fatal(http.ListenAndServe(httpPort, router))
}
