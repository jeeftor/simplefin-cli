package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessAccountsHappyPath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`{"Accounts": [{"Org": {"Name": "TestOrg"}, "Name": "TestAccount", "Balance": "1000", "Currency": "USD"}]}`))
	}))
	defer server.Close()

	err := processAccounts(server.URL, "")
	assert.NoError(t, err)
}

func TestProcessAccountsWithInvalidURL(t *testing.T) {
	err := processAccounts("invalid-url", "")
	assert.Error(t, err)
}

func TestProcessAccountsWithNonOKStatusCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	err := processAccounts(server.URL, "")
	assert.Error(t, err)
}

func TestProcessAccountsWithInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`invalid-json`))
	}))
	defer server.Close()

	err := processAccounts(server.URL, "")
	assert.Error(t, err)
}

func TestProcessAccountsWithProxy(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`{"Accounts": [{"Org": {"Name": "TestOrg"}, "Name": "TestAccount", "Balance": "1000", "Currency": "USD"}]}`))
	}))
	defer server.Close()

	proxyServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		http.Redirect(rw, req, server.URL, http.StatusFound)
	}))
	defer proxyServer.Close()

	err := processAccounts(server.URL, proxyServer.URL)
	assert.NoError(t, err)
}

func TestProcessAccountsWithInvalidProxyURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`{"Accounts": [{"Org": {"Name": "TestOrg"}, "Name": "TestAccount", "Balance": "1000", "Currency": "USD"}]}`))
	}))
	defer server.Close()

	err := processAccounts(server.URL, "invalid-url")
	assert.Error(t, err)
}

func TestProcessErrors(t *testing.T) {
	fd := FinancialData{
		Errors: []string{"Error for TestOrg"},
		Accounts: []Accounts{
			{
				Org: Org{
					Name: "TestOrg",
				},
			},
		},
	}

	processErrors(&fd)
	assert.True(t, fd.Accounts[0].PossibleError)
}

func TestProcessErrorsWithNoMatchingOrg(t *testing.T) {
	fd := FinancialData{
		Errors: []string{"Error for NonExistentOrg"},
		Accounts: []Accounts{
			{
				Org: Org{
					Name: "TestOrg",
				},
			},
		},
	}

	processErrors(&fd)
	assert.False(t, fd.Accounts[0].PossibleError)
}

func TestTable(t *testing.T) {
	data, err := ioutil.ReadFile("test1.json")
	if err != nil {
		t.Fatalf("Failed to read test1.json: %s", err)
	}

	// Unmarshal the JSON data into FinancialData struct
	var fd FinancialData
	err = json.Unmarshal(data, &fd)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %s", err)
	}

	processErrors(&fd)

	printFinancialData(fd)
}
