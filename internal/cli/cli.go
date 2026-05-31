package cli

import (
	"expense-tracker/internal/expense"
	"flag"
	"fmt"
	"sort"
	"strings"
	"time"
)

func Run(svc expense.Service, args []string) error {
	if len(args) == 0 {
		printHelp()
		return fmt.Errorf("no command provided")
	}

	if args[0] == "help" {
		printHelp()
		return nil
	}

	fs := flag.NewFlagSet("expense-tracker", flag.ContinueOnError)
	desc := fs.String("description", "", "Expense description")
	amount := fs.Float64("amount", 0, "Expense amount")
	cat := fs.String("category", string(expense.Unknown), "Expense category")
	date := fs.String("date", expense.Date(time.Now()).String(), "Date of expense (YYYY-MM-DD)")
	id := fs.Int("id", 0, "Expense ID")
	month := fs.Int("month", int(time.Now().Month()), "Month (1-12)")
	year := fs.Int("year", time.Now().Year(), "Year")

	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	parseDate := func() (expense.Date, error) {
		return expense.ParseDate(*date)
	}

	switch args[0] {
	case "add":
		expenseDate, err := parseDate()
		if err != nil {
			return err
		}
		if err := svc.AddExpense(float32(*amount), *desc, expense.Category(*cat), expenseDate); err != nil {
			return err
		}
		fmt.Println("Expense added successfully.")
		return nil
	case "delete":
		if *id == 0 {
			return fmt.Errorf("--id is required for delete")
		}
		if err := svc.DeleteExpense(*id); err != nil {
			return err
		}
		fmt.Println("Expense deleted successfully.")
		return nil
	case "update":
		if *id == 0 {
			return fmt.Errorf("--id is required for update")
		}
		expenseDate, err := parseDate()
		if err != nil {
			return err
		}
		if err := svc.UpdateExpense(*id, float32(*amount), *desc, expense.Category(*cat), expenseDate); err != nil {
			return err
		}
		fmt.Println("Expense updated successfully.")
		return nil
	case "list":
		expenses, err := svc.ListAll()
		if err != nil {
			return err
		}
		return svc.PrintExpenses(expenses)
	case "summary":
		monthSet := false
		fs.Visit(func(f *flag.Flag) {
			if f.Name == "month" {
				monthSet = true
			}
		})
		var expenses []expense.Expense
		var err error
		if monthSet {
			expenses, err = svc.ListByMonth(time.Month(*month), time.Now().Year())
			if err != nil {
				return err
			}
			fmt.Printf("Total expenses for %s: €%.2f\n", time.Month(*month), svc.TotalExpense(expenses))
		} else {
			expenses, err = svc.ListAll()
			if err != nil {
				return err
			}
			fmt.Printf("Total expenses: €%.2f\n", svc.TotalExpense(expenses))
		}
		return nil
	case "list-category":
		expenses, err := svc.ListByCategory(expense.Category(*cat))
		if err != nil {
			return err
		}
		return svc.PrintExpenses(expenses)
	case "list-month":
		expenses, err := svc.ListByMonth(time.Month(*month), *year)
		if err != nil {
			return err
		}
		return svc.PrintExpenses(expenses)
	case "list-year":
		expenses, err := svc.ListByYear(*year)
		if err != nil {
			return err
		}
		return svc.PrintExpenses(expenses)
	default:
		return fmt.Errorf("unknown command %q. Run 'expense-tracker help' for usage", args[0])
	}
}

func printHelp() {
	cats := expense.ValidCategories()
	names := make([]string, len(cats))
	for i, c := range cats {
		names[i] = string(c)
	}
	sort.Strings(names)

	fmt.Println("Usage: expense-tracker <command> [flags]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  add           Add a new expense          (--description, --amount, --category, --date)")
	fmt.Println("  update        Update an expense          (--id, --description, --amount, --category, --date)")
	fmt.Println("  delete        Delete an expense          (--id)")
	fmt.Println("  list          List all expenses")
	fmt.Println("  list-category List expenses by category  (--category)")
	fmt.Println("  list-month    List expenses by month     (--month, --year)")
	fmt.Println("  list-year     List expenses by year      (--year)")
	fmt.Println("  summary       Total of all expenses")
	fmt.Println("  summary       Total for a month          (--month)")
	fmt.Println("  export        Export expenses to CSV")
	fmt.Println("  help          Show this help message")
	fmt.Println()
	fmt.Println("Valid categories:")
	fmt.Println(" ", strings.Join(names, ", "))
}
