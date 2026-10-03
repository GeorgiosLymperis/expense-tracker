package storage

import (
	"cmp"
	"encoding/csv"
	"encoding/json"
	"expense-tracker/internal/expense"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"time"
)

type JSONRepo struct {
	file *os.File
}

func NewJSONRepo(f *os.File) *JSONRepo {
	return &JSONRepo{file: f}
}

type expenseRecord struct {
	ID          int              `json:"id"`
	Amount      expense.Amount   `json:"amount"`
	Date        expense.Date     `json:"date"`
	Description string           `json:"description"`
	Category    expense.Category `json:"category"`
}

func newExpenseRecord(e *expense.Expense) expenseRecord {
	return expenseRecord{
		ID:          e.ID(),
		Amount:      e.ExpenseAmount(),
		Date:        e.ExpenseDate(),
		Description: e.ExpenseDescription(),
		Category:    e.ExpenseCategory(),
	}
}

func (rec expenseRecord) toExpense() (expense.Expense, error) {
	e, err := expense.NewExpense(rec.Amount, rec.Description, rec.Category, rec.Date)
	if err != nil {
		return expense.Expense{}, err
	}
	return e.WithID(rec.ID), nil
}

// nextID returns one more than the highest ID in use. Deleting other
// expenses never changes an expense's ID; only the highest ID can be reused
// after it is deleted.
func nextID(expenses []expenseRecord) int {
	maxID := 0
	for _, e := range expenses {
		maxID = max(maxID, e.ID)
	}
	return maxID + 1
}

// indexOfID returns the position of the expense with the given ID, or -1.
func indexOfID(expenses []expenseRecord, id int) int {
	return slices.IndexFunc(expenses, func(e expenseRecord) bool { return e.ID == id })
}

// assignMissingIDs gives an ID to every record without one, in file order.
// It migrates files written before IDs were stored and reports whether
// anything changed.
func assignMissingIDs(expenses []expenseRecord) bool {
	changed := false
	for i := range expenses {
		if expenses[i].ID == 0 {
			expenses[i].ID = nextID(expenses)
			changed = true
		}
	}
	return changed
}

func (r *JSONRepo) Add(e *expense.Expense) error {
	expenses, err := r.load()
	if err != nil {
		return err
	}

	newRecord := newExpenseRecord(e)
	newRecord.ID = nextID(expenses)
	expenses = append(expenses, newRecord)

	slices.SortFunc(expenses, func(a, b expenseRecord) int {
		return cmp.Compare(a.Date.String(), b.Date.String())
	})

	return r.save(expenses)
}

func (r *JSONRepo) Delete(id int) error {
	expenses, err := r.load()
	if err != nil {
		return err
	}

	i := indexOfID(expenses, id)
	if i == -1 {
		return fmt.Errorf("%w: ID %d", expense.ErrNotFound, id)
	}
	expenses = slices.Delete(expenses, i, i+1)

	return r.save(expenses)
}

// load reads all records from the JSON file, giving IDs to any records
// saved before IDs existed.
func (r *JSONRepo) load() ([]expenseRecord, error) {
	byteValue, err := os.ReadFile(r.file.Name())
	if err != nil {
		return nil, fmt.Errorf("Error reading file:%v", err)
	}

	var expenses []expenseRecord
	if len(byteValue) != 0 {
		if err := json.Unmarshal(byteValue, &expenses); err != nil {
			return nil, fmt.Errorf("Error unmarshalling:%v", err)
		}
	}

	if assignMissingIDs(expenses) {
		if err := r.save(expenses); err != nil {
			return nil, err
		}
	}
	return expenses, nil
}

