package main

type FinancialData struct {
	Errors      []string   `json:"errors"`
	Accounts    []Accounts `json:"accounts"`
	XAPIMessage *[]string  `json:"x-api-message,omitempty"`
}
type Org struct {
	Domain  string `json:"domain"`
	Name    string `json:"name"`
	SfinURL string `json:"sfin-url"`
	URL     string `json:"url"`
}
type Transactions struct {
	ID          string `json:"id"`
	Posted      int    `json:"posted"`
	Amount      string `json:"amount"`
	Description string `json:"description"`
}
type Extra map[string]interface{}

type Accounts struct {
	Org              Org            `json:"org"`
	ID               string         `json:"id"`
	Name             string         `json:"name"`
	Currency         string         `json:"currency"`
	Balance          string         `json:"balance"`
	AvailableBalance string         `json:"available-balance"`
	BalanceDate      int            `json:"balance-date"`
	Transactions     []Transactions `json:"transactions"`
	Extra            Extra          `json:"extra,omitempty"`
	Holdings         []any          `json:"holdings"`
	PossibleError    bool           `json:"possible_error,omitempty"`
}
