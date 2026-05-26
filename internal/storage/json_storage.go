package storage

import (
	"cmp"
	"encoding/json"
	"expense-tracker/internal/expense"
	"fmt"
	"os"
	"slices"
	"time"
)

type JSONRepo struct {
	file *os.File
}

func NewJSONRepo(f *os.File) JSONRepo {
	return JSONRepo{file: f}
}

type expenseRecord struct {
	Amount      float32          `json:"amount"`
	Date        expense.Date     `json:"date"`
	Description string           `json:"description"`
	Category    expense.Category `json:"category"`
}

func NewExpenseRecord(e *expense.Expense) expenseRecord {
	return expenseRecord{
		Amount:      e.ExpenseAmount(),
		Date:        e.ExpenseDate(),
		Description: e.ExpenseDescription(),
		Category:    e.ExpenseCategory(),
	}
}

func (r JSONRepo) Save(e *expense.Expense) error {
	var expenses []expenseRecord

	file, err := OpenJsonRepoFile(r.file.Name(), &expenses)
	if err != nil {
		return err
	}
	defer file.Close()

	newRecord := NewExpenseRecord(e)
	expenses = append(expenses, newRecord)

	slices.SortFunc(expenses, func(a, b expenseRecord) int {
		return cmp.Compare(a.Date.String(), b.Date.String())
	})

	return UpdateJsonRepoFileOk(file, &expenses)
}

func (r JSONRepo) Delete(id int) error {
	var expenses []expenseRecord

	file, err := OpenJsonRepoFile(r.file.Name(), &expenses)
	if err != nil {
		return err
	}
	defer file.Close()

	if id > len(expenses) || id < 1 {
		return fmt.Errorf("Invalid ID. Please provide a valid expense ID [1, %v].", len(expenses))
	}
	expenses = slices.Delete(expenses, id-1, id)

	return UpdateJsonRepoFileOk(file, &expenses)
}

func OpenJsonRepoFile(name string, r *[]expenseRecord) (*os.File, error) {
	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("Error opening file:%v", err)
	}

	byteValue, err := os.ReadFile(file.Name())
	if err != nil {
		return nil, fmt.Errorf("Error reading file:%v", err)
	}

	if len(byteValue) != 0 {
		if err := json.Unmarshal(byteValue, &r); err != nil {
			return nil, fmt.Errorf("Error unmarshalling:%v", err)
		}
	}
	return file, nil
}

func UpdateJsonRepoFileOk(file *os.File, r *[]expenseRecord) error {
	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("Error truncating file:%v", err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("Error seeking file:%v", err)
	}

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(r); err != nil {
		return fmt.Errorf("Problem in encoding")
	}

	return nil
}

func (r JSONRepo) Update(id int, e *expense.Expense) error {
	var expenses []expenseRecord

	file, err := OpenJsonRepoFile(r.file.Name(), &expenses)
	if err != nil {
		return err
	}
	defer file.Close()

	if id > len(expenses) || id < 1 {
		return fmt.Errorf("Invalid ID. Please provide a valid expense ID [1, %v].", len(expenses))
	}

	updateRecord := NewExpenseRecord(e)
	expenses = slices.Delete(expenses, id-1, id)
	expenses = slices.Insert(expenses, id-1, updateRecord)

	return UpdateJsonRepoFileOk(file, &expenses)
}

func (r JSONRepo) FindByCategory(c expense.Category) ([]expense.Expense, error){
	var expenses []expenseRecord

	file, err := OpenJsonRepoFile(r.file.Name(), &expenses)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var expensesInCategory []expense.Expense
	for _, e := range expenses {
		if e.Category == c {
			rec, _ := expense.NewExpense(
				e.Amount,
				e.Description,
				e.Category,
				e.Date)
			expensesInCategory = append(expensesInCategory, rec)
		}
	}
	if len(expensesInCategory) == 0 {
		return expensesInCategory, fmt.Errorf("Category %v not present", c)
	}

	return expensesInCategory, nil
}

func (r JSONRepo) FindByMonth(m time.Month) ([]expense.Expense, error) {
	var expenses []expenseRecord

	file, err := OpenJsonRepoFile(r.file.Name(), &expenses)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var expensesInMonth []expense.Expense
	for _, e := range expenses {
		if e.Date.Month() == m {
			rec, _ := expense.NewExpense(
				e.Amount,
				e.Description,
				e.Category,
				e.Date)
			expensesInMonth = append(expensesInMonth, rec)
		}
	}
	if len(expensesInMonth) == 0 {
		return expensesInMonth, fmt.Errorf("Month %v not present", m)
	}

	return expensesInMonth, nil
}
