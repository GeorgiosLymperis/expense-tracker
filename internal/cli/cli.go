package cli

import (
	"expense-tracker/internal/expense"
	"flag"
	"fmt"
	"time"
)

func Run(svc expense.Service, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no command provided")
	}

	fs := flag.NewFlagSet("expense-tracker", flag.ContinueOnError)
	desc := fs.String("description", "", "Expense description")
	amount := fs.Float64("amount", 0, "Expense amount")
	cat := fs.String("category", string(expense.UNKNOWN), "Expense category")
	date := fs.String("date", expense.Date(time.Now()).String(), "Date of expense")
	id := fs.Int("id", 1, "Expense id")

	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	expenseDate, err := expense.ParseDate(*date)
	if err != nil {
		return err
	}

	switch args[0] {
	case "add":
		return svc.SaveExpense(float32(*amount), *desc, expense.Category(*cat), expenseDate)
	case "delete":
		return svc.DeleteExpense(*id)
	case "update":
		return svc.UpdateExpense(*id, float32(*amount), *desc, expense.Category(*cat), expenseDate)
	case "list-all":
		expenses, err := svc.ListAll()
		if err != nil {
			return err
		}
		return svc.PrintExpenses(&expenses)
	case "list-category":
		expenses, err := svc.ListByCategory(expense.Category(*cat))
		if err != nil {
			return err
		}
		return svc.PrintExpenses(&expenses)
	case "list-month":
		expenses, err := svc.ListByMonth(expenseDate.Month())
		if err != nil {
			return err
		}
		return svc.PrintExpenses(&expenses)
	case "list-year":
		expenses, err := svc.ListByYear(expenseDate.Year())
		if err != nil {
			return err
		}
		return svc.PrintExpenses(&expenses)
	default:
		return fmt.Errorf("Invalid command: %v", args[0])
	}
}
