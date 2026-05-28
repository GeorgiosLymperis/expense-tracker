package cli

import (
	"expense-tracker/internal/expense"
	"os"
	"time"
)

type fakeService struct {
	// recorded call arguments
	addedAmount      float32
	addedDescription string
	addedCategory    expense.Category
	addedDate        expense.Date

	updatedID          int
	updatedAmount      float32
	updatedDescription string

	deletedID int

	listByMonthValue time.Month
	listByMonthYear  int
	listByCat        expense.Category
	listByYear       int

	printedExpenses []expense.Expense

	// configurable returns
	listResult []expense.Expense
	listErr    error
	cmdErr     error
}

func (f *fakeService) AddExpense(amount float32, desc string, cat expense.Category, date expense.Date) error {
	f.addedAmount = amount
	f.addedDescription = desc
	f.addedCategory = cat
	f.addedDate = date
	return f.cmdErr
}

func (f *fakeService) UpdateExpense(id int, amount float32, desc string, cat expense.Category, date expense.Date) error {
	f.updatedID = id
	f.updatedAmount = amount
	f.updatedDescription = desc
	return f.cmdErr
}

func (f *fakeService) DeleteExpense(id int) error {
	f.deletedID = id
	return f.cmdErr
}

func (f *fakeService) ListAll() ([]expense.Expense, error) {
	return f.listResult, f.listErr
}

func (f *fakeService) ListByCategory(cat expense.Category) ([]expense.Expense, error) {
	f.listByCat = cat
	return f.listResult, f.listErr
}

func (f *fakeService) ListByMonth(month time.Month, year int) ([]expense.Expense, error) {
	f.listByMonthValue = month
	f.listByMonthYear = year
	return f.listResult, f.listErr
}

func (f *fakeService) ListByYear(year int) ([]expense.Expense, error) {
	f.listByYear = year
	return f.listResult, f.listErr
}

func (f *fakeService) ExportCSV(file *os.File) error {
	return f.cmdErr
}

func (f *fakeService) TotalExpense(expenses []expense.Expense) float32 {
	var sum float32
	for _, e := range expenses {
		sum += e.ExpenseAmount()
	}
	return sum
}

func (f *fakeService) PrintExpenses(expenses []expense.Expense) error {
	f.printedExpenses = expenses
	return nil
}
