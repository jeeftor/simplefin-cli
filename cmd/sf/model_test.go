package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestFinancialDataFields(t *testing.T) {
	ag := FinancialData{
		Errors: []string{"error1", "error2"},
		Accounts: []Accounts{
			{
				ID: "123",
			},
		},
		XAPIMessage: &[]string{"message1", "message2"},
	}

	assert.Equal(t, []string{"error1", "error2"}, ag.Errors)
	assert.Equal(t, "123", ag.Accounts[0].ID)
	assert.Equal(t, &[]string{"message1", "message2"}, ag.XAPIMessage)
}

func TestOrgFields(t *testing.T) {
	org := Org{
		Domain:  "example.com",
		Name:    "Example",
		SfinURL: "sfin.example.com",
		URL:     "www.example.com",
	}

	assert.Equal(t, "example.com", org.Domain)
	assert.Equal(t, "Example", org.Name)
	assert.Equal(t, "sfin.example.com", org.SfinURL)
	assert.Equal(t, "www.example.com", org.URL)
}

func TestTransactionsFields(t *testing.T) {
	trans := Transactions{
		ID:          "trans1",
		Posted:      123456,
		Amount:      "100.00",
		Description: "Test transaction",
	}

	assert.Equal(t, "trans1", trans.ID)
	assert.Equal(t, 123456, trans.Posted)
	assert.Equal(t, "100.00", trans.Amount)
	assert.Equal(t, "Test transaction", trans.Description)
}

func TestAccountsFields(t *testing.T) {
	acc := Accounts{
		Org: Org{
			Domain: "example.com",
		},
		ID:               "acc1",
		Name:             "Account 1",
		Currency:         "USD",
		Balance:          "1000.00",
		AvailableBalance: "900.00",
		BalanceDate:      123456789,
		Transactions: []Transactions{
			{
				ID: "trans1",
			},
		},
		Extra: map[string]interface{}{
			"key": "value",
		},
	}

	assert.Equal(t, "example.com", acc.Org.Domain)
	assert.Equal(t, "acc1", acc.ID)
	assert.Equal(t, "Account 1", acc.Name)
	assert.Equal(t, "USD", acc.Currency)
	assert.Equal(t, "1000.00", acc.Balance)
	assert.Equal(t, "900.00", acc.AvailableBalance)
	assert.Equal(t, 123456789, acc.BalanceDate)
	assert.Equal(t, "trans1", acc.Transactions[0].ID)
	assert.Equal(t, map[string]interface{}{"key": "value"}, acc.Extra)
}

func TestUnmarshalFinancialData(t *testing.T) {
	// Read the JSON file
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

	assert.Equal(t, len(fd.Accounts), 12)
	// Test the fields

}
