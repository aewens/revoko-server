package main

import (
	"fmt"
	"testing"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

func TestReadConfig(t *testing.T) {
	testFile := "/tmp/revoko_test_config.json"
	testPort := 3030
	testUser := "test"
	testPassword := "test"
	testJSON := `{
		"port": %d,
		"database": {
			"user": "%s",
			"password": "%s"
		}
	}`

	testData := []byte(fmt.Sprintf(testJSON, testPort, testUser, testPassword))
	err := ioutil.WriteFile(testFile, testData, 0644)
	
	if err != nil {
		t.Fatalf("Could not write temporary config file %s", testFile)
	}

	config, err := ReadConfig(testFile); if err != nil {
		t.Errorf("Could not read config file %s", testFile)
	}

	configPort := config.Port
	if configPort != testPort {
		t.Errorf("Parsed port as '%v' instead of '%d'", configPort, testPort)
	}

	configUser := config.User
	if configUser != testUser {
		t.Errorf("Parsed user as '%v' instead of '%s'", configUser, testUser)
	}

	configPassword := config.Password
	if configPassword != testPassword {
		t.Errorf("Parsed password as '%v' instead of '%s'", configPassword,
			testPassword)
	}
}

func TestGetEntries(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/entries", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetEntries)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("/api/entries got status code '%v' instead of '%v'", status,
			http.StatusOK)
	}

	testEntries := "[{\"ID\":1,\"ParentID\":0,\"Value\":\"test\"}]\n"
	rrEntries := rr.Body.String()
	if rrEntries != testEntries {
		t.Errorf("/api/entries returned '%v' instead of '%v'", rrEntries,
			testEntries)
	}
}
