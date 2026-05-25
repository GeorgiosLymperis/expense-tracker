package storage

import (
	"cmp"
	"encoding/json"
	"expense-tracker/internal/expense"
	"fmt"
	"os"
	"slices"
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
	file, err := os.OpenFile(r.file.Name(), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("Error opening file:%v", err)
	}
	defer file.Close()

	byteValue, err := os.ReadFile(file.Name())
	if err != nil {
		return fmt.Errorf("Error reading file:%v", err)
	}

	var expenses []expenseRecord
	json.Unmarshal(byteValue, &expenses)

	newRecord := NewExpenseRecord(e)
	expenses = append(expenses, newRecord)

	slices.SortFunc(expenses, func(a, b expenseRecord) int {
		return cmp.Compare(a.Date.String(), b.Date.String())
	})
	

	encoder := json.NewEncoder(file)
	err = encoder.Encode(expenses)
	if err != nil {
		return fmt.Errorf("Problem in encoding")
	}

	return nil
}

func (r JSONRepo) Delete(id int) error {
	file, err := os.OpenFile(r.file.Name(), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("Error opening file:%v", err)
	}
	defer file.Close()

	byteValue, err := os.ReadFile(file.Name())
	if err != nil {
		return fmt.Errorf("Error reading file:%v", err)
	}

	var expenses []expenseRecord
	json.Unmarshal(byteValue, &expenses)

	if id > len(expenses) || id < 1 {
		return fmt.Errorf("Invalid ID. Please provide a valid expense ID [1, %v].", len(expenses))
	}
	expenses = slices.Delete(expenses, id-1, id)

	if err:= file.Truncate(0); err != nil {
		return fmt.Errorf("Error truncating file:%v", err)
	}
	
	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("Error seeking file:%v", err)
	}

	encoder := json.NewEncoder(file)
	
	if err = encoder.Encode(expenses); err != nil {
		return fmt.Errorf("Problem in encoding")
	}

	return nil
}

// func openFile

// func (r JSONRepo) FindByCategory(c expense.Category) ([]expense.Expense, error){

// }
