package main

import (
	"expense-tracker/internal/cli"
	"expense-tracker/internal/expense"
	"expense-tracker/internal/storage"
	"os"
)

const jsonFile = "expenses.json"
const csvFile = "expenses.csv"

func main() {
	repoFile, err := os.OpenFile(jsonFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer repoFile.Close()
	svc := expense.NewService(storage.NewJSONRepo(repoFile))
	if os.Args[1] == "export" {
		f, err := os.Create(csvFile)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		if err := svc.ExportCSV(f); err != nil {
			panic(err)
		}
		return
	}
	if err := cli.Run(svc, os.Args[1:]); err != nil {
		panic(err)
	}
}
