package main

import (
	"os"
	"fmt"
	"log"
	"flag"
	"errors"
	"syscall"
	"net/http"
	"io/ioutil"
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

func Welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

func GetEntries(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(testEntries)
}

func cleanup() {
	fmt.Println("Closing server")
}

type Config struct {
	Port int
	DBUri string
	DBUser string
	DBPassword string
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

	rawPort, ok := raw["port"]; if !ok {
		log.Fatalf("Missing 'port' from config file: %s", path)
	}

	port, valid := rawPort.(float64); if !valid {
		log.Fatalf("Value of 'port' is not a number: %v", port)
	}

	rawDatabase, ok := raw["database"]; if !ok {
		log.Fatalf("Missing 'database' from config file: %s", path)
	}

	database, valid := rawDatabase.(map[string]interface{}); if !valid {
		log.Fatalf("Value of 'port' is not a dictionary: %v", rawDatabase)
	}

	rawUri, ok := database["uri"]; if !ok {
		log.Fatalf("Missing 'uri' from database in config file: %s", path)
	}

	uri, valid := rawUri.(string); if !valid {
		log.Fatalf("Value of 'uri' is not a string: %v", rawUri)
	}

	rawUser, ok := database["user"]; if !ok {
		log.Fatalf("Missing 'user' from database in config file: %s", path)
	}

	user, valid := rawUser.(string); if !valid {
		log.Fatalf("Value of 'user' is not a string: %v", rawUser)
	}

	rawPassword, ok := database["password"]; if !ok {
		log.Fatalf("Missing 'password' from database in config file: %s", path)
	}

	password, valid := rawPassword.(string); if !valid {
		log.Fatalf("Value of 'password' is not a string: %v", rawPassword)
	}

	config := &Config{
		Port: int(port),
		DBUri: uri,
		DBUser: user,
		DBPassword: password,
	}

	return config, err
}

func StartServer(config *Config) {
	// Setup routes for API
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api", Welcome)
	router.HandleFunc("/api/entries", GetEntries).Methods("GET")

	// Start HTTP server
	httpPort := fmt.Sprintf(":%d", config.Port)
	fmt.Printf("Starting server on port %d\n", config.Port)
	log.Fatal(http.ListenAndServe(httpPort, router))
}

func HandleSigterm() {
	sigterm := make(chan os.Signal)
	signal.Notify(sigterm, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigterm
		cleanup()
		os.Exit(1)
	}()
}

func DatabaseGet(config *Config, path string) (string, error) {
	client := &http.Client{}
	dbPath := fmt.Sprintf("%s%s", config.DBUri, path)
	req, err := http.NewRequest("GET", dbPath, nil); if err != nil {
		log.Fatal(err)
		return "", err
	}

	req.SetBasicAuth(config.DBUser, config.DBPassword)
	resp, err := client.Do(req); if err != nil {
		log.Fatal(err)
		return "", err
	}

	if status := resp.StatusCode; status != http.StatusOK {
		errMessage := "%s got status code '%v' instead of '%v'"
		err := errors.New(fmt.Sprintf(errMessage, dbPath, status,http.StatusOK))
		log.Fatal(err)
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body); if err != nil {
		log.Fatal(err)
		return "", err
	}

	return string(body), nil
}

func AllDatabases(config *Config) ([]string, error) {
	var dbs []string

	resp, err := DatabaseGet(config, "/_all_dbs"); if err == nil {
		json.Unmarshal([]byte(resp), &dbs)
	}

	return dbs, err
}

func ParseFlags() map[string]interface{} {
	configFlag := flag.String("config", "", "Path to config file")
	flag.Parse()

	if len(*configFlag) == 0 {
		log.Fatal("ERROR: Config flag is missing!")
		os.Exit(2)
	}

	flags := make(map[string]interface{})
	flags["config"] = *configFlag

	return flags
}

func main() {
	HandleSigterm()

	flags := ParseFlags()
	configFlag, ok := flags["config"]; if !ok {
		log.Fatal("ERROR: Config flag is missing!")
		os.Exit(2)
	}

	configPath, ok := configFlag.(string); if !ok {
		log.Fatalf("Config flag is not a string: %v", configFlag)
	}

	config, err := ReadConfig(configPath); if err != nil {
		log.Fatal(err)
		os.Exit(3)
	}

	dbs, err := AllDatabases(config); if err != nil {
		log.Fatal(err)
		os.Exit(4)
	}

	fmt.Println("Databases:")
	for _, db := range dbs {
		fmt.Printf("- %s\n", db)
	}
	fmt.Printf("\n")

	StartServer(config)
}
