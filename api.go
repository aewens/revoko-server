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

var testEntries = Entries{
	{
		ID:			1,	
		ParentID:	0,
		Value:		"test",
	},
};

type Config struct {
	Port int
}

func Welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

func GetEntries(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(testEntries)
}

func cleanup() {
	fmt.Println("Closing server")
}

func ReadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	//buf, err := ioutil.ReadAll(file)
	//if err != nil {
	//	log.Fatal(err)
	//}

	raw := make(map[string]interface{})
	json.NewDecoder(file).Decode(&raw)

	port, valid := raw["port"].(float64)
	if !valid {
		fmt.Printf("%v | %T\n", raw, raw["port"])
		port = float64(3030)
	}

	config := &Config{Port: int(port)}

	return config, err
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
	//portFlag := flag.Int("port", 3030, "Port to serve HTTP server")
	configFlag := flag.String("config", "", "Path to config file")
	flag.Parse()

	if len(*configFlag) == 0 {
		log.Fatal("ERROR: Config flag is missing!")
		os.Exit(2)
	}

	config, err := ReadConfig(*configFlag)
	if err != nil {
		log.Fatal(err)
		os.Exit(3)
	}

	// Setup routes for API
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api", Welcome)
	router.HandleFunc("/api/entries", GetEntries).Methods("GET")

	// Start HTTP server
	httpPort := fmt.Sprintf(":%d", config.Port)
	fmt.Printf("Starting server on port %d\n", config.Port)
	log.Fatal(http.ListenAndServe(httpPort, router))
}
