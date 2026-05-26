package expense

import "time"

type Repository interface {
	Save(e *Expense) error
	Delete(id int) error
	Update(id int, e *Expense) error
	FindByDate(date Date) ([]Expense, error)
	FindByMonth(month time.Month) ([]Expense, error)
	FindByYear(year int) ([]Expense, error)
	FindByCategory(category Category) ([]Expense, error)
	Export(filetype string) error
}
