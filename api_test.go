package main

import (
	"fmt"
	"testing"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

func TestReadConfig(t *testing.T) {
	testPort := 3030
	testFile := "/tmp/revoko_test_config.json"
	testData := []byte(fmt.Sprintf("{\"port\": %d}", testPort))
	err := ioutil.WriteFile(testFile, testData, 0644)
	
	if err != nil {
		t.Fatalf("Could not write temporary config file %s", testFile)
	}

	config, err := ReadConfig(testFile); if err != nil {
		t.Errorf("Could not read config file %s", testFile)
	}

	configPort := config.Port
	if configPort != testPort {
		t.Errorf("Parsed port as '%d' instead of '%d'", configPort, testPort)
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
