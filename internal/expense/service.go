package expense

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Service interface {
	// commands
	AddExpense(amount float32, description string, category Category, date Date) error
	UpdateExpense(id int, amount float32, description string, category Category, date Date) error
	DeleteExpense(id int) error

	// queries
	ListAll() ([]Expense, error)
	ListByCategory(category Category) ([]Expense, error)
	ListByMonth(month time.Month, year int) ([]Expense, error)
	ListByYear(year int) ([]Expense, error)

	// exports
	ExportCSV(file *os.File) error

	// helpers
	TotalExpense(expense []Expense) float32
	PrintExpenses(expenses []Expense) error
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return service{repo: r}
}

func validateCategory(category Category) error {
	if !category.IsValid() {
		cats := ValidCategories()
		names := make([]string, len(cats))
		for i, c := range cats {
			names[i] = string(c)
		}
		return fmt.Errorf("invalid category %q. Valid categories: %s", category, strings.Join(names, ", "))
	}
	return nil
}

func (s service) AddExpense(amount float32, description string, category Category, date Date) error {
	if err := validateCategory(category); err != nil {
		return err
	}
	e, err := NewExpense(amount, description, category, date)
	if err != nil {
		return err
	}
	return s.repo.Add(&e)
}

func (s service) UpdateExpense(id int, amount float32, description string, category Category, date Date) error {
	if err := validateCategory(category); err != nil {
		return err
	}
	e, err := NewExpense(amount, description, category, date)
	if err != nil {
		return err
	}
	return s.repo.Update(id, &e)
}

func (s service) DeleteExpense(id int) error {
	return s.repo.Delete(id)
}

func (s service) ListAll() ([]Expense, error) {
	return s.repo.ListAll()
}

func (s service) ListByCategory(category Category) ([]Expense, error) {
	return s.repo.ListByCategory(category)
}

func (s service) ListByMonth(month time.Month, year int) ([]Expense, error) {
	return s.repo.ListByMonth(month, year)
}

func (s service) ListByYear(year int) ([]Expense, error) {
	return s.repo.ListByYear(year)
}

func (s service) ExportCSV(file *os.File) error {
	return s.repo.ExportCSV(file)
}

func (s service) TotalExpense(e []Expense) float32 {
	var sum float32 = 0
	for _, expense := range e {
		sum += expense.ExpenseAmount()
	}
	return sum
}

func (s service) PrintExpenses(expenses []Expense) error {
	if len(expenses) == 0 {
		fmt.Println("No expenses found.")
		return nil
	}

	const sep = "+-----+------------+----------------------+----------------------+----------+"
	fmt.Println()
	fmt.Println(sep)
	fmt.Println("| ID  | Date       | Category             | Description          |   Amount |")
	fmt.Println(sep)
	for _, e := range expenses {
		fmt.Printf("| %03d | %10s | %-20s | %-20s | €%7.2f |\n",
			e.ID(), e.ExpenseDate(),
			truncate(string(e.ExpenseCategory()), 20),
			truncate(e.ExpenseDescription(), 20),
			e.ExpenseAmount())
	}
	fmt.Println(sep)
	fmt.Printf("| %64s €%7.2f |\n", "Total:", s.TotalExpense(expenses))
	fmt.Println(sep)
	fmt.Println()
	return nil
}

func truncate(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max-3]) + "..."
}
