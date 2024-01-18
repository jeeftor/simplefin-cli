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
			err := processAccounts(baseURL, proxyURL)
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

func processAccounts(baseURL string, proxyURL string) error {
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

	printFinancialData(financialData)
	return nil
}

//
//var orgColors = map[string]text.Colors{
//	"Citibank":             text.Colors{text.FgBlue},
//	"Fidelity Investments": text.Colors{text.FgGreen},
//	"Hanscom Federal CU":   text.Colors{text.FgRed},
//	"Wealthfront":          text.Colors{text.FgYellow},
//	// Add more organizations and colors as needed
//}

func processErrors(financialData *FinancialData) {
	for _, errorString := range financialData.Errors {
		for i, account := range financialData.Accounts {
			if strings.Contains(errorString, account.Org.Name) {
				financialData.Accounts[i].PossibleError = true
			}
		}
	}
}

func printFinancialData(financialData FinancialData) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Bank", "Account Name", "Balance", "Currency"})

	for _, account := range financialData.Accounts {
		status := ""
		if account.PossibleError {
			status = "⚠"
		}

		t.AppendRow([]interface{}{
			account.Org.Name,
			account.Name,
			account.Balance,
			account.Currency,
			status,
		})
	}

	t.SortBy([]table.SortBy{
		{Name: "Org Name", Mode: table.Asc},
		{Name: "Account Name", Mode: table.Asc},
	})

	t.SetColumnConfigs([]table.ColumnConfig{
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

	// t.SetStyle(table.Style{
	//     Name: "Custom Style",
	//     Box: table.StyleBox{
	//         MiddleHorizontal: "-", // Use '-' for horizontal lines
	//         MiddleVertical:   "|", // Use '|' for vertical lines
	//         MiddleSeparator:  "+", // Use '+' for the middle cross-section
	//     },
	//     Color: table.ColorOptions{
	//         Header: text.Colors{text.Bold, text.FgCyan},
	//         Row:    text.Colors{text.FgHiYellow},
	//     },
	// })

	t.SetStyle(table.StyleLight) // This style includes lines between columns
	//t.SetStyle(table.StyleColoredBright)

	t.Render()
}
