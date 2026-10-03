package storage

import (
	"encoding/csv"
	"encoding/json"
	"expense-tracker/internal/expense"
	"io"
	"os"
	"testing"
	"time"
)

func compareRecExpense(t *testing.T, rec expenseRecord, e expense.Expense) {
	if rec.Amount != e.ExpenseAmount() {
		t.Errorf("Expense amount should be %v. Got %v", e.ExpenseAmount(), rec.Amount)
	}
	if rec.Date != e.ExpenseDate() {
		t.Errorf("Expense date should be %v. Got %v", e.ExpenseDate(), rec.Date)
	}
	if rec.Description != e.ExpenseDescription() {
		t.Errorf("Expense description should be %v. Got %v", e.ExpenseDescription(), rec.Description)
	}
	if rec.Category != e.ExpenseCategory() {
		t.Errorf("Expense category should be %v. Got %v", e.ExpenseCategory(), rec.Category)
	}
}

func tempJSONFile(t *testing.T) *os.File {
	t.Helper()
	f, err := os.CreateTemp("", "example-*.json")
	if err != nil {
		t.Fatal("could not create temporary file")
	}
	t.Cleanup(func() {
		f.Close()
		os.Remove(f.Name())
	})
	return f
}

func tempCSVFile(t *testing.T) *os.File {
	t.Helper()
	f, err := os.CreateTemp("", "example-*.csv")
	if err != nil {
		t.Fatal("could not create temporary file")
	}
	t.Cleanup(func() {
		f.Close()
		os.Remove(f.Name())
	})
	return f
}

func TestNewJSONRepo(t *testing.T) {
	f := tempJSONFile(t)
	repo := NewJSONRepo(f)
	if repo.file != f {
		t.Errorf("Repo file should be the same with temp")
	}
}

func TestJSONRepoAdd(t *testing.T) {
	f := tempJSONFile(t)
	e, _ := expense.NewExpense(100.50, "EFKA", expense.Business,
		expense.Date(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)))

	repo := NewJSONRepo(f)
	if err := repo.Add(&e); err != nil {
		t.Fatalf("Could not Add expense %v", e)
	}

	byteValue, _ := io.ReadAll(f)
	var expenses []expenseRecord
	json.Unmarshal(byteValue, &expenses)

	compareRecExpense(t, expenses[0], e)
}

func TestJSONRepoAddOlderYearIsLast(t *testing.T) {
	f := tempJSONFile(t)
	old, _ := expense.NewExpense(100.50, "EFKA", expense.Business,
		expense.Date(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)))
	newer, _ := expense.NewExpense(60.50, "Coursera", expense.Education,
		expense.Date(time.Date(2023, 2, 3, 0, 0, 0, 0, time.UTC)))

	repo := NewJSONRepo(f)
	repo.Add(&newer)
	repo.Add(&old)

	byteValue, _ := io.ReadAll(f)
	var expenses []expenseRecord
	json.Unmarshal(byteValue, &expenses)

	compareRecExpense(t, expenses[1], newer)
}

func TestJSONRepoAddOlderMonthIsLast(t *testing.T) {
	f := tempJSONFile(t)
	old, _ := expense.NewExpense(100.50, "EFKA", expense.Business,
		expense.Date(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)))
	newer, _ := expense.NewExpense(60.50, "Coursera", expense.Education,
		expense.Date(time.Date(2023, 2, 3, 0, 0, 0, 0, time.UTC)))

	repo := NewJSONRepo(f)
	repo.Add(&newer)
	repo.Add(&old)

	byteValue, _ := io.ReadAll(f)
	var expenses []expenseRecord
	json.Unmarshal(byteValue, &expenses)

	compareRecExpense(t, expenses[1], newer)
}

func TestJSONRepoAddOlderDayIsLast(t *testing.T) {
	f := tempJSONFile(t)
	old, _ := expense.NewExpense(100.50, "EFKA", expense.Business,
		expense.Date(time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)))
	newer, _ := expense.NewExpense(60.50, "Coursera", expense.Education,
		expense.Date(time.Date(2023, 2, 3, 0, 0, 0, 0, time.UTC)))

	repo := NewJSONRepo(f)
	repo.Add(&newer)
	repo.Add(&old)

	byteValue, _ := io.ReadAll(f)
	var expenses []expenseRecord
	json.Unmarshal(byteValue, &expenses)

	compareRecExpense(t, expenses[1], newer)
}

