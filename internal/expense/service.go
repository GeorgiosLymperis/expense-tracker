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
	for i, e := range expenses {
		fmt.Printf("| ID: %d | %v | %v | %v | %.2f |\n", i+1, e.ExpenseDate(),
			e.ExpenseCategory(), e.ExpenseDescription(), e.ExpenseAmount())
	}

	fmt.Println()
	fmt.Println("Total Amount: ", s.TotalExpense(expenses))
	return nil
}
