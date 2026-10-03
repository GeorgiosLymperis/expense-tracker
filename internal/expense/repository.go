package expense

import (
	"errors"
	"os"
	"time"
)

// ErrNotFound is returned when no expense has the requested ID.
var ErrNotFound = errors.New("expense not found")

type Repository interface {
	Add(e *Expense) error
	Get(id int) (Expense, error)
	Delete(id int) error
	Update(id int, e *Expense) error
	ListAll() ([]Expense, error)
	ListByDate(date Date) ([]Expense, error)
	ListByMonth(month time.Month, year int) ([]Expense, error)
	ListByYear(year int) ([]Expense, error)
	ListByCategory(category Category) ([]Expense, error)
	ExportCSV(file *os.File) error
}
