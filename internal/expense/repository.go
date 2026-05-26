package expense

import (
	"os"
	"time"
)

type Repository interface {
	Save(e *Expense) error
	Delete(id int) error
	Update(id int, e *Expense) error
	ListAll() ([]Expense, error)
	FindByDate(date Date) ([]Expense, error)
	FindByMonth(month time.Month) ([]Expense, error)
	FindByYear(year int) ([]Expense, error)
	FindByCategory(category Category) ([]Expense, error)
	ExportCSV(file *os.File) error
	GetFile() *os.File
}
