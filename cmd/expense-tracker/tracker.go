package main

import (
	"expense-tracker/internal/cli"
	"expense-tracker/internal/expense"
	"expense-tracker/internal/httpapi"
	"expense-tracker/internal/storage"
	"fmt"
	"os"
)

const jsonFile = "expenses.json"
const csvFile = "expenses.csv"

func main() {
	repoFile, err := os.OpenFile(jsonFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	defer repoFile.Close()
	svc := expense.NewService(storage.NewJSONRepo(repoFile))
	if len(os.Args) < 2 {
		panic("no command provided")
	}
	if os.Args[1] == "export" {
		f, err := os.Create(csvFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		defer f.Close()
		if err := svc.ExportCSV(f); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		return
	}
	if os.Args[1] == "serve" {
		if err := httpapi.Run(svc, "localhost:8080"); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		return
	}
	if err := cli.Run(svc, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