func TestJSONRepoDelete(t *testing.T) {
	f := tempJSONFile(t)
	old, _ := expense.NewExpense(100.50, "EFKA", expense.Business,
		expense.Date(time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)))
	newer, _ := expense.NewExpense(60.50, "Coursera", expense.Education,
		expense.Date(time.Date(2023, 2, 3, 0, 0, 0, 0, time.UTC)))

	repo := NewJSONRepo(f)
	repo.Add(&newer) // ID 1
	repo.Add(&old)   // ID 2, but sorted first

	if err := repo.Delete(2); err != nil {
		t.Fatalf("Could not delete expense: %v", err)
	}

	byteValue, _ := io.ReadAll(f)
	var expenses []expenseRecord
	json.Unmarshal(byteValue, &expenses)

	compareRecExpense(t, expenses[0], newer)
}

func TestJSONRepoUpdate(t *testing.T) {
	f := tempJSONFile(t)
	e, _ := expense.NewExpense(100.50, "EFKA", expense.Business,
		expense.Date(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)))
	updated, _ := expense.NewExpense(25.00, "NETFLIX", expense.Entertainment,
		expense.Date(time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC)))

	repo := NewJSONRepo(f)
	repo.Add(&e)
	if err := repo.Update(1, &updated); err != nil {
		t.Errorf("Error in updating")
	}

	byteValue, _ := io.ReadAll(f)
	var expenses []expenseRecord
	json.Unmarshal(byteValue, &expenses)

	compareRecExpense(t, expenses[0], updated)
}

func TestJSONRepoListByDate(t *testing.T) {
	f := tempJSONFile(t)
	e1, _ := expense.NewExpense(100.50, "EFKA", expense.Business,
		expense.Date(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)))
	e2, _ := expense.NewExpense(60.50, "Coursera", expense.Education,
		expense.Date(time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)))
	e3, _ := expense.NewExpense(110.50, "BANK", expense.Business,
		expense.Date(time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)))

	repo := NewJSONRepo(f)
	repo.Add(&e2)
	repo.Add(&e1)
	repo.Add(&e3)

	date, err := time.Parse("2006-01-02", "2023-02-01")
	if err != nil {
		t.Fatalf("Could not parse date: %v", err)
	}
	expenses, err := repo.ListByDate(expense.Date(date))
	if err != nil {
		t.Errorf("Error in finding by date: %v", err)
	}

	compareRecExpense(t, newExpenseRecord(&e2), expenses[0])
	compareRecExpense(t, newExpenseRecord(&e3), expenses[1])
}

func TestJSONRepoListByMonth(t *testing.T) {
	f := tempJSONFile(t)
	e1, _ := expense.NewExpense(100.50, "EFKA", expense.Business,
		expense.Date(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)))
	e2, _ := expense.NewExpense(60.50, "Coursera", expense.Education,
		expense.Date(time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)))
	e3, _ := expense.NewExpense(110.50, "BANK", expense.Business,
		expense.Date(time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)))

	repo := NewJSONRepo(f)
	repo.Add(&e2)
	repo.Add(&e1)
	repo.Add(&e3)

	expenses, err := repo.ListByMonth(2, 2025)
	if err != nil {
		t.Errorf("Error in finding by month")
	}

	compareRecExpense(t, newExpenseRecord(&e3), expenses[0])
}

func TestJSONRepoListByYear(t *testing.T) {
	f := tempJSONFile(t)
	e1, _ := expense.NewExpense(100.50, "EFKA", expense.Business,
		expense.Date(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)))
	e2, _ := expense.NewExpense(60.50, "Coursera", expense.Education,
		expense.Date(time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)))
	e3, _ := expense.NewExpense(110.50, "BANK", expense.Business,
		expense.Date(time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)))

	repo := NewJSONRepo(f)
	repo.Add(&e2)
	repo.Add(&e1)
	repo.Add(&e3)

	expenses, err := repo.ListByYear(2023)
	if err != nil {
		t.Errorf("Error in finding by year")
	}

	compareRecExpense(t, newExpenseRecord(&e1), expenses[0])
	compareRecExpense(t, newExpenseRecord(&e2), expenses[1])
}

