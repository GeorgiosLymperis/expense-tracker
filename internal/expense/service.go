package expense

import (
	"os"
	"time"
)

type Service interface {
	// commands
	SaveExpense(amount float32, description string, category Category, date Date) error
	UpdateExpense(id int, amount float32, description string, category Category, date Date) error
	DeleteExpense(id int) error

	// queries
	ListAll() ([]Expense, error)
	ListByCategory(category Category) ([]Expense, error)
	ListByMonth(month time.Month) ([]Expense, error)
	ListByYear(year int) ([]Expense, error)

	// exports
	ExportCSV(file *os.File) error

	// helpers
	TotalExpense(expense *[]Expense) float32
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return service{repo: r}
}

func (s service) SaveExpense(amount float32, description string, category Category, date Date) error {
	return nil
}

func (s service) UpdateExpense(id int, amount float32, description string, category Category, date Date) error {
	return nil
}

func (s service) DeleteExpense(id int) error {
	return nil
}

func (s service) ListAll() ([]Expense, error) {
	return []Expense{}, nil
}

func (s service) ListByCategory(category Category) ([]Expense, error) {
	return []Expense{}, nil
}

func (s service) ListByMonth(month time.Month) ([]Expense, error) {
	return []Expense{}, nil
}

func (s service) ListByYear(year int) ([]Expense, error) {
	return []Expense{}, nil
}

func (s service) ExportCSV(file *os.File) error {
	return nil
}

func (s service) TotalExpense(e *[]Expense) float32 {
	return 0.0
}