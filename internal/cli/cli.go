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
	amount := fs.String("amount", "0", "Expense amount, e.g. 12.50")
	cat := fs.String("category", string(expense.Unknown), "Expense category")
	date := fs.String("date", expense.Date(time.Now()).String(), "Date of expense (YYYY-MM-DD)")
	id := fs.Int("id", 0, "Expense ID")
	month := fs.Int("month", int(time.Now().Month()), "Month (1-12)")
	year := fs.Int("year", time.Now().Year(), "Year")
	last := fs.Int("last", 0, "Show only the N most recent expenses (0 = all)")

	if err := fs.Parse(args[1:]); err != nil {
		return err
	}
	if *last < 0 {
		return fmt.Errorf("--last must be 0 or more")
	}

	parseDate := func() (expense.Date, error) {
		return expense.ParseDate(*date)
	}
	parseAmount := func() (expense.Amount, error) {
		return expense.ParseAmount(*amount)
	}
	// printList prints the expenses, keeping only the last N when --last is set.
	// Expenses are stored in date order, so these are the most recent ones.
	printList := func(expenses []expense.Expense) error {
		total := len(expenses)
		if *last > 0 && *last < total {
			expenses = expenses[total-*last:]
		}
		if err := svc.PrintExpenses(expenses); err != nil {
			return err
		}
		if len(expenses) < total {
			fmt.Printf("Showing the last %d of %d expenses.\n", len(expenses), total)
		}
		return nil
	}

	switch args[0] {
	case "add":
		expenseDate, err := parseDate()
		if err != nil {
			return err
		}
		expenseAmount, err := parseAmount()
		if err != nil {
			return err
		}
		if err := svc.AddExpense(expenseAmount, *desc, expense.Category(*cat), expenseDate); err != nil {
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
		// Start from the stored expense and change only the flags given.
		current, err := svc.GetExpense(*id)
		if err != nil {
			return err
		}
		set := setFlags(fs)
		newAmount := current.ExpenseAmount()
		if set["amount"] {
			if newAmount, err = parseAmount(); err != nil {
				return err
			}
		}
		newDesc := current.ExpenseDescription()
		if set["description"] {
			newDesc = *desc
		}
		newCat := current.ExpenseCategory()
		if set["category"] {
			newCat = expense.Category(*cat)
		}
		newDate := current.ExpenseDate()
		if set["date"] {
			if newDate, err = parseDate(); err != nil {
				return err
			}
		}
		if err := svc.UpdateExpense(*id, newAmount, newDesc, newCat, newDate); err != nil {
			return err
		}
		fmt.Println("Expense updated successfully.")
		return nil
	case "list":
		expenses, err := svc.ListAll()
		if err != nil {
			return err
		}
		return printList(expenses)
	case "summary":
		monthSet := setFlags(fs)["month"]
		var expenses []expense.Expense
		var err error
		if monthSet {
			expenses, err = svc.ListByMonth(time.Month(*month), time.Now().Year())
			if err != nil {
				return err
			}
			fmt.Printf("Total expenses for %s: €%s\n", time.Month(*month), svc.TotalExpense(expenses))
		} else {
			expenses, err = svc.ListAll()
			if err != nil {
				return err
			}
			fmt.Printf("Total expenses: €%s\n", svc.TotalExpense(expenses))
		}
		return nil
	case "list-category":
		expenses, err := svc.ListByCategory(expense.Category(*cat))
		if err != nil {
			return err
		}
		return printList(expenses)
	case "list-month":
		expenses, err := svc.ListByMonth(time.Month(*month), *year)
		if err != nil {
			return err
		}
		return printList(expenses)
	case "list-year":
		expenses, err := svc.ListByYear(*year)
		if err != nil {
			return err
		}
		return printList(expenses)
	default:
		return fmt.Errorf("unknown command %q. Run 'expense-tracker help' for usage", args[0])
	}
}

// setFlags returns the names of the flags given on the command line.
func setFlags(fs *flag.FlagSet) map[string]bool {
	set := map[string]bool{}
	fs.Visit(func(f *flag.Flag) {
		set[f.Name] = true
	})
	return set
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
	fmt.Println("  update        Update an expense          (--id, then any of --description, --amount, --category, --date)")
	fmt.Println("  delete        Delete an expense          (--id)")
	fmt.Println("  list          List all expenses          (--last)")
	fmt.Println("  list-category List expenses by category  (--category, --last)")
	fmt.Println("  list-month    List expenses by month     (--month, --year, --last)")
	fmt.Println("  list-year     List expenses by year      (--year, --last)")
	fmt.Println("  summary       Total of all expenses")
	fmt.Println("  summary       Total for a month          (--month)")
	fmt.Println("  export        Export expenses to CSV")
	fmt.Println("  help          Show this help message")
	fmt.Println()
	fmt.Println("Valid categories:")
	fmt.Println(" ", strings.Join(names, ", "))
}