func TestJSONRepoListByCategory(t *testing.T) {
	f := tempJSONFile(t)
	e1, _ := expense.NewExpense(100.50, "EFKA", expense.Business,
		expense.Date(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)))
	e2, _ := expense.NewExpense(60.50, "Coursera", expense.Education,
		expense.Date(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)))
	e3, _ := expense.NewExpense(110.50, "BANK", expense.Business,
		expense.Date(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)))

	repo := NewJSONRepo(f)
	repo.Add(&e2)
	repo.Add(&e1)
	repo.Add(&e3)

	expenses, err := repo.ListByCategory(expense.Business)
	if err != nil {
		t.Errorf("Error in finding by category: %v", err)
	}

	compareRecExpense(t, newExpenseRecord(&e1), expenses[0])
	compareRecExpense(t, newExpenseRecord(&e3), expenses[1])
}

func TestJSONRepoListByCategoryNotFound(t *testing.T) {
	f := tempJSONFile(t)
	e1, _ := expense.NewExpense(100.50, "EFKA", expense.Business,
		expense.Date(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)))

	repo := NewJSONRepo(f)
	repo.Add(&e1)

	expenses, err := repo.ListByCategory(expense.Entertainment)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if len(expenses) != 0 {
		t.Errorf("Expected no expenses, got %v", len(expenses))
	}
}

func TestJSONRepoExport(t *testing.T) {
	f := tempJSONFile(t)
	e, _ := expense.NewExpense(100.50, "EFKA", expense.Business,
		expense.Date(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)))

	repo := NewJSONRepo(f)
	repo.Add(&e)

	c := tempCSVFile(t)
	if err := repo.ExportCSV(c); err != nil {
		t.Fatalf("Could not export to CSV: %v", err)
	}
	c.Seek(0, 0)
	r := csv.NewReader(c)

	header, err := r.Read()
	if err != nil {
		t.Fatalf("Could not read header: %v", err)
	}
	if header[0] != "Date" {
		t.Errorf("Expected header to be Date, got %v", header[0])
	}
	if header[1] != "Category" {
		t.Errorf("Expected header to be Category, got %v", header[1])
	}
	if header[2] != "Description" {
		t.Errorf("Expected header to be Description, got %v", header[2])
	}
	if header[3] != "Amount" {
		t.Errorf("Expected header to be Amount, got %v", header[3])
	}

	row, err := r.Read()
	if err != nil {
		t.Fatalf("Could not read row: %v", err)
	}
	if row[0] != "2023-01-01" {
		t.Errorf("Expected date 2023-01-01, got %v", row[0])
	}
	if row[1] != "Business" {
		t.Errorf("Expected category Business, got %v", row[1])
	}
	if row[2] != "EFKA" {
		t.Errorf("Expected description EFKA, got %v", row[2])
	}
	if row[3] != "100.5" {
		t.Errorf("Expected amount 100.5, got %v", row[3])
	}
}

func TestListAll(t *testing.T) {
	f := tempJSONFile(t)
	e1, _ := expense.NewExpense(100.50, "EFKA", expense.Business,
		expense.Date(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)))
	e2, _ := expense.NewExpense(60.50, "Coursera", expense.Education,
		expense.Date(time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)))
	e3, _ := expense.NewExpense(110.50, "BANK", expense.Business,
		expense.Date(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)))

	repo := NewJSONRepo(f)
	repo.Add(&e2)
	repo.Add(&e1)
	repo.Add(&e3)

	expenses, err := repo.ListAll()
	if err != nil {
		t.Errorf("Error in listing all")
	}

	compareRecExpense(t, newExpenseRecord(&e3), expenses[0])
	compareRecExpense(t, newExpenseRecord(&e2), expenses[1])
	compareRecExpense(t, newExpenseRecord(&e1), expenses[2])
}

