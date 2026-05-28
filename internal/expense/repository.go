package expense

import (
	"os"
	"time"
)

type Repository interface {
	Add(e *Expense) error
	Delete(id int) error
	Update(id int, e *Expense) error
	ListAll() ([]Expense, error)
	ListByDate(date Date) ([]Expense, error)
	ListByMonth(month time.Month, year int) ([]Expense, error)
	ListByYear(year int) ([]Expense, error)
	ListByCategory(category Category) ([]Expense, error)
	ExportCSV(file *os.File) error
}