// save replaces the JSON file atomically: it writes a temporary file next to
// it and renames it into place, so a crash mid-write leaves the old file
// intact instead of an empty or half-written one.
func (r *JSONRepo) save(expenses []expenseRecord) error {
	if expenses == nil {
		expenses = []expenseRecord{}
	}
	name := r.file.Name()

	tmp, err := os.CreateTemp(filepath.Dir(name), filepath.Base(name)+".tmp-*")
	if err != nil {
		return fmt.Errorf("Error creating temp file:%v", err)
	}
	defer os.Remove(tmp.Name()) // no-op once the rename has succeeded

	if err := json.NewEncoder(tmp).Encode(expenses); err != nil {
		tmp.Close()
		return fmt.Errorf("Problem in encoding: %v", err)
	}
	if err := tmp.Chmod(0644); err != nil {
		tmp.Close()
		return fmt.Errorf("Error setting file permissions:%v", err)
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return fmt.Errorf("Error syncing file:%v", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("Error closing file:%v", err)
	}
	if err := os.Rename(tmp.Name(), name); err != nil {
		return fmt.Errorf("Error replacing file:%v", err)
	}
	return nil
}

func (r *JSONRepo) Get(id int) (expense.Expense, error) {
	expenses, err := r.load()
	if err != nil {
		return expense.Expense{}, err
	}

	i := indexOfID(expenses, id)
	if i == -1 {
		return expense.Expense{}, fmt.Errorf("%w: ID %d", expense.ErrNotFound, id)
	}
	return expenses[i].toExpense()
}

func (r *JSONRepo) Update(id int, e *expense.Expense) error {
	expenses, err := r.load()
	if err != nil {
		return err
	}

	i := indexOfID(expenses, id)
	if i == -1 {
		return fmt.Errorf("%w: ID %d", expense.ErrNotFound, id)
	}

	expenses[i] = newExpenseRecord(e)
	expenses[i].ID = id
	return r.save(expenses)
}

func (r *JSONRepo) ListByCategory(c expense.Category) ([]expense.Expense, error) {
	expenses, err := r.load()
	if err != nil {
		return nil, err
	}

	var expensesInCategory []expense.Expense
	for _, e := range expenses {
		if e.Category == c {
			rec, err := e.toExpense()
			if err != nil {
				return nil, err
			}
			expensesInCategory = append(expensesInCategory, rec)
		}
	}
	if len(expensesInCategory) == 0 {
		return expensesInCategory, fmt.Errorf("Category %v not present", c)
	}

	return expensesInCategory, nil
}

func (r *JSONRepo) ListByMonth(m time.Month, y int) ([]expense.Expense, error) {
	expenses, err := r.load()
	if err != nil {
		return nil, err
	}

	var expensesInMonth []expense.Expense
	for _, e := range expenses {
		if e.Date.Month() == m && e.Date.Year() == y {
			rec, err := e.toExpense()
			if err != nil {
				return nil, err
			}
			expensesInMonth = append(expensesInMonth, rec)
		}
	}
	if len(expensesInMonth) == 0 {
		return expensesInMonth, fmt.Errorf("Month %v not present", m)
	}

	return expensesInMonth, nil
}

func (r *JSONRepo) ListByYear(y int) ([]expense.Expense, error) {
	expenses, err := r.load()
	if err != nil {
		return nil, err
	}

	var expensesInYear []expense.Expense
	for _, e := range expenses {
		if e.Date.Year() == y {
			rec, err := e.toExpense()
			if err != nil {
				return nil, err
			}
			expensesInYear = append(expensesInYear, rec)
		}
	}
	if len(expensesInYear) == 0 {
		return expensesInYear, fmt.Errorf("Year %v not present", y)
	}

	return expensesInYear, nil
}

func (r *JSONRepo) ListByDate(d expense.Date) ([]expense.Expense, error) {
	expenses, err := r.load()
	if err != nil {
		return nil, err
	}

	var expensesInDate []expense.Expense
	for _, e := range expenses {
		if e.Date == d {
			rec, err := e.toExpense()
			if err != nil {
				return nil, err
			}
			expensesInDate = append(expensesInDate, rec)
		}
	}
	if len(expensesInDate) == 0 {
		return expensesInDate, fmt.Errorf("Date %v not present", d)
	}

	return expensesInDate, nil
}

func (r *JSONRepo) ExportCSV(f *os.File) error {
	expenses, err := r.load()
	if err != nil {
		return err
	}

	w := csv.NewWriter(f)
	if err := w.Write([]string{"Date", "Category", "Description", "Amount"}); err != nil {
		return err
	}

	for _, e := range expenses {
		if err := w.Write([]string{e.Date.String(), string(e.Category),
			e.Description, e.Amount.String()}); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}

func (r *JSONRepo) ListAll() ([]expense.Expense, error) {
	expenses, err := r.load()
	if err != nil {
		return nil, err
	}

	var allExpenses []expense.Expense
	for _, e := range expenses {
		rec, err := e.toExpense()
		if err != nil {
			return nil, err
		}
		allExpenses = append(allExpenses, rec)
	}
	return allExpenses, nil
}