func readRecords(t *testing.T, f *os.File) []expenseRecord {
	t.Helper()
	byteValue, err := os.ReadFile(f.Name())
	if err != nil {
		t.Fatalf("Could not read file: %v", err)
	}
	var expenses []expenseRecord
	if err := json.Unmarshal(byteValue, &expenses); err != nil {
		t.Fatalf("Could not unmarshal: %v", err)
	}
	return expenses
}

func TestJSONRepoIDsSurviveSortingAndDeletes(t *testing.T) {
	f := tempJSONFile(t)
	e1, _ := expense.NewExpense(10, "A", expense.Food,
		expense.Date(time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)))
	e2, _ := expense.NewExpense(20, "B", expense.Food,
		expense.Date(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)))
	e3, _ := expense.NewExpense(30, "C", expense.Food,
		expense.Date(time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)))

	repo := NewJSONRepo(f)
	repo.Add(&e1) // ID 1
	repo.Add(&e2) // ID 2, sorted before A
	if err := repo.Delete(1); err != nil {
		t.Fatalf("Could not delete expense: %v", err)
	}
	repo.Add(&e3) // ID 3

	got := map[string]int{}
	for _, rec := range readRecords(t, f) {
		got[rec.Description] = rec.ID
	}
	want := map[string]int{"B": 2, "C": 3}
	if len(got) != len(want) || got["B"] != want["B"] || got["C"] != want["C"] {
		t.Errorf("IDs should be %v. Got %v", want, got)
	}
}

func TestJSONRepoUpdateKeepsID(t *testing.T) {
	f := tempJSONFile(t)
	e1, _ := expense.NewExpense(10, "A", expense.Food,
		expense.Date(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)))
	e2, _ := expense.NewExpense(20, "B", expense.Food,
		expense.Date(time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)))
	updated, _ := expense.NewExpense(25, "B2", expense.Food,
		expense.Date(time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)))

	repo := NewJSONRepo(f)
	repo.Add(&e1)
	repo.Add(&e2)
	if err := repo.Update(2, &updated); err != nil {
		t.Fatalf("Could not update expense: %v", err)
	}

	expenses, err := repo.ListAll()
	if err != nil {
		t.Fatalf("Could not list expenses: %v", err)
	}
	for _, e := range expenses {
		if e.ExpenseDescription() == "B2" && e.ID() != 2 {
			t.Errorf("Updated expense should keep ID 2. Got %v", e.ID())
		}
	}
}

func TestJSONRepoUnknownID(t *testing.T) {
	f := tempJSONFile(t)
	e, _ := expense.NewExpense(10, "A", expense.Food, expense.Date(time.Now()))

	repo := NewJSONRepo(f)
	repo.Add(&e)
	if err := repo.Delete(5); err == nil {
		t.Error("Expected error deleting unknown ID, got nil")
	}
	if err := repo.Update(5, &e); err == nil {
		t.Error("Expected error updating unknown ID, got nil")
	}
}

func TestJSONRepoMigratesFileWithoutIDs(t *testing.T) {
	f := tempJSONFile(t)
	old := `[{"amount":1,"date":"2023-01-01","description":"A","category":"Food"},` +
		`{"amount":2,"date":"2023-01-02","description":"B","category":"Food"}]`
	if _, err := f.WriteString(old); err != nil {
		t.Fatalf("Could not write file: %v", err)
	}

	repo := NewJSONRepo(f)
	expenses, err := repo.ListAll()
	if err != nil {
		t.Fatalf("Could not list expenses: %v", err)
	}
	if expenses[0].ID() != 1 || expenses[1].ID() != 2 {
		t.Errorf("IDs should be 1, 2. Got %v, %v", expenses[0].ID(), expenses[1].ID())
	}

	recs := readRecords(t, f)
	if recs[0].ID != 1 || recs[1].ID != 2 {
		t.Errorf("Migrated IDs should be saved to the file. Got %v, %v", recs[0].ID, recs[1].ID)
	}
}
