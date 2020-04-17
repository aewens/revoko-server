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
	testDBUri := "http://0.0.0.0:5984"
	testDBUser := "test"
	testDBPassword := "test"
	testJSON := `{
		"port": %d,
		"database": {
			"uri": "%s",
			"user": "%s",
			"password": "%s"
		}
	}`

	testData := []byte(fmt.Sprintf(testJSON, testPort, testDBUri, testDBUser,
		testDBPassword))
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

	DBUri := config.DBUri
	if DBUri != testDBUri {
		t.Errorf("Parsed database uri as '%v' instead of '%s'", DBUri,
			testDBUri)
	}

	DBUser := config.DBUser
	if DBUser != testDBUser {
		t.Errorf("Parsed database user as '%v' instead of '%s'", DBUser,
			testDBUser)
	}

	DBPassword := config.DBPassword
	if DBPassword != testDBPassword {
		t.Errorf("Parsed password as '%v' instead of '%s'", DBPassword,
			testDBPassword)
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
