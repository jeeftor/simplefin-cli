package main

import (
	"encoding/json"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
)

var (
	version   = "dev"  // Default value for development
	commit    = "none" // Short commit hash
	buildTime = "unknown"
)

func main() {
	app := &cli.App{
		Name:  "sf",
		Usage: "A CLI interface to simplefin.org",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "url",
				Usage:    "Your specific SimpleFIN Access URL",
				EnvVars:  []string{"SF_URL"},
				Required: true,
			},
			&cli.StringFlag{
				Name:    "proxy",
				Usage:   "Set the proxy URL",
				EnvVars: []string{"SF_PROXY"},
			},
			&cli.StringFlag{
				Name:    "out",
				Usage:   "Output filename for JSON results",
				EnvVars: []string{"SF_OUT"},
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "version",
				Usage: "Print the version",
				Action: func(c *cli.Context) error {
					fmt.Printf("Version: %s\nCommit: %s\nBuilt at: %s\n", version, commit, buildTime)
					return nil
				},
			},
		},
		Action: func(c *cli.Context) error {
			baseURL := c.String("url")
			proxyURL := c.String("proxy")
			outFilename := c.String("out")
			err := processAccounts(baseURL, proxyURL, outFilename)
			if err != nil {
				return err
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func processAccounts(baseURL string, proxyURL string, outFilename string) error {
	accountsURL := baseURL + "/accounts"

	client := &http.Client{}
	if proxyURL != "" {
		proxy, err := url.Parse(proxyURL)
		if err != nil {
			return fmt.Errorf("invalid proxy URL: %v", err)
		}
		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxy)}
	}

	resp, err := client.Get(accountsURL)
	if err != nil {
		return fmt.Errorf("error making GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK HTTP status code: %d", resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}
	var financialData FinancialData
	err = json.Unmarshal(bodyBytes, &financialData)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	// Process the FinancialData as needed
	//fmt.Printf("Received financial data: %+v\n", financialData)
	processErrors(&financialData)

	if outFilename != "" {
		// Write the JSON results to the specified file
		err := writeJSONToFile(financialData, outFilename)
		if err != nil {
			return err
		}
	}

	printFinancialData(financialData)
	return nil
}

//	var orgColors = map[string]text.Colors{
//		"Citibank":             text.Colors{text.FgBlue},
//		"Fidelity Investments": text.Colors{text.FgGreen},
//		"Hanscom Federal CU":   text.Colors{text.FgRed},
//		"Wealthfront":          text.Colors{text.FgYellow},
//		// Add more organizations and colors as needed
//	}
func writeJSONToFile(data FinancialData, filename string) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return err
	}

	fmt.Printf("JSON results written to: %s\n", filename)
	return nil
}
func processErrors(financialData *FinancialData) {
	for _, errorString := range financialData.Errors {
		for i, account := range financialData.Accounts {
			if strings.Contains(errorString, account.Org.Name) {
				financialData.Accounts[i].PossibleError = true
			}
		}
	}
}

func sortAccountsByOrgNameAndAccountName(accounts []Accounts) []Accounts {
	// Define a custom sorting function
	customSort := func(accounts []Accounts, less func(i, j int) bool) {
		sort.SliceStable(accounts, less)
	}

	// Sort by Org.Name and then by Name
	customSort(accounts, func(i, j int) bool {
		if accounts[i].Org.Name != accounts[j].Org.Name {
			return accounts[i].Org.Name < accounts[j].Org.Name
		}
		return accounts[i].Name < accounts[j].Name
	})

	return accounts
}

func printFinancialData(financialData FinancialData) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Bank", "Account Name", "Balance", "Currency"})

	sortedAccounts := sortAccountsByOrgNameAndAccountName(financialData.Accounts)
	prevOrgName := sortedAccounts[0].Org.Name

	for _, account := range sortedAccounts {
		status := ""
		if account.PossibleError {
			status = "⚠"
		}
		// Check if the organization name has changed
		if account.Org.Name != prevOrgName {
			t.AppendSeparator()
			prevOrgName = account.Org.Name
		}
		t.AppendRow([]interface{}{
			account.Org.Name,
			account.Name,
			account.Balance,
			account.Currency,
			status,
		})
	}

	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true, Align: text.AlignCenter},
		{Number: 3, Align: text.AlignRight},
	})
	//t.SetStyle(table.Style{
	//	Name: "colorful",
	//	Box:  table.StyleBoxLight,
	//	Color: table.ColorOptions{
	//		Header: text.Colors{text.Bold, text.FgCyan},
	//		Row:    text.Colors{text.FgHiYellow},
	//	},
	//})

	//t.SetStyle(table.Style{})

	t.SetStyle(table.StyleLight) // This style includes lines between columns
	//t.SetStyle(table.StyleColoredBright)

	t.Render()
}
